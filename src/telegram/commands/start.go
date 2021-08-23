package commands

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type StartCommand struct {
	B *tb.Bot // Bot instance
}

// The entry point for the incoming update
func (s *StartCommand) Handler(m *tb.Message) {
	// Ignore channels and groups
	_, err := s.B.Reply(m, "Hello!")
	if err != nil {
		log.Errorf("Couldn't sent the start command response: %v", err)
	}
}
