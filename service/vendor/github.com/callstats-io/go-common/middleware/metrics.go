package middleware

import (
	"net/http"

	"github.com/callstats-io/go-common/metrics"
)

// Metrics registers labeled HTTP request metrics instrumentation for a given http.HandlerFunc
func Metrics(handlerLabel string, constLabels ...*metrics.Label) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return metrics.InstrumentHandler(handlerLabel, constLabels, h)
	}
}
