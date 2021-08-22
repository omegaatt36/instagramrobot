package telegram

import (
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (t *Bot) getMiddleware() *tb.MiddlewarePoller {
	poller := &tb.LongPoller{
		Timeout:        15 * time.Second,
		AllowedUpdates: []string{"message"},
	}

	return tb.NewMiddlewarePoller(poller, func(update *tb.Update) bool {
		if !update.Message.Private() {
			// TODO: add keyboard to start the bot private chat
			t.b.Reply(update.Message, "I'm limited to private chats!")
			t.b.Leave(update.Message.Chat)
			return false
		}

		for _, entity := range update.Message.Entities {
			if entity.Type == tb.EntityURL {
				link := substr(update.Message.Text, entity.Offset, entity.Length)

				// validate extracted link
				validLink, err := url.ParseRequestURI(link)
				if err != nil {
					t.b.Reply(update.Message, fmt.Sprintf("Invalid link '%v'", link))
				}

				log.Infof("valid link: %v", validLink)
				// TODO: process downloading the link
			}
		}

		return true
	})
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
