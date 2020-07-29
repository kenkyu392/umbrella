package umbrella

import (
	"net/http"
	"testing"
)

func TestRecover(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("recover test")
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(Recover(nil)(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusInternalServerError; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}
