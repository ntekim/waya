package domain

import (
	"errors"
	"time"
)

// PayoutStatus Enum
const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusSuccess    = "SUCCESS"
	StatusFailed     = "FAILED"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidCurrency   = errors.New("invalid currency pair")
)

// Payout represents a single money transfer
type Payout struct {
	ID           string
	BatchID      string
	ReferenceID  string
	
	RecipientName  string // "John Doe"
	RecipientPhone string // "+234..."
	RecipientEmail string // "john@example.com"
	RecipientTag   string // Optional: If sending to Afriex Wallet directly
	
	CountryCode    string // "NG", "GH", "KE"
	BankCode       string // "033" (UBA)
	AccountNumber  string // "2039..."
	BankName       string // "United Bank for Africa"
	// -----------------------------

	Amount       int64  // Cents
	Currency     string // "NGN"
	Status       string
	ErrorMessage string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Batch represents a bulk transfer request
type Batch struct {
	ID          string
	TotalAmount int64
	TotalCount  int
	Status      string
	Payouts     []Payout
}