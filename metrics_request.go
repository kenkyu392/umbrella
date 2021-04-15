package umbrella

import (
	"time"
)

// RequestMetrics ...
type RequestMetrics struct {
	StartTime                   time.Time `json:"startTime"`
	EndTime                     time.Time `json:"endTime"`
	Method                      string    `json:"method"`
	Status                      int       `json:"status"`
	UserAgent                   string    `json:"userAgent"`
	Referer                     string    `json:"referer"`
	GoroutinesCount             int64     `json:"goroutinesCount"`
	RequestDurationNanoseconds  int64     `json:"requestDurationNanoseconds"`
	RequestDurationMilliseconds int64     `json:"requestDurationMilliseconds"`
	RequestBytesCount           int64     `json:"requestBytesCount"`
	ResponseBytesCount          int64     `json:"responseBytesCount"`
}

// Clone returns a new RequestMetrics with the same value.
func (r *RequestMetrics) Clone() *RequestMetrics {
	return &RequestMetrics{
		StartTime:                   r.StartTime,
		EndTime:                     r.EndTime,
		Method:                      r.Method,
		Status:                      r.Status,
		UserAgent:                   r.UserAgent,
		Referer:                     r.Referer,
		GoroutinesCount:             r.GoroutinesCount,
		RequestDurationNanoseconds:  r.RequestDurationNanoseconds,
		RequestDurationMilliseconds: r.RequestDurationMilliseconds,
		RequestBytesCount:           r.RequestBytesCount,
		ResponseBytesCount:          r.ResponseBytesCount,
	}
}
