package umbrella

// MetricsRecorderOption ...
type MetricsRecorderOption func(mr *MetricsRecorder)

// WithRequestMetricsHookFunc sets the hook function to be called for each request.
// The hook function receives the request metrics as an argument.
func WithRequestMetricsHookFunc(fn func(*RequestMetrics)) MetricsRecorderOption {
	return func(mr *MetricsRecorder) {
		if fn != nil {
			mr.requestMetricsHookFunc = fn
		}
	}
}
