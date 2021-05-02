package umbrella

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	const limit = 3
	interval := time.Second / limit
	times := make([]time.Time, 0)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		times = append(times, now)
		t.Logf("%v", now)
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(RateLimit(limit)(handler))
	defer teardown()

	var testCases = []struct {
		want    time.Time
		timeout time.Duration
		err     error
	}{
		{
			timeout: time.Second,
			err:     nil,
		},
		{
			timeout: time.Second,
			err:     nil,
		},
		{
			timeout: time.Millisecond * 100,
			err:     context.DeadlineExceeded,
		},
		{
			timeout: time.Second,
			err:     nil,
		},
	}

	for i := 0; i < len(testCases); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), testCases[i].timeout)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, httpServer.URL, nil)
		_, err := httpClient.Do(req)
		if got := errors.Unwrap(err); (testCases[i].err != nil || got != nil) && !errors.Is(got, testCases[i].err) {
			t.Errorf("want:%v got:%v", testCases[i].err, got)
		}
		cancel()
	}

	for i, t1 := range times {
		if len(times) > i+1 {
			t2 := times[i+1]
			t.Logf("%v >= %v: want:%v+ got: %v",
				t1.Format("15:04:05.000000"),
				t2.Format("15:04:05.000000"),
				interval, t2.Sub(t1))
			if got := t1.Sub(t2); got >= interval {
				t.Errorf("want:%v+ got: %v", interval, got)
			}
		}
	}
}

func TestRateLimitPerIP(t *testing.T) {
	const limit = 3
	interval := time.Second / limit
	times := make([]time.Time, 0)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		times = append(times, now)
		t.Logf("%v", now)
		w.WriteHeader(http.StatusOK)
	})

	teardown := setup(RateLimitPerIP(limit)(handler))
	defer teardown()

	var testCases = []struct {
		want    time.Time
		timeout time.Duration
		err     error
	}{
		{
			timeout: time.Second,
			err:     nil,
		},
		{
			timeout: time.Second,
			err:     nil,
		},
		{
			timeout: time.Millisecond * 100,
			err:     context.DeadlineExceeded,
		},
		{
			timeout: time.Second,
			err:     nil,
		},
	}

	for i := 0; i < len(testCases); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), testCases[i].timeout)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, httpServer.URL, nil)
		_, err := httpClient.Do(req)
		if got := errors.Unwrap(err); (testCases[i].err != nil || got != nil) && !errors.Is(got, testCases[i].err) {
			t.Errorf("want:%v got:%v", testCases[i].err, got)
		}
		cancel()
	}

	for i, t1 := range times {
		if len(times) > i+1 {
			t2 := times[i+1]
			t.Logf("%v >= %v: want:%v+ got: %v",
				t1.Format("15:04:05.000000"),
				t2.Format("15:04:05.000000"),
				interval, t2.Sub(t1))
			if got := t1.Sub(t2); got >= interval {
				t.Errorf("want:%v+ got: %v", interval, got)
			}
		}
	}
}
