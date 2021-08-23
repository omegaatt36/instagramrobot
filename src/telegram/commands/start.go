package commands

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Start struct {
	B *tb.Bot
}

func (s *Start) Get(m *tb.Message) {
	// Ignore channels and groups
	s.B.Reply(m, "Hello!")
}
