package umbrella

import "net/http"

// Reject provides middleware that returns badStatus if the result of f is false.
func Reject(badStatus int, f func(r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !f(r) {
				http.Error(w, http.StatusText(badStatus), badStatus)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
