package umbrella

import (
	"net/http"
	"path"
	"strings"
)

// Static calls StaticFS internally.
func Static(m ServeMux, pattern, root string, middlewares ...func(http.Handler) http.Handler) {
	StaticFS(m, pattern, http.Dir(root), middlewares...)
}

// StaticFS adds an endpoint for static files to ServeMux.
func StaticFS(m ServeMux, pattern string, fs http.FileSystem, middlewares ...func(http.Handler) http.Handler) {
	pattern2, handler := StaticHandlerFS(pattern, fs)
	if pattern != pattern2 {
		m.Handle(pattern, http.RedirectHandler(pattern2, http.StatusMovedPermanently))
	}
	m.Handle(pattern2, Use(middlewares...)(handler))
}

// StaticHandler calls StaticFileHandlerFS internally.
func StaticHandler(pattern, root string) (string, http.Handler) {
	return StaticHandlerFS(pattern, http.Dir(root))
}

// StaticHandlerFS returns a normalized pattern and Handler for static file delivery.
func StaticHandlerFS(pattern string, fs http.FileSystem) (string, http.Handler) {
	if pattern != "/" && pattern[len(pattern)-1] != '/' {
		pattern += "/"
	}
	return pattern, http.StripPrefix(pattern, http.FileServer(fileSystem{fs: fs}))
}

type fileSystem struct {
	fs http.FileSystem
}

// Open implements the FileSystem and opens index.html with a request to the directory.
func (fs fileSystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	if s, _ := f.Stat(); s.IsDir() {
		index := path.Join(strings.TrimSuffix(name, "/"), "index.html")
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}
	return f, nil
}
