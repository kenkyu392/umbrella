package umbrella

import (
	"net/http"
	"testing"
)

func TestClickjacking(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("deny", func(t *testing.T) {
		teardown := setup(Clickjacking("DENY")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("X-Frame-Options"), "deny"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("sameorigin", func(t *testing.T) {
		teardown := setup(Clickjacking("SAMEORIGIN")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("X-Frame-Options"), "sameorigin"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("empty", func(t *testing.T) {
		teardown := setup(Clickjacking("")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("X-Frame-Options"), "deny"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
