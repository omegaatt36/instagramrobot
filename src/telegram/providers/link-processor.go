package providers

import (
	"fmt"
	"net/url"

	"github.com/feelthecode/instagramrobot/src/instagram"
	"github.com/feelthecode/instagramrobot/src/telegram/utils"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type linkProcessor struct {
	bot *tb.Bot
	msg *tb.Message
}

// NewLinkProcessor constructor
func NewLinkProcessor(bot *tb.Bot, msg *tb.Message) linkProcessor {
	return linkProcessor{
		bot: bot,
		msg: msg,
	}
}

// ProcessLink will process a single link
func (l *linkProcessor) ProcessLink(link string) {
	// Convert link to URL object
	url, err := url.ParseRequestURI(link)

	// Validate URL
	if err != nil {
		utils.ReplyError(l.bot, l.msg, fmt.Sprintf("I couldn't parse the [%v] link.", link))
		return
	}

	// Validate HOST in the URL (only instagram.com is allowed)
	if url.Host != "instagram.com" && url.Host != "www.instagram.com" {
		utils.ReplyError(l.bot, l.msg, fmt.Sprintf("I can only process links from [instagram.com] not [%v].", url.Host))
		return
	}

	// Extract shortcode
	shortcode, err := instagram.ExtractShortcodeFromLink(url.Path)
	if err != nil {
		utils.ReplyError(l.bot, l.msg, err.Error())
		return
	}

	log.WithFields(log.Fields{
		"chat_id":   l.msg.Sender.ID,
		"shortcode": shortcode,
	}).Infof("Processing link")

	// Process fetching the shortcode from Instagram
	response, err := instagram.GetPostWithCode(shortcode)
	if err != nil {
		utils.ReplyError(l.bot, l.msg, err.Error())
		return
	}

	mediaSender := MediaSender{
		bot:   l.bot,
		msg:   l.msg,
		media: response,
	}
	mediaSender.Send()
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
