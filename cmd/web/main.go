package main

import (
	"context"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/app/web"
	"github.com/omegaatt36/instagramrobot/health"
	"github.com/omegaatt36/instagramrobot/logging"
	"github.com/urfave/cli/v2"
)

// Main is the entry point of the application.
func Main(ctx context.Context) {
	logging.Init(false)

	go health.StartServer()

	// if err := bot.Register(config.BotToken()); err != nil {
	// 	logging.Fatalf("couldn't register the Telegram bot: %v", err)
	// }

	stopped := web.NewServer().Start(ctx)

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
