package main

import (
	"gifka-bot/internal/app"
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("can't initialize logger", err)
	}

	app.Run(logger)
}
