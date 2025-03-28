//go:generate go-enum -f=$GOFILE

package domain

// Source indicates the origin platform of the media (e.g., Instagram, Threads).
// ENUM(
// Instagram
// Threads
// )
type Source string

// MediaItem represents a single item (photo or video) within a post.
type MediaItem struct {
	IsVideo bool `json:"is_video"`
	// URL is the direct URL to the media content.
	URL string `json:"url"`
}

// Media represents a single post, which might contain multiple media items.
type Media struct {
	// ShortCode is the unique identifier for the post on its platform.
	ShortCode string
	// Caption is the text content associated with the post.
	Caption string
	// IsVideo indicates if the primary media item is a video (relevant for single-item posts).
	IsVideo bool
	// URL is the direct URL to the primary media item (relevant for single-item posts).
	URL string
	// Items holds multiple media items if the post is a carousel or album.
	Items []*MediaItem
	// Source specifies the platform (Instagram or Threads).
	Source Source
}

// MediaSender defines the interface for sending media (e.g., to a Telegram chat).
type MediaSender interface {
	// Send sends the media content (photos/videos).
	Send(*Media) error
	SendCaption(*Media) error
}
