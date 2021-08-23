package utils

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

func ReplyError(b *tb.Bot, to *tb.Message, text string) {
	_, err := b.Reply(to, fmt.Sprintf("⚠️ *Oops, ERROR!*\n\n`%v`", text), tb.ModeMarkdown)
	if err != nil {
		log.WithFields(log.Fields{
			"chat_id": to.Chat.ID,
		}).Errorf("Couldn't reply the Error: %v", err)
	}
}
