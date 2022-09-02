package telegram

import (
	"time"

	"github.com/omegaatt36/instagramrobot/src/config"
	"github.com/omegaatt36/instagramrobot/src/telegram/commands"
	"github.com/omegaatt36/instagramrobot/src/telegram/events"
	"github.com/omegaatt36/instagramrobot/src/telegram/middleware"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

var b *tb.Bot

// Register will generate a fresh Telegram bot instance
// and registers it's handler logics
func Register() error {
	poller := &tb.LongPoller{
		Timeout:        15 * time.Second,
		AllowedUpdates: []string{"message"},
	}

	// Generate middleware
	m := middleware.Middleware{
		B: b,
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:   viper.GetString("BOT_TOKEN"),
		Poller:  tb.NewMiddlewarePoller(poller, m.GetFilter),
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
	text := events.TextHandler(b)
	b.Handle(tb.OnText, text.Handler)
}

// Start brings bot into motion by consuming incoming updates
func Start() {
	log.Warn("Telegram bot starting")
	b.Start()
}
