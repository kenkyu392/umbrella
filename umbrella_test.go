package umbrella

import (
	"net/http"
	"net/http/httptest"
)

var (
	httpServer *httptest.Server
	httpClient *http.Client
)

func setup(handler http.Handler) func() {
	httpClient = new(http.Client)
	httpServer = httptest.NewServer(handler)
	return httpServer.Close
}
