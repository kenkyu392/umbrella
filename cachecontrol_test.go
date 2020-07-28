package umbrella

import (
	"net/http"
	"testing"
)

func TestCacheControl(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("custom", func(t *testing.T) {
		teardown := setup(CacheControl("public", "max-age=86400", "s-maxage=86400")(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("Cache-Control"), "public,max-age=86400,s-maxage=86400"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("empty", func(t *testing.T) {
		teardown := setup(CacheControl()(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("Cache-Control"), ""; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})

	t.Run("no-cache", func(t *testing.T) {
		teardown := setup(NoCache()(handler))
		defer teardown()

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		if got, want := resp.Header.Get("Cache-Control"), "private,no-cache,no-store,must-revalidate,max-age=0,proxy-revalidate,s-maxage=0"; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		_ = resp.Body.Close()
	})
}
