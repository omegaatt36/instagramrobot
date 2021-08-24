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
	goQuery "github.com/PuerkitoBio/goquery"
	"github.com/feelthecode/instagramrobot/src/instagram/response"
	"github.com/feelthecode/instagramrobot/src/instagram/transform"
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
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return transform.Media{}, fmt.Errorf("failed to initial HTTP request to the Instagram: %v", err)
	}
	// Set random chrome user-agent
	req.Header.Set("User-Agent", browser.Chrome())

	res, err := client.Do(req)
	if err != nil {
		return transform.Media{}, fmt.Errorf("failed to send HTTP request to the Instagram: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return transform.Media{}, fmt.Errorf("request failed with %v HTTP error", res.Status)
	}

	// Load the HTML document
	doc, err := goQuery.NewDocumentFromReader(res.Body)
	if err != nil {
		return transform.Media{}, fmt.Errorf("couldn't initial document parser: %v", err)
	}

	embedResponse := response.EmbedResponse{}

	// For each script item found
	doc.Find("script").Each(func(i int, s *goQuery.Selection) {
		checkScriptForJSON(&embedResponse, s.Text())
	})

	// If the method one which is JSON parsing didn't fail
	if !embedResponse.IsEmpty() {
		// Transform the Embed response and return
		return transform.FromEmbedResponse(embedResponse), nil
	}

	// TODO: Try second method by HTML body parsing

	// If every two methods have failed, then return an error
	return transform.Media{}, errors.New("failed to fetch the post\nthe page might be \"private\", or\nthe link is completely wrong")
}

// ExtractShortcodeFromLink will extract the media shortcode from a URL link or path
func ExtractShortcodeFromLink(link string) (string, error) {
	values := regexp.MustCompile(`(p|tv|reel)\/([A-Za-z0-9-_]+)`).FindStringSubmatch(link)
	if len(values) != 3 {
		return "", errors.New("couldn't extract the media shortcode from the link")
	}
	// return shortcode
	return values[2], nil
}

func checkScriptForJSON(embedResponse *response.EmbedResponse, scriptContent string) {
	validateText := "window.__additionalDataLoaded('extra',"
	if !strings.Contains(scriptContent, validateText) {
		return
	}
	res := strings.Replace(scriptContent, validateText, "", 1)
	res = strings.Replace(res, ");", "", 1)

	_ = json.Unmarshal([]byte(res), &embedResponse)
}
