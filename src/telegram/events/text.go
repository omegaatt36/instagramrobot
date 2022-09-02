package events

import (
	"github.com/omegaatt36/instagramrobot/src/helpers"
	"github.com/omegaatt36/instagramrobot/src/telegram/providers"
	"github.com/omegaatt36/instagramrobot/src/telegram/utils"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TextHandler constructor
func TextHandler(bot *tb.Bot) textHandler {
	return textHandler{
		bot: bot,
	}
}

type textHandler struct {
	bot *tb.Bot // Bot instance
}

// Handler is the entry point for the incoming update
func (l *textHandler) Handler(m *tb.Message) {
	links := helpers.ExtractLinksFromString(m.Text)
	// Send proper error if text has no link inside
	if len(links) == 0 {
		utils.ReplyError(l.bot, m, "Invalid command,\nPlease send the Instagram post link.")
		return
	}
	l.processLinks(links, m)
}

// Gets list of links from user message text
// and processes each one of them one by one.
func (l *textHandler) processLinks(links []string, m *tb.Message) {
	for index, link := range links {
		linkProcessor := providers.NewLinkProcessor(l.bot, m)

		if spam := linkProcessor.CheckIndexForSpam(index); spam {
			break
		}

		linkProcessor.ProcessLink(link)
	}
}
