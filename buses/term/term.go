package term

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/mmcclimon/marvin"
)

type Term struct {
	name marvin.BusName
}

func Assemble(name marvin.BusName, cfg map[string]any) (marvin.Bus, error) {
	return &Term{name}, nil
}

func (b *Term) Run(ctx context.Context, comm marvin.BusBundle) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down term bus")
			return nil

		case reply := <-comm.Replies:
			b.SendMessage(ctx, nil, reply.Text)
			fmt.Print("> ")

		default:
			fmt.Print("> ")
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			switch {
			case errors.Is(err, io.EOF):
				slog.Info("caught EOF, shutting down bus")
				return marvin.ErrShuttingDown

			case err != nil:
				comm.Errors <- err

			case text == "error":
				comm.Errors <- errors.New("induced error")

			default:
				event := b.eventFromText(text)
				comm.Events <- event
				<-event.Done()
			}
		}
	}
}

func (b *Term) eventFromText(text string) marvin.Event {
	ev := marvin.NewEvent(b)
	ev.Text = text
	return ev
}

func (b *Term) Name() marvin.BusName { return b.name }

func (b *Term) SendMessage(_ context.Context, _ any, text string) {
	fmt.Printf("| %s\n", text)
}
