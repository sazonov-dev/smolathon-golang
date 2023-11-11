package main

import (
	"context"
	"log"

	"github.com/smolathon/internal/app"
	"github.com/smolathon/internal/config"
	"github.com/smolathon/pkg/logging"
)

func main() {
	log.Print("config")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)
	ctx := context.Background()
	a, err := app.NewApp(cfg, &logger, ctx)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Running Application")
	a.Run()
}
