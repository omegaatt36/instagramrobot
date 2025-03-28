package domain

// InstagramFetcher defines the interface for fetching Instagram posts.
type InstagramFetcher interface {
	// GetPostWithCode fetches media information for a given Instagram post shortcode.
	GetPostWithCode(code string) (Media, error)
}
