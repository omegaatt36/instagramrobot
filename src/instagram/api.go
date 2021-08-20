package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/feelthecode/instagramrobot/src/instagram/response"
	"github.com/gocolly/colly"
)

func GetPostWithCode(code string) (string, error) {
	c := colly.NewCollector()

	res := ""

	// Find and visit all links
	c.OnHTML("script", func(e *colly.HTMLElement) {
		validateText := "window.__additionalDataLoaded('extra',"
		if !strings.Contains(e.Text, validateText) {
			return
		}
		res = strings.Replace(e.Text, validateText, "", 1)
		res = strings.Replace(res, ");", "", 1)
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	// TODO: replace with debugger
	// 	fmt.Println("Visiting", r.URL)
	// })

	// c.OnResponse(func(r *colly.Response) {
	// 	// TODO: replace with debugger
	// 	fmt.Println("Response", r.Headers)
	// })

	c.Visit(fmt.Sprintf("https://www.instagram.com/p/%v/embed/captioned/", code))

	if res == "" {
		return "", errors.New("couldn't parse Instagram JSON")
	}

	var result response.EmbedResponse

	json.Unmarshal([]byte(res), &result)

	fmt.Printf("%+v\n", result)

	return res, nil
}
