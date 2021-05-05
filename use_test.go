package umbrella

import (
	"net/http"
	"reflect"
	"testing"
)

func TestUse(t *testing.T) {
	done := map[string]bool{
		"root:before": false,
		"root:after":  false,
		"mw1:before":  false,
		"mw1:after":   false,
		"mw2:before":  false,
		"mw2:after":   false,
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done["root:before"] = true
		if want := map[string]bool{
			"mw1:before":  true,
			"mw2:before":  true,
			"root:before": true,
			"root:after":  false,
			"mw2:after":   false,
			"mw1:after":   false,
		}; !reflect.DeepEqual(done, want) {
			t.Errorf("\ngot: %#v \nwant: %#v", done, want)
		}
		w.WriteHeader(http.StatusOK)
		done["root:after"] = true
		if want := map[string]bool{
			"mw1:before":  true,
			"mw2:before":  true,
			"root:before": true,
			"root:after":  true,
			"mw2:after":   false,
			"mw1:after":   false,
		}; !reflect.DeepEqual(done, want) {
			t.Errorf("\ngot: %#v \nwant: %#v", done, want)
		}
	})

	mw1 := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			done["mw1:before"] = true
			if want := map[string]bool{
				"mw1:before":  true,
				"mw2:before":  false,
				"root:before": false,
				"root:after":  false,
				"mw2:after":   false,
				"mw1:after":   false,
			}; !reflect.DeepEqual(done, want) {
				t.Errorf("\ngot: %#v \nwant: %#v", done, want)
			}
			next.ServeHTTP(w, r)
			done["mw1:after"] = true
			if want := map[string]bool{
				"mw1:before":  true,
				"mw2:before":  true,
				"root:before": true,
				"root:after":  true,
				"mw2:after":   true,
				"mw1:after":   true,
			}; !reflect.DeepEqual(done, want) {
				t.Errorf("\ngot: %#v \nwant: %#v", done, want)
			}
		}
		return http.HandlerFunc(fn)
	}

	mw2 := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			done["mw2:before"] = true
			if want := map[string]bool{
				"mw1:before":  true,
				"mw2:before":  true,
				"root:before": false,
				"root:after":  false,
				"mw2:after":   false,
				"mw1:after":   false,
			}; !reflect.DeepEqual(done, want) {
				t.Errorf("\ngot: %#v \nwant: %#v", done, want)
			}
			next.ServeHTTP(w, r)
			done["mw2:after"] = true
			if want := map[string]bool{
				"mw1:before":  true,
				"mw2:before":  true,
				"root:before": true,
				"root:after":  true,
				"mw2:after":   true,
				"mw1:after":   false,
			}; !reflect.DeepEqual(done, want) {
				t.Errorf("\ngot: %#v \nwant: %#v", done, want)
			}
		}
		return http.HandlerFunc(fn)
	}

	mw := Use(mw1, mw2)

	teardown := setup(mw(handler))
	defer teardown()

	req, _ := http.NewRequest(http.MethodGet, httpServer.URL, nil)
	resp, _ := httpClient.Do(req)
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
	if want := map[string]bool{
		"mw1:before":  true,
		"mw2:before":  true,
		"root:before": true,
		"root:after":  true,
		"mw2:after":   true,
		"mw1:after":   true,
	}; !reflect.DeepEqual(done, want) {
		t.Errorf("\ngot: %#v \nwant: %#v", done, want)
	}
	_ = resp.Body.Close()
}
