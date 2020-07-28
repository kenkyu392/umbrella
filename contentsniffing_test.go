package umbrella

import (
	"net/http"
	"testing"
)

func TestContentSniffing(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(ContentSniffing()(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	if got, want := resp.Header.Get("X-Content-Type-Options"), "nosniff"; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}
