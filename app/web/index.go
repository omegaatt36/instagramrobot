package web

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/omegaatt36/instagramrobot/appmodule/instagram"
	"github.com/omegaatt36/instagramrobot/logging"
)

//go:embed index.html
var indexHTML embed.FS

// media is a helper struct for passing data to the HTML template.
type media struct {
	// URL is the base64 encoded data URL for the media.
	URL string
	// IsVideo indicates if the media is a video.
	IsVideo bool
}

// index handles requests to the root path ("/"), rendering the main index.html page.
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if err := s.indexPage.Execute(w, nil); err != nil {
		logging.ErrorCtx(r.Context(), err)
	}
}

// parse handles POST requests to "/parse/" (typically from HTMX).
// It extracts the Instagram URL from the form, fetches the media using the
// Instagram fetcher, encodes the media to base64 data URLs, and renders the
// "medias" template block with the results.
func (s *Server) parse(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")

	shortCode, err := instagram.ExtractShortCodeFromLink(url)
	if err != nil {
		logging.ErrorCtx(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	domainMedia, err := instagram.NewExtractor().GetPostWithCode(shortCode)
	if err != nil {
		logging.ErrorCtx(r.Context(), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	medias := make([]media, 0)
	// Check if media has no child item
	if len(domainMedia.Items) == 0 {
		m, err := encodeMediaToBase64(domainMedia.URL, domainMedia.IsVideo)
		if err != nil {
			logging.ErrorCtx(r.Context(), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		medias = append(medias, m)
	} else {
		for _, item := range domainMedia.Items {
			m, err := encodeMediaToBase64(item.URL, item.IsVideo)
			if err != nil {
				logging.ErrorCtx(r.Context(), err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			medias = append(medias, m)
		}
	}

	captionWithBreaks := strings.ReplaceAll(domainMedia.Caption, "\n", "<br>")

	if err := s.indexPage.ExecuteTemplate(w, "medias", map[string]any{
		"Caption": captionWithBreaks,
		"Medias":  medias,
	}); err != nil {
		logging.ErrorCtx(r.Context(), err)
	}
}

// encodeMediaToBase64 downloads media content from a URL and encodes it into a
// base64 data URL string suitable for embedding directly in HTML.
// Note: This can be inefficient for large files, especially videos. Consider
// proxying or direct linking if performance becomes an issue.
func encodeMediaToBase64(link string, isVideo bool) (m media, err error) {
	bs, err := downloadImage(link)
	if err != nil {
		return media{}, err
	}

	base64Str := encodeToBase64(bs)

	mediaType := `data:image/jpeg;base64,%s`
	if isVideo {
		mediaType = `data:video/mp4;base64,%s`
	}

	m.IsVideo = isVideo
	m.URL = fmt.Sprintf(mediaType, base64Str)

	return
}

// downloadImage fetches the content of a given URL and returns it as a byte slice.
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

// encodeToBase64 encodes a byte slice into a standard base64 string.
func encodeToBase64(imageData []byte) string {
	// 將圖片數據編碼為 Base64 字符串
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	return base64Str
}
