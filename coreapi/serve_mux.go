package coreapi

import (
	"fmt"
	"net/http"

	"github.com/intercom/gocore/log"
	"github.com/intercom/gocore/metrics"
	"github.com/intercom/gocore/monitoring"
)

// ServeMux handles requests matching a given pattern, and routes them appropriately
type ServeMux struct {
	mux     *http.ServeMux
	logger  *log.CoreLogger
	metrics metrics.MetricsRecorder
	monitor monitoring.Monitor
}

// ServeMuxWithDefaults creates a ServeMux with a default CoreLogger, MetricsRecorder and Monitor, which provide
// a baseline for per-request setup.
func ServeMuxWithDefaults(logger *log.CoreLogger, metrics metrics.MetricsRecorder, monitor monitoring.Monitor) *ServeMux {
	logger.UseTimestamp(true) // we always want timestamp fields
	return &ServeMux{
		mux:     http.NewServeMux(),
		logger:  logger,
		metrics: metrics,
		monitor: monitor,
	}
}

// Handle a pattern match, sending to a ContextHandlerFunc
func (bh *ServeMux) Handle(pattern string, f ContextHandlerFunc) {
	bh.mux.Handle(pattern, bh.EndpointFor(pattern, f))
}

// ListenAndServe on the ServeMux
func (bh *ServeMux) ListenAndServe(host, port string) {
	handler := http.HandlerFunc(bh.mux.ServeHTTP)
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), handler)
}

// EndpointFor generates a ContextEndpoint that for a pattern, that has per-endpoint logging, metrics and monitoring
func (bh *ServeMux) EndpointFor(pattern string, f ContextHandlerFunc) *ContextEndpoint {
	return &ContextEndpoint{
		handlerFunc: f,
		logger:      bh.logger.SetStandardFields("url", pattern),
		metrics:     bh.metrics.WithTag("url", pattern),
		monitor:     bh.monitor,
	}
}

// ContextEndpoint holds references to per-endpoint logger, metric and monitor instances
type ContextEndpoint struct {
	handlerFunc ContextHandlerFunc
	logger      *log.CoreLogger
	metrics     metrics.MetricsRecorder
	monitor     monitoring.Monitor
}

// ServeHTTP allows the ContextEndpoint to act as a http.Handler
func (endpoint *ContextEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := &ContextHandler{
		Logger:      endpoint.logger,
		Metrics:     endpoint.metrics,
		Monitor:     endpoint.monitor,
		handlerFunc: endpoint.handlerFunc,
	}
	handler.ServeHTTP(w, r)
}

func (endpoint *ContextEndpoint) Metrics() metrics.MetricsRecorder {
	return endpoint.metrics
}
