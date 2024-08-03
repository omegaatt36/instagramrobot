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
	ThreadsFetcher   domain.ThreadsFetcher
	MediaSender      domain.MediaSender
}

// NewLinkProcessorRequest is the request for LinkProcessor.
type NewLinkProcessorRequest struct {
	InstagramFetcher domain.InstagramFetcher
	ThreadsFetcher   domain.ThreadsFetcher
	Sender           domain.MediaSender
}

// NewLinkProcessor constructor
func NewLinkProcessor(req NewLinkProcessorRequest) *LinkProcessor {
	return &LinkProcessor{
		InstagramFetcher: req.InstagramFetcher,
		ThreadsFetcher:   req.ThreadsFetcher,
		MediaSender:      req.Sender,
	}
}

// ProcessLink will process a single link
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
