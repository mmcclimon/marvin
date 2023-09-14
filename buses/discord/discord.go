package discord

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mitchellh/mapstructure"
	"github.com/mmcclimon/marvin"
	"github.com/mmcclimon/marvin/buses/discord/internal/discord"
)

type Discord struct {
	name    marvin.BusName
	discord *discord.Client
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

	return &Discord{
		name:    name,
		raw:     make(chan []byte),
		discord: discord.NewClient(cfg.Token),
	}, nil
}

func (d *Discord) Run(ctx context.Context, eventCh chan<- marvin.Event, errCh chan<- error) error {
	// TODO remove this constant
	url := "wss://gateway.discord.gg"
	if err := d.discord.Connect(ctx, url); err != nil {
		return err
	}

	evtCh := make(chan discord.GatewayEvent)
	go d.discord.Run(ctx, evtCh, errCh)

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down discord channel")
			return marvin.ErrShuttingDown
		case <-d.discord.C:
			err := d.discord.Err()
			slog.Info("caught fatal err from discord", "err", err)
			return err
		case evt := <-evtCh:
			fmt.Printf("%+v\n", evt)
		}
	}
}

func (d *Discord) SendMessage(text string) {
	fmt.Printf("TODO: send message %s\n", text)
}
