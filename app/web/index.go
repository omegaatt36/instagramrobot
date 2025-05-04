package web

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

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

// "medias" template block with the results or an error message.
func (s *Server) parse(w http.ResponseWriter, r *http.Request) {
	renderError := func(errMsg string) {
		logging.ErrorCtx(r.Context(), errMsg)

		data := map[string]any{
			"Error": errMsg,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if err := s.indexPage.ExecuteTemplate(w, "medias", data); err != nil {
			logging.ErrorCtx(r.Context(), fmt.Errorf("failed to execute error template: %w", err))
			http.Error(w, "An internal error occurred while trying to display the error message.", http.StatusInternalServerError)
		}
	}

	postURL := r.PostFormValue("url")
	if postURL == "" {
		renderError("Please enter an Instagram URL.")
		return
	}

	domainMedia, err := s.linkProcessor.ProcessLink(postURL)
	if err != nil {
		renderError(fmt.Sprintf("Failed to process link '%s': %v", postURL, err))
		return
	}

	medias := make([]media, 0)
	if len(domainMedia.Items) == 0 && domainMedia.URL != "" {
		m, err := encodeMediaToBase64(domainMedia.URL, domainMedia.IsVideo)
		if err != nil {
			renderError(fmt.Sprintf("Failed to encode media from %s: %v", domainMedia.URL, err))
			return
		}
		medias = append(medias, m)
	} else if len(domainMedia.Items) > 0 {
		for i, item := range domainMedia.Items {
			m, err := encodeMediaToBase64(item.URL, item.IsVideo)
			if err != nil {
				renderError(fmt.Sprintf("Failed to encode media item #%d from %s: %v", i+1, item.URL, err))
				return
			}
			medias = append(medias, m)
		}
	} else if domainMedia.Caption == "" {
		renderError(fmt.Sprintf("Could not find any media or caption for the provided URL: %s", postURL))
		return
	}

	captionWithBreaks := strings.ReplaceAll(domainMedia.Caption, "\n", "<br>")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := s.indexPage.ExecuteTemplate(w, "medias", map[string]any{
		"Caption": captionWithBreaks,
		"Medias":  medias,
	}); err != nil {
		logging.ErrorCtx(r.Context(), fmt.Errorf("failed to execute success template: %w", err))
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
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

// encodeToBase64 encodes a byte slice into a standard base64 string.
func encodeToBase64(imageData []byte) string {
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	return base64Str
}
