package web

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/logging"
)

// Embed the entire directory.
//
//go:embed index.html
var indexHTML embed.FS

type Item struct {
	URL string
}

// handler function #1 - returns the index.html template, with film data
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	index, err := template.ParseFS(indexHTML, "index.html")
	if err != nil {
		logging.ErrorCtx(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := index.Execute(w, nil); err != nil {
		logging.ErrorCtx(r.Context(), err)
	}
}

// handler function #2 - returns the template block with the newly added film, as an HTMX response
func (s *Server) addFilm(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	index, err := template.ParseFS(indexHTML, "index.html")
	if err != nil {
		logging.ErrorCtx(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortCode, err := instagram.ExtractShortCodeFromLink(url)
	if err != nil {
		logging.ErrorCtx(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	media, err := instagram.NewInstagramFetcher().GetPostWithCode(shortCode)
	if err != nil {
		logging.ErrorCtx(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]Item, 0)
	// Check if media has no child item
	if len(media.Items) == 0 {
		link, err := parseURLToDom(media.Url, media.IsVideo)
		if err != nil {
			logging.ErrorCtx(r.Context(), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, Item{
			URL: link,
		})
	} else {
		for _, item := range media.Items {
			link, err := parseURLToDom(item.Url, item.IsVideo)
			if err != nil {
				logging.ErrorCtx(r.Context(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, Item{
				URL: link,
			})
		}
	}

	if err := index.ExecuteTemplate(w, "instagram-item-element", map[string]any{
		"Caption": media.Caption,
		"Items":   items,
	}); err != nil {
		logging.ErrorCtx(r.Context(), err)
	}
}

func parseURLToDom(link string, isVideo bool) (string, error) {
	bs, err := downloadImage(link)
	if err != nil {
		return "", nil
	}

	base64Str := encodeToBase64(bs)

	if isVideo {
		return fmt.Sprintf(`<video src='data:video/mp4;base64,%s' controls style="max-width: 300px; height: auto;"/>`, base64Str), nil
	}

	return fmt.Sprintf(`<img src='data:image/jpeg;base64,%s' style="max-width: 300px; height: auto;"/>`, base64Str), nil
}

func downloadImage(url string) ([]byte, error) {
	// 發送 GET 請求以下載圖片
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 讀取圖片的內容
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func encodeToBase64(imageData []byte) string {
	// 將圖片數據編碼為 Base64 字符串
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	return base64Str
}
