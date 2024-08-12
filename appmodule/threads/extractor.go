package threads

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/omegaatt36/instagramrobot/domain"
)

// InstagramFetcherRepo is the repository for fetching Instagram media.
type Extractor struct {
	client *http.Client
}

// NewExtractor will create a new instance of InstagramFetcherRepo.
func NewExtractor() domain.ThreadsFetcher {
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

func (repo *Extractor) GetPostWithURL(URL *url.URL) (media domain.Media, err error) {
	URL.RawQuery = ""

	URL = URL.JoinPath("embed")

	collector := colly.NewCollector()
	collector.SetClient(repo.client)

	// case single image or video
	collector.OnHTML("div.SingleInnerMediaContainer", func(e *colly.HTMLElement) {
		if src := e.ChildAttr("img", "src"); src != "" {
			media.URL = src
		}
		if src := e.ChildAttr("video > source", "src"); src != "" {
			media.URL = src
			media.IsVideo = true
		}
	})

	// case multiple image or video
	collector.OnHTML("div.MediaScrollImageContainer", func(e *colly.HTMLElement) {
		if src := e.ChildAttr("img", "src"); src != "" {
			media.Items = append(media.Items, &domain.MediaItem{
				URL: src,
			})
		}
		if src := e.ChildAttr("video > source", "src"); src != "" {
			media.Items = append(media.Items, &domain.MediaItem{
				URL:     src,
				IsVideo: true,
			})
		}
	})

	// caption
	collector.OnHTML("span.BodyTextContainer", func(e *colly.HTMLElement) {
		media.Caption = e.Text
	})

	if err := collector.Visit(URL.String()); err != nil {
		return domain.Media{}, fmt.Errorf("failed to send HTTP request to the Instagram: %v", err)
	}

	return
}
