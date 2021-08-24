package main

import (
	"github.com/feelthecode/instagramrobot/src/config"
	"github.com/feelthecode/instagramrobot/src/helpers"
	"github.com/feelthecode/instagramrobot/src/telegram"

	log "github.com/sirupsen/logrus"
)

func main() {
	helpers.RegisterLogger()
	config.Load()

	if config.IsDevelopment() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	bot := telegram.Bot{}
	if err := bot.Register(); err != nil {
		log.Fatalf("Couldn't register the Telegram bot: %v", err)
	}
	bot.Start()
}
