package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/app/bot"
	"github.com/omegaatt36/instagramrobot/app/bot/config"
	"github.com/omegaatt36/instagramrobot/health"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Main is the entry point of the application.
func Main(ctx context.Context) {
	logging.Init(!config.IsLocal())

	go health.StartServer()

	if err := bot.Register(config.BotToken()); err != nil {
		logging.Fatalf("couldn't register the Telegram bot: %v", err)
	}

	// Wait for the bot to stop and for the context cancellation signal.
	stopped := bot.Start(ctx)

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
