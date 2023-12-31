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

func (d *Discord) Run(ctx context.Context, comm marvin.BusBundle) error {
	if err := d.discord.Connect(ctx); err != nil {
		return err
	}

	msgCh := make(chan discord.Message)
	go d.discord.Run(ctx, msgCh, comm.Errors)

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("shutting down discord channel")
			return marvin.ErrShuttingDown

		case <-d.discord.Fatal():
			err := d.discord.Err
			d.logger.Warn("fatal err from discord", "err", err)
			return err

		case err := <-d.discord.Errors():
			d.logger.Warn("caught error from discord client", "err", err)

		case reply := <-comm.Replies:
			d.SendMessage(ctx, reply.Address, reply.Text)

		case msg := <-msgCh:
			if msg.Author.IsBot {
				continue
			}

			evt := d.eventFromMessage(msg)
			comm.Events <- evt
		}
	}
}

func (d *Discord) eventFromMessage(msg discord.Message) marvin.Event {
	ev := marvin.NewEvent(d)
	ev.Text = d.discord.DecodeFormatting(msg)
	ev.Address = msg.ChannelID
	return ev
}

func (d *Discord) Name() marvin.BusName { return d.name }

func (d *Discord) SendMessage(ctx context.Context, address any, text string) {
	url := discord.URLFor("/channels/%s/messages", address)
	res, err := d.discord.Post(ctx, url, map[string]string{"content": text})

	if err != nil {
		d.logger.Warn("bad message post", "err", err)
		return
	}

	d.discord.CheckAPIResponse(res)
}
