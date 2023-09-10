package marvin

import "context"

type BusAssembler func(arbitraryConfig) (Bus, error)

type Bus interface {
	Run(context.Context, chan<- Event, chan<- error) error
}

func (m *Marvin) wrapBusFunc(
	ctx context.Context,
	base func(context.Context, chan<- Event, chan<- error) error,
) func() error {
	return func() error {
		return base(ctx, m.events, m.errs)
	}
}
