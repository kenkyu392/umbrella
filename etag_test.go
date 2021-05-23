package umbrella

import (
	"net/http"
	"testing"
	"time"
)

func TestEtag(t *testing.T) {
	data := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
	<circle cx="50" cy="50" r="40" stroke="#6a737d" stroke-width="4" fill="#1b1f23" />
	</svg>`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", time.Now().Add(time.Second*60).Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
	teardown := setup(ETag()(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	if got, want := resp.Header.Get("ETag"), `"7e36ab9325b08dcb99adb7894770c4b8"`; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()

	req.Header.Set("If-None-Match", resp.Header.Get("ETag"))
	resp, _ = httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusNotModified; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}
