package umbrella

import (
	"net/http"
)

// Debug provides middleware that executes the handler only if d is true.
// If d is false, it will return 404.
func Debug(d bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if d {
				next.ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
