package usecase

import (
	"context"
	"errors"
	"time"

	"data-bus/internal/entity"
	"data-bus/pkg/externalservice"
	"github.com/rs/zerolog/log"
)

type BatchUseCase struct {
	repo    BatchRepo
	service Service
}

// NewBatchUseCase returns new BatchUseCase instance
func NewBatchUseCase(repo BatchRepo, service Service) *BatchUseCase {
	return &BatchUseCase{
		repo:    repo,
		service: service,
	}
}

// Process handles items
func (uc *BatchUseCase) Process(ctx context.Context, items []externalservice.Item) error {
	const logTrace = "BatchUseCase.Process"

	if len(items) == 0 {
		return nil
	}

	// get service limits; todo: may need cache
	n, _ := uc.service.GetLimits()

	// compose batches of items
	uLen := uint64(len(items))
	for i, j := uint64(0), n; i < uLen; i, j = j, j+n {
		if j > uLen {
			j = uLen
		}

		batch := make(entity.Batch, 0, n)
		batch = append(batch, items[i:j]...)
		if err := uc.repo.AddBatch(ctx, batch); err != nil {
			log.Error().Str("log_trace", logTrace).Err(err).Msg("items batch adding error")
			return errors.New("items batch adding error")
		}
	}

	return nil
}

// SendNextBatch sends next in line batch to the processing service
func (uc *BatchUseCase) SendNextBatch(ctx context.Context) error {
	const logTrace = "BatchUseCase.SendNextBatch"

	// get service limits; todo: may need cache
	_, p := uc.service.GetLimits()
	lastReqTime, err := uc.repo.GetLastRequestTime(ctx)
	if err != nil {
		log.Error().Str("log_trace", logTrace).Err(err).Msg("last request time getting error")
		return errors.New("last request time getting error")
	}

	// check if time is within the limits
	lastReqTime = lastReqTime.Add(p)
	if lastReqTime.After(time.Now()) {
		return nil
	}

	// check if batch for processing exists
	batch, err := uc.repo.GetFirstBatchItems(ctx)
	if err != nil {
		log.Error().Str("log_trace", logTrace).Err(err).Msg("batch getting error")
		return errors.New("batch getting error")
	}
	if len(batch) == 0 {
		return nil
	}

	if err := uc.service.Process(ctx, externalservice.Batch(batch)); err != nil {
		log.Error().Str("log_trace", logTrace).Err(err).Msg("batch processing error")
		if addBatchErr := uc.repo.AddBatch(ctx, batch); addBatchErr != nil {
			log.Error().Str("log_trace", logTrace).Err(err).Msg("batch adding error")
		}
		return errors.New("batch adding error")
	}

	if err := uc.repo.SaveLastRequestTime(ctx, time.Now()); err != nil {
		log.Error().Str("log_trace", logTrace).Err(err).Msg("last request time updating error")
		return errors.New("last request time updating error")
	}

	return nil
}
