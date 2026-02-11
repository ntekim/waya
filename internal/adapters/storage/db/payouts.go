package db

import (
	"context"
	"database/sql"
	"fmt"

	// "fmt"
	// "waya/internal/adapters/storage/db"
	"waya/internal/core/domain"
	"waya/internal/core/ports"
)

type SQLiteRepo struct {
	q *Queries
}

// Ensure SQLiteRepo implements PaymentRepository
var _ ports.PaymentRepository = (*SQLiteRepo)(nil)

func NewRepository(database *Database) *SQLiteRepo {
	return &SQLiteRepo{
		q: database.Q,
	}
}

func (r *SQLiteRepo) SavePayout(ctx context.Context, p domain.Payout) error {
    // Handle Nullable Strings for SQLC
    email := sql.NullString{String: p.RecipientEmail, Valid: p.RecipientEmail != ""}
    tag := sql.NullString{String: p.RecipientTag, Valid: p.RecipientTag != ""}
    bankCode := sql.NullString{String: p.BankCode, Valid: p.BankCode != ""}
    accNum := sql.NullString{String: p.AccountNumber, Valid: p.AccountNumber != ""}
    bankName := sql.NullString{String: p.BankName, Valid: p.BankName != ""}
    batchID := sql.NullString{String: p.BatchID, Valid: p.BatchID != ""}

    _, err := r.q.CreatePayout(ctx, CreatePayoutParams{
        ID:             p.ID,
        BatchID:        batchID,
        ReferenceID:    p.ReferenceID,
        RecipientName:  p.RecipientName,
        RecipientPhone: p.RecipientPhone,
        RecipientEmail: email,
        RecipientTag:   tag,
        CountryCode:    p.CountryCode,
        BankCode:       bankCode,
        AccountNumber:  accNum,
        BankName:       bankName,
        Amount:         p.Amount,
        Currency:       p.Currency,
        Status:         p.Status,
    })
    return err
}

// Also update GetPayout and ListPayouts to map back from DB to Domain!
func (r *SQLiteRepo) GetPayout(ctx context.Context, id string) (*domain.Payout, error) {
    row, err := r.q.GetPayout(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, return nil without error
		}
		return nil, err
	}

	// Map DB model to Domain model, handling NULLs
    return &domain.Payout{
        ID:             row.ID,
        BatchID:        row.BatchID.String,
        ReferenceID:    row.ReferenceID,
        RecipientName:  row.RecipientName,
        RecipientPhone: row.RecipientPhone,
        RecipientEmail: row.RecipientEmail.String,
        CountryCode:    row.CountryCode,
        BankCode:       row.BankCode.String,
        AccountNumber:  row.AccountNumber.String,
        Amount:         row.Amount,
        Currency:       row.Currency,
        Status:         row.Status,
        // ... rest
    }, nil
}

func (r *SQLiteRepo) UpdatePayoutStatus(ctx context.Context, id string, status string, errMsg string) error {
	return r.q.UpdatePayoutStatus(ctx, UpdatePayoutStatusParams{
		ID:     id,
		Status: status,
		ErrorMessage: sql.NullString{
			String: errMsg,
			Valid:  errMsg != "",
		},
	})
}

func (r *SQLiteRepo) ListPayouts(ctx context.Context, limit int) ([]domain.Payout, error) {
	rows, err := r.q.ListPayouts(ctx)
	if err != nil {
		return nil, err
	}

    var payouts []domain.Payout
    for _, row := range rows {
        payouts = append(payouts, domain.Payout{
            ID:             row.ID,
            BatchID:        row.BatchID.String,
            ReferenceID:    row.ReferenceID,
            RecipientName:  row.RecipientName,
            RecipientPhone: row.RecipientPhone,
            RecipientEmail: row.RecipientEmail.String,
            CountryCode:    row.CountryCode,
            BankCode:       row.BankCode.String,
            AccountNumber:  row.AccountNumber.String,
            BankName:       row.BankName.String,
            Amount:         row.Amount,
            Currency:       row.Currency,
            Status:         row.Status,
        })
    }

	return payouts, nil
}

// New, efficient method to get payouts by BatchID
func (r *SQLiteRepo) ListPayoutsByBatchID(ctx context.Context, batchID string) ([]domain.Payout, error) {
    // Call the SQLC-generated function directly
    rows, err := r.q.ListPayoutsByBatchID(ctx, sql.NullString{String: batchID, Valid: true})
    if err != nil {
        return nil, err
    }

    // Map rows to domain.Payout (you must write the mapping logic here)
    var payouts []domain.Payout
    for _, row := range rows {
        payouts = append(payouts, domain.Payout{
            ID:             row.ID,
            BatchID:        row.BatchID.String,
            ReferenceID:    row.ReferenceID,
            RecipientName:  row.RecipientName,
            RecipientPhone: row.RecipientPhone,
            RecipientEmail: row.RecipientEmail.String,
            CountryCode:    row.CountryCode,
            BankCode:       row.BankCode.String,
            AccountNumber:  row.AccountNumber.String,
            BankName:       row.BankName.String,
            Amount:         row.Amount,
            Currency:       row.Currency,
            Status:         row.Status,
        })
    }

    if len(payouts) == 0 {
        return nil, fmt.Errorf("not found") // Return error on empty set
    }
    return payouts, nil
}