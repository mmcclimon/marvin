package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type arbitraryJSON = map[string]any

type HelloData struct {
	HeartbeatInterval int `mapstructure:"heartbeat_interval"`
}

func (c *Client) doHello(ctx context.Context, event GatewayEvent) (*GatewayEvent, error) {
	var data HelloData
	if err := mapstructure.Decode(event.Data, &data); err != nil {
		return nil, fmt.Errorf("bad hello decode: %w", err)
	}

	c.logger.Debug("got hello data", "interval", data.HeartbeatInterval)

	interval := time.Duration(data.HeartbeatInterval) * time.Millisecond
	go c.runHeartbeatLoop(ctx, interval)

	return nil, nil
}

func (c *Client) sendHeartbeat(ctx context.Context, seq *int) {
	outgoing := map[string]any{
		"op": Heartbeat,
		"d":  seq,
	}

	c.logger.Debug("will send heartbeat", "data", outgoing)
	data, _ := json.Marshal(outgoing)

	c.acked = false
	c.write(ctx, data)
}

func (c *Client) ackHeartbeat(ctx context.Context, event GatewayEvent) (*GatewayEvent, error) {
	c.acked = true
	return nil, nil
}

const intents = GuildMessages | GuildMessageReactions | DirectMessages | DirectMessageReactions

func (c *Client) doIdentify(ctx context.Context) {
	outgoing := arbitraryJSON{
		"op": Identify,
		"d": arbitraryJSON{
			"token":   c.token,
			"intents": intents,
			"properties": arbitraryJSON{
				"os":      "macos",
				"browser": "marvin",
				"device":  "marvin",
			},
			"presence": arbitraryJSON{
				"status":     "online",
				"activities": []string{},
				"afk":        false,
			},
		},
	}

	data, _ := json.Marshal(outgoing)

	c.logger.Debug("will identify")
	c.write(ctx, data)
}
