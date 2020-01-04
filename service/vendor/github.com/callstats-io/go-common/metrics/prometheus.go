package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusEndpoint wraps the promhttp.Handler()
func PrometheusEndpoint() http.HandlerFunc {
	return promhttp.Handler().(http.HandlerFunc)
}

// PrometheusEndpointWithoutCompression wraps the promhttp.Handler() without response body compression
func PrometheusEndpointWithoutCompression() http.HandlerFunc {
	return promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: true}).(http.HandlerFunc)
}

// ResetRegistry resets the prometheus registry. Mainly used in tests.
func ResetRegistry() {
	newRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = newRegistry
	prometheus.DefaultGatherer = newRegistry
}
