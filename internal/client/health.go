package client

import "context"

type HealthResponse struct {
	Status string `json:"status"`
}

func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var out HealthResponse
	if err := c.DoJSON(ctx, "GET", "/api/v1/health", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
