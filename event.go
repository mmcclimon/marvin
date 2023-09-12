package marvin

import (
	"sync"
	"time"
)

const (
	eventTimeout = 250 * time.Millisecond
)

var nextEventID struct {
	mu sync.Mutex
	id uint64
}

type Event struct {
	Text      string
	SourceBus Bus
	id        uint64
	watchdog  *time.Timer
	done      chan struct{}
}

func NewEvent(source Bus) Event {
	nextEventID.mu.Lock()
	id := nextEventID.id
	nextEventID.id++
	nextEventID.mu.Unlock()

	evt := Event{
		id:        id,
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
	close(e.done)
}
