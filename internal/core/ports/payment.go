package ports

import (
	"context"
	"waya/internal/adapters/payments/afriex"
	"waya/internal/core/domain"
)

// PaymentRepository defines how we store data (Database Port)
type PaymentRepository interface {
	SavePayout(ctx context.Context, payout domain.Payout) error
	GetPayout(ctx context.Context, id string) (*domain.Payout, error)
	UpdatePayoutStatus(ctx context.Context, id string, status string, errMsg string) error
	ListPayouts(ctx context.Context, limit int) ([]domain.Payout, error)
	ListPayoutsByBatchID(ctx context.Context, batchID string) ([]domain.Payout, error)
}

// AfriexGateway defines how we talk to the outside world (API Port)
type AfriexGateway interface {
	GetCustomerByEmail(ctx context.Context, email string) (string, error) // Returns customerID
    FindPaymentMethod(ctx context.Context, customerID, accountNumber string) (string, error) // Returns paymentMethodID
	// Step 1: Onboard
	CreateCustomer(ctx context.Context, req afriex.CreateCustomerRequest) (string, error)
	
	// Step 2: Link Bank/Wallet
	CreatePaymentMethod(ctx context.Context, req afriex.CreatePaymentMethodRequest) (string, error)
	
	// Step 3: Pay
	CreateTransaction(ctx context.Context, req afriex.CreateTransactionRequest) (*afriex.TransactionResponse, error)
	
	// Utils
	GetRates(ctx context.Context, base, symbols string) (*afriex.RateResponse, error)
}

type ExternalClientNotifier interface{
	NotifyBatchCompletion(ctx context.Context, batchID string, payouts []domain.Payout) error
}

// Data structures specifically for the Afriex Port
type AfriexTransferRequest struct {
	Reference    string
	RecipientTag string
	Amount       int64
	Currency     string
}

type AfriexTransferResponse struct {
	TransactionID string
	Status        string
	Fee           int64
}