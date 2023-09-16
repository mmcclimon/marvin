package discord

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/mmcclimon/marvin"
	"github.com/mmcclimon/marvin/buses/discord/internal/discord"
)

type Discord struct {
	name    marvin.BusName
	discord *discord.Client
	logger  *slog.Logger
	raw     chan []byte
}

type config struct {
	Token string `mapstructure:"api_token"`
}

func Assemble(name marvin.BusName, rawConfig map[string]any) (marvin.Bus, error) {
	var cfg config
	if err := mapstructure.Decode(rawConfig, &cfg); err != nil {
		return nil, fmt.Errorf("bad config for %s bus: %w", name, err)
	}

	logger := slog.Default().With("bus", name)

	return &Discord{
		name:    name,
		raw:     make(chan []byte),
		discord: discord.NewClient(logger, cfg.Token),
		logger:  logger,
	}, nil
}

func (d *Discord) Run(ctx context.Context, eventCh chan<- marvin.Event, errCh chan<- error) error {
	// TODO remove this constant
	url := "wss://gateway.discord.gg"
	if err := d.discord.Connect(ctx, url); err != nil {
		return err
	}

	evtCh := make(chan discord.Message)
	go d.discord.Run(ctx, evtCh, errCh)

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("shutting down discord channel")
			return marvin.ErrShuttingDown
		case <-d.discord.C:
			err := d.discord.Err()
			d.logger.Info("caught fatal err from discord", "err", err)
			return err
		case msg := <-evtCh:
			if msg.Author.IsBot {
				continue
			}

			evt := d.eventFromMessage(msg)
			eventCh <- evt
		}
	}
}

func (d *Discord) eventFromMessage(msg discord.Message) marvin.Event {
	ev := marvin.NewEvent(d)
	ev.Text = msg.Content
	return ev
}

func (d *Discord) SendMessage(ctx context.Context, text string) {
	url := discord.URLFor("/channels/%s/messages", "1152064873161834579") // lol
	res, err := d.discord.Post(ctx, url, map[string]string{"content": text})

	if err != nil {
		d.logger.Warn("bad message post", "err", err)
		return
	}

	defer res.Body.Close()
	_, _ = io.Copy(os.Stdout, res.Body)
}
