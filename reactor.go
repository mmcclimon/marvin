package marvin

import "context"

type ReactorName string

type ReactorAssembler func(ReactorName, arbitraryConfig) (Reactor, error)

type ReactorBundle struct {
	Events  <-chan Event
	Replies chan<- Reply
	Errors  chan<- error
}

type Reactor interface {
	Run(context.Context, ReactorBundle) error
}

func (h *Hub) wrapReactorFunc(
	ctx context.Context,
	base func(context.Context, ReactorBundle) error,
	bundle ReactorBundle,
) func() error {
	return func() error { return base(ctx, bundle) }
}
