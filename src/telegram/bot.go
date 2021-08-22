package telegram

import (
	"fmt"

	"github.com/feelthecode/instagramrobot/src/config"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	b *tb.Bot
}

func (t *Bot) Register() error {
	b, err := tb.NewBot(tb.Settings{
		Token:   config.C.BOT_TOKEN,
		Poller:  t.getMiddleware(),
		Verbose: config.IsDevelopment(),
	})
	if err != nil {
		return fmt.Errorf("couldn't create the Telegram bot instance: %v", err)
	}
	t.b = b
	log.WithFields(log.Fields{
		"id":       b.Me.ID,
		"username": b.Me.Username,
		"title":    b.Me.FirstName,
	}).Info("Telegram bot instance created")

	t.registerCommands()

	// TODO: set bot commands

	return nil
}

func (t *Bot) registerCommands() {
	t.b.Handle("/start", t.start)
}

func (t *Bot) Start() {
	log.Warn("Telegram bot started")
	t.b.Start()
}
