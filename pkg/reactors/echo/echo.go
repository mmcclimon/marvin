package echo

import (
	"context"
	"fmt"
	"log"

	"github.com/mmcclimon/marvin/pkg/marvin"
)

type Echo struct{}

func Assemble(cfg any) (marvin.Reactor, error) {
	return &Echo{}, nil
}

func (r *Echo) Run(ctx context.Context, eventCh <-chan marvin.Event, errCh chan<- error) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("shutting down echo reactor")
			return nil
		case event := <-eventCh:
			fmt.Printf("echo: >>> %s <<<\n", event.Text)
		}
	}
}
