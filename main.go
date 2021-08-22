package main

import (
	"github.com/feelthecode/instagramrobot/src/config"
	"github.com/feelthecode/instagramrobot/src/telegram"
	"github.com/feelthecode/instagramrobot/src/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	utils.RegisterLogger()
	config.Load()

	if config.IsDevelopment() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	bot := telegram.Bot{}
	bot.Register()
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
