package telegram

import (
	"fmt"
	"net/url"
	"regexp"

	tb "gopkg.in/tucnak/telebot.v2"

	log "github.com/sirupsen/logrus"
)

func (t *Bot) links(m *tb.Message) {
	// Ignore channels and groups
	r := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	matches := r.FindAllString(m.Text, -1)
	for index, link := range matches {
		// Protect user from sending bulk links in a single message.
		// TODO: load from Config
		AllowedLinksPerMessage := 3
		if index == AllowedLinksPerMessage {
			// TODO: error helper
			t.ReplyError(m, fmt.Sprintf("I can't process more than %c links per message.", AllowedLinksPerMessage))
			break
		}

		// Convert link to URL object
		url, err := url.ParseRequestURI(link)

		// Validate URL
		if err != nil || url == nil {
			// TODO: error helper
			t.ReplyError(m, fmt.Sprintf("I couldn't parse the [%v] link.", link))
			continue
		}

		// Validate HOST in the URL (only instagram.com is allowed)
		if url.Host != "instagram.com" {
			t.ReplyError(m, fmt.Sprintf("I can only process links from [instagram.com] not [%v].", url.Host))
			continue
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
		t.b.Reply(m, fmt.Sprintf("processing path %v", url.Path))
	}
}
