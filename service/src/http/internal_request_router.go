package http

import (
	"encoding/json"
	"net/http"
	"time"

	"context"

	"github.com/callstats-io/ai-decision/service/src/config"
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/metrics"
	"github.com/callstats-io/go-common/response"
)

var startTime = time.Now().Format(time.RFC3339)

// service states
const (
	StatusUp   = "UP"
	StatusDown = "DOWN"
)

// StatusChecker defines the function type for status check functions.
// If any StatusChecker returns an error, status is set as "DOWN".
type StatusChecker func(ctx context.Context) error

// Status implements stats response
type Status struct {
	Status    string `json:"status"`           // "UP" if all connections to dependent services are ok, "DOWN" otherwise
	BuildTime string `json:"imageBuildTime"`   // Build time
	StartTime string `json:"serviceStartTime"` // Service instance start time
	Version   string `json:"serviceVersion"`   // version from Git
}

// InternalRequestRouter handles status requests
type InternalRequestRouter struct {
	statusHandler  http.Handler
	metricsHandler http.Handler
	checkers       []StatusChecker
}

// NewInternalRequestRouter returns a new status handler
func NewInternalRequestRouter(metricsHandler http.Handler, checkers ...StatusChecker) *InternalRequestRouter {
	sh := &InternalRequestRouter{
		statusHandler:  &statusHandler{checkers},
		metricsHandler: metricsHandler,
		checkers:       checkers,
	}
	return sh
}
func (s *InternalRequestRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.RequestURI {
	case "/status":
		s.statusHandler.ServeHTTP(w, r)
	case metrics.InternalMetricsPath:
		s.metricsHandler.ServeHTTP(w, r)
	default:
		response.NotFound(w, nil)
	}
}

type statusHandler struct {
	checkers []StatusChecker
}

// Status responds if the system is functional or not, TODO: healthcheck should be async
func (sh *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := &Status{
		Status:    StatusUp,
		BuildTime: config.ImageBuildTime,
		StartTime: startTime,
		Version:   config.ServiceVersion,
	}

	for _, check := range sh.checkers {
		if err := check(r.Context()); err != nil {
			log.FromContext(r.Context()).Error("failed status", log.Error(err))
			status.Status = StatusDown
			break
		}
	}

	payload, _ := json.Marshal(status)
	if status.Status == StatusUp {
		response.OK(w, payload)
	} else {
		response.RequiredServiceUnavailable(w, payload)
	}
}
