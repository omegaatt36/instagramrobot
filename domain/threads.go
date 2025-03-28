package domain

import "net/url"

// ThreadsFetcher defines the interface for fetching Threads posts.
type ThreadsFetcher interface {
	// GetPostWithURL fetches media information for a given Threads post URL.
	GetPostWithURL(*url.URL) (Media, error)
}
