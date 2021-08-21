package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (t *Bot) start(m *tb.Message) {
	// Ignore channels and groups
	if !m.Private() {
		t.b.Reply(m, "I'm limited to private messages!")
		t.b.Leave(m.Chat)
		return
	}
	t.b.Reply(m, "Hello!")
}
