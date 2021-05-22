package umbrella

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

type testStaticServeMux struct {
	handlers map[string]http.Handler
}

func (t *testStaticServeMux) Handle(pattern string, handler http.Handler) {
	t.handlers[pattern] = handler
}

func (t *testStaticServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	t.handlers[pattern] = http.HandlerFunc(handler)
}

type testStaticBrokenFS struct {
	fs   http.FileSystem
	open func(name string) (http.File, error)
}

func (t *testStaticBrokenFS) Open(name string) (http.File, error) {
	if t.open != nil {
		return t.open(name)
	}
	if strings.HasSuffix(name, "index.html") {
		return nil, errors.New("error")
	}
	return t.fs.Open(name)
}

func TestStatic(t *testing.T) {
	htmlData, err := os.ReadFile("./testdata/index.html")
	if err != nil {
		log.Fatal(err)
	}
	jsData, err := os.ReadFile("./testdata/js/app.js")
	if err != nil {
		log.Fatal(err)
	}

	m := &testStaticServeMux{handlers: make(map[string]http.Handler)}
	Static(m, "/static", "./testdata", ETag())

	testCases := []struct {
		pattern string
		code    int
		path    string
		body    []byte
		etag    string
	}{
		{
			pattern: "/static",
			code:    http.StatusMovedPermanently,
			path:    "/static",
			body:    nil,
			etag:    "",
		},
		{
			pattern: "/static/",
			code:    http.StatusOK,
			path:    "/static/",
			body:    htmlData,
			etag:    fmt.Sprintf(`"%x"`, md5.Sum(htmlData)),
		},
		{
			pattern: "/static/",
			code:    http.StatusOK,
			path:    "/static/js/app.js",
			body:    jsData,
			etag:    fmt.Sprintf(`"%x"`, md5.Sum(jsData)),
		},
	}
	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodGet, tc.path, nil)
		w := httptest.NewRecorder()
		m.handlers[tc.pattern].ServeHTTP(w, r)
		if got, want := w.Code, tc.code; got != want {
			t.Errorf("got: %d, want: %d", got, want)
		}
		if tc.code == http.StatusOK && !reflect.DeepEqual(w.Body.Bytes(), tc.body) {
			t.Errorf("%s: response body does not match", tc.path)
		}
		if got, want := w.Header().Get("ETag"), tc.etag; got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
}

func TestStaticHandler(t *testing.T) {
	htmlData, err := os.ReadFile("./testdata/index.html")
	if err != nil {
		log.Fatal(err)
	}
	jsData, err := os.ReadFile("./testdata/js/app.js")
	if err != nil {
		log.Fatal(err)
	}

	pattern, handler := StaticHandler("/static", "./testdata")
	if got, want := pattern, "/static/"; got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}

	testCases := []struct {
		code int
		path string
		body []byte
	}{
		{code: http.StatusOK, path: "/static/", body: htmlData},
		{code: http.StatusOK, path: "/static/js/app.js", body: jsData},
	}
	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodGet, tc.path, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if got, want := w.Code, tc.code; got != want {
			t.Errorf("got: %d, want: %d", got, want)
		}
		if tc.code == http.StatusOK && !reflect.DeepEqual(w.Body.Bytes(), tc.body) {
			t.Errorf("%s: response body does not match", tc.path)
		}
	}
}

func TestStaticHandlerFS(t *testing.T) {
	t.Run("case=broken-fs-a", func(t *testing.T) {
		fs := &testStaticBrokenFS{}
		fs.open = func(name string) (http.File, error) {
			return nil, errors.New("error")
		}
		pattern, handler := StaticHandlerFS("/static", fs)
		if got, want := pattern, "/static/"; got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
		r := httptest.NewRequest(http.MethodGet, "/static/js/app.js", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if got, want := w.Code, http.StatusInternalServerError; got != want {
			t.Errorf("got: %d, want: %d", got, want)
		}
	})

	t.Run("case=broken-fs-b", func(t *testing.T) {
		pattern, handler := StaticHandlerFS("/static", &testStaticBrokenFS{fs: http.Dir("./testdata")})
		if got, want := pattern, "/static/"; got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
		r := httptest.NewRequest(http.MethodGet, "/static/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if got, want := w.Code, http.StatusInternalServerError; got != want {
			t.Errorf("got: %d, want: %d", got, want)
		}
	})
}
