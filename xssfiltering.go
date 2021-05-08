package umbrella

import (
	"net/http"
)

// XSSFiltering provides middleware that enables the ability to stop a page from
// loading when a cross-site scripting attack is detected.
func XSSFiltering(opt string) func(http.Handler) http.Handler {
	return ResponseHeader(XSSFilteringHeaderFunc(opt))
}

// XSSFilteringHeaderFunc returns a HeaderFunc to enable XSS filtering.
func XSSFilteringHeaderFunc(opt string) HeaderFunc {
	if opt == "" {
		opt = "1; mode=block"
	}
	return func(header http.Header) {
		header.Set("X-XSS-Protection", opt)
	}
}
