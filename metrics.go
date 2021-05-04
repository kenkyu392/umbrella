package umbrella

import (
	"net/http"
)

// Metrics ...
type Metrics struct {
	UptimeDurationNanoseconds  int64 `json:"uptimeDurationNanoseconds"`
	UptimeDurationMilliseconds int64 `json:"uptimeDurationMilliseconds"`

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
		UptimeDurationNanoseconds:        m.UptimeDurationNanoseconds,
		UptimeDurationMilliseconds:       m.UptimeDurationMilliseconds,
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
