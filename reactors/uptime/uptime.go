package uptime

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/mmcclimon/marvin"
)

type Uptime struct {
	name  marvin.ReactorName
	start time.Time
}

func Assemble(name marvin.ReactorName, rawConfig map[string]any) (marvin.Reactor, error) {
	return &Uptime{name: name}, nil
}

func (r *Uptime) Run(ctx context.Context, comm marvin.ReactorBundle) error {
	r.start = time.Now()

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down uptime reactor")
			return nil

		case event := <-comm.Events:
			if strings.ToLower(event.Text) != "uptime" {
				continue
			}

			event.MarkHandled()

			uptime := time.Since(r.start)

			trunc := time.Second
			if uptime > time.Hour {
				trunc = time.Minute
			}

			comm.Replies <- event.Reply("Online for %s", uptime.Truncate(trunc))
		}
	}
}
