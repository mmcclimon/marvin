package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type HelloData struct {
	HeartbeatInterval int `mapstructure:"heartbeat_interval"`
}

func (c *Client) doHello(ctx context.Context, event GatewayEvent) (*GatewayEvent, error) {
	var data HelloData
	if err := mapstructure.Decode(event.Data, &data); err != nil {
		return nil, fmt.Errorf("bad hello decode: %w", err)
	}

	Logger.Debug("got hello data", "interval", data.HeartbeatInterval)

	interval := time.Duration(data.HeartbeatInterval) * time.Millisecond
	go c.runHeartbeatLoop(ctx, interval, event.Seq)

	return nil, nil
}

func (c *Client) sendHeartbeat(ctx context.Context, seq *int) {
	outgoing := map[string]any{
		"op": Heartbeat,
		"d":  seq,
	}

	Logger.Debug("will send heartbeat", "data", outgoing)
	data, _ := json.Marshal(outgoing)

	c.acked = false
	c.write(ctx, data)
}

func (c *Client) ackHeartbeat(ctx context.Context, event GatewayEvent) (*GatewayEvent, error) {
	c.acked = true
	return nil, nil
}
