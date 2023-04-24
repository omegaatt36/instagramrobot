package telegram

import (
	"context"
	"time"

	"github.com/omegaatt36/instagramrobot/src/config"
	"github.com/omegaatt36/instagramrobot/src/telegram/commands"
	"github.com/omegaatt36/instagramrobot/src/telegram/events"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

var b *telebot.Bot

// Register will generate a fresh Telegram bot instance
// and registers it's handler logics
func Register(botToken string) error {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:   botToken,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		Verbose: config.IsDevelopment(),
	})
	if err != nil {
		log.Error("Couldn't create the Telegram bot instance")
		log.Fatal(err)
	}
	b = bot
	log.WithFields(log.Fields{
		"id":       b.Me.ID,
		"username": b.Me.Username,
		"title":    b.Me.FirstName,
	}).Info("Telegram bot instance created")

	registerCommands()

	// TODO: set bot commands

	return nil
}

func registerCommands() {
	// Commands
	start := commands.Start{B: b}
	b.Handle("/start", start.Handler)

	// Events
	text := events.NewTextHandler(b)
	b.Handle(telebot.OnText, text.Handler)
}

// Start brings bot into motion by consuming incoming updates
func Start(ctx context.Context) {
	log.Warn("Telegram bot starting")
	go func() {
		b.Start()
	}()
	<-ctx.Done()
	b.Stop()

}
