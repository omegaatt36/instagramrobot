package utils

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func ReplyError(b *tb.Bot, to *tb.Message, text string) (*tb.Message, error) {
	return b.Reply(to, fmt.Sprintf("⚠️ *Oops, ERROR!*\n\n`%v`", text), tb.ModeMarkdown)
}
