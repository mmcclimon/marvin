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

func FromFile(path string) *Marvin {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)

	if err != nil {
		return &Marvin{err: err}
	}

	return cfg.Assemble()
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

	// We're also going to set up a signal channel, so we can shut down on
	// SIGINT or SIGKILL.
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		select {
		case sig := <-sigChan:
			log.Printf("caught %s, shutting down", sig)
			cancel()
		case <-ctx.Done():
			// just exit
		}
	}()

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

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case err := <-m.errs:
				log.Printf("oh no, an error: %s", err)
				cancel()
				return

			case event := <-m.events:
				log.Printf("dispatching event: %+v", event)
				for _, ch := range m.reactorChs {
					ch <- event
				}
			}
		}
	}()

	return eg.Wait()
}
