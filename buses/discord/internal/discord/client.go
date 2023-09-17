package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"nhooyr.io/websocket"
)

type Client struct {
	Err error // set on fatal errors

	// these are passed in and stashed
	token  string
	logger *slog.Logger

	// persistent state
	ws    *websocket.Conn
	state clientState

	// communication channels
	fatalNotifier chan struct{} // closed when we die, which sets .Err
	errors        chan error    // over which we send non-fatal errors
	reconnecting  chan struct{} // used to signal the heartbeat loop to shut down
}

type clientState struct {
	identified bool
	acked      bool
	seq        *int
	sessionID  string
	gatewayURL string
	resumeURL  string
}

func NewClient(logger *slog.Logger, token string) *Client {
	return &Client{
		token:         token,
		logger:        logger,
		fatalNotifier: make(chan struct{}),
		reconnecting:  make(chan struct{}),
		errors:        make(chan error),
	}
}

func (c *Client) Fatal() <-chan struct{} {
	return c.fatalNotifier
}

func (c *Client) Errors() <-chan error {
	return c.errors
}

func (c *Client) Connect(ctx context.Context) error {
	if err := c.loadGatewayURL(); err != nil {
		return err
	}

	conn, _, err := websocket.Dial(ctx, c.state.gatewayURL, nil)
	if err != nil {
		return fmt.Errorf("could not connect to websocket: %w", err)
	}

	c.ws = conn
	return nil
}

var errFrameNotText = errors.New("got binary websocket type")

func (c *Client) Run(ctx context.Context, dataCh chan<- Message, errCh chan<- error) {
	defer c.ws.Close(websocket.StatusNormalClosure, "so long")

	var closeErr websocket.CloseError

	for {
		// c.logger.Debug("will read from websocket")
		readCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		typ, data, err := c.ws.Read(readCtx)
		cancel()

		switch {
		case errors.Is(err, context.Canceled):
			// this is fine.
		case errors.As(err, &closeErr):
			err = c.maybeResume(ctx, closeErr)
			if err == nil {
				continue
			}

			fallthrough
		case err != nil:
			c.Err = fmt.Errorf("ws read error: %w", err)
			close(c.fatalNotifier)
			return
		case typ != websocket.MessageText:
			errCh <- errFrameNotText
		}

		select {
		case <-ctx.Done():
			return
		case <-c.fatalNotifier:
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

func (c *Client) handleFrame(ctx context.Context, data []byte) (*Message, error) {
	var event GatewayEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("bad frame from discord: %w", err)
	}

	if event.Seq != nil {
		c.state.seq = event.Seq
	}

	switch event.Op {
	case Hello:
		err := c.doHello(ctx, event)
		return nil, err

	case Heartbeat:
		c.sendHeartbeat(ctx, c.state.seq)
		return nil, nil

	case HeartbeatACK:
		if !c.state.identified {
			c.doIdentify(ctx)
			c.state.identified = true
		}

		err := c.ackHeartbeat(ctx, event)
		return nil, err

	case Reconnect:
		return nil, c.resume(ctx)

	case InvalidSession:
		canResume, ok := event.Data.(bool)
		if !ok {
			return nil, fmt.Errorf("got an InvalidSession op with bad data: %+v", event.Data)
		}

		if canResume {
			return nil, c.resume(ctx)
		}

		return nil, c.reconnect(ctx)

	case Dispatch:
		return c.dispatch(ctx, &event)

	default:
		c.logger.Debug("ignoring gateway event", "type", event.Op)
		fmt.Printf("  frame: %s\n", string(data))
		return nil, nil
	}
}

func reconnectContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 15*time.Second)
}

func (c *Client) dispatch(ctx context.Context, evt *GatewayEvent) (*Message, error) {
	switch evt.Type {
	case TypeReady:
		return nil, c.handleReady(evt)

	case TypeMessageCreate:
		return c.handleMessage(evt)

	case TypeResumed:
		c.logger.Debug("finished resuming")
		return nil, nil

	default:
		c.logger.Debug("ignoring dispatch event", "type", evt.Type)
	}

	return nil, nil
}

func (c *Client) runHeartbeatLoop(ctx context.Context, interval time.Duration) {
	// jitter := rand.Float64()
	jitter := 0.09
	firstInterval := time.Duration(float64(interval) * jitter)
	c.logger.Debug("waiting to send first heartbeat", "interval", firstInterval)

	timer := time.NewTimer(firstInterval)
	first := true

	shutdown := func(reason string) {
		c.logger.Debug("shutting down heartbeat loop", "reason", reason)
		timer.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			shutdown("context cancelled")
			return
		case <-c.reconnecting:
			shutdown("reconnecting")
			return

		case <-timer.C:
			if !c.state.acked && !first {
				c.errors <- fmt.Errorf("failed to receive ack for last heartbeat")

				if err := c.resume(ctx); err != nil {
					c.Err = fmt.Errorf("failed to reconnect after lost ack: %w", err)
					c.fatalNotifier <- struct{}{}
				}

				continue
			}

			c.sendHeartbeat(ctx, c.state.seq)
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
		c.errors <- fmt.Errorf("bad websocket write: %w", err)
	}
}
