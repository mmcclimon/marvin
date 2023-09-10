package marvin

import "context"

type BusAssembler func(cfg any) (Bus, error)

type Bus interface {
	Run(context.Context) error
}

func wrapBusFunc(ctx context.Context, base func(context.Context) error) func() error {
	return func() error {
		return base(ctx)
	}
}
