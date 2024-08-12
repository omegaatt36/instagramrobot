package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gocolly/colly/v2"

	"github.com/omegaatt36/instagramrobot/domain"
)

// Extractor is the implement for fetching Instagram media.
type Extractor struct {
	client *http.Client
}

// NewInstagramFetcher will create a new instance of InstagramFetcherRepo.
func NewInstagramFetcher() domain.InstagramFetcher {
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

// fromEmbedResponse will automatically transforms the EmbedResponse to the Media
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

// GetPostWithCode lets you to get information about specific Instagram post
// by providing its unique short code
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
		r := regexp.MustCompile(`\\\"gql_data\\\":([\s\S]*)\}\"\}\]\]\,\[\"NavigationMetrics`)
		match := r.FindStringSubmatch(e.Text)

		if len(match) < 2 {
			return
		}

		s := strings.ReplaceAll(match[1], `\"`, `"`)
		s = strings.ReplaceAll(s, `\\/`, `/`)
		s = strings.ReplaceAll(s, `\\`, `\`)

		err := json.Unmarshal([]byte(s), &embedResponse)
		if err != nil {
			log.Fatal(err)
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

	// if every two methods have failed, then return an error
	return domain.Media{}, errors.New("failed to fetch the post\nthe page might be \"private\", or\nthe link is completely wrong")
}

// ExtractShortCodeFromLink will extract the media short code from a URL link or path
func ExtractShortCodeFromLink(link string) (string, error) {
	values := regexp.MustCompile(`(p|tv|reel|reels\/videos)\/([A-Za-z0-9-_]+)`).FindStringSubmatch(link)
	if len(values) != 3 {
		return "", errors.New("couldn't extract the media short code from the link")
	}
	// return short code
	return values[2], nil
}
