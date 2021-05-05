package umbrella

import "net/http"

// Switch provides a middleware that executes the next handler if the result of
// f is true, and executes h if it is false.
func Switch(f func(r *http.Request) bool, h http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if f(r) {
				next.ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
