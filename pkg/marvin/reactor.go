package marvin

import "context"

type ReactorAssembler func(cfg any) (Reactor, error)

type Reactor interface {
	Run(context.Context) error
}

func wrapReactorFunc(ctx context.Context, base func(context.Context) error) func() error {
	return func() error {
		return base(ctx)
	}
}
