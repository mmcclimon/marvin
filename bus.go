package marvin

import "context"

type BusName string

type BusAssembler func(BusName, arbitraryConfig) (Bus, error)

type Bus interface {
	Run(context.Context, chan<- Event, chan<- error) error
	SendMessage(string)
}

func (m *Marvin) wrapBusFunc(
	ctx context.Context,
	base func(context.Context, chan<- Event, chan<- error) error,
) func() error {
	return func() error {
		return base(ctx, m.events, m.errs)
	}
}
