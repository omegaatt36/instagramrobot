package providers

import (
	"fmt"
	"net/url"

	"github.com/feelthecode/instagramrobot/src/telegram/utils"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type linkProcessor struct {
	bot *tb.Bot
	msg *tb.Message
}

func NewLinkProcessor(bot *tb.Bot, msg *tb.Message) linkProcessor {
	return linkProcessor{
		bot: bot,
		msg: msg,
	}
}

// Process a single link
func (l *linkProcessor) ProcessLink(link string) {
	// Convert link to URL object
	url, err := url.ParseRequestURI(link)

	// Validate URL
	if err != nil {
		utils.ReplyError(l.bot, l.msg, fmt.Sprintf("I couldn't parse the [%v] link.", link))
		return
	}

	// Validate HOST in the URL (only instagram.com is allowed)
	if url.Host != "instagram.com" {
		utils.ReplyError(l.bot, l.msg, fmt.Sprintf("I can only process links from [instagram.com] not [%v].", url.Host))
		return
	}

	log.Infof("link %+v", url.Path)

	// TODO: Validate URL path (only "/p/" or "/tv/" are acceptable)

	// TODO: Extract shortcode

	// TODO: Validate shortcode

	log.WithFields(log.Fields{
		"chat_id": l.msg.Sender.ID,
		"link":    url,
	}).Infof("Processing link")

	// TODO: process downloading the shortcode
	l.bot.Reply(l.msg, fmt.Sprintf("processing path %v", url.Path))
}

// Protect user from sending bulk links in a single message.
func (l *linkProcessor) CheckIndexForSpam(index int) bool {
	// TODO: load from Config
	AllowedLinksPerMessage := 3
	if index == AllowedLinksPerMessage {
		utils.ReplyError(l.bot, l.msg, fmt.Sprintf("can't process more than %c links per message", AllowedLinksPerMessage))
		return true
	}

	return false
}
