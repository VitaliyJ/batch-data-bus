package app

import (
	"errors"
	"fmt"

	"data-bus/internal/config"
	"data-bus/internal/logger"
	"data-bus/internal/repository"
	"data-bus/internal/server/http"
	"data-bus/internal/usecase"
	"data-bus/internal/worker"
	"data-bus/pkg/externalservice"
	"github.com/rs/zerolog/log"
)

// Run starts application
func Run(cfg *config.Config) error {
	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		return fmt.Errorf("application run: init logger error: %w", err)
	}

	// init repository
	inMemoryRepo := repository.NewInMemoryRepository()
	// init external service client
	externalService := externalservice.NewClient()

	// init use cases
	batchUseCase := usecase.NewBatchUseCase(
		inMemoryRepo,
		externalService,
	)

	// init workers
	batchWorker := worker.NewBatchWorker(batchUseCase)

	// run workers
	batchWorker.Start()
	defer func() {
		if err := batchWorker.Close(); err != nil {
			log.Error().Err(fmt.Errorf("application run: %w", err)).Msg("batch worker closing error")
		}
	}()

	httpServer := http.NewServer(
		&http.Config{HTTPPort: cfg.HTTP.Port},
		batchUseCase,
	)
	httpServer.Start()

	return errors.New("application stopped")
}
