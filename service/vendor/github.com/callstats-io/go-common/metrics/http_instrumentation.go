// A little modified version of https://github.com/orian/go-http-instrument/blob/master/instrumentation/instrument.go

package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Label is a representation of metrics labels as key-value pair
type Label struct {
	Key   string
	Value string
}

// NewAppLabels builds a slice of labels in app scope with keys "app" and "env"
func NewAppLabels(app, env string) []*Label {
	return []*Label{
		{Key: "app", Value: app},
		{Key: "env", Value: env},
	}
}

// InstrumentHandler wraps the given HTTP handler for instrumentation. It
// registers four metric collectors (if not already done) and reports HTTP
// metrics to the (newly or already) registered collectors: http_inflight_requests
// (GaugeVec), http_request_count (CounterVec), http_request_request_time
// (Histogram), http_response_size_bytes (Summary). Each has a set of arbitrary labels.
// http_request_count is a metric vector partitioned by HTTP method
// (label name "method") and HTTP status code (label name "status_code").
func InstrumentHandler(handlerName string, labels []*Label, handler http.Handler) http.HandlerFunc {
	return InstrumentHandlerFunc(handlerName, labels, handler.ServeHTTP)
}

// InstrumentHandlerFunc wraps the given function for instrumentation. It
// otherwise works in the same way as InstrumentHandler (and shares the same
// issues).
func InstrumentHandlerFunc(handlerName string, labels []*Label, handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	constLabels := prometheus.Labels{}
	for _, label := range labels {
		constLabels[label.Key] = label.Value
	}
	return InstrumentHandlerFuncWithOpts(
		handlerName,
		prometheus.Opts{
			Namespace:   "http",
			ConstLabels: constLabels,
		},
		handlerFunc,
	)
}

// InstrumentHandlerWithOpts works like InstrumentHandler (and shares the same
// issues) but provides more flexibility (at the cost of a more complex call
// syntax). As InstrumentHandler, this function registers four metric
// collectors, but it uses the provided SummaryOpts to create them. However, the
// fields "Name" and "Help" in the SummaryOpts are ignored. "Name" is replaced
// by "inflight_requests", "request_count", "request_time" and
// "response_size_bytes", respectively. "Help" is replaced by an appropriate
// help string. The names of the variable labels of the http_requests_total
// CounterVec are "method" (get, post, etc.), and "status_code" (HTTP status code).
//
// If InstrumentHandlerWithOpts is called as follows, it mimics exactly the
// behavior of InstrumentHandler:
//
//     prometheus.InstrumentHandlerWithOpts(
//         handlerName,
//         prometheus.SummaryOpts{
//              Subsystem:   "http",
//              ConstLabels: prometheus.Labels{"app": serviceName, "env": env},
//         },
//         handler,
//     )
//
// Technical detail: "request_count" is a CounterVec, not a SummaryVec, so it
// cannot use SummaryOpts. Instead, a CounterOpts struct is created internally,
// and all its fields are set to the equally named fields in the provided
// SummaryOpts.
func InstrumentHandlerWithOpts(handlerName string, opts prometheus.Opts, handler http.Handler) http.HandlerFunc {
	return InstrumentHandlerFuncWithOpts(handlerName, opts, handler.ServeHTTP)
}

// InstrumentHandlerFuncWithOpts works like InstrumentHandlerFunc (and shares
// the same issues) but provides more flexibility (at the cost of a more complex
// call syntax). See InstrumentHandlerWithOpts for details how the provided
// Opts are used.
func InstrumentHandlerFuncWithOpts(handlerName string, opts prometheus.Opts, handlerFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   opts.Namespace,
			Subsystem:   opts.Subsystem,
			Name:        MetricRequestCount,
			Help:        "Total number of HTTP requests made.",
			ConstLabels: opts.ConstLabels,
		},
		[]string{"handler"},
	)
	if err := prometheus.Register(requestCounter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			requestCounter = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			// This should never happen.
			// It will trigger if you try to instrument new handler metrics
			// with already existing metric name but with different help text or ConstLabels.
			panic(err)
		}
	}

	responseTimer := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   opts.Namespace,
			Subsystem:   opts.Subsystem,
			Name:        MetricResponseTime,
			Help:        "The HTTP response time in seconds.",
			ConstLabels: opts.ConstLabels,
			Buckets:     prometheus.ExponentialBuckets(0.005, 2, 12),
		},
		[]string{"handler", "status_code"},
	)
	if err := prometheus.Register(responseTimer); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			responseTimer = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			// This should never happen.
			// It will trigger if you try to instrument new handler metrics
			// with already existing metric name but with different help text or ConstLabels.
			panic(err)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter.WithLabelValues(handlerName).Inc()
		now := time.Now()

		rw := &responseWriterDelegator{ResponseWriter: w}
		handlerFunc(rw, r)

		elapsed := float64(time.Since(now)) / float64(time.Second)
		statusCode := sanitizeCode(rw.status)
		responseTimer.WithLabelValues(handlerName, statusCode).Observe(elapsed)
	})
}

type responseWriterDelegator struct {
	http.ResponseWriter

	handler, method string
	status          int
	written         int64
	wroteHeader     bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func sanitizeCode(s int) string {
	switch s {
	case 100:
		return "100"
	case 101:
		return "101"

	case 200:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"

	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 511:
		return "511"

	default:
		return strconv.Itoa(s)
	}
}
