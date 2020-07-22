# umbrella

[![test](https://github.com/kenkyu392/umbrella/workflows/test/badge.svg)](https://github.com/kenkyu392/umbrella)
[![codecov](https://codecov.io/gh/kenkyu392/umbrella/branch/master/graph/badge.svg)](https://codecov.io/gh/kenkyu392/umbrella)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-00ADD8?logo=go)](https://pkg.go.dev/github.com/kenkyu392/umbrella)
[![go report card](https://goreportcard.com/badge/github.com/kenkyu392/umbrella)](https://goreportcard.com/report/github.com/kenkyu392/umbrella)
[![license](https://img.shields.io/github/license/kenkyu392/umbrella.svg)](LICENSE)

This package provides middleware intended for use with various frameworks compatible with the standard `net/http` ecosystem.

## Installation

```
go get -u github.com/kenkyu392/umbrella
```

## Middlewares

### Context

Context is middleware that operates the context of request scope.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kenkyu392/umbrella"
)

type key struct{}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "context: %v", r.Context().Value(key{}))
	})

	m := http.NewServeMux()

	// You can embed the value in the request context.
	mw := umbrella.Context(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key{}, time.Now().UnixNano())
	})
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

### AllowUserAgent/DisallowUserAgent

Allow/DisallowUserAgent is middleware that performs authentication based on the request User-Agent.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ua: %v", r.UserAgent())
	})

	m := http.NewServeMux()

  // Only accessible in Firefox and Chrome.
  allows := umbrella.AllowUserAgent("Firefox", "Chrome")
	m.Handle("/allows",
		allows(handler),
  )

  // Not accessible in Edge and Internet Explorer.
  disallows := umbrella.DisallowUserAgent("Edg", "MSIE")
	m.Handle("/disallows",
		disallows(handler),
	)

	http.ListenAndServe(":3000", m)
}
```

## License

[MIT](LICENSE)
