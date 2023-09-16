package marvin

import "context"

type ReactorName string

type ReactorAssembler func(ReactorName, arbitraryConfig) (Reactor, error)

type Reactor interface {
	Run(context.Context, <-chan Event, chan<- error) error
}

func (h *Hub) wrapReactorFunc(
	ctx context.Context,
	base func(context.Context, <-chan Event, chan<- error) error,
	ch <-chan Event,
) func() error {
	return func() error {
		return base(ctx, ch, h.errs)
	}
}
