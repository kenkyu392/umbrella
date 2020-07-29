package umbrella

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

// Recover is a middleware that recovers from panic and records a
// stack trace and returns a 500 Internal Server Error status.
func Recover(out io.Writer) func(http.Handler) http.Handler {
	const size = 4 << 10
	if out == nil {
		out = os.Stderr
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, size)
					length := runtime.Stack(stack, true)
					fmt.Fprintf(out, "[RECOVER] %v\n%s\n", err, stack[:length])
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
