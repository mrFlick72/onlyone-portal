package expense

import (
	"context"
	"sync"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

// InternalEvent carries a budget expense whose stored tags must be durably
// reclassified to the UNKNOWN sentinel because a read resolved one of its tag
// keys to a tag that no longer exists in tag-api. It is produced on the
// repository read path and consumed by UpdateBudgetExpense.Listen, which
// persists the fix and re-emits the UPDATE event so analytics is reindexed.
type InternalEvent struct {
	Payload *BudgetExpense
	Ctx     context.Context
}

// reclassificationBusCapacity bounds how many pending reclassifications the bus
// buffers. When it is full Publish drops the event rather than blocking the
// read that produced it — a dropped reclassification is simply re-detected the
// next time the record is read, so best-effort delivery is acceptable.
const reclassificationBusCapacity = 256

// InternalEventBus is a bounded, non-blocking, in-process bus connecting the
// repository read path to the reclassification listener. It owns its channel:
// Publish never blocks the caller and Close stops the listener.
type InternalEventBus struct {
	events    chan InternalEvent
	done      chan struct{}
	closeOnce sync.Once
	logger    *logging.Logger
}

func NewEventBus() *InternalEventBus {
	return &InternalEventBus{
		events: make(chan InternalEvent, reclassificationBusCapacity),
		done:   make(chan struct{}),
		logger: logging.GetLoggerInstanceForComponentByTypeName("InternalEventBus"),
	}
}

// Publish enqueues an event without ever blocking the caller. The event is
// dropped (and logged) when the buffer is full, so a read is never stalled by a
// slow or stopped listener.
func (bus *InternalEventBus) Publish(event InternalEvent) {
	select {
	case bus.events <- event:
	default:
		bus.logger.LogErrorfFor("reclassification bus full; dropping event for budget expense %s", event.Payload.Id)
	}
}

// Events is the receive side consumed by the listener.
func (bus *InternalEventBus) Events() <-chan InternalEvent {
	return bus.events
}

// Done is closed by Close to signal the listener to stop.
func (bus *InternalEventBus) Done() <-chan struct{} {
	return bus.done
}

// Close signals the listener to stop after draining any still-buffered events
// (see UpdateBudgetExpense.Listen / drainAndStop). The events channel is
// intentionally left open so a read racing shutdown cannot panic on a send to a
// closed channel. Safe to call more than once.
func (bus *InternalEventBus) Close() {
	bus.closeOnce.Do(func() { close(bus.done) })
}

type BudgetExpenseEventPublisher interface {
	CreateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	UpdateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	DeleteBudgetExpense(ctx context.Context, expense BudgetExpense) error
}
