package umbrella

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	type key struct{}
	value := time.Now().UnixNano()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Context().Value(key{}), value; got != want {
			t.Errorf("got: %v, want: %v", got, want)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(Context(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key{}, value)
	})(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}
