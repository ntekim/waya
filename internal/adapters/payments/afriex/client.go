package afriex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"waya/internal/config"
	// "waya/internal/core/ports" // Don't import ports here to avoid cycle if types are in same package
    // Instead, make sure types.go is in THIS package (adapters/payment/afriex)
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg config.AfriexConfig) *Client {
	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: "https://staging.afx-server.com", // Force Staging for Hackathon
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generic Helper
func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, _ := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// LOGGING (Crucial for Hackathon debugging)
	slog.Debug("Afriex API", "path", path, "status", resp.StatusCode, "resp", string(respBody))

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil {
			return &apiErr
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		return json.Unmarshal(respBody, result)
	}
	return nil
}