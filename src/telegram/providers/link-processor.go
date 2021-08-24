package providers

import (
	"fmt"
	"net/url"
	"regexp"

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

	log.Infof("link %+v", url.Path)

	// TODO: move logic to the instagram package
	// Validate URL path (only "/p/" or "/tv/" are acceptable)
	// Validate shortcode
	re := regexp.MustCompile(`(p|tv|reel)\/([A-Za-z0-9-_]+)`)
	values := re.FindStringSubmatch(url.Path)
	log.Infof("%+v %v", values, len(values))
	if len(values) != 3 {
		utils.ReplyError(l.bot, l.msg, "The link structure is invalid")
		return
	}
	// Extract shortcode
	shortcode := values[2]

	log.WithFields(log.Fields{
		"chat_id":   l.msg.Sender.ID,
		"shortcode": shortcode,
	}).Infof("Processing link")

	// Process fetching the shortcode from Instagram
	ig := instagram.API{}
	response, err := ig.GetPostWithCode(shortcode)
	if err != nil {
		utils.ReplyError(l.bot, l.msg, err.Error())
		return
	}

	// TODO: Process response object and send proper content type
	_, _ = l.bot.Reply(l.msg, response.Url)
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
