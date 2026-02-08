package http

import (
	"context"
	"log/slog"
	"net/http"

	// "time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"waya/internal/core/domain"
	"waya/internal/core/services"
)

type PayoutHandler struct {
	service *services.PayoutService
}

func NewPayoutHandler(service *services.PayoutService) *PayoutHandler {
	return &PayoutHandler{service: service}
}

// @Summary Trigger Bulk Payout
// @Description Accepts a list of recipients, creates customers/payment methods on Afriex, and sends money.
// @Tags Payouts
// @Accept json
// @Produce json
// @Param request body BulkPayoutRequest true "The batch of payouts to process"
// @Success 202 {object} BulkPayoutResponse "Batch accepted for background processing"
// @Failure 400 {object} map[string]string "Invalid JSON or payload"
// @Router /payouts [post]
func (h *PayoutHandler) HandleBulkPayout(c echo.Context) error {
	var req BulkPayoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	// 1. Generate a Batch ID
	batchID := uuid.New().String()

	// 2. Map DTO -> Domain Entities
	var domainPayouts []domain.Payout
	for _, item := range req.Items {
		domainPayouts = append(domainPayouts, domain.Payout{
			ID:             uuid.New().String(), // We generate our own ID
			BatchID:        batchID,
			ReferenceID:    req.BatchReference + "-" + uuid.New().String()[0:8],
			
			// User Data
			RecipientName:  item.RecipientName,
			RecipientPhone: item.RecipientPhone,
			RecipientEmail: item.RecipientEmail,
			CountryCode:    item.CountryCode,
			
			// Bank Data
			BankCode:       item.BankCode,
			AccountNumber:  item.AccountNumber,
			
			// Money (Convert Float to Cents/Kobo)
			Amount:   int64(item.Amount * 100), 
			Currency: item.Currency,
		})
	}

	// 3. Call the Orchestrator (Async)
	// We use a goroutine here so the HTTP request returns immediately (202 Accepted)
	// while the heavy lifting happens in the background.
	go func() {
		// Create a background context since the request context will cancel when we return
		ctx := context.Background() 
		_ = h.service.ExecuteBatch(ctx, batchID, domainPayouts)
	}()

	return c.JSON(http.StatusAccepted, BulkPayoutResponse{
		BatchID: batchID,
		Status:  "PROCESSING",
		Message: "Batch accepted. Check status via /payouts/status/" + batchID,
	})
}

// @Summary Get Batch Status
// @Description Retrieves all payouts and status for a given batch ID.
// @Tags Payouts
// @Produce json
// @Param batch_id path string true "Unique ID of the payout batch"
// @Success 200 {object} domain.Batch "Returns batch details and list of payouts"
// @Failure 404 {object} map[string]string "Batch ID not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /payouts/{batch_id} [get]
func (h *PayoutHandler) GetBatchStatus(c echo.Context) error {
	batchID := c.Param("batch_id")

	// This is a simplified fetch for a hackathon.
	// In production, we'd query a dedicated Batch table.
	// For now, we query all payouts and filter by BatchID.
	
	ctx := c.Request().Context()
	
	// Since we don't have a ListByBatchID, we'll implement a mock-up/quick query
	// by fetching all and filtering in the service (not ideal, but fast for now)
	payouts, err := h.service.ListPayoutsByBatchID(ctx, batchID)
	
	if err != nil {
		if err.Error() == "not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Batch ID not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve batch status"})
	}

	// Create a mock Batch response
	totalAmount := int64(0)
	for _, p := range payouts {
		totalAmount += p.Amount
	}

	response := domain.Batch{
		ID: batchID,
		TotalAmount: totalAmount,
		TotalCount: len(payouts),
		Status: "COMPLETED", // Simplified status for demo
		Payouts: payouts,
	}

	return c.JSON(http.StatusOK, response)
}

// HandleAfriexWebhook processes incoming transaction status updates from Afriex
// @Summary Afriex Webhook Listener
// @Description Receives real-time transaction updates (e.g., SUCCESS/FAILED) from Afriex.
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Afriex Webhook Payload"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid Signature or Payload"
// @Router /webhooks/afriex [post]
func (h *PayoutHandler) HandleAfriexWebhook(c echo.Context) error {
	// 1. **CRITICAL:** SECURITY/Signature Verification (Mocked for Hackathon)
	// You should verify the 'x-webhook-signature' header here.
	
	// 2. Decode Payload
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Payload")
	}

	// 3. Process Event
	event := payload["event"].(string)
	data := payload["data"].(map[string]interface{})
	
	if event == "TRANSACTION.UPDATED" {
		txID := data["transactionId"].(string)
		status := data["status"].(string) // "SUCCESS", "FAILED"
		
		slog.Info("🔔 Webhook received", "tx_id", txID, "new_status", status)

		// 4. Update the DB (This requires an ID lookup, which we can mock)
		// We can't map Afriex TX ID to Waya Payout ID without a new DB column.
		// For the demo: just log and show the concept.
		
		// In a real app, you would: 
		// h.service.UpdateStatusByAfriexID(txID, status) 
	}

	// 5. Respond Immediately (200 OK)
	return c.String(http.StatusOK, "OK")
}

// @Summary List All Payouts
// @Description Retrieves a complete, paginated list of all payout records for the Waya Admin Dashboard.
// @Tags Payouts
// @Produce json
// @Success 200 {object} []domain.Payout "List of all payouts"
// @Failure 500 {object} map[string]string "Server error"
// @Router /payouts/all [get]
func (h *PayoutHandler) HandleListAllPayouts(c echo.Context) error {
    ctx := c.Request().Context()
    
    // Hardcode a high limit for the dashboard demo
    limit := 100 
    
    payouts, err := h.service.ListPayouts(ctx, limit)

    if err != nil {
        slog.Error("Failed to list all payouts", "err", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve payout history"})
    }

    // Return the list directly
    return c.JSON(http.StatusOK, payouts)
}