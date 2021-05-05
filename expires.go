package umbrella

import (
	"net/http"
	"time"
)

// Expires provides middleware for adding response expiration dates.
// The expiration time is set to the current time plus d.
func Expires(d time.Duration) func(http.Handler) http.Handler {
	return ResponseHeader(ExpiresHeaderFunc(d))
}

// ExpiresHeaderFunc returns a HeaderFunc that adds an expiration date.
// The expiration time is set to the current time plus d.
func ExpiresHeaderFunc(d time.Duration) HeaderFunc {
	return func(header http.Header) {
		header.Set("Expires", time.Now().Add(d).Format(http.TimeFormat))
	}
}
