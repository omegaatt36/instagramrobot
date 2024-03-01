package domain

// InstagramFetcher defines the contract for fetching Instagram post.
type InstagramFetcher interface {
	GetPostWithCode(code string) (Media, error)
}
