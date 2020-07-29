package umbrella

import (
	"net/http"
)

// Use creates a single middleware that executes multiple middlewares.
func Use(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := range middlewares {
			next = middlewares[len(middlewares)-1-i](next)
		}
		return next
	}
}
