package domain

import "net/url"

// ThreadsFetcher defines the contract for fetching Threads post.
type ThreadsFetcher interface {
	GetPostWithURL(*url.URL) (Media, error)
}
