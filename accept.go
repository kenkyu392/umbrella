package umbrella

import "net/http"

// AllowAccept is middleware that allows a request only if any of
// the specified strings is included in the Accept header.
// Returns 406 Not Acceptable status if the request has a type that is not allowed.
func AllowAccept(contentTypes ...string) func(http.Handler) http.Handler {
	return AllowHTTPHeader(http.StatusNotAcceptable, "Accept", contentTypes...)
}

// DisallowAccept is middleware that disallow a request only if any of
// the specified strings is included in the Accept header.
// Returns 406 Not Acceptable status if the request has a type that is not allowed.
func DisallowAccept(contentTypes ...string) func(http.Handler) http.Handler {
	return DisallowHTTPHeader(http.StatusNotAcceptable, "Accept", contentTypes...)
}
