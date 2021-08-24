package helpers

import "net/url"

// IsValidUrl checks if the String is a valid URL
// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
