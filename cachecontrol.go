package umbrella

import (
	"net/http"
	"strings"
)

// CacheControl adds the Cache-Control header.
func CacheControl(opts ...string) func(next http.Handler) http.Handler {
	return ResponseHeader(CacheControlHeaderFunc(opts...))
}

// NoCache adds the Cache-Control to disable the cache.
func NoCache() func(next http.Handler) http.Handler {
	return ResponseHeader(NoCacheHeaderFunc())
}

// CacheControlHeaderFunc returns a HeaderFunc that adds a
// Cache-Control header.
func CacheControlHeaderFunc(opts ...string) HeaderFunc {
	if len(opts) == 0 {
		return func(_ http.Header) {}
	}
	value := strings.Join(opts, ",")
	return func(header http.Header) {
		header.Set("Cache-Control", value)
	}
}

// NoCacheHeaderFunc returns the HeaderFunc to add the Cache-Control
// header that disables the cache.
func NoCacheHeaderFunc() HeaderFunc {
	return CacheControlHeaderFunc(
		"private", "no-cache", "no-store", "must-revalidate",
		"max-age=0", "proxy-revalidate", "s-maxage=0",
	)
}
