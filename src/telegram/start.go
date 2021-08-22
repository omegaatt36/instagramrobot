package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (t *Bot) start(m *tb.Message) {
	// Ignore channels and groups
	t.b.Reply(m, "Hello!")
}
