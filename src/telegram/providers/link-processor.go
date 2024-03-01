package providers

import (
	"net/url"

	"github.com/omegaatt36/instagramrobot/logging"
	"github.com/omegaatt36/instagramrobot/src/instagram"
	"github.com/omegaatt36/instagramrobot/src/telegram/utils"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

// LinkProcessor defines link processor.
type LinkProcessor struct {
	bot *telebot.Bot
	msg *telebot.Message
}

// NewLinkProcessor constructor
func NewLinkProcessor(bot *telebot.Bot, msg *telebot.Message) LinkProcessor {
	return LinkProcessor{
		bot: bot,
		msg: msg,
	}
}

// ProcessLink will process a single link
func (l *LinkProcessor) ProcessLink(link string) error {
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

	// Extract short code
	shortCode, err := instagram.ExtractShortCodeFromLink(url.Path)
	if err != nil {
		return err
	}

	logging.Infof("chatID(%d) shortcode(%s)", l.msg.Sender.ID, shortCode)

	// Process fetching the short code from Instagram
	response, err := instagram.GetPostWithCode(shortCode)
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
func (l *LinkProcessor) CheckIndexForSpam(index int) bool {
	// TODO: load from Config
	AllowedLinksPerMessage := 3
	if index == AllowedLinksPerMessage {
		logging.Errorf("can't process more than %c links per message", AllowedLinksPerMessage)
		return true
	}

	return false
}
