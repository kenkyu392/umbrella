package umbrella

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestStampede(t *testing.T) {

	t.Run("method=GET", func(t *testing.T) {
		var cached = false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cached {
				t.Fatal("cache not used")
			}
			cached = true
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%d", time.Now().UnixNano())
		})

		teardown := setup(Stampede(time.Second * 5)(handler))
		defer teardown()

		var results = make([][]byte, 0)
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
			resp, err := httpClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			results = append(results, raw)
			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
		}

		if !bytes.Equal(results[0], results[1]) {
			t.Fatalf("must be equal: %s, %s", results[0], results[1])
		}
	})

	t.Run("method=POST", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "%d", time.Now().UnixNano())
		})

		teardown := setup(Stampede(time.Second * 5)(handler))
		defer teardown()

		var results = make([][]byte, 0)
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest(http.MethodPost, httpServer.URL, nil)
			resp, err := httpClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			results = append(results, raw)
			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
		}

		if bytes.Equal(results[0], results[1]) {
			t.Fatalf("should not be equal: %s, %s", results[0], results[1])
		}
	})

	t.Run("expired", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%d", time.Now().UnixNano())
		})

		teardown := setup(Stampede(0)(handler))
		defer teardown()

		var results = make([][]byte, 0)
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
			resp, err := httpClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			results = append(results, raw)
			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
			time.Sleep(time.Microsecond)
		}

		if bytes.Equal(results[0], results[1]) {
			t.Fatalf("should not be equal: %s, %s", results[0], results[1])
		}
	})
}
