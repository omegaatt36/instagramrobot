package events

import (
	"fmt"
	"net/url"

	"github.com/feelthecode/instagramrobot/src/helpers"
	"github.com/feelthecode/instagramrobot/src/telegram/utils"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Text handler
type TextHandler struct {
	B *tb.Bot // Bot instance
}

// The entry point for the incoming update
func (l *TextHandler) Handler(m *tb.Message) {
	links := helpers.ExtractLinksFromString(m.Text)
	l.processLinks(links, m)
}

// Gets list of links from user message text
// and processes each one of them one by one.
func (l *TextHandler) processLinks(links []string, m *tb.Message) {
	for index, link := range links {
		if spam := l.checkIndexForSpam(index, m); spam {
			break
		}

		l.processLink(link, m)
	}
}

// Process a single link
func (l *TextHandler) processLink(link string, m *tb.Message) {
	// Convert link to URL object
	url, err := url.ParseRequestURI(link)

	// Validate URL
	if err != nil {
		utils.ReplyError(l.B, m, fmt.Sprintf("I couldn't parse the [%v] link.", link))
		return
	}

	// Validate HOST in the URL (only instagram.com is allowed)
	if url.Host != "instagram.com" {
		utils.ReplyError(l.B, m, fmt.Sprintf("I can only process links from [instagram.com] not [%v].", url.Host))
		return
	}

	log.Infof("link %+v", url.Path)

	// TODO: Validate URL path (only "/p/" or "/tv/" are acceptable)

	// TODO: Extract shortcode

	// TODO: Validate shortcode

	log.WithFields(log.Fields{
		"chat_id": m.Sender.ID,
		"link":    url,
	}).Infof("Processing link")

	// TODO: process downloading the shortcode
	l.B.Reply(m, fmt.Sprintf("processing path %v", url.Path))
}

// Protect user from sending bulk links in a single message.
func (l *TextHandler) checkIndexForSpam(index int, m *tb.Message) bool {
	// TODO: load from Config
	AllowedLinksPerMessage := 3
	if index == AllowedLinksPerMessage {
		utils.ReplyError(l.B, m, fmt.Sprintf("can't process more than %c links per message", AllowedLinksPerMessage))
		return true
	}
	return false
}
