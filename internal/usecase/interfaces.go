package usecase

import (
	"context"
	"time"

	"data-bus/internal/entity"
	"data-bus/pkg/externalservice"
)

//go:generate go run github.com/vektra/mockery/v2@v2.36.0 --name=BatchRepo
type BatchRepo interface {
	AddBatch(ctx context.Context, batch entity.Batch) error
	GetFirstBatchItems(ctx context.Context) (entity.Batch, error)
	SaveLastRequestTime(ctx context.Context, t time.Time) error
	GetLastRequestTime(ctx context.Context) (time.Time, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.0 --name=Service
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch externalservice.Batch) error
}
