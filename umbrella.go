package umbrella

import (
	"net/http"
	"time"
)

var (
	processStartTime time.Time
)

func init() {
	processStartTime = time.Now()
}

// ServeMux ...
type ServeMux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
