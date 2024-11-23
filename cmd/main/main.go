package main

import (
	"github.com/un1uckyyy/go-telegram-oidc/internal/config"
	"github.com/un1uckyyy/go-telegram-oidc/internal/tg"
	"github.com/un1uckyyy/go-telegram-oidc/pkg/logger"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		logger.ErrorLogger.Fatalf("failed load config: %v", err)
	}

	telegramService, err := tg.NewService()
	if err != nil {
		logger.ErrorLogger.Fatalf("failed init tg bot: %v", err)
	}

	logger.InfoLogger.Println("bot is starting...")
	telegramService.Start()
}
