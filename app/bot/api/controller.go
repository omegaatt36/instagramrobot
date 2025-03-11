package api

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/telebot.v4"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/appmodule/providers"
	"github.com/omegaatt36/instagramrobot/appmodule/telegram"
	"github.com/omegaatt36/instagramrobot/appmodule/threads"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Controller is the main controller for the bot.
type Controller struct {
	bot       *telebot.Bot // Bot instance
	urlParser *regexp.Regexp
}

// NewController creates a new Controller instance.
func NewController(b *telebot.Bot) *Controller {
	regex := `(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`
	r := regexp.MustCompile(regex)
	return &Controller{bot: b, urlParser: r}
}

// OnStart is the entry point for the incoming update
func (*Controller) OnStart(c telebot.Context) error {
	// Ignore channels and groups
	if c.Chat().Type != telebot.ChatPrivate {
		return nil
	}

	if err := c.Reply("Hello! I'm instagram keeper, post some instagram public post/reels to me."); err != nil {
		return fmt.Errorf("couldn't sent the start command response: %w", err)
	}

	return nil
}

func (x *Controller) extractLinksFromString(input string) []string {
	return x.urlParser.FindAllString(input, -1)
}

func (x *Controller) OnText(c telebot.Context) error {
	links := x.extractLinksFromString(c.Message().Text)

	// Send proper error if text has no link inside
	if len(links) == 0 {
		if c.Chat().Type != telebot.ChatPrivate {
			return nil
		}

		err := errors.New("Invalid command,\nPlease send the Instagram post link.")
		logging.Error(fmt.Errorf("OnText.replyError: %w", err))
		return x.replyError(c, "Invalid command,\nPlease send the Instagram post link.")
	}

	if err := x.processLinks(links, c.Message()); err != nil {
		if c.Chat().Type != telebot.ChatPrivate {
			return nil
		}

		logging.Error(fmt.Errorf("OnText.processLinks: %w", err))
		return x.replyError(c, err.Error())
	}

	return nil
}

// Gets list of links from user message text
// and processes each one of them one by one.
func (x *Controller) processLinks(links []string, m *telebot.Message) error {
	const maxLinksPerMessage = 3

	linkProcessor := providers.NewLinkProcessor(providers.NewLinkProcessorRequest{
		InstagramFetcher: instagram.NewInstagramFetcher(),
		ThreadsFetcher:   threads.NewExtractor(),
		Sender:           telegram.NewMediaSender(x.bot, m),
	})

	for index, link := range links {
		if index == maxLinksPerMessage {
			logging.Errorf("can't process more than %c links per message", maxLinksPerMessage)
			break
		}

		if err := linkProcessor.ProcessLink(link); err != nil {
			logging.Error(fmt.Errorf("processLinks.ProcessLink: %w", err))
			continue // 繼續處理下一個 link
		}
	}

	return nil
}

func (*Controller) replyError(c telebot.Context, text string) error {
	_, err := c.Bot().Reply(c.Message(), fmt.Sprintf("⚠️ *Oops, ERROR!*\n\n`%v`", text), telebot.ModeMarkdown)
	if err != nil {
		return fmt.Errorf("couldn't reply the Error, chat_id: %d, err: %w", c.Chat().ID, err)
	}

	return nil
}
