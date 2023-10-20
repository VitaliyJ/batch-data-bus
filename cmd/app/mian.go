package main

import (
	"log"

	"data-bus/internal/app"
	"data-bus/internal/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	if err := app.Run(cfg); err != nil {
		log.Fatalf("Application start error: %s", err)
	}
}
