package echo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/mmcclimon/marvin"
)

type Echo struct {
	name marvin.ReactorName
	config
}

type config struct {
	ShouldUpper bool `mapstructure:"upper"`
}

func Assemble(name marvin.ReactorName, rawConfig map[string]any) (marvin.Reactor, error) {
	var cfg config
	if err := mapstructure.Decode(rawConfig, &cfg); err != nil {
		return nil, fmt.Errorf("bad config for %s reactor: %w", name, err)
	}

	return &Echo{
		name:   name,
		config: cfg,
	}, nil
}

func (r *Echo) Run(ctx context.Context, eventCh <-chan marvin.Event, errCh chan<- error) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("shutting down echo reactor")
			return nil
		case event := <-eventCh:
			text := event.Text
			if r.ShouldUpper {
				text = strings.ToUpper(text)
			}

			fmt.Printf("echo: >>> %s <<<\n", text)
		}
	}
}
