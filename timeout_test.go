package umbrella

import (
	"net/http"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 10):
		}
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(Timeout(time.Second * 2)(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusGatewayTimeout; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}
