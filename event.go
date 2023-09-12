package marvin

import (
	"sync"
	"time"
)

const (
	eventTimeout = 250 * time.Millisecond
)

type Event struct {
	Text      string
	SourceBus Bus
	id        uint64
	watchdog  *time.Timer
	done      chan struct{}
}

func NewEvent(source Bus) Event {
	evt := Event{
		id:        nextID(),
		SourceBus: source,
		done:      make(chan struct{}),
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
	return e.done
}

func (e *Event) Reply(text string) {
	e.SourceBus.SendMessage(text)
	// lol this panics if you call Reply() more than once to an event
	close(e.done)
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
