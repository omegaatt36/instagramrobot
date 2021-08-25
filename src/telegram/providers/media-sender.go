package providers

import (
	"github.com/feelthecode/instagramrobot/src/instagram/transform"
	"github.com/feelthecode/instagramrobot/src/telegram/cache"
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
	// Hit cache
	if cachedMedia, found := cache.Get(m.media.Shortcode); found {
		if len(cachedMedia.Items) == 0 {
			if cachedMedia.IsVideo {
				_, err := m.bot.Send(m.msg.Chat, &tb.Video{
					File:    tb.File{FileID: cachedMedia.FileID},
					Caption: caption,
				})
				if err == nil {
					log.Debugf("Sent single video with shortcode [%v] from cache", m.media.Shortcode)
					m.sendCaption(cachedMedia.Caption)
					return
				}
			} else {
				_, err := m.bot.Send(m.msg.Chat, &tb.Photo{
					File:    tb.File{FileID: cachedMedia.FileID},
					Caption: caption,
				})
				if err == nil {
					log.Debugf("Sent single photo with shortcode [%v] from cache", m.media.Shortcode)
					m.sendCaption(cachedMedia.Caption)
					return
				}
			}
		} else {
			// Send album from cache
			_, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromCachedMedia(cachedMedia.Items))
			if err == nil {
				log.Debugf("Sent album with shortcode [%v] from cache", m.media.Shortcode)
				m.sendCaption(cachedMedia.Caption)
				return
			}
		}
	}

	// Check if media has no child item
	if len(m.media.Items) == 0 {
		m.sendSingleMedia()
	} else {
		m.sendNestedMedia()
	}
}

func (m *MediaSender) sendSingleMedia() {
	if m.media.IsVideo {
		msg, err := m.bot.Send(m.msg.Chat, &tb.Video{
			File:    tb.FromURL(m.media.Url),
			Caption: caption,
		})
		if err != nil {
			log.Errorf("couldn't send the single video: %v", err)
			return
		}
		log.Debugf("Sent single video with shortcode [%v]", m.media.Shortcode)

		// Cache file_id
		cache.Set(m.media.Shortcode, cache.Media{
			IsVideo: m.media.IsVideo,
			FileID:  msg.Video.File.FileID,
			Caption: m.media.Caption,
		})
	} else {
		msg, err := m.bot.Send(m.msg.Chat, &tb.Photo{
			File:    tb.FromURL(m.media.Url),
			Caption: caption,
		})
		if err != nil {
			log.Errorf("couldn't send the single photo: %v", err)
			return
		}
		log.Debugf("Sent single photo with shortcode [%v]", m.media.Shortcode)

		// Cache file_id
		cache.Set(m.media.Shortcode, cache.Media{
			IsVideo: m.media.IsVideo,
			FileID:  msg.Photo.File.FileID,
			Caption: m.media.Caption,
		})
	}

	m.sendCaption(m.media.Caption)
}

func (m *MediaSender) sendNestedMedia() {
	message, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromMedia())
	if err != nil {
		log.Errorf("couldn't send the single video: %v", err)
	}

	log.Debugf("Sent album with shortcode [%v]", m.media.Shortcode)

	// Start caching process
	var media cache.Media
	media.Caption = m.media.Caption

	for _, msg := range message {
		if msg.Video != nil {
			media.Items = append(media.Items, cache.MediaItem{
				IsVideo: true,
				FileID:  msg.Video.File.FileID,
			})
		} else {
			media.Items = append(media.Items, cache.MediaItem{
				IsVideo: false,
				FileID:  msg.Photo.File.FileID,
			})
		}
	}
	cache.Set(m.media.Shortcode, media)

	m.sendCaption(m.media.Caption)
}

func (m *MediaSender) generateAlbumFromCachedMedia(items []cache.MediaItem) tb.Album {
	var album tb.Album

	for _, media := range items {
		if media.IsVideo {
			album = append(album, &tb.Video{
				File: tb.File{FileID: media.FileID},
			})
		} else {
			album = append(album, &tb.Photo{
				File: tb.File{FileID: media.FileID},
			})
		}
	}

	return album
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
