package umbrella

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestMetricsRecorder(t *testing.T) {
	t.Run("case=ok", func(t *testing.T) {
		mr := NewMetricsRecorder()
		mw := mr.Middleware()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				mr.Handler(w, r)
				return
			}
			fn := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				switch r.URL.Path {
				case "/1xx":
					w.WriteHeader(http.StatusSwitchingProtocols)
				case "/2xx":
					w.WriteHeader(http.StatusOK)
				case "/3xx":
					w.WriteHeader(http.StatusMultipleChoices)
				case "/4xx":
					w.WriteHeader(http.StatusBadRequest)
				case "/5xx":
					w.WriteHeader(http.StatusInternalServerError)
				default:
					w.WriteHeader(http.StatusOK)
				}
			}
			mw(http.HandlerFunc(fn)).ServeHTTP(w, r)
		})

		teardown := setup(handler)
		defer teardown()

		for _, path := range []string{"/5xx", "/4xx", "/3xx", "/2xx", "/1xx"} {
			func() {
				req, err := http.NewRequest(http.MethodGet, httpServer.URL+path, nil)
				if err != nil {
					t.Error(err)
				}
				resp, err := httpClient.Do(req)
				if err != nil {
					t.Error(err)
				}
				resp.Body.Close()
			}()
		}

		req, _ := http.NewRequest(http.MethodGet, httpServer.URL+"/metrics", nil)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		got := &Metrics{}
		if err := json.Unmarshal(body, got); err != nil {
			t.Error(err)
		}
		want := mr.Metrics()
		want.GoroutinesTotalCount = got.GoroutinesTotalCount
		want.MaxGoroutinesCount = got.MaxGoroutinesCount
		want.MinGoroutinesCount = got.MinGoroutinesCount
		want.AvgGoroutinesCount = got.AvgGoroutinesCount
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %#v \nwant: %#v", got, want)
		}
		resp.Body.Close()

		m := mr.Metrics()
		if got := m.RequestsTotalCount; got != 5 {
			t.Errorf("got: %v, want: 5", got)
		}
		if got := m.MethodCount[http.MethodGet]; got != 5 {
			t.Errorf("got: %v, want: 4", got)
		}
		if got := m.StatusCount[http.StatusOK]; got != 1 {
			t.Errorf("got: %v, want: 1", got)
		}
		if got := m.StatusClassCount["1xx"]; got != 1 {
			t.Errorf("got: %v, want: 1", got)
		}
		if got := m.StatusClassCount["2xx"]; got != 1 {
			t.Errorf("got: %v, want: 1", got)
		}
		if got := m.StatusClassCount["3xx"]; got != 1 {
			t.Errorf("got: %v, want: 1", got)
		}
		if got := m.StatusClassCount["4xx"]; got != 1 {
			t.Errorf("got: %v, want: 1", got)
		}
		if got := m.StatusClassCount["5xx"]; got != 1 {
			t.Errorf("got: %v, want: 1", got)
		}
	})

	t.Run("case=error", func(t *testing.T) {
		mr := NewMetricsRecorder()
		mw := mr.Middleware()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = r.Body.Close()
			mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(w, r)
		})

		teardown := setup(handler)
		defer teardown()

		reqBody := bytes.NewBuffer([]byte("1234567890"))
		req, _ := http.NewRequest(http.MethodPost, httpServer.URL, reqBody)
		resp, _ := httpClient.Do(req)
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
		resp.Body.Close()

		got := mr.Metrics()
		if want := newMetrics(); !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot: %#v \nwant: %#v", got, want)
		}
	})
}
