package marvin

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"golang.org/x/sync/errgroup"
)

type Marvin struct {
	err      error
	buses    map[string]Bus
	reactors map[string]Reactor
}

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
		sig := <-sigChan
		log.Printf("caught %s, shutting down", sig)
		cancel()
	}()

	for name, bus := range m.buses {
		log.Printf("starting bus %s", name)
		eg.Go(wrapBusFunc(ctx, bus.Run))
	}

	for name, reactor := range m.reactors {
		log.Printf("starting reactor %s", name)
		eg.Go(wrapBusFunc(ctx, reactor.Run))
	}

	return eg.Wait()
}
