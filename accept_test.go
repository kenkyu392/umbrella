package umbrella

import (
	"net/http"
	"testing"
)

func TestAllowAccept(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(AllowAccept("application/json", "text/json")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Accept", testContentTypeJSON)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("406", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Accept", testContentTypeHTML)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusNotAcceptable; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}

func TestDisallowAccept(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(DisallowAccept("text/html")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Accept", testContentTypeJSON)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("406", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("Accept", testContentTypeHTML)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusNotAcceptable; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
