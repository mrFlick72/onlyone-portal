package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGenerationJob struct {
	runs chan context.Context
}

func (f *fakeGenerationJob) Execute(ctx context.Context) error {
	f.runs <- ctx
	<-ctx.Done()
	return ctx.Err()
}

// The job runs once as soon as the scheduler starts (so a restart doubles as
// a manual trigger), and Dispose cancels the in-flight run's context.
func TestConfigureRunsTheJobImmediatelyAndDisposeCancelsIt(t *testing.T) {
	job := &fakeGenerationJob{runs: make(chan context.Context, 1)}
	uut := NewGocronGenerationConfigurer(job, time.Hour)

	require.NoError(t, uut.Configure())

	var runCtx context.Context
	select {
	case runCtx = <-job.runs:
	case <-time.After(5 * time.Second):
		t.Fatal("expected the job to run immediately on Configure")
	}

	disposeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, uut.Dispose(disposeCtx))

	select {
	case <-runCtx.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("expected Dispose to cancel the running job's context")
	}
}

func TestDisposeWithoutConfigureIsANoOp(t *testing.T) {
	uut := NewGocronGenerationConfigurer(&fakeGenerationJob{runs: make(chan context.Context, 1)}, time.Hour)

	assert.NoError(t, uut.Dispose(context.Background()))
}

func TestName(t *testing.T) {
	uut := NewGocronGenerationConfigurer(&fakeGenerationJob{}, time.Hour)

	assert.Equal(t, "scheduled-expense-generation", uut.Name())
}
