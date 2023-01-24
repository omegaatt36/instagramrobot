package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/omegaatt36/instagramrobot/src/instagram/response"
	"github.com/omegaatt36/instagramrobot/src/instagram/transform"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gocolly/colly/v2"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
)

// GetPostWithCode lets you to get information about specific Instagram post
// by providing its unique shortcode
func GetPostWithCode(code string) (transform.Media, error) {
	// TODO: validate code

	URL := fmt.Sprintf("https://www.instagram.com/p/%v/embed/captioned/", code)

	var embeddedMediaImage string
	var embedResponse = response.EmbedResponse{}
	collector := colly.NewCollector()
	collector.SetClient(client)

	collector.OnHTML("img.EmbeddedMediaImage", func(e *colly.HTMLElement) {
		embeddedMediaImage = e.Attr("src")
	})

	collector.OnHTML("script", func(e *colly.HTMLElement) {
		r := regexp.MustCompile(`window\.__additionalDataLoaded\(\'extra\',([\s\S]*)\);`)
		match := r.FindStringSubmatch(e.Text)

		if len(match) < 2 {
			return
		}

		_ = json.Unmarshal([]byte(match[1]), &embedResponse)
	})

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", browser.Chrome())
	})

	if err := collector.Visit(URL); err != nil {
		return transform.Media{}, fmt.Errorf("failed to send HTTP request to the Instagram: %v", err)
	}

	// If the method one which is JSON parsing didn't fail
	if !embedResponse.IsEmpty() {
		// Transform the Embed response and return
		return transform.FromEmbedResponse(embedResponse), nil
	}

	if embeddedMediaImage != "" {
		return transform.Media{
			Url: embeddedMediaImage,
		}, nil
	}

	// If every two methods have failed, then return an error
	return transform.Media{}, errors.New("failed to fetch the post\nthe page might be \"private\", or\nthe link is completely wrong")

}

// ExtractShortcodeFromLink will extract the media shortcode from a URL link or path
func ExtractShortcodeFromLink(link string) (string, error) {
	values := regexp.MustCompile(`(p|tv|reel|reels\/videos)\/([A-Za-z0-9-_]+)`).FindStringSubmatch(link)
	if len(values) != 3 {
		return "", errors.New("couldn't extract the media shortcode from the link")
	}
	// return shortcode
	return values[2], nil
}
