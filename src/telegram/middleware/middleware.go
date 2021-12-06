package middleware

import (
	"github.com/feelthecode/instagramrobot/src/telegram/utils"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Middleware struct
type Middleware struct {
	B *tb.Bot
}

// GetFilter will return the middleware filter function
func (m *Middleware) GetFilter(update *tb.Update) bool {
	if !update.Message.Private() {
		utils.ReplyError(m.B, update.Message, "I'm limited to private chats!")
		if err := m.B.Leave(update.Message.Chat); err != nil {
			log.WithField("chat_id", update.Message.Chat.ID).Errorf("Couldn't leave the chat: %v", err)
		}
		return false
	}

	return true
}
