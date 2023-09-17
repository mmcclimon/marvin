package marvin

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	eventTimeout = 250 * time.Millisecond
)

type Event struct {
	Text      string
	SourceBus BusName
	Address   any
	id        uint64
	watchdog  *time.Timer

	// look, this is super weird, but I just want a done channel
	ctx    context.Context
	cancel context.CancelFunc
}

type Reply struct {
	Bus     BusName
	Address any
	Text    string
}

func NewEvent(source Bus) Event {
	ctx, cancel := context.WithCancel(context.Background())
	evt := Event{
		id:        nextID(),
		SourceBus: source.Name(),
		ctx:       ctx,
		cancel:    cancel,
	}

	return evt
}

func (e *Event) setWatchdog(ch chan<- Reply) {
	e.watchdog = time.AfterFunc(eventTimeout, func() {
		ch <- e.Reply("does not compute")
	})
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

func (e *Event) Reply(format string, args ...any) Reply {
	e.cancel()
	return Reply{
		Bus:     e.SourceBus,
		Address: e.Address,
		Text:    fmt.Sprintf(format, args...),
	}
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
