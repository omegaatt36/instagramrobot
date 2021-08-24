package providers

import (
	"github.com/feelthecode/instagramrobot/src/instagram/transform"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type MediaSender struct {
	bot   *tb.Bot
	msg   *tb.Message
	media transform.Media
}

const (
	caption = "🤖 Downloaded via @InstagramRobot"
)

// Send will start to process Media and eventually send it to the Telegram chat
func (m *MediaSender) Send() {
	// TODO: hit cache

	// Check if media has no child item
	if len(m.media.Items) == 0 {
		m.sendSingleMedia()
	} else {
		m.sendNestedMedia()
	}
}

func (m *MediaSender) sendSingleMedia() {
	var message *tb.Message

	if m.media.IsVideo {
		msg, err := m.bot.Send(m.msg.Chat, &tb.Video{
			File:    tb.FromURL(m.media.Url),
			Caption: caption,
		})
		if err != nil {
			log.Errorf("couldn't send the single video: %v", err)
			return
		}
		message = msg
	} else {
		msg, err := m.bot.Send(m.msg.Chat, &tb.Photo{
			File:    tb.FromURL(m.media.Url),
			Caption: caption,
		})
		if err != nil {
			log.Errorf("couldn't send the single photo: %v", err)
			return
		}
		message = msg
	}

	// TODO: cache the file_id

	m.sendCaption(message)
}

func (m *MediaSender) sendNestedMedia() {
	message, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromMedia())
	if err != nil {
		log.Errorf("couldn't send the single video: %v", err)
	}
	// TODO: cache album file_ids

	m.sendCaption(&message[0])
}

func (m *MediaSender) generateAlbumFromMedia() tb.Album {
	var album tb.Album

	for _, media := range m.media.Items {
		if m.media.IsVideo {
			album = append(album, &tb.Video{
				File: tb.FromURL(media.Url),
			})
		} else {
			album = append(album, &tb.Photo{
				File: tb.FromURL(media.Url),
			})
		}
	}

	return album
}

func (m *MediaSender) sendCaption(replyTo *tb.Message) {
	// If caption is empty, ignore sending it
	if m.media.Caption == "" {
		return
	}
	// TODO: cache caption
	// TODO: chunk caption if the length is above the Telegram limit
	_, _ = m.bot.Reply(replyTo, m.media.Caption)
}
