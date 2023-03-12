package validation

import "net/url"

// IsValidUrl verifies that the string passed to it is a valid URL.
// It expects the URL Scheme to be either http or https, and for it to have a non empty host.
func IsValidUrl(urlString string) bool {
	url, err := url.ParseRequestURI(urlString)
	return err == nil && (url.Scheme == "http" || url.Scheme == "https") && url.Host != ""
}
