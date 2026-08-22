● Nothing is ever pushed into done — that's the key. The signal is closing the channel, not sending a value.

Look at Close():                                                                                                                                                                                                                                                                                         
func (bus *InternalEventBus) Close() {                                                                                                                                                                                                                                                                   
bus.closeOnce.Do(func() { close(bus.done) })   // close, not send                                                                                                                                                                                                                                  
}

There is no bus.done <- something anywhere. The mechanism is Go's "signal by close" idiom:

- Receiving from an open channel with nothing to give blocks. So while running, case <-action.EventBus.Done(): never fires, and the select keeps servicing Events().
- Receiving from a closed channel never blocks — it returns immediately (the zero value, with ok == false). So the instant close(bus.done) runs, <-action.EventBus.Done() becomes ready and that select case fires.

That's why the type is chan struct{}: struct{} is a zero-size value — you're not carrying data, you only care about the open→closed transition. The value received is meaningless; the closing is the message.

Two properties that make close (not send) the right choice here:
1. Broadcast — closing unblocks every receiver at once. A single send would only wake one. So if there were several Listen goroutines, one Close() stops them all.
2. Idempotent-safe via closeOnce — closing an already-closed channel panics, so sync.Once guarantees Close() can be called more than once safely (e.g. a double shutdown).

So the flow is: main defers Close → close(bus.done) → the closed channel makes <-Done() return instantly → Listen drains and exits. No message is ever enqueued into done; the closed state is the signal.