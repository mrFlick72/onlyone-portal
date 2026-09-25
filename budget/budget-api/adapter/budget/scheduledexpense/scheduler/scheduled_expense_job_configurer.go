package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

// Job is what the configurer runs on each tick —
// scheduledexpense.ScheduledExpenseJob in production, a fake in tests.
type Job interface {
	Execute(ctx context.Context) error
}

// ScheduledExpenseJobConfigurer runs the ScheduledExpenseJob in-process on a
// gocron v2 scheduler, as a golang-web-framework WebServerConfigurer
// (registered with RegisterConfigurer): Configure starts it, Dispose shuts it
// down within the provisioner's shutdown budget. See ADR 0005.
//
// The job runs once immediately on start — a restart doubles as a manual
// trigger — then every interval. Running more often than daily is safe:
// LastEvaluatedDate makes each day evaluated exactly once. Singleton mode
// keeps a slow run from overlapping the next. Single replica only: there is
// no distributed lock (ADR 0005).
type ScheduledExpenseJobConfigurer struct {
	job       Job
	interval  time.Duration
	scheduler gocron.Scheduler
	logger    *logging.Logger
}

func NewScheduledExpenseJobConfigurer(job Job, interval time.Duration) *ScheduledExpenseJobConfigurer {
	return &ScheduledExpenseJobConfigurer{
		job:      job,
		interval: interval,
		logger:   logging.GetLoggerInstanceForComponentByType(&ScheduledExpenseJobConfigurer{}),
	}
}

func (c *ScheduledExpenseJobConfigurer) Name() string {
	return "scheduled-expense-job"
}

func (c *ScheduledExpenseJobConfigurer) Configure() error {
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		return fmt.Errorf("create scheduler: %w", err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(c.interval),
		// gocron passes a ctx it cancels on shutdown; the engine checks it
		// between definitions and between days.
		gocron.NewTask(func(ctx context.Context) {
			if err := c.job.Execute(ctx); err != nil {
				c.logger.LogErrorfFor("scheduled expense job run failed: %v", err)
			}
		}),
		gocron.WithName(c.Name()),
		gocron.WithStartAt(gocron.WithStartImmediately()),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		_ = scheduler.Shutdown()
		return fmt.Errorf("schedule scheduled expense job: %w", err)
	}

	scheduler.Start()
	c.scheduler = scheduler
	c.logger.LogInfofFor("scheduled expense job started, every %s", c.interval)
	return nil
}

func (c *ScheduledExpenseJobConfigurer) Dispose(ctx context.Context) error {
	if c.scheduler == nil {
		return nil
	}
	return c.scheduler.ShutdownWithContext(ctx)
}
