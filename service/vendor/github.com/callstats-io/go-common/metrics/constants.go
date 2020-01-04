package metrics

// Constants
const (
	InternalMetricsPath = "/internal/metrics"

	// metric names
	MetricRequestCount       = "request_count"
	MetricResponseTime       = "response_time_seconds"
	MetricStreamResponseTime = "stream_response_time_seconds"
)

// Slice can't be made constant but variable works just fine
var (
	ResponseTimeBuckets = []float64{10, 25, 50, 100, 250, 500, 1000, 2000, 4000, 8000}
)
