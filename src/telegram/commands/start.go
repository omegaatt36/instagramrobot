package commands

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Start struct {
	B *tb.Bot // Bot instance
}

// The entry point for the incoming update
func (s *Start) Handler(m *tb.Message) {
	// Ignore channels and groups
	s.B.Reply(m, "Hello!")
}
