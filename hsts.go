package umbrella

import (
	"fmt"
	"net/http"
)

// HSTS adds the Strict-Transport-Security header.
// Proper use of this header will mitigate stripping attacks.
func HSTS(maxAge int, opt string) func(http.Handler) http.Handler {
	return ResponseHeader(HSTSHeaderFunc(maxAge, opt))
}

// HSTSHeaderFunc returns a HeaderFunc that adds a
// Strict-Transport-Security header.
func HSTSHeaderFunc(maxAge int, opt string) HeaderFunc {
	if maxAge < 0 {
		maxAge = 31536000 // 365 days
	}
	value := fmt.Sprintf("max-age=%d", maxAge)
	if opt == "includeSubDomains" || opt == "preload" {
		value = fmt.Sprintf("%s; %s", value, opt)
	}
	return func(header http.Header) {
		header.Set("Strict-Transport-Security", value)
	}
}
