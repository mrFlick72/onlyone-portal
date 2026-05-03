package server

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// callRecorder logs Configure / Dispose calls in the order they happen,
// shared across fakes so tests can assert ordering across configurers.
type callRecorder struct {
	mu    sync.Mutex
	calls []string
}

func (r *callRecorder) record(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, s)
}

func (r *callRecorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.calls))
	copy(out, r.calls)
	return out
}

// fakeConfigurer implements WebServerConfigurer with hooks for tests:
// canned Configure / Dispose errors, and an option to make Dispose block
// until the caller-supplied ctx is done — used to assert deadline semantics.
type fakeConfigurer struct {
	name          string
	rec           *callRecorder
	configureErr  error
	disposeErr    error
	disposeBlocks bool
}

func (f *fakeConfigurer) Name() string { return f.name }

func (f *fakeConfigurer) Configure() error {
	f.rec.record("configure:" + f.name)
	return f.configureErr
}

func (f *fakeConfigurer) Dispose(ctx context.Context) error {
	f.rec.record("dispose:" + f.name)
	if f.disposeBlocks {
		<-ctx.Done()
		return ctx.Err()
	}
	return f.disposeErr
}

// withFakes builds a provisioner pre-populated with the given fakes and a
// non-nil engine, bypassing ConfigureEngine. Lets tests focus on Shutdown
// semantics without booting the real OTel / OAuth2 wiring.
func withFakes(fakes ...WebServerConfigurer) *WebServerProvisioner {
	gin.SetMode(gin.TestMode)
	return &WebServerProvisioner{
		engine:      gin.New(),
		configurers: append([]WebServerConfigurer(nil), fakes...),
	}
}

// Disposal must follow registration order so that, e.g., a configurer that
// depends on a downstream one (OAuth2 holding a ctx derived from OTel) tears
// down predictably.
func TestShutdown_DisposesInRegistrationOrder(t *testing.T) {
	rec := &callRecorder{}
	a := &fakeConfigurer{name: "a", rec: rec}
	b := &fakeConfigurer{name: "b", rec: rec}
	c := &fakeConfigurer{name: "c", rec: rec}
	wsp := withFakes(a, b, c)

	err := wsp.Shutdown(context.Background())
	require.NoError(t, err)

	assert.Equal(t,
		[]string{"dispose:a", "dispose:b", "dispose:c"},
		rec.snapshot(),
	)
}

// Calling Shutdown twice must not double-Dispose. The slice reset inside
// Shutdown gives us this for free; the test pins the contract.
func TestShutdown_IsIdempotent(t *testing.T) {
	rec := &callRecorder{}
	a := &fakeConfigurer{name: "a", rec: rec}
	wsp := withFakes(a)

	require.NoError(t, wsp.Shutdown(context.Background()))
	require.NoError(t, wsp.Shutdown(context.Background()))

	assert.Equal(t, []string{"dispose:a"}, rec.snapshot(),
		"Dispose should fire exactly once even with two Shutdown calls")
}

// After Shutdown the provisioner should be back to its zero state so that a
// fresh ConfigureEngine builds a clean engine, and a retried Shutdown
// iterates an empty slice.
func TestShutdown_ResetsEngineAndConfigurers(t *testing.T) {
	wsp := withFakes(&fakeConfigurer{name: "a", rec: &callRecorder{}})
	require.NotNil(t, wsp.engine)
	require.Len(t, wsp.configurers, 1)

	require.NoError(t, wsp.Shutdown(context.Background()))

	assert.Nil(t, wsp.engine, "engine should be cleared so a retry rebuilds")
	assert.Nil(t, wsp.configurers, "configurer slice should be cleared so a retry doesn't double-dispose")
}

// One Dispose failure must not stop the others, and every error must surface
// to the caller so the process exit code reflects a partial shutdown.
func TestShutdown_AggregatesDisposeErrorsViaErrorsJoin(t *testing.T) {
	rec := &callRecorder{}
	errA := errors.New("a-fail")
	errC := errors.New("c-fail")
	a := &fakeConfigurer{name: "a", rec: rec, disposeErr: errA}
	b := &fakeConfigurer{name: "b", rec: rec} // succeeds
	c := &fakeConfigurer{name: "c", rec: rec, disposeErr: errC}
	wsp := withFakes(a, b, c)

	err := wsp.Shutdown(context.Background())
	require.Error(t, err)

	assert.ErrorIs(t, err, errA, "joined error should wrap a's failure")
	assert.ErrorIs(t, err, errC, "joined error should wrap c's failure")

	assert.Equal(t,
		[]string{"dispose:a", "dispose:b", "dispose:c"},
		rec.snapshot(),
		"all three Dispose calls should run despite errors")
}

// The whole point of Batch 2: a Dispose that would otherwise block forever
// (e.g. OTel exporter unable to reach its collector) must return as soon as
// the caller-supplied ctx expires. If this regresses, processes hang on
// shutdown until SIGKILL.
func TestShutdown_HonorsContextDeadline(t *testing.T) {
	rec := &callRecorder{}
	blocker := &fakeConfigurer{name: "blocker", rec: rec, disposeBlocks: true}
	wsp := withFakes(blocker)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := wsp.Shutdown(ctx)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded,
		"the joined error should wrap context.DeadlineExceeded from the blocked Dispose")
	assert.Less(t, elapsed, 500*time.Millisecond,
		"Shutdown should return shortly after ctx deadline, not block indefinitely")
}

// A blocker followed by a normal configurer: once the deadline expires, the
// remaining configurers still get a chance to Dispose, but they receive an
// already-done ctx. This documents the shared-budget behaviour discussed in
// review: configurers later in the slice get whatever budget the earlier ones
// left behind.
func TestShutdown_RemainingConfigurersStillDisposedAfterDeadline(t *testing.T) {
	rec := &callRecorder{}
	blocker := &fakeConfigurer{name: "blocker", rec: rec, disposeBlocks: true}
	after := &fakeConfigurer{name: "after", rec: rec}
	wsp := withFakes(blocker, after)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	require.Error(t, wsp.Shutdown(ctx))
	assert.Equal(t,
		[]string{"dispose:blocker", "dispose:after"},
		rec.snapshot(),
		"Dispose loop should continue past a deadline-exceeded configurer")
}

// Edge case: zero configurers should be a clean no-op. Important because the
// "shutdown after a clean shutdown" path goes through this branch (slice was
// reset to nil on the first call).
func TestShutdown_NoConfigurersReturnsNil(t *testing.T) {
	wsp := &WebServerProvisioner{}
	assert.NoError(t, wsp.Shutdown(context.Background()))
}

// ConfigureEngine memoises the engine — the second call must return the same
// pointer without re-running configurers (which would double-register every
// middleware on the chain).
func TestConfigureEngine_ReturnsCachedEngineOnSecondCall(t *testing.T) {
	gin.SetMode(gin.TestMode)
	existing := gin.New()
	wsp := &WebServerProvisioner{engine: existing}

	got := wsp.ConfigureEngine()

	assert.Same(t, existing, got,
		"second call should return the cached engine without rebuilding the middleware stack")
}
