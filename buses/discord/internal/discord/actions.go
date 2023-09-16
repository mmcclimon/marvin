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

func (c *Client) doHello(ctx context.Context, event GatewayEvent) error {
	var data HelloData
	if err := mapstructure.Decode(event.Data, &data); err != nil {
		return fmt.Errorf("bad hello decode: %w", err)
	}

	c.logger.Debug("got hello data", "interval", data.HeartbeatInterval)

	interval := time.Duration(data.HeartbeatInterval) * time.Millisecond
	go c.runHeartbeatLoop(ctx, interval)

	return nil
}

func (c *Client) sendHeartbeat(ctx context.Context, seq *int) {
	pretty := "<nil>"
	if seq != nil {
		pretty = fmt.Sprint(*seq)
	}
	c.logger.Debug("will send heartbeat", "data", pretty)

	outgoing := map[string]any{
		"op": Heartbeat,
		"d":  seq,
	}
	data, _ := json.Marshal(outgoing)

	c.state.acked = false
	c.write(ctx, data)
}

func (c *Client) ackHeartbeat(ctx context.Context, event GatewayEvent) error {
	c.state.acked = true
	return nil
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

func (c *Client) handleReady(event *GatewayEvent) error {
	var ready Ready
	err := mapstructure.Decode(event.Data, &ready)
	if err != nil {
		return fmt.Errorf("failed to decode ready event: %w", err)
	}

	c.state.resumeURL = ready.ResumeGatewayURL
	c.state.sessionID = ready.SessionID
	return nil
}

func (c *Client) handleMessage(event *GatewayEvent) (*Message, error) {
	var message Message

	err := mapstructure.Decode(event.Data, &message)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &message, nil
}
