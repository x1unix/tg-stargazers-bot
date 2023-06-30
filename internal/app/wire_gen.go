// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"github.com/x1unix/tg-stargazers-bot/internal/handlers/chat"
	"github.com/x1unix/tg-stargazers-bot/internal/repository"
	"github.com/x1unix/tg-stargazers-bot/internal/services"
	"github.com/x1unix/tg-stargazers-bot/internal/services/auth"
	"github.com/x1unix/tg-stargazers-bot/internal/services/bot"
	"github.com/x1unix/tg-stargazers-bot/internal/services/feedback"
)

// Injectors from wire.go:

// BuildService constructs service instance with app dependencies using Wire.
func BuildService() (*Service, error) {
	commandFlags := config.ReadCommandFlags()
	configConfig, err := config.FromEnv(commandFlags)
	if err != nil {
		return nil, err
	}
	logger, err := provideLogger(configConfig)
	if err != nil {
		return nil, err
	}
	botConfig := provideBotConfig(configConfig)
	urlBuilder := provideURLBuilder(configConfig)
	resolvedAuthConfig, err := provideAuthConfig(configConfig)
	if err != nil {
		return nil, err
	}
	cmdable, err := provideRedis(configConfig)
	if err != nil {
		return nil, err
	}
	tokenRepository := repository.NewTokenRepository(cmdable)
	service := auth.NewService(logger, resolvedAuthConfig, tokenRepository)
	gitHubTokensRepository := repository.NewGitHubTokensRepository(cmdable)
	hookRepository := repository.NewHookRepository(cmdable)
	gitHubService := provideGitHubService(configConfig, service, gitHubTokensRepository, hookRepository, urlBuilder)
	handlers := chat.NewHandlers(logger, urlBuilder, gitHubService, service)
	eventHandler := services.NewEventRouter(handlers)
	botService, err := bot.NewService(logger, botConfig, eventHandler)
	if err != nil {
		return nil, err
	}
	notificationsService := feedback.NewNotificationsService(botService)
	appService := NewService(logger, configConfig, botService, service, gitHubService, notificationsService)
	return appService, nil
}
