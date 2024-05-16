package domain

// MediaItem contains information about the Instagram post
// which is similar to the Instagram Media struct
type MediaItem struct {
	Id        string `json:"id"`
	ShortCode string `json:"shortcode"`
	Type      string `json:"type"`
	IsVideo   bool   `json:"is_video"`
	Url       string `json:"url"`
}

// Media which contains a single Instagram post
type Media struct {
	Id        string      `json:"id"`
	ShortCode string      `json:"shortcode"`
	Type      string      `json:"type"`
	Comments  uint64      `json:"comments_count"`
	Likes     uint64      `json:"likes_count"`
	Caption   string      `json:"caption"`
	IsVideo   bool        `json:"is_video"`
	Url       string      `json:"url"`
	Items     []MediaItem `json:"items"`
	TakenAt   int64       `json:"taken_at"` // Timestamp
}

// MediaSender defines the contract for sending media to the Telegram chat.
type MediaSender interface {
	Send(*Media) error
	SendCaption(*Media) error
}
