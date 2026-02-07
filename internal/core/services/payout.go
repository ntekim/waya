package services

import (
	// "context"
	// "fmt"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	// "sync"
	// "time"

	// "waya/internal/core/domain"
	"waya/internal/adapters/payments/afriex"
	"waya/internal/core/domain"
	"waya/internal/core/ports"
)

type PayoutService struct {
	repo    ports.PaymentRepository
	gateway ports.AfriexGateway
	logger  *slog.Logger
}

func NewPayoutService(repo ports.PaymentRepository, gateway ports.AfriexGateway, logger *slog.Logger) *PayoutService {
	return &PayoutService{
		repo:    repo,
		gateway: gateway,
		logger:  logger,
	}
}

// ExecuteBatch is the "Money Maker" function.
// It takes a list of requests, saves them to DB, and fires them in parallel.
func (s *PayoutService) ExecuteBatch(ctx context.Context, batchID string, payouts []domain.Payout) error {
	slog.Info("🚀 Starting Batch Execution", "batch_id", batchID, "count", len(payouts))

	// 1. Validation & Persistence Loop
	for _, p := range payouts {
		p.BatchID = batchID
		p.Status = domain.StatusPending
		p.CreatedAt = time.Now()
		
		if err := s.repo.SavePayout(ctx, p); err != nil {
			return fmt.Errorf("failed to save payout %s: %w", p.ID, err)
		}
	}

	// 2. Parallel Execution (The "Orchestration")
	// We use a WaitGroup to handle concurrency
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent requests to avoid Rate Limits (429)

	for _, p := range payouts {
		wg.Add(1)
		
		go func(payout domain.Payout) {
			defer wg.Done()
			
			// Acquire token (rate limiting)
			semaphore <- struct{}{} 
			defer func() { <-semaphore }()

			s.processSinglePayout(context.Background(), payout)
		}(p)
	}

	// Don't wait for completion if you want non-blocking API response.
	// For Hackathon demo, waiting is safer to show results immediately.
	wg.Wait()
	
	slog.Info("✅ Batch Execution Complete", "batch_id", batchID)
	return nil
}

func (s *PayoutService) processSinglePayout(ctx context.Context, p domain.Payout) {
	s.repo.UpdatePayoutStatus(ctx, p.ID, domain.StatusProcessing, "")

	// --- STEP 1: CREATE CUSTOMER ---
	// In a real app, you'd check DB if customer exists. For Hackathon, strict upsert or create.
	custID, err := s.gateway.CreateCustomer(ctx, afriex.CreateCustomerRequest{
		FullName:    p.RecipientName, // Need to add this to Domain Payout struct
		Email:       "temp_" + p.RecipientTag + "@waya.com", // Fake email if not provided
		Phone:       p.RecipientPhone, // Need to add this
		CountryCode: p.CountryCode,    // Need to add this (e.g. "NG")
	})
	if err != nil {
		s.handleError(ctx, p, "Failed to create customer", err)
		return
	}

	// --- STEP 2: CREATE PAYMENT METHOD (BANK) ---
	pmID, err := s.gateway.CreatePaymentMethod(ctx, afriex.CreatePaymentMethodRequest{
		Channel:       "BANK_ACCOUNT", // Or MOBILE_MONEY based on logic
		CustomerID:    custID,
		AccountName:   p.RecipientName,
		AccountNumber: p.AccountNumber, // Add to Domain
		CountryCode:   p.CountryCode,
		Institution: afriex.Institution{
			InstitutionCode: p.BankCode, // Add to Domain
		},
	})
	if err != nil {
		s.handleError(ctx, p, "Failed to link bank account", err)
		return
	}

	// --- STEP 3: SEND MONEY ---
	// Convert int64 cents to string "100.50"
	amountStr := fmt.Sprintf("%.2f", float64(p.Amount)/100.0)

	txResp, err := s.gateway.CreateTransaction(ctx, afriex.CreateTransactionRequest{
		CustomerID:          custID,
		DestinationID:       pmID,
		SourceCurrency:      "USD", // We pay in USD
		DestinationCurrency: p.Currency,
		DestinationAmount:   amountStr,
		Meta: map[string]string{
			"narration": "Waya Payout - " + p.BatchID,
		},
	})

	if err != nil {
		s.handleError(ctx, p, "Transaction failed", err)
		return
	}

	// Success!
	slog.Info("💰 Paid!", "tx_id", txResp.Data.TransactionID)
	s.repo.UpdatePayoutStatus(ctx, p.ID, domain.StatusSuccess, "")
}

func (s *PayoutService) handleError(ctx context.Context, p domain.Payout, msg string, err error) {
	slog.Error(msg, "id", p.ID, "err", err)
	s.repo.UpdatePayoutStatus(ctx, p.ID, domain.StatusFailed, fmt.Sprintf("%s: %v", msg, err))
}

// ListPayoutsByBatchID fetches all payouts belonging to a single batch.
func (s *PayoutService) ListPayoutsByBatchID(ctx context.Context, batchID string) ([]domain.Payout, error) {
	// WARNING: This is inefficient for a production system but works for a hackathon
	// because we don't want to create a new SQL query file and a new SQLC call.
	
	// In a real app, you would add a SQL query like:
	// -- name: ListPayoutsByBatchID :many SELECT * FROM payouts WHERE batch_id = ? ORDER BY created_at DESC;
	

	// For now, let's just use the ListPayouts and filter manually.
	allPayouts, err := s.repo.ListPayoutsByBatchID(ctx, batchID)
	if err != nil {
		return nil, err
	}
	
	var batchPayouts []domain.Payout
	for _, p := range allPayouts {
		if p.BatchID == batchID {
			batchPayouts = append(batchPayouts, p)
		}
	}

	if len(batchPayouts) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return batchPayouts, nil
}