package umbrella

import (
	"net/http"
	"testing"
)

func TestAllowUserMethod(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(AllowMethod("get")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("405", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusMethodNotAllowed; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}

func TestDisallowMethod(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(DisallowMethod("post")(handler))
	defer teardown()

	t.Run("200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("405", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusMethodNotAllowed; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
