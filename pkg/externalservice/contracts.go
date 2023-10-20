package externalservice

import "errors"

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}
