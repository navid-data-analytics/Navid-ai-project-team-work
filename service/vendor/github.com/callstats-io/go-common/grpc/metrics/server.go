package metrics

import (
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

var (
	requestCounter      *prometheus.CounterVec
	responseTimer       *prometheus.HistogramVec
	streamRequestsTimer *prometheus.HistogramVec
)

// Register initializes gRPC metrics and registers them to Prometheus.
func Register(ctx context.Context) error {
	return RegisterWithHistogram(ctx, prometheus.ExponentialBuckets(0.005, 2, 12))
}

// RegisterWithHistogram initializes gRPC metrics with customized histogram and registers them to Prometheus.
func RegisterWithHistogram(ctx context.Context, histogramBucket []float64) error {
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc",
			Name:      metrics.MetricRequestCount,
			Help:      "Total number of RPCs started on the server.",
		},
		[]string{LabelServiceName, LabelMethod},
	)
	var err error
	err = prometheus.Register(requestCounter)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			requestCounter = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			// This should never happen.
			// It triggers only if you have already have registered a metric
			// with the same name and different help text or ConstLabels.
			return err
		}
	}

	responseTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "grpc",
			Name:      metrics.MetricResponseTime,
			Help:      "The gRPC response time in seconds.",
			Buckets:   histogramBucket,
		},
		[]string{LabelServiceName, LabelMethod, LabelStatusCode},
	)
	err = prometheus.Register(responseTimer)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			responseTimer = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			// This should never happen.
			// It triggers only if you have already have registered a metric
			// with the same name and different help text or ConstLabels.
			return err
		}
	}

	streamRequestsTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "grpc",
			Name:      metrics.MetricStreamResponseTime,
			Help:      "The gRPC response time in seconds for stream requests.",
			Buckets:   histogramBucket,
		},
		[]string{LabelServiceName, LabelMethod, LabelStatusCode, LabelStreamDirection},
	)
	err = prometheus.Register(streamRequestsTimer)
	if err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			streamRequestsTimer = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			// This should never happen.
			// It triggers only if you have already have registered a metric
			// with the same name and different help text or ConstLabels.
			return err
		}
	}

	log.FromContextWithPackageName(ctx, "go-common/metrics/grpc").Info("registered gRPC metrics")
	return nil
}
