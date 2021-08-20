package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/feelthecode/instagramrobot/src/instagram/response"
	"github.com/gocolly/colly"
)

type API struct {
	result response.EmbedResponse
	body   []byte
}

func (a *API) GetPostWithCode(code string) (response.EmbedResponse, error) {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("script", func(e *colly.HTMLElement) {
		validateText := "window.__additionalDataLoaded('extra',"
		if !strings.Contains(e.Text, validateText) {
			return
		}
		res := strings.Replace(e.Text, validateText, "", 1)
		res = strings.Replace(res, ");", "", 1)

		json.Unmarshal([]byte(res), &a.result)
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	// TODO: replace with debugger
	// 	fmt.Println("Visiting", r.URL)
	// })

	c.OnResponse(func(r *colly.Response) {
		// TODO: replace with debugger
		// fmt.Println("Response", r.Headers)
		a.body = r.Body
	})

	c.Visit(fmt.Sprintf("https://www.instagram.com/p/%v/embed/captioned/", code))

	if !a.result.IsEmpty() {
		fmt.Printf("%+v\n", a.result)

		return a.result, nil
	}

	// Try second method by HTML parse body
	return a.result, errors.New("couldn't parse Instagram JSON")
}
