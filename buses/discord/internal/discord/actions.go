package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"nhooyr.io/websocket"
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

func (c *Client) maybeResume(ctx context.Context, err websocket.CloseError) error {
	switch err.Code {
	case AuthenticationFailed, InvalidShard, ShardingRequired,
		InvalidAPIVersion, InvalidIntent, DisallowedIntent:
		return err
	}

	return c.resume(ctx)
}

func (c *Client) resume(ctx context.Context) error {
	rctx, cancel := reconnectContext(ctx)
	defer cancel()

	c.logger.Info("resuming websocket connection")
	conn, _, err := websocket.Dial(rctx, c.state.resumeURL, nil)
	if err != nil {
		return fmt.Errorf("could not reconnect to websocket: %w", err)
	}

	c.ws = conn

	data, _ := json.Marshal(arbitraryJSON{
		"op": Resume,
		"d": GatewayResume{
			Token:     c.token,
			SessionID: c.state.sessionID,
			Seq:       c.state.seq,
		},
	})

	c.write(rctx, data)
	return nil
}

func (c *Client) reconnect(ctx context.Context) error {
	c.reconnecting <- struct{}{}
	c.state = clientState{
		gatewayURL: c.state.gatewayURL,
	}

	rctx, cancel := reconnectContext(ctx)
	defer cancel()

	c.logger.Info("reconnecting websocket connection")
	return c.Connect(rctx)
}

func (c *Client) handleMessage(event *GatewayEvent) (*Message, error) {
	var message Message

	err := mapstructure.Decode(event.Data, &message)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &message, nil
}
