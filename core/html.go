package core

import (
	"regexp"
	"strings"
)

func rewriteLinks(htmlContent []byte, proxyDomain string) string {
	re := regexp.MustCompile(`href=["'](http://([a-zA-Z0-9\-]+)\.onion[^"']*)["']`)

	result := re.ReplaceAllStringFunc(string(htmlContent), func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match // No change if pattern doesn't match
		}
		originalURL := submatches[1]
		onionSubdomain := submatches[2]

		newURL := strings.Replace(originalURL, onionSubdomain+".onion", onionSubdomain+"."+proxyDomain, 1)
		return strings.Replace(match, originalURL, newURL, 1)
	})

	return result
}
