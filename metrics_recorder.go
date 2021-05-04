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

// MetricsRecorder provides features for recording and retrieving metrics.
type MetricsRecorder struct {
	m   *Metrics
	rwm *sync.RWMutex

	requestMetricsHookFunc func(*RequestMetrics)
}

// NewMetricsRecorder creates and returns a new MetricsRecorder.
func NewMetricsRecorder(opts ...MetricsRecorderOption) *MetricsRecorder {
	mr := &MetricsRecorder{
		m:                      newMetrics(),
		rwm:                    new(sync.RWMutex),
		requestMetricsHookFunc: func(rm *RequestMetrics) {},
	}
	for _, opt := range opts {
		opt(mr)
	}
	return mr
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
			startTime := time.Now()
			next.ServeHTTP(rec, r)
			endTime := time.Now()
			d := endTime.Sub(startTime)

			// Start recording metrics.
			mr.rwm.Lock()

			uptime := startTime.Sub(processStartTime)
			mr.m.UptimeDurationNanoseconds = uptime.Nanoseconds()
			mr.m.UptimeDurationMilliseconds = uptime.Milliseconds()

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

			// Pass the request metrics to the hook function.
			mr.requestMetricsHookFunc(&RequestMetrics{
				StartTime:                   startTime,
				EndTime:                     endTime,
				Method:                      r.Method,
				Status:                      rec.Code,
				UserAgent:                   r.UserAgent(),
				Referer:                     r.Referer(),
				GoroutinesCount:             goroutines,
				RequestDurationNanoseconds:  ns,
				RequestDurationMilliseconds: ms,
				RequestBytesCount:           requestBytesCount,
				ResponseBytesCount:          responseBytesCount,
			})

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
