package umbrella

import (
	"net/http"
	"testing"
)

func TestRealIP(t *testing.T) {
	t.Run("case=X-Forwarded-For", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.RemoteAddr, "100.100.100.100"; got != want {
				t.Errorf("got: %v, want: %v", got, want)
			}
			w.WriteHeader(http.StatusOK)
		})

		teardown := setup(RealIP()(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("X-Real-IP", "101.101.101.101")
		req.Header.Set("X-Forwarded-For", "127.0.0.1, 100.100.100.100, 192.168.0.4, localhost")
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("case=X-Forwarded-For", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.RemoteAddr, "101.101.101.101"; got != want {
				t.Errorf("got: %v, want: %v", got, want)
			}
			w.WriteHeader(http.StatusOK)
		})

		teardown := setup(RealIP()(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("X-Real-IP", "101.101.101.101")
		req.Header.Set("X-Forwarded-For", "127.0.0.1, 192.168.0.4, localhost")
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("case=X-Real-IP", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.RemoteAddr, "100.100.100.100"; got != want {
				t.Errorf("got: %v, want: %v", got, want)
			}
			w.WriteHeader(http.StatusOK)
		})

		teardown := setup(RealIP()(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("X-Real-IP", "100.100.100.100")
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("case=empty", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		teardown := setup(RealIP()(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		req.Header.Set("X-Real-IP", "127.0.0.1")
		req.Header.Set("X-Forwarded-For", "192.168.0.3, localhost, 127.0.0.1")
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
