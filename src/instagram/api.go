package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	goQuery "github.com/PuerkitoBio/goquery"
	"github.com/feelthecode/instagramrobot/src/instagram/response"
	"github.com/feelthecode/instagramrobot/src/instagram/transform"
)

type API struct {
	embedResponse response.EmbedResponse
}

// TODO change response type
func (a *API) GetPostWithCode(code string) (transform.Media, error) {
	// TODO: validate code

	URL := fmt.Sprintf("https://www.instagram.com/p/%v/embed/captioned/", code)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		// TODO: return error
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", browser.Chrome())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// TODO: use logger
		// TODO: return error
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// TODO: use logger
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return transform.Media{}, fmt.Errorf("Instagram returned %v error.", res.Status)
	}

	// Load the HTML document
	doc, err := goQuery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// For each script item found
	doc.Find("script").Each(func(i int, s *goQuery.Selection) {
		a.checkScriptForJSON(s.Text())
	})

	// If the method one which is JSON parsing didn't fail
	if !a.embedResponse.IsEmpty() {
		// Transform the Embed response and return
		return transform.FromEmbedResponse(a.embedResponse), nil
	}

	// TODO: Try second method by HTML body parsing

	// If every two methods have failed, then return an error
	return transform.Media{}, errors.New("failed to fetch the post\nthe page might be \"private\", or\nthe link is completely wrong")
}

func (a *API) checkScriptForJSON(scriptContent string) {
	validateText := "window.__additionalDataLoaded('extra',"
	if !strings.Contains(scriptContent, validateText) {
		return
	}
	res := strings.Replace(scriptContent, validateText, "", 1)
	res = strings.Replace(res, ");", "", 1)

	_ = json.Unmarshal([]byte(res), &a.embedResponse)
}
