package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const urlBase = "https://discord.com/api/v10"

var httpClient = http.Client{Timeout: 5 * time.Second}

func URLFor(endpoint string, args ...any) string {
	return fmt.Sprintf(urlBase+endpoint, args...)
}

func (c *Client) Post(ctx context.Context, url string, data any) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("bad json encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("bad request creation: %w", err)
	}

	req.Header.Add("Authorization", "Bot "+c.token)
	req.Header.Add("Content-Type", "application/json")

	return httpClient.Do(req)
}
