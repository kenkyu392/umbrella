package umbrella

import (
	"net/http"
	"strings"
)

// AllowMethod is a middleware that returns a 405 Method Not Allowed
// status if the request method is not one of the given methods.
func AllowMethod(methods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for _, method := range methods {
				if strings.EqualFold(r.Method, method) {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return http.HandlerFunc(fn)
	}
}

// DisallowMethod is a middleware that returns a 405 Method Not Allowed
// status if the request method is one of the given methods.
func DisallowMethod(methods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for _, method := range methods {
				if strings.EqualFold(r.Method, method) {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
