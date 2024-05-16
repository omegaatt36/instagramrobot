package providers

import (
	"fmt"
	"net/url"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/domain"
)

// LinkProcessor is the implementation of LinkProcessor.
type LinkProcessor struct {
	InstagramFetcher domain.InstagramFetcher
	MediaSender      domain.MediaSender
}

// NewLinkProcessor constructor
func NewLinkProcessor(fetcher domain.InstagramFetcher, sender domain.MediaSender) *LinkProcessor {
	return &LinkProcessor{
		InstagramFetcher: fetcher,
		MediaSender:      sender,
	}
}

// ProcessLink will process a single link
func (processor *LinkProcessor) ProcessLink(link string) error {
	// Convert link to URL object
	url, err := url.ParseRequestURI(link)

	// Validate URL
	if err != nil {
		return fmt.Errorf("I couldn't parse the [%v] link: %w", link, err)
	}

	// Validate HOST in the URL (only instagram.com is allowed)
	if url.Host != "instagram.com" && url.Host != "www.instagram.com" {
		return fmt.Errorf("can only process links from [instagram.com] not [%s], %w", url.Host, ErrInvalidHost)
	}

	// Extract short code
	shortCode, err := instagram.ExtractShortCodeFromLink(url.Path)
	if err != nil {
		return fmt.Errorf("couldn't extract the short code from the link [%s]: %w", link, err)
	}

	// Process fetching the short code from Instagram
	response, err := processor.InstagramFetcher.GetPostWithCode(shortCode)
	if err != nil {
		return fmt.Errorf("couldn't fetch the post with short code [%s]: %w", shortCode, err)
	}

	if err := processor.MediaSender.Send(&response); err != nil {
		return fmt.Errorf("couldn't send the media with short code [%s], %w", shortCode, err)
	}

	return nil
}
