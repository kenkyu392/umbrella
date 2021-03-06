# umbrella

[![test](https://github.com/kenkyu392/umbrella/workflows/test/badge.svg)](https://github.com/kenkyu392/umbrella)
[![codecov](https://codecov.io/gh/kenkyu392/umbrella/branch/master/graph/badge.svg)](https://codecov.io/gh/kenkyu392/umbrella)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-00ADD8?logo=go)](https://pkg.go.dev/github.com/kenkyu392/umbrella)
[![go report card](https://goreportcard.com/badge/github.com/kenkyu392/umbrella)](https://goreportcard.com/report/github.com/kenkyu392/umbrella)
[![license](https://img.shields.io/github/license/kenkyu392/umbrella.svg)](LICENSE)

This package provides 20+ middleware that can be used in various frameworks compatible with Go standard `net/http` ecosystem.

## Installation

```
go get -u github.com/kenkyu392/umbrella
```

## Usage

The middleware provided in this package can be used by intuitively combining individual functions.  
Documentation for all middleware can be found [here](#middleware).

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kenkyu392/umbrella"
)

func init() {
	os.Setenv("DEBUG", "true")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := http.NewServeMux()

	// Create a MetricsRecorder and start recording using the middleware.
	mr := umbrella.NewMetricsRecorder(
		umbrella.WithRequestMetricsHookFunc(func(rm *umbrella.RequestMetrics) {
			// You can use hook functions to output request metrics to a log
			// or send them to a monitoring service.
			raw, err := json.Marshal(rm)
			if err != nil {
				log.Println(err)
			}
			log.Printf("%s", raw)
		}),
	)

	// Create the middleware.
	mw := umbrella.Use(
		// Recover panic.
		umbrella.Recover(os.Stderr),
		// Enables the recording of metrics.
		mr.Middleware(),
		// Only accessible in Firefox and Chrome.
		umbrella.AllowUserAgent("Firefox", "Chrome"),
		// Limit the display of iframe to mitigate clickjacking attacks.
		umbrella.Clickjacking("deny"),
		// It implements a countermeasure for Content-Type snuffing vulnerability,
		// which is a problem in old Internet Explorer, for example.
		umbrella.ContentSniffing(),
		// Enable browser cache for 120s.
		umbrella.CacheControl("public", "max-age=120", "s-maxage=120"),
		// Set the timeout at 800ms.
		umbrella.Timeout(time.Millisecond*800),
		// Enable caching for 2s.
		umbrella.Stampede(time.Second*2),
		// Limit the number of requests from the same IP address to 10 per second.
		umbrella.RateLimitPerIP(10),
	)

	// Enable the handler for debugging using the value set in the environment variable.
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	debugMW := umbrella.Use(
		umbrella.Recover(os.Stderr),
		umbrella.Debug(debug),
	)

	m.Handle("/search", mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		log.Printf("Search: %s", q)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Search: %s", q)
	})))

	// You can use MetricsRecorder.Handler to view the metrics.
	// ~$ curl -s http://localhost:3000/metrics | jq .
	m.Handle("/metrics", debugMW(http.HandlerFunc(mr.Handler)))

	http.ListenAndServe(":3000", m)
}
```

This is an example of middleware used to deliver static files such as images.

```go
mw := umbrella.Use(
  // Sets the ETag.
  umbrella.ETag(),
  // Add headers for basic cache control, etc.
  umbrella.ResponseHeader(
    umbrella.ClickjackingHeaderFunc("deny"),
    umbrella.ContentSniffingHeaderFunc(),
    umbrella.XSSFilteringHeaderFunc("1; mode=block"),
    umbrella.CacheControlHeaderFunc("public", "max-age=86400", "no-transform"),
    umbrella.ExpiresHeaderFunc(time.Second*86400),
  ),
  // Set the timeout at 2s.
  umbrella.Timeout(time.Second*2),
  // Enable caching for 60s.
  umbrella.Stampede(time.Second*60),
)
```

## Middleware

| Middleware | Description |
| ---------- | ----------- |
| [Use](#use)                                                                  | Creates a single middleware that executes multiple middleware. |
| [RealIP](#realip)                                                            | Override the RemoteAddr in http.Request with an X-Forwarded-For or X-Real-IP header. |
| [Recover](#recover)                                                          | Recover from panic and record a stack trace and return a 500 Internal Server Error status. |
| [Timeout](#timeout)                                                          | Timeout cancels the context at the given time. |
| [Context](#context)                                                          | Context is middleware that manipulates request scope context. |
| [Stampede](#stampede)                                                        | Stampede provides a simple cache middleware that is valid for a specified amount of time. |
| [RateLimit/RateLimitPerIP](#ratelimitratelimitperip)                         | RateLimit provides middleware that limits the number of requests processed per second. |
| [MetricsRecorder](#metricsrecorder)                                          | MetricsRecorder provides simple metrics such as request/response size and request duration. |
| [HSTS](#hsts)                                                                | HSTS adds the Strict-Transport-Security header. |
| [Clickjacking](#clickjacking)                                                | Clickjacking mitigates clickjacking attacks by limiting the display of iframe. |
| [ContentSniffing](#contentsniffing)                                          | ContentSniffing adds a header for Content-Type sniffing vulnerability countermeasures. |
| [XSSFiltering](#xssfiltering)                                                | XSSFiltering provides middleware that enables the ability to stop a page from loading when a cross-site scripting attack is detected. |
| [CacheControl/NoCache](#cachecontrolnocache)                                 | CacheControl/NoCache adds the Cache-Control header. |
| [AllowUserAgent/DisallowUserAgent](#allowuseragentdisallowuseragent)         | Allow/DisallowUserAgent middleware controls the request based on the User-Agent header of the request. |
| [AllowContentType/DisallowContentType](#allowcontenttypedisallowcontenttype) | Allow/DisallowContentType middleware controls the request based on the Content-Type header of the request. |
| [AllowAccept/DisallowAccept](#allowacceptdisallowaccept)                     | Allow/DisallowAccept middleware controls the request based on the Accept header of the request. |
| [AllowMethod/DisallowMethod](#allowmethoddisallowmethod)                     | Create an access control using the request method. |
| [RequestHeader/ResponseHeader](#requestheaderresponseheader)                 | Request/ResponseHeader is middleware that edits request and response headers. |
| [Debug](#debug)                                                               | Debug provides middleware that executes the handler only if d is true. |
| [Switch](#switch)                                                             | Switch provides a middleware that executes the next handler if the result of f is true, and executes h if it is false. |
| [ETag](#etag)                                                                 | ETag provides middleware that calculates MD5 from the response data and sets it in the ETag header. |
| [Expires](#expires)                                                           | Expires provides middleware for adding response expiration dates. |
| [Static](#static)                                                             | Provides a handler to deliver static files. |

### Use

Creates a single middleware that executes multiple middleware.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for k := range w.Header() {
			fmt.Fprintf(w, "%s: %s\n", k, w.Header().Get(k))
		}
	})

	m := http.NewServeMux()

	// Creates a single middleware that executes multiple middleware.
	mw := umbrella.Use(
		umbrella.AllowUserAgent("Firefox", "Chrome"),
		umbrella.Clickjacking("deny"),
		umbrella.ContentSniffing(),
		umbrella.NoCache(),
		umbrella.Timeout(time.Millisecond*800),
	)
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### RealIP

Override the RemoteAddr in http.Request with an X-Forwarded-For or X-Real-IP header.

<details>
<summary><b><i>Example :</i></b></summary>

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
		// If an X-Forwarded-For or X-Real-IP header is received,
		// RemoteAddr will be overwritten.
		fmt.Fprintf(w, "RemoteAddr: %v\n", r.RemoteAddr)
		r.Write(w)
	})

	m := http.NewServeMux()

	mw := umbrella.RealIP()
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Recover

Recover from panic and record a stack trace and return a 500 Internal Server Error status.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		if now.Unix()%2 == 0 {
			panic(fmt.Sprintf("panic: %v\n", now))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Time: %v\n", now)
		r.Write(w)
	})

	m := http.NewServeMux()

	// If you give nil, it will be output to os.Stderr.
	mw := umbrella.Recover(nil)
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Timeout

Timeout cancels the context at the given time.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/kenkyu392/umbrella"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := time.Millisecond * time.Duration(rand.Intn(500)+500)
		ctx := r.Context()
		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "duration: %v", d)
	})

	m := http.NewServeMux()

	// This handler times out in 800ms.
	mw := umbrella.Timeout(time.Millisecond * 800)
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Stampede

Stampede provides a simple cache middleware that is valid for a specified amount of time.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kenkyu392/umbrella"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		log.Printf("Search: %s", q)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Search: %s", q)
	})

	m := http.NewServeMux()

	// Create a middleware with a cache that expires in 5 seconds.
	mw := umbrella.Stampede(time.Second * 5)
	m.Handle("/search", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### RateLimit/RateLimitPerIP

RateLimit provides middleware that limits the number of requests processed per second.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kenkyu392/umbrella"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		log.Printf("Search: %s", q)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Search: %s", q)
	})

	m := http.NewServeMux()

	// Create a middleware that processes 5 requests per second.
	mw := umbrella.RateLimit(5) // or umbrella.RateLimitPerIP(5)
	m.Handle("/search", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### MetricsRecorder

MetricsRecorder provides simple metrics such as request/response size and request duration.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kenkyu392/umbrella"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Search: %s", q)
	})

	m := http.NewServeMux()

	// Create a MetricsRecorder and start recording using the middleware.
	mr := umbrella.NewMetricsRecorder(
		umbrella.WithRequestMetricsHookFunc(func(rm *umbrella.RequestMetrics) {
			// You can use hook functions to output request metrics to a log
			// or send them to a monitoring service.
			raw, err := json.Marshal(rm)
			if err != nil {
				log.Println(err)
			}
			log.Printf("%s", raw)
		}),
	)
	m.Handle("/search", mr.Middleware()(handler))
	// You can use MetricsRecorder.Handler to view the metrics.
	// ~$ curl -s http://localhost:3000/metrics | jq .
	m.Handle("/metrics", http.HandlerFunc(mr.Handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### HSTS

HSTS adds the Strict-Transport-Security header.  
Proper use of this header will mitigate stripping attacks.

<details>
<summary><b><i>Example :</i></b></summary>

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
		fmt.Fprintf(w,
			"Strict-Transport-Security: %v",
			w.Header().Get("Strict-Transport-Security"),
		)
	})

	m := http.NewServeMux()

	// Tells the browser to use HTTPS instead of HTTP to connect to a domain
	// (including subdomains).
	mw := umbrella.HSTS(60, "includeSubDomains")
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Clickjacking

Clickjacking mitigates clickjacking attacks by limiting the display of iframe.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"net/http"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// This iframe is not displayed.
		w.Write([]byte(`<iframe src="https://www.google.com/"></iframe>`))
	})

	m := http.NewServeMux()

	// Limit the display of iframe to mitigate clickjacking attacks.
	mw := umbrella.Clickjacking("deny")
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### ContentSniffing

ContentSniffing adds a header for Content-Type sniffing vulnerability countermeasures.

<details>
<summary><b><i>Example :</i></b></summary>

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
		fmt.Fprintf(w,
			"X-Content-Type-Options: %v",
			w.Header().Get("X-Content-Type-Options"),
		)
	})

	m := http.NewServeMux()

	// It implements a countermeasure for Content-Type snuffing vulnerability,
	// which is a problem in old Internet Explorer, for example.
	mw := umbrella.ContentSniffing()
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### XSSFiltering

XSSFiltering provides middleware that enables the ability to stop a page from loading when a cross-site scripting attack is detected.

<details>
<summary><b><i>Example :</i></b></summary>

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
		fmt.Fprintf(w,
			"X-XSS-Protection: %v",
			w.Header().Get("X-XSS-Protection"),
		)
	})

	m := http.NewServeMux()

	// Enable the filter function (stop page loading) for cross-site scripting (XSS).
	// This header is for browsers that do not support Content-Security-Policy.
	mw := umbrella.XSSFiltering("0") // If empty, "1; mode=block" will be set.
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### CacheControl/NoCache

CacheControl/NoCache adds the Cache-Control header.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"github.com/kenkyu392/umbrella"
)

func main() {
	data := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
	<circle cx="50" cy="50" r="40" stroke="#6a737d" stroke-width="4" fill="#1b1f23" />
	</svg>`)
	etag := fmt.Sprintf(`"%x"`, md5.Sum(data))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if match := r.Header.Get("If-None-Match"); strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

	m := http.NewServeMux()

	// Enable browser cache for 2 days.
	// mw := umbrella.NoCache()
	mw := umbrella.CacheControl("public", "max-age=172800", "s-maxage=172800")
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Context

Context is middleware that manipulates request scope context.

<details>
<summary><b><i>Example :</i></b></summary>

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

</details>


### AllowUserAgent/DisallowUserAgent

Allow/DisallowUserAgent middleware controls the request based on the User-Agent header of the request.

<details>
<summary><b><i>Example :</i></b></summary>

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

</details>


### AllowContentType/DisallowContentType

Allow/DisallowContentType middleware controls the request based on the Content-Type header of the request.

<details>
<summary><b><i>Example :</i></b></summary>

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

	allows := umbrella.AllowContentType(
		"application/json", "text/json",
		"application/xml", "text/xml",
	)
	disallows := umbrella.DisallowContentType(
		"text/plain", "application/octet-stream",
	)

	// Only accessible in JSON and XML.
	m.Handle("/allows",
		allows(handler),
	)
	// Not accessible in Plain text and Binary data.
	m.Handle("/disallows",
		disallows(handler),
	)

	http.ListenAndServe(":3000", m)
}
```

</details>


### AllowAccept/DisallowAccept

Allow/DisallowAccept middleware controls the request based on the Accept header of the request.

<details>
<summary><b><i>Example :</i></b></summary>

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

	allows := umbrella.AllowAccept(
		"application/json", "text/json",
	)
	disallows := umbrella.DisallowAccept(
		"text/plain", "text/html",
	)

	// Only accessible in JSON.
	m.Handle("/allows",
		allows(handler),
	)
	// Not accessible in Plain text and HTML data.
	m.Handle("/disallows",
		disallows(handler),
	)

	http.ListenAndServe(":3000", m)
}
```

</details>


### AllowMethod/DisallowMethod

Create an access control using the request method.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"net/http"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Write(w)
	})

	m := http.NewServeMux()

	// Create an access control using the request method.
	mw1 := umbrella.DisallowMethod(http.MethodGet)
	mw2 := umbrella.AllowMethod(http.MethodGet)
	m.Handle("/mw1", mw1(handler))
	m.Handle("/mw2", mw2(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### RequestHeader/ResponseHeader

Request/ResponseHeader is middleware that edits request and response headers.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/kenkyu392/umbrella"
)

const Version = "1.0.0"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "version:%s request:%s response:%s",
			w.Header().Get("X-Version"),
			r.Header.Get("X-Request-Id"),
			w.Header().Get("X-Response-Id"),
		)
	})

	m := http.NewServeMux()

	// You can embed values in request and response headers.
	mw1 := umbrella.RequestHeader(
		func(h http.Header) {
			h.Set("X-Request-Id", generateID())
		},
	)
	mw2 := umbrella.ResponseHeader(
		func(h http.Header) {
			h.Set("X-Response-Id", generateID())
		},
		umbrella.AddHeaderFunc("X-Version", Version),
	)
	m.Handle("/", mw1(mw2(handler)))

	http.ListenAndServe(":3000", m)
}

func generateID() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var b strings.Builder
	for i := 0; i < 6; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return fmt.Sprintf("%s%x", b.String(), time.Now().UnixNano())
}
```

</details>


### Debug

Debug provides middleware that executes the handler only if d is true.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Debug handler implementation.
		w.WriteHeader(http.StatusOK)
	})

	m := http.NewServeMux()

	// Enable the handler for debugging using the value set in the environment variable.
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	mw := umbrella.Debug(debug)
	m.Handle("/debug", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Switch

Switch provides a middleware that executes the next handler if the result of f is true, and executes h if it is false.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"net/http"

	"github.com/kenkyu392/umbrella"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})

	m := http.NewServeMux()

	mw := umbrella.Switch(func(r *http.Request) bool {
		// e.g. checking the prerequisites for a request such as authentication information...
		return true
	}, forbidden)

	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### ETag

ETag provides middleware that calculates MD5 from the response data and sets it in the ETag header.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kenkyu392/umbrella"
)

func main() {
	data := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
	<circle cx="50" cy="50" r="40" stroke="#6a737d" stroke-width="4" fill="#1b1f23" />
	</svg>`)
	etag := fmt.Sprintf(`"%x"`, md5.Sum(data))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if match := r.Header.Get("If-None-Match"); strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

	m := http.NewServeMux()

	mw := umbrella.Use(
		// Automatically calculates and sets the ETag.
		umbrella.ETag(),
		umbrella.ResponseHeader(
			umbrella.CacheControlHeaderFunc("public", "max-age=86400", "no-transform"),
			umbrella.ExpiresHeaderFunc(time.Second*86400),
		),
	)
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


### Expires

Expires provides middleware for adding response expiration dates.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kenkyu392/umbrella"
)

func main() {
	data := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
	<circle cx="50" cy="50" r="40" stroke="#6a737d" stroke-width="4" fill="#1b1f23" />
	</svg>`)
	etag := fmt.Sprintf(`"%x"`, md5.Sum(data))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if match := r.Header.Get("If-None-Match"); strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("ETag", etag)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	})

	m := http.NewServeMux()

	// Enable browser cache for 2 days.
	mw := umbrella.Expires(time.Second * 172800)
	m.Handle("/", mw(handler))

	http.ListenAndServe(":3000", m)
}
```

</details>


## Static

Provides a handler to deliver static files.

<details>
<summary><b><i>Example :</i></b></summary>

```go
package main

import (
	"net/http"
	"time"

	"github.com/kenkyu392/umbrella"
)

func main() {
	m := http.NewServeMux()

	// Add a handler to ServeMux to deliver static files.
	umbrella.Static(m, "/static", "./static",
		// Any middleware can be passed as a variable-length argument.
		// * Sets the ETag.
		umbrella.ETag(),
		// * Add headers for basic cache control, etc.
		umbrella.ResponseHeader(
			umbrella.ClickjackingHeaderFunc("deny"),
			umbrella.ContentSniffingHeaderFunc(),
			umbrella.XSSFilteringHeaderFunc("1; mode=block"),
			umbrella.CacheControlHeaderFunc("public", "max-age=86400", "no-transform"),
			umbrella.ExpiresHeaderFunc(time.Second*86400),
		),
		// * Set the timeout at 2s.
		umbrella.Timeout(time.Second*2),
		// * Enable caching for 60s.
		umbrella.Stampede(time.Second*60),
	)
	// You can also use StaticHandler to create a http.Handler.
	// m.Handle(umbrella.StaticHandler("/static", "./static"))

	http.ListenAndServe(":3000", m)
}
```

</details>


## License

[MIT](LICENSE)
