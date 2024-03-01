package telegram

import (
	"context"
	"time"

	"github.com/omegaatt36/instagramrobot/config"
	"github.com/omegaatt36/instagramrobot/logging"
	"github.com/omegaatt36/instagramrobot/src/telegram/commands"
	"github.com/omegaatt36/instagramrobot/src/telegram/events"
	"gopkg.in/telebot.v3"
)

var b *telebot.Bot

// Register will generate a fresh Telegram bot instance
// and registers it's handler logics
func Register(botToken string) error {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:   botToken,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		Verbose: !config.IsProduction(),
	})
	if err != nil {
		logging.Error("Couldn't create the Telegram bot instance")
		logging.Fatal(err)
	}
	b = bot

	logging.Info("Telegram bot instance created")
	logging.Infof("Bot info: id(%d) username(%s) title(%s)",
		b.Me.ID, b.Me.Username, b.Me.FirstName)

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
	logging.Info("Telegram bot starting")
	go func() {
		b.Start()
		<-ctx.Done()
		b.Stop()
	}()
}
