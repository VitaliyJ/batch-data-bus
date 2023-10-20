package externalservice

import (
	"context"
	"fmt"
	"time"
)

type Client struct{}

// NewClient returns new Client instance
func NewClient() *Client {
	return &Client{}
}

// GetLimits returns requests limits
// TODO need implementation
func (c Client) GetLimits() (n uint64, p time.Duration) {
	return 4, time.Second * 20
}

// Process sends batch of items
// TODO need implementation
func (c Client) Process(_ context.Context, batch Batch) error {
	fmt.Println("processed batch; len:", len(batch), "time:", time.Now())
	return nil
}
