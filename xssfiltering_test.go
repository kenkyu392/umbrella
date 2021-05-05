package umbrella

import (
	"net/http"
	"testing"
)

func TestXSSFiltering(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("case=disable", func(t *testing.T) {
		teardown := setup(XSSFiltering("0")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("X-XSS-Protection"), "0"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("case=empty", func(t *testing.T) {
		teardown := setup(XSSFiltering("")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("X-XSS-Protection"), "1; mode=block"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
