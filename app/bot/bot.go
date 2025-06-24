package bot

import (
	"context"
	"time"

	"gopkg.in/telebot.v4"

	"github.com/omegaatt36/instagramrobot/app/bot/api"
	"github.com/omegaatt36/instagramrobot/app/bot/config"
	"github.com/omegaatt36/instagramrobot/logging"
)

// TelegramBot wraps the telebot.Bot instance and provides lifecycle management.
type TelegramBot struct {
	bot *telebot.Bot
}

// NewTelegramBot creates and configures a new TelegramBot instance.
// It sets up polling options and logging verbosity based on the environment.
func NewTelegramBot(botToken string) (*TelegramBot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:   botToken,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		Verbose: !config.IsProduction(),
	})
	if err != nil {
		logging.Error("Couldn't create the Telegram bot instance")
		return nil, err
	}

	telegramBot := &TelegramBot{bot: bot}

	logging.Info("Telegram bot instance created")
	logging.Infof("Bot info: id(%d) username(%s) title(%s)",
		bot.Me.ID, bot.Me.Username, bot.Me.FirstName)

	telegramBot.registerCommands()

	// TODO: set bot commands

	return telegramBot, nil
}

// registerCommands sets up the handlers for different Telegram events (commands, text messages).
func (tb *TelegramBot) registerCommands() {
	x := api.NewController(tb.bot)

	tb.bot.Handle("/start", x.OnStart)
	tb.bot.Handle(telebot.OnText, x.OnText)
}

// Start brings bot into motion by consuming incoming updates.
// It initiates the bot's polling process in a separate goroutine.
// It returns a channel that will be closed when the bot stops.
// It listens to the provided context for cancellation signals to gracefully stop the bot.
func (tb *TelegramBot) Start(ctx context.Context) <-chan struct{} {
	logging.Info("Telegram bot starting")
	closeChain := make(chan struct{})
	go tb.bot.Start()
	go func() {
		defer func() {
			logging.Info("Telegram bot stopped")
			closeChain <- struct{}{}
			close(closeChain)
		}()

		<-ctx.Done()
		tb.bot.Stop()
	}()

	return closeChain
}
