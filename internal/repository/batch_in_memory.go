package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"data-bus/internal/entity"
)

var ErrInterrupted = errors.New("interrupted")

// node of list
type node struct {
	batch entity.Batch
	next  *node
}

// InMemory represents items storage
type InMemory struct {
	firstNode       *node
	lastNode        *node
	mu              *sync.Mutex
	lastRequestTime time.Time
}

// NewInMemoryRepository returns new InMemory repository instance
func NewInMemoryRepository() *InMemory {
	return &InMemory{
		mu: &sync.Mutex{},
	}
}

// GetFirstBatchItems returns first batch of items
func (r *InMemory) GetFirstBatchItems(ctx context.Context) (entity.Batch, error) {
	if err := ctx.Err(); err != nil {
		return nil, ErrInterrupted
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.firstNode == nil {
		return nil, nil
	}

	firstNode := r.firstNode
	r.firstNode = firstNode.next

	if r.firstNode == nil {
		r.lastNode = nil
	}

	return firstNode.batch, nil
}

// AddBatch adds batch of items to storage
func (r *InMemory) AddBatch(ctx context.Context, batch entity.Batch) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return ErrInterrupted
	}

	n := &node{batch: batch}
	if r.firstNode == nil {
		// set first node if doesn't exists
		r.firstNode = n
		r.lastNode = n
	} else {
		// update last node
		r.lastNode.next = n
		r.lastNode = n
	}

	return nil
}

// SaveLastRequestTime stores last request time to storage
func (r *InMemory) SaveLastRequestTime(ctx context.Context, t time.Time) error {
	if err := ctx.Err(); err != nil {
		return ErrInterrupted
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastRequestTime = t
	return nil
}

// GetLastRequestTime returns last request time
func (r *InMemory) GetLastRequestTime(ctx context.Context) (time.Time, error) {
	if err := ctx.Err(); err != nil {
		return time.Time{}, ErrInterrupted
	}

	return r.lastRequestTime, nil
}
