package umbrella

import (
	"net/http"
	"testing"
)

func TestAllowUserAgent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(AllowUserAgent("Chrome")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("User-Agent", chromeUserAgent)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("403", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("User-Agent", firefoxUserAgent)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusForbidden; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}

func TestDisallowUserAgent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(DisallowUserAgent("Firefox")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("User-Agent", chromeUserAgent)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("403", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("User-Agent", firefoxUserAgent)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusForbidden; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
