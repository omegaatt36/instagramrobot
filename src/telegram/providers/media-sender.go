package providers

import (
	"github.com/omegaatt36/instagramrobot/src/instagram/transform"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"

	log "github.com/sirupsen/logrus"
)

type MediaSender struct {
	bot   *telebot.Bot
	msg   *telebot.Message
	media transform.Media
}

const (
	caption = "🤖 Downloaded via @InstagramKeeperBot"
)

// Send will start to process Media and eventually send it to the Telegram chat
func (m *MediaSender) Send() error {
	var err error
	// Check if media has no child item
	if len(m.media.Items) == 0 {
		err = m.sendSingleMedia()
	} else {
		err = m.sendNestedMedia()
	}

	return err
}

func (m *MediaSender) sendSingleMedia() error {
	if m.media.IsVideo {
		if _, err := m.bot.Send(m.msg.Chat, &telebot.Video{
			File:    telebot.FromURL(m.media.Url),
			Caption: caption,
		}); err != nil {
			return errors.Wrap(err, "couldn't send the single video")
		}

		log.Debugf("Sent single video with short code [%v]", m.media.Shortcode)
	} else {
		if _, err := m.bot.Send(m.msg.Chat, &telebot.Photo{
			File:    telebot.FromURL(m.media.Url),
			Caption: caption,
		}); err != nil {
			return errors.Wrap(err, "couldn't send the single photo")
		}

		log.Debugf("Sent single photo with short code [%v]", m.media.Shortcode)
	}

	return m.sendCaption(m.media.Caption)
}

func (m *MediaSender) sendNestedMedia() error {
	_, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromMedia())
	if err != nil {
		return errors.Wrap(err, "couldn't send the nested media")
	}
	return m.sendCaption(m.media.Caption)
}

func (m *MediaSender) generateAlbumFromMedia() telebot.Album {
	var album telebot.Album

	for _, media := range m.media.Items {
		if media.IsVideo {
			album = append(album, &telebot.Video{
				File: telebot.FromURL(media.Url),
			})
		} else {
			album = append(album, &telebot.Photo{
				File: telebot.FromURL(media.Url),
			})
		}
	}

	return album
}

func (m *MediaSender) sendCaption(caption string) error {
	// If caption is empty, ignore sending it
	if m.media.Caption == "" {
		return nil
	}
	// TODO: chunk caption if the length is above the Telegram limit
	_, err := m.bot.Reply(m.msg, caption)
	return err
}
