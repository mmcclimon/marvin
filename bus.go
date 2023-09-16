package marvin

import "context"

type BusName string

type BusAssembler func(BusName, arbitraryConfig) (Bus, error)

type BusBundle struct {
	Events  chan<- Event
	Replies <-chan Reply
	Errors  chan<- error
}

type Bus interface {
	Run(context.Context, BusBundle) error
	SendMessage(ctx context.Context, address any, text string)
}

func (h *Hub) wrapBusFunc(
	ctx context.Context,
	base func(context.Context, BusBundle) error,
) func() error {
	bb := BusBundle{
		Events:  h.events,
		Replies: h.replies,
		Errors:  h.errs,
	}

	return func() error {
		return base(ctx, bb)
	}
}
