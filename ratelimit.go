package umbrella

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimit provides middleware that limits the number of requests processed per second.
func RateLimit(rl int) func(http.Handler) http.Handler {
	var l = rate.NewLimiter(rate.Limit(rl), 1)
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			waitRateLimit(l, next, w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// RateLimitPerIP provides middleware that limits the number of requests processed per second per IP.
func RateLimitPerIP(rl int) func(http.Handler) http.Handler {
	var m sync.Map
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ip := realIP(r)
			if v, ok := m.Load(ip); ok {
				if l, ok := v.(*rate.Limiter); ok {
					waitRateLimit(l, next, w, r)
					return
				}
			}
			m.Store(ip, rate.NewLimiter(rate.Limit(rl), 1))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func waitRateLimit(l *rate.Limiter, next http.Handler, w http.ResponseWriter, r *http.Request) {
	if err := l.Wait(r.Context()); err != nil {
		log.Printf("ratelimit.error: %#v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	next.ServeHTTP(w, r)
}
