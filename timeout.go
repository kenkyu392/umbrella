package umbrella

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Timeout cancels the context at the given time.
// Returns 504 Gateway Timeout status if a timeout occurs.
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer func() {
				cancel()
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
