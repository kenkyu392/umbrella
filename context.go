package umbrella

import (
	"context"
	"net/http"
)

// ContextFunc ...
type ContextFunc func(ctx context.Context) context.Context

// Context is middleware that edits the context of the request.
func Context(f ContextFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = f(ctx)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
