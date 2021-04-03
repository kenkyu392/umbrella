package umbrella

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// Stampede provides a simple cache middleware that is valid for a specified amount of time.
// It uses singleflight for caching to prevent thundering-herd and cache-stampede.
// If this middleware is requested at the same time, it executes the handler
// only once and shares the execution result with all requests.
func Stampede(d time.Duration) func(http.Handler) http.Handler {
	type cache struct {
		time   time.Time
		code   int
		header http.Header
		body   []byte
	}
	var sm sync.Map
	var group singleflight.Group
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// If the method is not GET or HEAD, call the handler.
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			path := r.URL.String()
			v, _, _ := group.Do(path, func() (interface{}, error) {
				if v, ok := sm.Load(path); ok {
					if c, ok := v.(*cache); ok {
						if c.time.After(time.Now()) {
							return c, nil
						}
					}
					sm.Delete(path)
				}
				// Use ResponseRecorder to record the results.
				rec := httptest.NewRecorder()
				next.ServeHTTP(rec, r)
				c := &cache{
					time:   time.Now().Add(d),
					code:   rec.Code,
					header: rec.Header(),
					body:   rec.Body.Bytes(),
				}
				sm.Store(path, c)
				return c, nil
			})

			c := v.(*cache)
			for k, v := range c.header {
				w.Header()[k] = v
			}
			w.WriteHeader(c.code)
			_, _ = w.Write(c.body)
		}
		return http.HandlerFunc(fn)
	}
}
