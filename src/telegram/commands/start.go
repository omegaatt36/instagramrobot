package commands

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type StartCommand struct {
	B *tb.Bot // Bot instance
}

// The entry point for the incoming update
func (s *StartCommand) Handler(m *tb.Message) {
	// Ignore channels and groups
	s.B.Reply(m, "Hello!")
}
