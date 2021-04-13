package umbrella

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"time"
)

// Metrics ...
type Metrics struct {
	GoroutinesTotalCount int64 `json:"goroutinesTotalCount"`
	MaxGoroutinesCount   int64 `json:"maxGoroutinesCount"`
	MinGoroutinesCount   int64 `json:"minGoroutinesCount"`
	AvgGoroutinesCount   int64 `json:"avgGoroutinesCount"`

	RequestsTotalCount int64 `json:"requestsTotalCount"`

	TotalRequestDurationNanoseconds int64 `json:"totalRequestDurationNanoseconds"`
	MaxRequestDurationNanoseconds   int64 `json:"maxRequestDurationNanoseconds"`
	MinRequestDurationNanoseconds   int64 `json:"minRequestDurationNanoseconds"`
	AvgRequestDurationNanoseconds   int64 `json:"avgRequestDurationNanoseconds"`

	TotalRequestDurationMilliseconds int64 `json:"totalRequestDurationMilliseconds"`
	MaxRequestDurationMilliseconds   int64 `json:"maxRequestDurationMilliseconds"`
	MinRequestDurationMilliseconds   int64 `json:"minRequestDurationMilliseconds"`
	AvgRequestDurationMilliseconds   int64 `json:"avgRequestDurationMilliseconds"`

	MaxRequestBytesCount  int64 `json:"maxRequestBytesCount"`
	MinRequestBytesCount  int64 `json:"minRequestBytesCount"`
	MaxResponseBytesCount int64 `json:"maxResponseBytesCount"`
	MinResponseBytesCount int64 `json:"minResponseBytesCount"`

	MethodCount      map[string]int64 `json:"methodCount"`
	StatusCount      map[int]int64    `json:"statusCount"`
	StatusClassCount map[string]int64 `json:"statusClassCount"`
}

// Clone returns a new Metrics with the same value.
func (m *Metrics) Clone() *Metrics {
	m2 := &Metrics{
		GoroutinesTotalCount:             m.GoroutinesTotalCount,
		MaxGoroutinesCount:               m.MaxGoroutinesCount,
		MinGoroutinesCount:               m.MinGoroutinesCount,
		AvgGoroutinesCount:               m.AvgGoroutinesCount,
		RequestsTotalCount:               m.RequestsTotalCount,
		TotalRequestDurationNanoseconds:  m.TotalRequestDurationNanoseconds,
		MaxRequestDurationNanoseconds:    m.MaxRequestDurationNanoseconds,
		MinRequestDurationNanoseconds:    m.MinRequestDurationNanoseconds,
		AvgRequestDurationNanoseconds:    m.AvgRequestDurationNanoseconds,
		TotalRequestDurationMilliseconds: m.TotalRequestDurationMilliseconds,
		MaxRequestDurationMilliseconds:   m.MaxRequestDurationMilliseconds,
		MinRequestDurationMilliseconds:   m.MinRequestDurationMilliseconds,
		AvgRequestDurationMilliseconds:   m.AvgRequestDurationMilliseconds,
		MaxRequestBytesCount:             m.MaxRequestBytesCount,
		MinRequestBytesCount:             m.MinRequestBytesCount,
		MaxResponseBytesCount:            m.MaxResponseBytesCount,
		MinResponseBytesCount:            m.MinResponseBytesCount,
		MethodCount:                      make(map[string]int64),
		StatusCount:                      make(map[int]int64),
		StatusClassCount:                 make(map[string]int64),
	}
	for k, v := range m.MethodCount {
		m2.MethodCount[k] = v
	}
	for k, v := range m.StatusCount {
		m2.StatusCount[k] = v
	}
	for k, v := range m.StatusClassCount {
		m2.StatusClassCount[k] = v
	}
	return m2
}

func newMetrics() *Metrics {
	m := &Metrics{
		MethodCount: map[string]int64{
			http.MethodGet:     0,
			http.MethodHead:    0,
			http.MethodPost:    0,
			http.MethodPut:     0,
			http.MethodPatch:   0,
			http.MethodDelete:  0,
			http.MethodConnect: 0,
			http.MethodOptions: 0,
			http.MethodTrace:   0,
		},
		StatusCount: make(map[int]int64),
		StatusClassCount: map[string]int64{
			"1xx": 0,
			"2xx": 0,
			"3xx": 0,
			"4xx": 0,
			"5xx": 0,
		},
	}
	for i := 0; i < 600; i++ {
		if http.StatusText(i) != "" {
			m.StatusCount[i] = 0
		}
	}
	return m
}

// MetricsRecorder provides features for recording and retrieving metrics.
type MetricsRecorder struct {
	m   *Metrics
	rwm *sync.RWMutex
}

// NewMetricsRecorder creates and returns a new MetricsRecorder.
func NewMetricsRecorder() *MetricsRecorder {
	return &MetricsRecorder{
		m:   newMetrics(),
		rwm: new(sync.RWMutex),
	}
}

// Metrics ...
func (mr *MetricsRecorder) Metrics() *Metrics {
	return mr.m.Clone()
}

// Middleware records metrics.
// If an error occurs, call the next handler without recording any metrics.
func (mr *MetricsRecorder) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("metrics.middleware.error: %#v", err)
				next.ServeHTTP(w, r)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// Use ResponseRecorder to record the results.
			rec := httptest.NewRecorder()
			now := time.Now()
			next.ServeHTTP(rec, r)
			d := time.Since(now)

			// Start recording metrics.
			mr.rwm.Lock()

			mr.m.RequestsTotalCount++

			// Measure the body size of the request/response.
			requestBytesCount := int64(len(requestBody))
			if mr.m.MaxRequestBytesCount < requestBytesCount || mr.m.MaxRequestBytesCount == 0 {
				mr.m.MaxRequestBytesCount = requestBytesCount
			}
			if mr.m.MinRequestBytesCount > requestBytesCount || mr.m.MinRequestBytesCount == 0 {
				mr.m.MinRequestBytesCount = requestBytesCount
			}

			responseBody := rec.Body.Bytes()
			responseBytesCount := int64(len(responseBody))
			if mr.m.MaxResponseBytesCount < responseBytesCount || mr.m.MaxResponseBytesCount == 0 {
				mr.m.MaxResponseBytesCount = responseBytesCount
			}
			if mr.m.MinResponseBytesCount > responseBytesCount || mr.m.MinResponseBytesCount == 0 {
				mr.m.MinResponseBytesCount = responseBytesCount
			}

			// Measure request status and methods.
			mr.m.MethodCount[r.Method]++
			mr.m.StatusCount[rec.Code]++
			switch {
			case 100 <= rec.Code && rec.Code < 200:
				// Informational
				mr.m.StatusClassCount["1xx"]++
			case 200 <= rec.Code && rec.Code < 300:
				// Successful
				mr.m.StatusClassCount["2xx"]++
			case 300 <= rec.Code && rec.Code < 400:
				// Redirection
				mr.m.StatusClassCount["3xx"]++
			case 400 <= rec.Code && rec.Code < 500:
				// Client Error
				mr.m.StatusClassCount["4xx"]++
			case 500 <= rec.Code && rec.Code < 600:
				// Server Error
				mr.m.StatusClassCount["5xx"]++
			}

			// Measure the duration of the request.
			ns := d.Nanoseconds()
			ms := d.Milliseconds()
			mr.m.TotalRequestDurationNanoseconds += ns
			mr.m.TotalRequestDurationMilliseconds += ms

			mr.m.AvgRequestDurationNanoseconds = mr.m.TotalRequestDurationNanoseconds / mr.m.RequestsTotalCount
			if mr.m.MaxRequestDurationNanoseconds < ns || mr.m.MaxRequestDurationNanoseconds == 0 {
				mr.m.MaxRequestDurationNanoseconds = ns
			}
			if mr.m.MinRequestDurationNanoseconds > ns || mr.m.MinRequestDurationNanoseconds == 0 {
				mr.m.MinRequestDurationNanoseconds = ns
			}

			mr.m.AvgRequestDurationMilliseconds = mr.m.TotalRequestDurationMilliseconds / mr.m.RequestsTotalCount
			if mr.m.MaxRequestDurationMilliseconds < ms || mr.m.MaxRequestDurationMilliseconds == 0 {
				mr.m.MaxRequestDurationMilliseconds = ms
			}
			if mr.m.MinRequestDurationMilliseconds > ms || mr.m.MinRequestDurationMilliseconds == 0 {
				mr.m.MinRequestDurationMilliseconds = ms
			}
			if mr.m.MinRequestDurationMilliseconds > ms || mr.m.MinRequestDurationMilliseconds == 0 {
				mr.m.MinRequestDurationMilliseconds = ms
			}

			goroutines := int64(runtime.NumGoroutine())
			mr.m.GoroutinesTotalCount += goroutines
			mr.m.AvgGoroutinesCount = mr.m.GoroutinesTotalCount / mr.m.RequestsTotalCount
			if mr.m.MaxGoroutinesCount < goroutines || mr.m.MaxGoroutinesCount == 0 {
				mr.m.MaxGoroutinesCount = goroutines
			}
			if mr.m.MinGoroutinesCount > goroutines || mr.m.MinGoroutinesCount == 0 {
				mr.m.MinGoroutinesCount = goroutines
			}

			mr.rwm.Unlock()

			// Write response.
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			_, _ = w.Write(responseBody)
		}
		return http.HandlerFunc(fn)
	}
}

// Handler returns metrics in JSON.
func (mr *MetricsRecorder) Handler(w http.ResponseWriter, r *http.Request) {
	raw, _ := json.MarshalIndent(mr.m, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}
