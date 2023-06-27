package main

import (
	"fmt"
	"os"

	"github.com/x1unix/tg-stargazers-bot/internal/app"
	"github.com/x1unix/tg-stargazers-bot/internal/config"
	"go.uber.org/zap"
)

func main() {
	ctx, cancelFn := config.NewApplicationContext()
	defer cancelFn()

	svc, err := app.BuildService()
	if err != nil {
		die("Error:", err)
	}

	if err := svc.Start(ctx); err != nil {
		zap.L().Fatal("failed to start service", zap.Error(err))
	}
}

func die(args ...any) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}
