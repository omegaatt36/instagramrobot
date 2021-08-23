package telegram

import (
	"time"

	"github.com/feelthecode/instagramrobot/src/config"
	"github.com/feelthecode/instagramrobot/src/telegram/commands"
	"github.com/feelthecode/instagramrobot/src/telegram/events"
	"github.com/feelthecode/instagramrobot/src/telegram/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	b *tb.Bot
}

func (t *Bot) Register() error {
	poller := &tb.LongPoller{
		Timeout:        15 * time.Second,
		AllowedUpdates: []string{"message"},
	}
	// Generate middleware
	m := middleware.Middleware{
		B: t.b,
	}

	b, err := tb.NewBot(tb.Settings{
		Token:   viper.GetString("BOT_TOKEN"),
		Poller:  tb.NewMiddlewarePoller(poller, m.Get),
		Verbose: config.IsDevelopment(),
	})
	if err != nil {
		log.Error("Couldn't create the Telegram bot instance")
		log.Fatal(err)
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
	// Commands
	start := commands.StartCommand{B: t.b}
	t.b.Handle("/start", start.Handler)

	// Events
	text := events.TextHandler(t.b)
	t.b.Handle(tb.OnText, text.Handler)
}

func (t *Bot) Start() {
	log.Warn("Telegram bot starting")
	t.b.Start()
}
