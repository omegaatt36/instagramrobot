package telegram

import (
	"fmt"

	"gopkg.in/telebot.v3"

	"github.com/omegaatt36/instagramrobot/domain"
	"github.com/omegaatt36/instagramrobot/logging"
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

// Send will start to process Media and eventually send it to the Telegram chat
func (m *MediaSender) Send(media *domain.Media) error {
	logging.Infof("chatID(%d) source(%s) short code(%s)", m.msg.Sender.ID, media.Source, media.ShortCode)

	var fnSend func(*domain.Media) error

	// Check if media has no child item
	if len(media.Items) == 0 {
		fnSend = m.sendSingleMedia
	} else {
		fnSend = m.sendNestedMedia
	}

	if err := fnSend(media); err != nil {
		return fmt.Errorf("sent the media failed, %w", err)
	}

	return m.SendCaption(media)
}

func (m *MediaSender) sendSingleMedia(media *domain.Media) error {
	if media.URL == "" {
		return nil
	}

	mediaToSend := convertMediaToInputtable(media)

	if _, err := m.bot.Send(m.msg.Chat, mediaToSend); err != nil {
		return fmt.Errorf("couldn't send the %s photo, %w", mediaToSend.MediaType(), err)
	}

	logging.Debugf("Sent single %s with short code [%v]", mediaToSend.MediaType(), media.ShortCode)

	return nil
}

func (m *MediaSender) sendNestedMedia(media *domain.Media) error {
	var album telebot.Album

	for _, media := range media.Items {
		album = append(album, convertMediaItemToInputtable(media))
	}

	_, err := m.bot.SendAlbum(m.msg.Chat, album)
	if err != nil {
		return fmt.Errorf("couldn't send the nested media, %w", err)
	}

	return nil
}

// SendCaption will send the caption to the chat.
func (m *MediaSender) SendCaption(media *domain.Media) error {
	// If caption is empty, ignore sending it
	if media.Caption == "" {
		return nil
	}

	// shrink media caption below 4096 characters
	if len(media.Caption) > 4096 {
		media.Caption = media.Caption[:4096]
	}

	_, err := m.bot.Reply(m.msg, media.Caption)
	return err
}
