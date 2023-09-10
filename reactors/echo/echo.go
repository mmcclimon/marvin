package echo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mmcclimon/marvin"
)

type Echo struct {
	name        marvin.ReactorName
	shouldUpper bool
}

func Assemble(name marvin.ReactorName, cfg map[string]any) (marvin.Reactor, error) {
	echo := Echo{name: name}
	upper, ok := cfg["upper"]
	if ok {
		switch val := upper.(type) {
		case bool:
			echo.shouldUpper = val
		default:
			return nil, fmt.Errorf("bad 'upper' key for echo reactor: %v", upper)
		}
	}

	return &echo, nil
}

func (r *Echo) Run(ctx context.Context, eventCh <-chan marvin.Event, errCh chan<- error) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("shutting down echo reactor")
			return nil
		case event := <-eventCh:
			text := event.Text
			if r.shouldUpper {
				text = strings.ToUpper(text)
			}

			fmt.Printf("echo: >>> %s <<<\n", text)
		}
	}
}
