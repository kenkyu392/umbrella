package umbrella

import (
	"net/http"
	"testing"
)

const (
	testContentTypeJSON = "application/json"
	testContentTypeHTML = "text/html; charset=utf-8"
)

func TestAllowContentType(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(AllowContentType("application/json", "text/json")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Content-Type", testContentTypeJSON)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("415", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Content-Type", testContentTypeHTML)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusUnsupportedMediaType; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}

func TestDisallowContentType(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(DisallowContentType("text/html")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Content-Type", testContentTypeJSON)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("415", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Content-Type", testContentTypeHTML)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusUnsupportedMediaType; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
