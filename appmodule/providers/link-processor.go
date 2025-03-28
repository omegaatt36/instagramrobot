package providers

import (
	"fmt"
	"net/url"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/domain"
)

// LinkProcessor orchestrates the fetching and sending of media based on the input link.
type LinkProcessor struct {
	// InstagramFetcher fetches content from Instagram.
	InstagramFetcher domain.InstagramFetcher
	// ThreadsFetcher fetches content from Threads.
	ThreadsFetcher domain.ThreadsFetcher
	// MediaSender sends the processed media to the destination (e.g., Telegram).
	MediaSender domain.MediaSender
}

// NewLinkProcessorRequest contains the dependencies required to create a LinkProcessor.
type NewLinkProcessorRequest struct {
	// InstagramFetcher fetches content from Instagram.
	InstagramFetcher domain.InstagramFetcher
	// ThreadsFetcher fetches content from Threads.
	ThreadsFetcher domain.ThreadsFetcher
	// Sender sends the processed media.
	Sender domain.MediaSender
}

// NewLinkProcessor creates a new instance of LinkProcessor with the provided dependencies.
func NewLinkProcessor(req NewLinkProcessorRequest) *LinkProcessor {
	return &LinkProcessor{
		InstagramFetcher: req.InstagramFetcher,
		ThreadsFetcher:   req.ThreadsFetcher,
		MediaSender:      req.Sender,
	}
}

// ProcessLink takes a URL string, determines the source (Instagram or Threads),
// fetches the media content using the appropriate fetcher, and then sends it
// using the configured MediaSender.
func (processor *LinkProcessor) ProcessLink(link string) error {
	// Convert link to URL object
	url, err := url.ParseRequestURI(link)

	// Validate URL
	if err != nil {
		return fmt.Errorf("couldn't parse the [%v] link: %w", link, err)
	}

	var media = domain.Media{}

	switch url.Host {
	case "instagram.com", "www.instagram.com":
		shortCode, err := instagram.ExtractShortCodeFromLink(url.Path)
		if err != nil {
			return fmt.Errorf("couldn't extract the short code from the link [%s]: %w", link, err)
		}

		// Process fetching the short code from Instagram
		media, err = processor.InstagramFetcher.GetPostWithCode(shortCode)
		if err != nil {
			return fmt.Errorf("couldn't fetch the post with short code [%s]: %w", shortCode, err)
		}

		media.Source = domain.SourceInstagram
	case "www.threads.net":
		media, err = processor.ThreadsFetcher.GetPostWithURL(url)
		if err != nil {
			return fmt.Errorf("couldn't fetch the post with URL [%s]: %w", url, err)
		}
		media.Source = domain.SourceThreads
	default:
		return fmt.Errorf("can only process links from [instagram.com] or [www.threads.net] not [%s]", url.Host)
	}

	// Extract short code
	if err := processor.MediaSender.Send(&media); err != nil {
		return fmt.Errorf("couldn't send the media, %w", err)
	}

	return nil
}
