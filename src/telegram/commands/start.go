package commands

import (
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

// Start command
type Start struct {
	B *telebot.Bot // Bot instance
}

// Handler is the entry point for the incoming update
func (s *Start) Handler(c telebot.Context) error {
	// Ignore channels and groups
	if c.Chat().Type != telebot.ChatPrivate {
		return nil
	}

	if err := c.Reply("Hello!"); err != nil {
		return errors.Wrap(err, "Couldn't sent the start command response")
	}

	return nil
}
