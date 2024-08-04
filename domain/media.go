//go:generate go-enum -f=$GOFILE

package domain

// Source defines the source of the media.
// ENUM(
// Instagram
// Threads
// )
type Source string

// MediaItem contains media information
type MediaItem struct {
	IsVideo bool   `json:"is_video"`
	URL     string `json:"url"`
}

// Media which contains a single post
type Media struct {
	ShortCode string
	Caption   string
	IsVideo   bool
	URL       string
	Items     []MediaItem
	Source    Source
}

// MediaSender defines the contract for sending media to the Telegram chat.
type MediaSender interface {
	Send(*Media) error
	SendCaption(*Media) error
}
