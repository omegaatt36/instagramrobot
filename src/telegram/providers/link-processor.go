package providers

import (
	"net/url"

	"github.com/omegaatt36/instagramrobot/src/instagram"
	"github.com/omegaatt36/instagramrobot/src/telegram/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type linkProcessor struct {
	bot *telebot.Bot
	msg *telebot.Message
}

// NewLinkProcessor constructor
func NewLinkProcessor(bot *telebot.Bot, msg *telebot.Message) linkProcessor {
	return linkProcessor{
		bot: bot,
		msg: msg,
	}
}

// ProcessLink will process a single link
func (l *linkProcessor) ProcessLink(link string) error {
	// Convert link to URL object
	url, err := url.ParseRequestURI(link)

	// Validate URL
	if err != nil {
		return errors.Wrapf(err, "I couldn't parse the [%v] link.", link)
	}

	// Validate HOST in the URL (only instagram.com is allowed)
	if url.Host != "instagram.com" && url.Host != "www.instagram.com" {
		return errors.Wrapf(utils.ErrInvalidHost, "can only process links from [instagram.com] not [%s]", url.Host)
	}

	// Extract shortcode
	shortcode, err := instagram.ExtractShortcodeFromLink(url.Path)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"chat_id":   l.msg.Sender.ID,
		"shortcode": shortcode,
	}).Infof("Processing link")

	// Process fetching the shortcode from Instagram
	response, err := instagram.GetPostWithCode(shortcode)
	if err != nil {
		return err
	}

	mediaSender := MediaSender{
		bot:   l.bot,
		msg:   l.msg,
		media: response,
	}

	return mediaSender.Send()
}

// Protect user from sending bulk links in a single message.
func (l *linkProcessor) CheckIndexForSpam(index int) bool {
	// TODO: load from Config
	AllowedLinksPerMessage := 3
	if index == AllowedLinksPerMessage {
		log.Errorf("can't process more than %c links per message", AllowedLinksPerMessage)
		return true
	}

	return false
}
