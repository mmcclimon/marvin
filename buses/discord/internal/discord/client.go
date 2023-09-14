package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

const urlBase = "https://discord.com/api/v10"

var httpClient = http.Client{Timeout: 5 * time.Second}

type Client struct {
	C      <-chan struct{}
	ch     chan struct{}
	token  string
	ws     *websocket.Conn
	err    error
	acked  bool
	logger *slog.Logger
}

func NewClient(logger *slog.Logger, token string) *Client {
	ch := make(chan struct{})
	return &Client{
		token:  token,
		C:      ch,
		ch:     ch,
		logger: logger,
	}
}

func (c *Client) Err() error {
	return c.err
}

func (c *Client) Connect(ctx context.Context, wssURL string) error {
	if wssURL == "" {
		resp, err := httpClient.Get(urlBase + "/gateway")
		if err != nil {
			return fmt.Errorf("could not fetch gateway: %w", err)
		}

		defer resp.Body.Close()

		var data struct{ URL string }
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("could not read gateway response: %w", err)
		}

		wssURL = data.URL
		fmt.Println(wssURL)
	}

	conn, _, err := websocket.Dial(ctx, wssURL, nil)
	if err != nil {
		return fmt.Errorf("could not connect to websocket: %w", err)
	}

	c.ws = conn
	return nil
}

var errFrameNotText = errors.New("got binary websocket type")

func (c *Client) Run(ctx context.Context, dataCh chan<- GatewayEvent, errCh chan<- error) {
	defer c.ws.Close(websocket.StatusNormalClosure, "so long")

	for {
		c.logger.Debug("will read from websocket")
		readCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		typ, data, err := c.ws.Read(readCtx)
		cancel()

		switch {
		case errors.Is(err, context.Canceled):
			// this is fine.
		case err != nil:
			c.err = fmt.Errorf("ws read error: %w", err)
			close(c.ch)
			return
		case typ != websocket.MessageText:
			errCh <- errFrameNotText
		}

		select {
		case <-ctx.Done():
			return
		default:
			evt, err := c.handleFrame(ctx, data)

			switch {
			case err != nil:
				errCh <- err
			case evt != nil:
				dataCh <- *evt
			}
		}
	}
}

func (c *Client) handleFrame(ctx context.Context, data []byte) (*GatewayEvent, error) {
	fmt.Printf("FRAME: %s\n", string(data))

	var event GatewayEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("bad frame from discord: %w", err)
	}

	switch event.Op {
	case Hello:
		return c.doHello(ctx, event)
	case HeartbeatACK:
		return c.ackHeartbeat(ctx, event)
	default:
		fmt.Printf("ignoring gateway event: %+v\n", event)
	}

	return &event, nil
}

func (c *Client) runHeartbeatLoop(ctx context.Context, interval time.Duration, seq *int) {
	jitter := rand.Float64()
	firstInterval := time.Duration(float64(interval) * jitter)
	c.logger.Debug("waiting to send first heartbeat", "interval", firstInterval)

	timer := time.NewTimer(firstInterval)
	first := true

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("shutting down heartbeat loop")
			timer.Stop()
			return

		case <-timer.C:
			if !c.acked && !first {
				// TODO: handle this somehow
				c.logger.Warn("failed to receive ack for last heartbeat")
				return
			}
			c.sendHeartbeat(ctx, seq)
			timer.Reset(interval)
			first = false
		}
	}
}

func (c *Client) write(ctx context.Context, data []byte) {
	writeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := c.ws.Write(writeCtx, websocket.MessageText, data)
	if err != nil {
		// TODO: handle this somehow
		c.logger.Debug("bad websocket write", "err", err)
	}
}
