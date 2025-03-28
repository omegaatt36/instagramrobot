package api

import (
	"fmt"
	"regexp"

	"gopkg.in/telebot.v4"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/appmodule/providers"
	"github.com/omegaatt36/instagramrobot/appmodule/telegram"
	"github.com/omegaatt36/instagramrobot/appmodule/threads"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Controller handles incoming Telegram updates and orchestrates responses.
type Controller struct {
	bot *telebot.Bot // bot is the Telegram bot instance.
	// urlParser is a pre-compiled regex for finding URLs in text messages.
	urlParser *regexp.Regexp
}

// NewController creates a new Controller instance.
// NewController creates a new instance of the Controller.
// It initializes the URL parsing regex.
func NewController(b *telebot.Bot) *Controller {
	regex := `(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`
	r := regexp.MustCompile(regex)
	return &Controller{bot: b, urlParser: r}
}

// OnStart is the entry point for the incoming update
// OnStart handles the /start command, sending a welcome message in private chats.
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

// extractLinksFromString uses the pre-compiled regex to find all URLs in a given string.
func (x *Controller) extractLinksFromString(input string) []string {
	return x.urlParser.FindAllString(input, -1)
}

// OnText handles incoming text messages. It extracts URLs from the message text
// and processes them. If no links are found, it sends an error message in private chats.
func (x *Controller) OnText(c telebot.Context) error {
	links := x.extractLinksFromString(c.Message().Text)

	// Send proper error if text has no link inside
	if len(links) == 0 {
		if c.Chat().Type != telebot.ChatPrivate {
			return nil
		}

		logging.Error("onText.replyError: invalid command, please send the Instagram post link.")
		return x.replyError(c, "Invalid command,\nPlfmtease send the Instagram post link.")
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

// processLinks takes a list of URL strings and processes each one using LinkProcessor.
// It limits the number of links processed per message.
func (x *Controller) processLinks(links []string, m *telebot.Message) error {
	const maxLinksPerMessage = 3

	linkProcessor := providers.NewLinkProcessor(providers.NewLinkProcessorRequest{
		InstagramFetcher: instagram.NewExtractor(),
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

// replyError sends a formatted error message back to the user as a reply,
// using Markdown for formatting.
func (*Controller) replyError(c telebot.Context, text string) error {
	_, err := c.Bot().Reply(c.Message(), fmt.Sprintf("⚠️ *Oops, ERROR!*\n\n`%v`", text), telebot.ModeMarkdown)
	if err != nil {
		return fmt.Errorf("couldn't reply the Error, chat_id: %d, err: %w", c.Chat().ID, err)
	}

	return nil
}
