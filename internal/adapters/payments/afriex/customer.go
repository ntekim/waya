package afriex

import "context"

func (c *Client) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (string, error) {
	var resp CustomerResponse
	err := c.do(ctx, "POST", "/api/v1/customer", req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Data.CustomerID, nil
}