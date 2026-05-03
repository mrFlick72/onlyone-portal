package otel

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
)

// recordingShutdown returns a ShutdownFunc that flips the supplied flag when
// invoked. Used to assert which providers had their shutdown called by the
// rollback (and which did not).
func recordingShutdown(flag *atomic.Bool) ShutdownFunc {
	return func(_ context.Context) error {
		flag.Store(true)
		return nil
	}
}

func successfulSetup(shut ShutdownFunc) providerSetupFn {
	return func(_ context.Context, _ otelConfig, _ *resource.Resource) (ShutdownFunc, error) {
		return shut, nil
	}
}

func failingSetup(err error) providerSetupFn {
	return func(_ context.Context, _ otelConfig, _ *resource.Resource) (ShutdownFunc, error) {
		return nil, err
	}
}

// neverCalled fails the test if its setup runs. We use it to assert that
// early failures stop the chain — later providers should never be reached.
func neverCalled(t *testing.T) providerSetupFn {
	t.Helper()
	return func(_ context.Context, _ otelConfig, _ *resource.Resource) (ShutdownFunc, error) {
		t.Fatalf("provider setup ran but should have been skipped after a prior failure")
		return nil, nil
	}
}

func enabledCfg() otelConfig {
	return otelConfig{Enabled: true, ServiceName: "test-service"}
}

// The whole point of the partial-setup rollback: when step 2 fails, the
// already-installed step 1 must have its shutdown invoked before Setup
// returns — otherwise its goroutines (BatchSpanProcessor etc.) leak.
func TestSetup_RollsBackPriorProvidersOnFailure(t *testing.T) {
	var traceShutdownCalled atomic.Bool
	var loggerWasCalled atomic.Bool

	providers := []providerSetup{
		{"trace", successfulSetup(recordingShutdown(&traceShutdownCalled))},
		{"metrics", failingSetup(errors.New("forced metrics failure"))},
		{"logs", func(_ context.Context, _ otelConfig, _ *resource.Resource) (ShutdownFunc, error) {
			loggerWasCalled.Store(true)
			return noopShutdown, nil
		}},
	}

	shutdown, err := setupWithProviders(context.Background(), enabledCfg(), providers)

	require.Error(t, err)
	assert.Nil(t, shutdown, "Setup should not return a shutdown when it failed")
	assert.True(t, traceShutdownCalled.Load(),
		"trace's shutdown must be invoked by the rollback so its goroutines stop")
	assert.False(t, loggerWasCalled.Load(),
		"logs setup must not run after metrics failed")
}

// On full success: all three setups ran, NONE of their shutdowns ran (rollback
// is skipped), and the returned shutdown function still works when the caller
// invokes it on graceful exit.
func TestSetup_NoRollbackOnSuccess(t *testing.T) {
	var traceShut, metricsShut, logsShut atomic.Bool

	providers := []providerSetup{
		{"trace", successfulSetup(recordingShutdown(&traceShut))},
		{"metrics", successfulSetup(recordingShutdown(&metricsShut))},
		{"logs", successfulSetup(recordingShutdown(&logsShut))},
	}

	shutdown, err := setupWithProviders(context.Background(), enabledCfg(), providers)

	require.NoError(t, err)
	require.NotNil(t, shutdown)
	assert.False(t, traceShut.Load(), "rollback must not run on success")
	assert.False(t, metricsShut.Load(), "rollback must not run on success")
	assert.False(t, logsShut.Load(), "rollback must not run on success")

	// The caller's normal Shutdown should still flush every provider.
	require.NoError(t, shutdown(context.Background()))
	assert.True(t, traceShut.Load())
	assert.True(t, metricsShut.Load())
	assert.True(t, logsShut.Load())
}

// First-step failure: nothing was installed, so rollback has nothing to do
// and the chain stops immediately. Pins the "empty rollback is a no-op"
// behaviour so a future change can't introduce an NPE here.
func TestSetup_FirstStepFails_NoRollbackNeeded(t *testing.T) {
	providers := []providerSetup{
		{"trace", failingSetup(errors.New("forced trace failure"))},
		{"metrics", neverCalled(t)},
		{"logs", neverCalled(t)},
	}

	shutdown, err := setupWithProviders(context.Background(), enabledCfg(), providers)

	require.Error(t, err)
	assert.Nil(t, shutdown)
}

// Disabled config short-circuits before any provider runs. The TestMain in
// tracing_test.go points at application.yml which has otel.enabled: false,
// so this also documents what the production short-circuit looks like.
func TestSetup_Disabled_InstallsNothing(t *testing.T) {
	providers := []providerSetup{
		{"trace", neverCalled(t)},
		{"metrics", neverCalled(t)},
		{"logs", neverCalled(t)},
	}

	shutdown, err := setupWithProviders(context.Background(), otelConfig{Enabled: false}, providers)

	require.NoError(t, err)
	require.NotNil(t, shutdown, "should return a usable noop shutdown, never nil")
	assert.NoError(t, shutdown(context.Background()))
}
