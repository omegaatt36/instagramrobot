package main

import (
	"github.com/omegaatt36/instagramrobot/src/config"
	"github.com/omegaatt36/instagramrobot/src/health"
	"github.com/omegaatt36/instagramrobot/src/helpers"
	"github.com/omegaatt36/instagramrobot/src/telegram"

	log "github.com/sirupsen/logrus"
)

func main() {
	go health.StartServer()

	helpers.RegisterLogger()
	config.Load()

	if config.IsDevelopment() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := telegram.Register(); err != nil {
		log.Fatalf("Couldn't register the Telegram bot: %v", err)
	}
	telegram.Start()
}
