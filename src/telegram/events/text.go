package events

import (
	"github.com/omegaatt36/instagramrobot/logging"
	"github.com/omegaatt36/instagramrobot/src/helpers"
	"github.com/omegaatt36/instagramrobot/src/telegram/providers"
	"github.com/omegaatt36/instagramrobot/src/telegram/utils"
	"gopkg.in/telebot.v3"
)

// NewTextHandler constructor
func NewTextHandler(bot *telebot.Bot) TextHandler {
	return TextHandler{
		bot: bot,
	}
}

type TextHandler struct {
	bot *telebot.Bot // Bot instance
}

// Handler is the entry point for the incoming update
func (l *TextHandler) Handler(c telebot.Context) error {
	links := helpers.ExtractLinksFromString(c.Message().Text)
	// Send proper error if text has no link inside
	if len(links) == 0 {
		if c.Chat().Type != telebot.ChatPrivate {
			return nil
		}

		logging.Error("Invalid command,\nPlease send the Instagram post link.")
		return utils.ReplyError(c, "Invalid command,\nPlease send the Instagram post link.")
	}

	if err := l.processLinks(links, c.Message()); err != nil {
		if c.Chat().Type != telebot.ChatPrivate {
			return nil
		}

		logging.Error(err)
		return utils.ReplyError(c, err.Error())
	}

	return nil
}

// Gets list of links from user message text
// and processes each one of them one by one.
func (l *TextHandler) processLinks(links []string, m *telebot.Message) error {
	for index, link := range links {
		linkProcessor := providers.NewLinkProcessor(l.bot, m)

		if spam := linkProcessor.CheckIndexForSpam(index); spam {
			break
		}

		if err := linkProcessor.ProcessLink(link); err != nil {
			return err
		}
	}
	return nil
}
