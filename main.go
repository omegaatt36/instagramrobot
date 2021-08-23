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

	// ig := instagram.API{}
	// code := "CSft2G5pFgr"
	// response, err := ig.GetPostWithCode(code)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	return
	// }

	// fmt.Printf("%+v\n", response)
}
