package umbrella

import (
	"net/http"
	"testing"
)

func TestSwitch(t *testing.T) {
	hA := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	hB := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	teardown := setup(Switch(func(r *http.Request) bool {
		return r.Method == http.MethodGet
	}, hB)(hA))
	defer teardown()

	var testCases = []struct {
		want   int
		method string
	}{
		{
			want:   http.StatusOK,
			method: http.MethodGet,
		},
		{
			want:   http.StatusCreated,
			method: http.MethodPost,
		},
	}
	for _, testCase := range testCases {
		req, _ := http.NewRequest(testCase.method, httpServer.URL, nil)
		resp, _ := httpClient.Do(req)
		if got := resp.StatusCode; got != testCase.want {
			t.Errorf("method:%s, got: %v, want: %v", testCase.method, got, testCase.want)
		}
		_ = resp.Body.Close()
	}
}
