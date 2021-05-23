package umbrella

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// ETag provides middleware that calculates MD5 from the response data and sets
// it in the ETag header.
func ETag() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)
			body := rec.Body.Bytes()
			etag := fmt.Sprintf(`"%x"`, md5.Sum(body))
			w.Header().Set("ETag", etag)
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.WriteHeader(rec.Code)
			_, _ = w.Write(body)
		}
		return http.HandlerFunc(fn)
	}
}
