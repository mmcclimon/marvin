package term

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mmcclimon/marvin/pkg/marvin"
)

type Term struct{}

func Assemble(cfg any) (marvin.Bus, error) {
	return &Term{}, nil
}

func (b *Term) Run(ctx context.Context, eventCh chan<- marvin.Event, errCh chan<- error) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down term bus")
			return nil

		default:
			fmt.Print("> ")
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			switch {
			case errors.Is(err, io.EOF):
				log.Println("caught EOF, shutting down bus")
				return marvin.ErrShuttingDown

			case err != nil:
				errCh <- err

			case text == "error":
				errCh <- errors.New("induced error")

			default:
				eventCh <- marvin.Event{Text: text}
			}
		}
	}
}
