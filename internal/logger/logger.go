package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// InitLogger inits global logger
// TODO may need use interface instead of global logger
func InitLogger(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("logget init error: %w", err)
	}

	zerolog.SetGlobalLevel(lvl)

	return nil
}
