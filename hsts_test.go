package umbrella

import (
	"net/http"
	"testing"
)

func TestHSTS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("includeSubDomains", func(t *testing.T) {
		teardown := setup(HSTS(60, "includeSubDomains")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("Strict-Transport-Security"), "max-age=60; includeSubDomains"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("preload", func(t *testing.T) {
		teardown := setup(HSTS(60, "preload")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("Strict-Transport-Security"), "max-age=60; preload"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("empty", func(t *testing.T) {
		teardown := setup(HSTS(-1, "")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("Strict-Transport-Security"), "max-age=31536000"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
