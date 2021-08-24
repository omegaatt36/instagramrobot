package helpers

import "regexp"

// ExtractLinksFromString lets you to extract HTTP links from a string
func ExtractLinksFromString(input string) []string {
	regex := `(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`
	r := regexp.MustCompile(regex)
	return r.FindAllString(input, -1)
}
