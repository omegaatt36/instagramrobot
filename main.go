package main

import (
	"context"

	"github.com/omegaatt36/instagramrobot/app"
	"github.com/omegaatt36/instagramrobot/src/config"
	"github.com/omegaatt36/instagramrobot/src/health"
	"github.com/omegaatt36/instagramrobot/src/helpers"
	"github.com/omegaatt36/instagramrobot/src/telegram"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Main(ctx context.Context) {
	go health.StartServer()

	helpers.RegisterLogger()

	if config.IsDevelopment() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := telegram.Register(config.BotToken()); err != nil {
		log.Fatalf("Couldn't register the Telegram bot: %v", err)
	}

	telegram.Start(ctx)
}

func main() {
	app := app.App{
		Main:  Main,
		Flags: []cli.Flag{},
	}

	app.Run()
}
