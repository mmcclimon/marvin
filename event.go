package marvin

import (
	"context"
	"sync"
	"time"
)

const (
	eventTimeout = 250 * time.Millisecond
)

type Event struct {
	Text      string
	SourceBus Bus
	Address   any
	id        uint64
	watchdog  *time.Timer

	// look, this is super weird, but I just want a done channel
	ctx    context.Context
	cancel context.CancelFunc
}

func NewEvent(source Bus) Event {
	ctx, cancel := context.WithCancel(context.Background())
	evt := Event{
		id:        nextID(),
		SourceBus: source,
		ctx:       ctx,
		cancel:    cancel,
	}

	evt.watchdog = time.AfterFunc(eventTimeout, func() {
		evt.Reply("does not compute")
	})

	return evt
}

func (e *Event) ID() uint64 {
	return e.id
}

func (e *Event) MarkHandled() {
	e.watchdog.Stop()
}

func (e *Event) Done() <-chan struct{} {
	return e.ctx.Done()
}

func (e *Event) Reply(text string) {
	e.SourceBus.SendMessage(context.TODO(), e.Address, text)
	e.cancel()
}

var nextEventID struct {
	mu sync.Mutex
	id uint64
}

func nextID() uint64 {
	nextEventID.mu.Lock()
	defer nextEventID.mu.Unlock()

	nextEventID.id++
	return nextEventID.id
}
