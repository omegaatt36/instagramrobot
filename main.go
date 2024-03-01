package main

import (
	"context"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/config"
	"github.com/omegaatt36/instagramrobot/logging"
	"github.com/omegaatt36/instagramrobot/src/health"
	"github.com/omegaatt36/instagramrobot/src/telegram"
	"github.com/urfave/cli/v2"
)

// Main is the entry point of the application.
func Main(ctx context.Context) {
	logging.Init()

	go health.StartServer()

	if err := telegram.Register(config.BotToken()); err != nil {
		logging.Fatalf("couldn't register the Telegram bot: %v", err)
	}

	stopped := telegram.Start(ctx)

	<-stopped
	<-ctx.Done()
	logging.Info("Shutting down")
}

func main() {
	app := app.App{
		Main:  Main,
		Flags: []cli.Flag{},
	}

	app.Run()
}
