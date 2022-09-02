package providers

import (
	"github.com/omegaatt36/instagramrobot/src/instagram/transform"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type MediaSender struct {
	bot   *tb.Bot
	msg   *tb.Message
	media transform.Media
}

const (
	caption = "🤖 Downloaded via @InstagramKeeperBot"
)

// Send will start to process Media and eventually send it to the Telegram chat
func (m *MediaSender) Send() {
	// Check if media has no child item
	if len(m.media.Items) == 0 {
		m.sendSingleMedia()
	} else {
		m.sendNestedMedia()
	}
}

func (m *MediaSender) sendSingleMedia() {
	if m.media.IsVideo {
		if _, err := m.bot.Send(m.msg.Chat, &tb.Video{
			File:    tb.FromURL(m.media.Url),
			Caption: caption,
		}); err != nil {
			log.Errorf("couldn't send the single video: %v", err)
			return
		}
		log.Debugf("Sent single video with shortcode [%v]", m.media.Shortcode)

	} else {
		if _, err := m.bot.Send(m.msg.Chat, &tb.Photo{
			File:    tb.FromURL(m.media.Url),
			Caption: caption,
		}); err != nil {
			log.Errorf("couldn't send the single photo: %v", err)
			return
		}
		log.Debugf("Sent single photo with shortcode [%v]", m.media.Shortcode)
	}

	m.sendCaption(m.media.Caption)
}

func (m *MediaSender) sendNestedMedia() {
	_, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromMedia())
	if err != nil {
		log.Errorf("couldn't send the single video: %v", err)
	}
	m.sendCaption(m.media.Caption)
}

func (m *MediaSender) generateAlbumFromMedia() tb.Album {
	var album tb.Album

	for _, media := range m.media.Items {
		if media.IsVideo {
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

func (m *MediaSender) sendCaption(caption string) {
	// If caption is empty, ignore sending it
	if m.media.Caption == "" {
		return
	}
	// TODO: chunk caption if the length is above the Telegram limit
	_, _ = m.bot.Reply(m.msg, caption)
}
