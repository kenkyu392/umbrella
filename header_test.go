package umbrella

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestRequestHeader(t *testing.T) {
	value := strconv.Itoa(int(time.Now().UnixNano()))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Request-ID"), value; got != want {
			t.Errorf("got: %v, want: %v", got, want)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(RequestHeader(func(header http.Header) {
		header.Set("Request-ID", value)
	})(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}

func TestResponseHeader(t *testing.T) {
	value := strconv.Itoa(int(time.Now().UnixNano()))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := w.Header().Get("Request-ID"), value; got != want {
			t.Errorf("got: %v, want: %v", got, want)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(ResponseHeader(func(header http.Header) {
		header.Set("Request-ID", value)
	})(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	if got, want := resp.Header.Get("Request-ID"), value; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	_ = resp.Body.Close()
}
