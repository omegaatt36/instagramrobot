package main

import (
	"github.com/feelthecode/instagramrobot/src/config"
	"github.com/feelthecode/instagramrobot/src/helpers"
	"github.com/feelthecode/instagramrobot/src/telegram"
	"github.com/feelthecode/instagramrobot/src/telegram/cache"

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

	// Register Telegram cache
	cache.Register()

	if err := telegram.Register(); err != nil {
		log.Fatalf("Couldn't register the Telegram bot: %v", err)
	}
	telegram.Start()
}
