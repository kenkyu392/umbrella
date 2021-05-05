package umbrella

import (
	"net/http"
	"testing"
	"time"
)

func TestExpires(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	teardown := setup(Expires(time.Minute)(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	expires, err := time.Parse(http.TimeFormat, resp.Header.Get("Expires"))
	if err != nil {
		t.Fatal(err)
	}
	if time.Now().After(expires) {
		t.Fatal("the expiration date must be newer than the now")
	}
	_ = resp.Body.Close()
}
