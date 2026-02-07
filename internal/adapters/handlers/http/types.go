package http

// BulkPayoutRequest is what the Frontend/User sends us
type BulkPayoutRequest struct {
	BatchReference string       `json:"batch_reference" example:"JAN_SALARY_2025"`
	Items          []PayoutItem `json:"items"`
}

type PayoutItem struct {
	// User Details (Step 1 of Afriex Flow)
	RecipientName  string `json:"recipient_name" example:"Emeka Okonkwo"`
	RecipientPhone string `json:"recipient_phone" example:"+2348012345678"`
	RecipientEmail string `json:"recipient_email" example:"emeka@example.com"`
	CountryCode    string `json:"country_code" example:"NG"`

	// Bank Details (Step 2 of Afriex Flow)
	BankCode      string `json:"bank_code" example:"033"`       // e.g., UBA
	AccountNumber string `json:"account_number" example:"2000012345"`
	
	// Transaction Details (Step 3 of Afriex Flow)
	Amount   float64 `json:"amount" example:"5000.00"` // User sends float, we convert to cents
	Currency string  `json:"currency" example:"NGN"`
}

type BulkPayoutResponse struct {
	BatchID string `json:"batch_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}