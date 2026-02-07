package afriex

// import "time"

// --- ERROR HANDLING ---
type APIError struct {
	Code    string                 `json:"code"`
	ErrorMsg string                `json:"error"`
	Details map[string]interface{} `json:"details"`
}

func (e *APIError) Error() string {
	return e.Code + ": " + e.ErrorMsg
}

// --- 1. CUSTOMER ---
type CreateCustomerRequest struct {
	FullName    string                 `json:"fullName"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone"`
	CountryCode string                 `json:"countryCode"` // e.g., "NG", "US"
	Kyc         map[string]interface{} `json:"kyc"`         // Empty for now
	Meta        map[string]interface{} `json:"meta"`
}

type CustomerResponse struct {
	Data struct {
		CustomerID string `json:"customerId"`
		Email      string `json:"email"`
	} `json:"data"`
}

// --- 2. PAYMENT METHOD ---
type Institution struct {
	InstitutionName string `json:"institutionName,omitempty"`
	InstitutionCode string `json:"institutionCode,omitempty"` // Bank Code
}

type CreatePaymentMethodRequest struct {
	Channel       string      `json:"channel"`       // "BANK_ACCOUNT", "MOBILE_MONEY"
	CustomerID    string      `json:"customerId"`
	AccountName   string      `json:"accountName"`
	AccountNumber string      `json:"accountNumber"`
	CountryCode   string      `json:"countryCode"`
	Institution   Institution `json:"institution"`
	// We can leave recipient/transaction fields empty for basic cases
}

type PaymentMethodResponse struct {
	Data struct {
		PaymentMethodID string `json:"paymentMethodId"`
		Channel         string `json:"channel"`
	} `json:"data"`
}

// --- 3. TRANSACTION (The Payout) ---
type CreateTransactionRequest struct {
	CustomerID          string            `json:"customerId"`
	DestinationAmount   string            `json:"destinationAmount"` // Docs say string! "100.50"
	DestinationCurrency string            `json:"destinationCurrency"`
	SourceCurrency      string            `json:"sourceCurrency"`
	DestinationID       string            `json:"destinationId"` // From PaymentMethod
	Meta                map[string]string `json:"meta"`
}

type TransactionResponse struct {
	Data struct {
		TransactionID     string `json:"transactionId"`
		Status            string `json:"status"`
		SourceAmount      string `json:"sourceAmount"`
		DestinationAmount string `json:"destinationAmount"`
	} `json:"data"`
}

// --- 4. RATES ---
type RateResponse struct {
	Rates     map[string]map[string]string `json:"rates"`
	UpdatedAt int64                        `json:"updatedAt"`
}