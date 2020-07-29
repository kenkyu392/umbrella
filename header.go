package umbrella

import (
	"net/http"
	"strings"
)

// HeaderFunc ...
type HeaderFunc func(header http.Header)

// RequestHeader is middleware that edits the header of the request.
func RequestHeader(fs ...HeaderFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for _, f := range fs {
				f(r.Header)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// ResponseHeader is middleware that edits the header of the response.
func ResponseHeader(fs ...HeaderFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for _, f := range fs {
				f(w.Header())
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// AllowHTTPHeader is middleware that allows a request only when one
// of the specified strings is included in the specified request header.
func AllowHTTPHeader(badStatus int, name string, values ...string) func(http.Handler) http.Handler {
	key := http.CanonicalHeaderKey(name)
	list := make([]string, len(values))
	for i, v := range values {
		list[i] = strings.ToLower(v)
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			value := strings.ToLower(r.Header.Get(key))
			for _, v := range list {
				if strings.Contains(value, v) {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.WriteHeader(badStatus)
		}
		return http.HandlerFunc(fn)
	}
}

// DisallowHTTPHeader is middleware that disallows a request only when one
// of the specified strings is included in the specified request header.
func DisallowHTTPHeader(badStatus int, name string, values ...string) func(http.Handler) http.Handler {
	key := http.CanonicalHeaderKey(name)
	list := make([]string, len(values))
	for i, v := range values {
		list[i] = strings.ToLower(v)
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			value := strings.ToLower(r.Header.Get(key))
			for _, v := range list {
				if strings.Contains(value, v) {
					w.WriteHeader(badStatus)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
