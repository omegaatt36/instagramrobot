package telegram

import (
	"fmt"

	"github.com/omegaatt36/instagramrobot/domain"
	"github.com/omegaatt36/instagramrobot/logging"
	"gopkg.in/telebot.v3"
)

type MediaSender struct {
	bot *telebot.Bot
	msg *telebot.Message
}

// NewMediaSender creates a new MediaSenderImpl instance
func NewMediaSender(bot *telebot.Bot, msg *telebot.Message) domain.MediaSender {
	return &MediaSender{
		bot: bot,
		msg: msg,
	}
}

const (
	caption = "🤖 Downloaded via @InstagramKeeperBot"
)

// Send will start to process Media and eventually send it to the Telegram chat
func (m *MediaSender) Send(media *domain.Media) error {
	logging.Infof("chatID(%d) source(%s) short code(%s)", m.msg.Sender.ID, media.Source, media.ShortCode)

	// Check if media has no child item
	if len(media.Items) == 0 {
		return m.sendSingleMedia(media)
	}

	return m.sendNestedMedia(media)
}

func (m *MediaSender) sendSingleMedia(media *domain.Media) error {
	if media.IsVideo {
		if _, err := m.bot.Send(m.msg.Chat, &telebot.Video{
			File:    telebot.FromURL(media.Url),
			Caption: caption,
		}); err != nil {
			return fmt.Errorf("couldn't send the single video, %w", err)
		}

		logging.Debugf("Sent single video with short code [%v]", media.ShortCode)
	} else {
		if _, err := m.bot.Send(m.msg.Chat, &telebot.Photo{
			File:    telebot.FromURL(media.Url),
			Caption: caption,
		}); err != nil {
			return fmt.Errorf("couldn't send the single photo, %w", err)
		}

		logging.Debugf("Sent single photo with short code [%v]", media.ShortCode)
	}

	return m.SendCaption(media)
}

func (m *MediaSender) sendNestedMedia(media *domain.Media) error {
	_, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromMedia(media))
	if err != nil {
		return fmt.Errorf("couldn't send the nested media, %w", err)
	}
	return m.SendCaption(media)
}

func (m *MediaSender) generateAlbumFromMedia(media *domain.Media) telebot.Album {
	var album telebot.Album

	for _, media := range media.Items {
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

// SendCaption will send the caption to the chat.
func (m *MediaSender) SendCaption(media *domain.Media) error {
	// If caption is empty, ignore sending it
	if media.Caption == "" {
		return nil
	}
	// TODO: chunk caption if the length is above the Telegram limit
	_, err := m.bot.Reply(m.msg, media.Caption)
	return err
}
