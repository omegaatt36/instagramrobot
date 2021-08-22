package telegram

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (t *Bot) ReplyError(to *tb.Message, text string) (*tb.Message, error) {
	return t.b.Reply(to, fmt.Sprintf("⚠️ *Oops, ERROR!*\n\n`%v`", text), tb.ModeMarkdown)
}
