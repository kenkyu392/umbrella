package umbrella

import (
	"net/http"
)

// Debug provides middleware that executes the handler only if b is true.
// If b is false, it will return 404.
func Debug(b bool) func(http.Handler) http.Handler {
	return Switch(func(r *http.Request) bool {
		return b
	}, http.NotFoundHandler())
}
