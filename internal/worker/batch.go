package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const batchWorkerTickDuration = 5 * time.Second

type BatchUseCase interface {
	SendNextBatch(ctx context.Context) error
}

type BatchWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
	uc     BatchUseCase
}

// NewBatchWorker returns new BatchWorker instance
func NewBatchWorker(uc BatchUseCase) *BatchWorker {
	return &BatchWorker{
		uc: uc,
	}
}

// Start runs batch worker
func (w *BatchWorker) Start() {
	w.ctx, w.cancel = context.WithCancel(context.Background())

	log.Info().Str("tick_duration", batchWorkerTickDuration.String()).Msg("starting batch worker")

	go func() {
		ticker := time.NewTicker(batchWorkerTickDuration)
		for {
			select {
			case <-ticker.C:
				go w.processBatch()
			case <-w.ctx.Done():
				log.Info().Msg("batch worker: context is done, stopping")
				return
			}
		}
	}()
}

// Close finalizes batch worker
func (w *BatchWorker) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}

func (w *BatchWorker) processBatch() {
	if err := w.uc.SendNextBatch(w.ctx); err != nil {
		log.Error().Time("send_time", time.Now()).Err(fmt.Errorf("BatchWorker.processBatch: batch sending error: %w", err)).Msg("batch sending error")
	}
}
