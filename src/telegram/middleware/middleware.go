package middleware

import (
	"github.com/feelthecode/instagramrobot/src/telegram/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Middleware struct {
	B *tb.Bot
}

func (m *Middleware) Get(update *tb.Update) bool {
	if !update.Message.Private() {
		utils.ReplyError(m.B, update.Message, "I'm limited to private chats!")
		m.B.Leave(update.Message.Chat)
		return false
	}

	return true
}
