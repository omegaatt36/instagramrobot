package telegram

import (
	"github.com/omegaatt36/instagramrobot/domain"
	"gopkg.in/telebot.v3"
)

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
