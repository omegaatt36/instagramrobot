package providers

import (
	"net/url"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/domain"
	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "I couldn't parse the [%v] link.", link)
	}

	// Validate HOST in the URL (only instagram.com is allowed)
	if url.Host != "instagram.com" && url.Host != "www.instagram.com" {
		return errors.Wrapf(ErrInvalidHost, "can only process links from [instagram.com] not [%s]", url.Host)
	}

	// Extract short code
	shortCode, err := instagram.ExtractShortCodeFromLink(url.Path)
	if err != nil {
		return errors.Wrapf(err, "couldn't extract the short code from the link [%s]", link)
	}

	// Process fetching the short code from Instagram
	response, err := processor.InstagramFetcher.GetPostWithCode(shortCode)
	if err != nil {
		return errors.Wrapf(err, "couldn't fetch the post with short code [%s]", shortCode)
	}

	return errors.Wrapf(processor.MediaSender.Send(&response), "couldn't send the media with short code [%s]", shortCode)
}
