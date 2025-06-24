package main

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/app/bot"
	"github.com/omegaatt36/instagramrobot/app/bot/config"
	"github.com/omegaatt36/instagramrobot/health"
	"github.com/omegaatt36/instagramrobot/logging"
)

var botToken string

// Main is the entry point of the application.
func Main(ctx context.Context, _ *cli.Command) error {
	logging.Init(!config.IsLocal())

	go health.StartServer()

	telegramBot, err := bot.NewTelegramBot(botToken)
	if err != nil {
		// Use Errorf instead of Fatalf to allow potential cleanup via After hooks
		logging.Errorf("couldn't create the Telegram bot: %v", err)
		return err
	}

	stopped := telegramBot.Start(ctx)

	select {
	case <-stopped:
		logging.Info("Bot stopped normally.")
	case <-ctx.Done():
		logging.Info("Shutdown signal received, waiting for bot to stop...")
		// Wait for bot to finish stopping after context cancellation
		<-stopped
		logging.Info("Bot stopped after signal.")
	}

	logging.Info("Shutting down")
	return nil
}

func main() {
	application := app.App{
		Main: Main,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bot-token",
				Sources:     cli.EnvVars("BOT_TOKEN"),
				Destination: &botToken,
				Required:    true,
			},
		},
	}

	application.Run()
}
