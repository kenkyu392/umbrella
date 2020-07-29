package umbrella

import "net/http"

// AllowContentType is middleware that allows a request only if any of
// the specified strings is included in the Content-Type header.
// Returns 415 Unsupported Media Type status if the request has a type that is not allowed.
func AllowContentType(userAgents ...string) func(http.Handler) http.Handler {
	return AllowHTTPHeader(http.StatusUnsupportedMediaType, "Content-Type", userAgents...)
}

// DisallowContentType is middleware that disallow a request only if any of
// the specified strings is included in the Content-Type header.
// Returns 415 Unsupported Media Type status if the request has a type that is not allowed.
func DisallowContentType(userAgents ...string) func(http.Handler) http.Handler {
	return DisallowHTTPHeader(http.StatusUnsupportedMediaType, "Content-Type", userAgents...)
}
