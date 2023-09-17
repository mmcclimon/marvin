package eject

import (
	"context"
	"log/slog"
	"regexp"
	"time"

	"github.com/mmcclimon/marvin"
)

var prefix = regexp.MustCompile(`(?i)^eject warp core\s*`)

type Eject struct {
	name marvin.ReactorName
}

func Assemble(name marvin.ReactorName, rawConfig map[string]any) (marvin.Reactor, error) {
	return &Eject{name}, nil
}

func (r *Eject) Run(ctx context.Context, comm marvin.ReactorBundle) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down eject reactor")
			return nil

		case event := <-comm.Events:
			if !prefix.MatchString(event.Text) {
				continue
			}

			event.MarkHandled()
			comm.Replies <- event.Reply("so long!")

			time.Sleep(2 * time.Second)
			return marvin.ErrShuttingDown
		}
	}
}
