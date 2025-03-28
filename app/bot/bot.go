package bot

import (
	"context"
	"time"

	"gopkg.in/telebot.v4"

	"github.com/omegaatt36/instagramrobot/app/bot/api"
	"github.com/omegaatt36/instagramrobot/app/bot/config"
	"github.com/omegaatt36/instagramrobot/logging"
)

// b holds the global bot instance. Consider dependency injection for better testability.
var b *telebot.Bot

// Register will generate a fresh Telegram bot instance
// and registers it's handler logics
// Register creates and configures the global Telegram bot instance `b`.
// It sets up polling options and logging verbosity based on the environment.
// It fatally logs errors during bot initialization.
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

// registerCommands sets up the handlers for different Telegram events (commands, text messages).
func registerCommands() {
	x := api.NewController(b)

	b.Handle("/start", x.OnStart)
	b.Handle(telebot.OnText, x.OnText)
}

// Start brings bot into motion by consuming incoming updates
// Start initiates the bot's polling process in a separate goroutine.
// It returns a channel that will be closed when the bot stops.
// It listens to the provided context for cancellation signals to gracefully stop the bot.
func Start(ctx context.Context) <-chan struct{} {
	logging.Info("Telegram bot starting")
	closeChain := make(chan struct{})
	go b.Start()
	go func() {
		defer func() {
			logging.Info("Telegram bot stopped")
			closeChain <- struct{}{}
			close(closeChain)
		}()

		<-ctx.Done()
		b.Stop()
	}()

	return closeChain
}
