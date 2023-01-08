package utils

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

// Errors
var (
	ErrInvalidHost = errors.New("invalid host")
)

// ReplyError will sends the error to specific Telegram chat
// with a pre-defined structure
func ReplyError(c telebot.Context, text string) error {
	_, err := c.Bot().Reply(c.Message(), fmt.Sprintf("⚠️ *Oops, ERROR!*\n\n`%v`", text), telebot.ModeMarkdown)
	if err != nil {
		return errors.Wrapf(err, "Couldn't reply the Error, chat_id: %d", c.Chat().ID)
	}

	return nil

}
