package afriex

import (
	"context"
	"fmt"
)

func (c *Client) GetRates(ctx context.Context, base, symbols string) (*RateResponse, error) {
	var resp RateResponse
	// Query params handling manually for speed
	path := fmt.Sprintf("/v2/public/rates?base=%s&symbols=%s", base, symbols)
	err := c.do(ctx, "GET", path, nil, &resp)
	return &resp, err
}