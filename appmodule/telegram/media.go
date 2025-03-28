package telegram

import (
	"gopkg.in/telebot.v4"

	"github.com/omegaatt36/instagramrobot/domain"
)

// convertMediaToInputtable converts a domain.Media (representing a single item post)
// into a telebot.Inputtable (Photo or Video) suitable for sending.
func convertMediaToInputtable(media *domain.Media) telebot.Inputtable {
	if media.IsVideo {
		return &telebot.Video{
			File: telebot.FromURL(media.URL),
		}
	}

	return &telebot.Photo{
		File: telebot.FromURL(media.URL),
	}
}

// convertMediaItemToInputtable converts a domain.MediaItem (from a carousel)
// into a telebot.Inputtable (Photo or Video) suitable for sending within an album.
func convertMediaItemToInputtable(mediaItem *domain.MediaItem) telebot.Inputtable {
	if mediaItem.IsVideo {
		return &telebot.Video{
			File: telebot.FromURL(mediaItem.URL),
		}
	}

	return &telebot.Photo{
		File: telebot.FromURL(mediaItem.URL),
	}
}
