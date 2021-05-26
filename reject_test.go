package umbrella

import (
	"net/http"
	"testing"
)

func TestReject(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(Reject(http.StatusForbidden, func(r *http.Request) bool {
		return r.Header.Get("X-API-Key") != ""
	})(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("X-API-Key", "abcdefghijklmnopqrstuvwxyz")
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("403", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusForbidden; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
