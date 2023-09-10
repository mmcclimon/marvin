package echo

import (
	"context"
	"log"

	"github.com/mmcclimon/marvin/pkg/marvin"
)

type Echo struct{}

func Assemble(cfg any) (marvin.Reactor, error) {
	return &Echo{}, nil
}

func (r *Echo) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("shutting down echo reactor")
			return nil
		default:
			// do nothing
		}
	}
}
