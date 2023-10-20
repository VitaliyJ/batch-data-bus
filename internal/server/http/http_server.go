package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"data-bus/pkg/externalservice"
	"data-bus/pkg/httpserver"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HTTPPort string
}

type Server struct {
	cfg *Config
	uc  BatchUseCase
}

type BatchUseCase interface {
	Process(ctx context.Context, items []externalservice.Item) error
}

// NewServer returns new Server instance
func NewServer(
	cfg *Config,
	uc BatchUseCase,
) *Server {
	return &Server{
		cfg: cfg,
		uc:  uc,
	}
}

// Start starts the server
func (s *Server) Start() {
	httpServer := httpserver.New(s.newRouter(), httpserver.Port(s.cfg.HTTPPort))
	log.Info().Str("port", s.cfg.HTTPPort).Msg("http server started")

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Str("signal", s.String()).Msg("application was interrupted")
	case err := <-httpServer.Notify():
		log.Error().Err(fmt.Errorf("httpServer.Notify: %w", err)).Msg("http server error")
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		log.Error().Err(fmt.Errorf("httpServer.Shutdown: %w", err)).Msg("http server shutdown error")
	}
}
