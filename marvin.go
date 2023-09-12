package marvin

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"golang.org/x/sync/errgroup"
)

type Marvin struct {
	err      error
	buses    map[BusName]Bus
	reactors map[ReactorName]Reactor
	events   chan Event
	errs     chan error

	reactorChs []chan Event
}

var ErrShuttingDown = errors.New("shutting down")

func FromFile(path string, registry Registry) *Marvin {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)

	if err != nil {
		return &Marvin{err: err}
	}

	return cfg.Assemble(registry)
}

func (m *Marvin) Run() error {
	if m.err != nil {
		return m.err
	}

	// Alright, so we're gonna set up a context here, and then an error group,
	// which will run all the channels and reactors, so we'll shut down if any
	// of them error out.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	m.startComponents(ctx, eg)
	go m.sigChan(ctx, cancel)
	go m.ioLoop(ctx)

	return eg.Wait()
}

func (m *Marvin) sigChan(ctx context.Context, cancel context.CancelFunc) {
	// We're also going to set up a signal channel, so we can shut down on
	// SIGINT or SIGKILL.
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	select {
	case sig := <-sigChan:
		log.Printf("caught %s, shutting down", sig)
		cancel()
	case <-ctx.Done():
		// just exit
	}
}

func (m *Marvin) startComponents(ctx context.Context, eg *errgroup.Group) {
	for name, bus := range m.buses {
		log.Printf("starting bus %s", name)
		eg.Go(m.wrapBusFunc(ctx, bus.Run))
	}

	for name, reactor := range m.reactors {
		log.Printf("starting reactor %s", name)

		ch := make(chan Event)
		m.reactorChs = append(m.reactorChs, ch)

		eg.Go(m.wrapReactorFunc(ctx, reactor.Run, ch))
	}
}

func (m *Marvin) ioLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case err := <-m.errs:
			log.Printf("caught non-fatal error: %s", err)

		case event := <-m.events:
			log.Printf("dispatching event: id=%d, text=%s", event.ID(), event.Text)
			for _, ch := range m.reactorChs {
				ch <- event
			}
		}
	}
}
