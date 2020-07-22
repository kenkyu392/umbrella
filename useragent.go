package umbrella

import "net/http"

const (
	chromeUserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.XX (KHTML, like Gecko) Chrome/84.0.XXXX.XX Safari/537.XX"
	firefoxUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/XXXXXXXX Firefox/76.XX"
)

// AllowUserAgent is middleware that allows a request only if any of
// the specified strings is included in the User-Agent header.
func AllowUserAgent(userAgents ...string) func(next http.Handler) http.Handler {
	return AllowHTTPHeader(http.StatusForbidden, "User-Agent", userAgents...)
}

// DisallowUserAgent is middleware that disallow a request only if any of
// the specified strings is included in the User-Agent header.
func DisallowUserAgent(userAgents ...string) func(next http.Handler) http.Handler {
	return DisallowHTTPHeader(http.StatusForbidden, "User-Agent", userAgents...)
}
