package cache

type Media struct {
	IsVideo bool
	FileID  string
	Caption string
	Items   []MediaItem
}

type MediaItem struct {
	IsVideo bool
	FileID  string
}
