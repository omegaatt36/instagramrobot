package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gocolly/colly/v2"

	"github.com/omegaatt36/instagramrobot/domain"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Extractor implements the domain.InstagramFetcher interface.
type Extractor struct {
	// client is the HTTP client used for making requests.
	client *http.Client
}

var (
	// Regular expression to extract embedded JSON data from Instagram embed pages.
	gqlDataRegex = regexp.MustCompile(`\\\"gql_data\\\":([\s\S]*)\}\"\}\]\]\,\[\"NavigationMetrics`)
	// Regular expression to extract shortcode from Instagram URLs.
	shortCodeRegex = regexp.MustCompile(`(p|tv|reel|reels\/videos)\/([A-Za-z0-9-_]+)`)
)

// NewExtractor creates a new instance of Extractor with a configured HTTP client.
func NewExtractor() domain.InstagramFetcher {
	return &Extractor{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

// fromEmbedResponse transforms the internal EmbedResponse structure to the domain.Media model.
func fromEmbedResponse(embed EmbedResponse) domain.Media {
	media := domain.Media{
		ShortCode: embed.Media.ShortCode,
		URL:       embed.ExtractMediaURL(),
		IsVideo:   embed.IsVideo(),
		Caption:   embed.GetCaption(),
	}

	for _, item := range embed.Media.SliderItems.Edges {
		media.Items = append(media.Items, &domain.MediaItem{
			IsVideo: item.Node.IsVideo,
			URL:     item.Node.ExtractMediaURL(),
		})
	}

	return media
}

// GetPostWithCode fetches media information for a given Instagram post shortcode.
// It attempts to parse embedded JSON data from the embed page.
// If JSON parsing fails, it falls back to extracting the cover photo URL.
func (repo *Extractor) GetPostWithCode(code string) (domain.Media, error) {
	URL := fmt.Sprintf("https://www.instagram.com/p/%v/embed/captioned/", code)

	var coverPhoto string
	var embedResponse = EmbedResponse{}
	collector := colly.NewCollector()
	collector.SetClient(repo.client)

	collector.OnHTML("img.EmbeddedMediaImage", func(e *colly.HTMLElement) {
		coverPhoto = e.Attr("src")
	})

	collector.OnHTML("script", func(e *colly.HTMLElement) {
		match := gqlDataRegex.FindStringSubmatch(e.Text)

		if len(match) < 2 {
			return
		}

		// Note: This string replacement is fragile and might break if Instagram changes the embedding format.
		s := strings.ReplaceAll(match[1], `\"`, `"`)
		s = strings.ReplaceAll(s, `\\/`, `/`)
		s = strings.ReplaceAll(s, `\\`, `\`)

		// Use a local error variable instead of log.Fatal
		parseErr := json.Unmarshal([]byte(s), &embedResponse)
		if parseErr != nil {
			// Log the error but don't terminate the application.
			// The function will continue and might fall back to cover photo extraction.
			logging.Errorf("Failed to unmarshal GQL data for code %s: %v", code, parseErr)
			// Clear embedResponse to ensure fallback logic triggers correctly if needed
			embedResponse = EmbedResponse{}
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", browser.Random())
	})

	if err := collector.Visit(URL); err != nil {
		return domain.Media{}, fmt.Errorf("failed to send HTTP request to the Instagram: %v", err)
	}

	// If the method one which is JSON parsing didn't fail
	if !embedResponse.IsEmpty() {
		// Transform the Embed response and return
		return fromEmbedResponse(embedResponse), nil
	}

	if coverPhoto != "" {
		return domain.Media{
			URL:     coverPhoto,
			Caption: "can only fetch the cover photo",
		}, nil
	}

	// If both JSON parsing and cover photo extraction failed, return an error.
	return domain.Media{}, errors.New("failed to fetch the post\nthe page might be \"private\", or\nthe link is completely wrong")
}

// ExtractShortCodeFromLink extracts the media shortcode from an Instagram URL path.
// It supports various URL formats like /p/, /tv/, /reel/, /reels/videos/.
func ExtractShortCodeFromLink(link string) (string, error) {
	values := shortCodeRegex.FindStringSubmatch(link)
	if len(values) != 3 {
		return "", errors.New("couldn't extract the media short code from the link")
	}
	// return short code
	return values[2], nil
}
