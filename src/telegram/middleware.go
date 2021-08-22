package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (t *Bot) middleware(update *tb.Update) bool {
	if !update.Message.Private() {
		t.ReplyError(update.Message, "I'm limited to private chats!")
		t.b.Leave(update.Message.Chat)
		return false
	}

	return true
}
