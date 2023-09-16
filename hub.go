package marvin

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Hub struct {
	buses    map[BusName]Bus
	reactors map[ReactorName]Reactor
	events   chan Event
	replies  chan Reply
	errs     chan error

	reactorChs []chan Event
}

func New() *Hub {
	return &Hub{
		events:   make(chan Event),
		errs:     make(chan error),
		reactors: make(map[ReactorName]Reactor),
		buses:    make(map[BusName]Bus),
		replies:  make(chan Reply),
	}
}

func (h *Hub) Run() error {
	// Alright, so we're gonna set up a context here, and then an error group,
	// which will run all the channels and reactors, so we'll shut down if any
	// of them error out.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	h.startComponents(ctx, eg)
	go h.sigChan(ctx, cancel)
	go h.ioLoop(ctx)

	return eg.Wait()
}

func (h *Hub) startComponents(ctx context.Context, eg *errgroup.Group) {
	for name, bus := range h.buses {
		slog.Info("starting bus", "name", name)
		eg.Go(h.wrapBusFunc(ctx, bus.Run))
	}

	for name, reactor := range h.reactors {
		slog.Info("starting reactor", "name", name)

		ch := make(chan Event)
		h.reactorChs = append(h.reactorChs, ch)

		eg.Go(h.wrapReactorFunc(ctx, reactor.Run, ch))
	}
}

func (h *Hub) sigChan(ctx context.Context, cancel context.CancelFunc) {
	// We're also going to set up a signal channel, so we can shut down on
	// SIGINT or SIGKILL.
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		slog.Info("shutting after catching signal", "signal", sig)
		cancel()
	case <-ctx.Done():
		// just exit
	}
}

func (h *Hub) ioLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case err := <-h.errs:
			slog.Debug("caught non-fatal error", "err", err)

		case event := <-h.events:
			slog.LogAttrs(ctx, slog.LevelDebug,
				"dispatching event",
				slog.Uint64("id", event.ID()),
				slog.String("text", event.Text),
			)

			for _, ch := range h.reactorChs {
				ch <- event
			}
		case reply := <-h.replies:
			reply.Event.Reply(reply.Text)
		}

	}
}
