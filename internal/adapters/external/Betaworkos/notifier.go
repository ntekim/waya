package betaworkos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"waya/internal/config"
	"waya/internal/core/domain" // Use domain model for the payload
)

type Notifier struct {
	webhookURL string
	httpClient *http.Client
}

func NewNotifier(cfg config.WayaConfig) *Notifier {
	return &Notifier{
		webhookURL: cfg.BETAWORKOSWebhookURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// NotifyBatchCompletion sends the final batch status to the client's webhook URL
func (n *Notifier) NotifyBatchCompletion(ctx context.Context, batchID string, payouts []domain.Payout) error {
	slog.Info("🔔 Attempting to notify client system (BetaWorkOS)", "batch_id", batchID, "url", n.webhookURL)
	
	if n.webhookURL == "" {
		slog.Warn("Skipping client notification: CLIENT_WEBHOOK_URL is not set.")
		return nil
	}

	// 1. Prepare Payload (use a clean struct for the final result)
	payload := map[string]any{
		"event": "WAYA.BATCH_COMPLETED",
		"batch_id": batchID,
		"timestamp": time.Now().UTC(),
		"data": map[string]any{
			"total_count": len(payouts),
			"payouts": payouts,
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	// 2. Make the HTTP Call
	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create client notification request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Optional: Add a signature header for real security: req.Header.Set("X-Waya-Signature", ...)

	resp, err := n.httpClient.Do(req)
	if err != nil {
		slog.Error("❌ Client notification FAILED (Client system down?)", "err", err, "url", n.webhookURL)
		// For a real system, you would retry this later.
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		slog.Error("❌ Client notification failed: received non-2xx status", "status", resp.StatusCode)
		return fmt.Errorf("client returned status code %d", resp.StatusCode)
	}
	
	slog.Info("✅ Successfully notified client system.")
	return nil
}