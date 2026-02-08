package afriex

import "context"

func (c *Client) FindPaymentMethod(ctx context.Context, customerID, accountNumber string) (string, error) {
	return "", nil
}


func (c *Client) CreatePaymentMethod(ctx context.Context, req CreatePaymentMethodRequest) (string, error) {
	var resp PaymentMethodResponse
	err := c.do(ctx, "POST", "/api/v1/payment-method", req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Data.PaymentMethodID, nil
}

func (c *Client) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*TransactionResponse, error) {
	var resp TransactionResponse
	err := c.do(ctx, "POST", "/api/v1/transaction", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}