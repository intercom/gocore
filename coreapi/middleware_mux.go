package coreapi

import (
	"net/http"

	"github.com/intercom/gocore/log"
	"github.com/intercom/gocore/metrics"
	"github.com/intercom/gocore/monitoring"
)

// MiddlewareMux wraps the default router with a couple of methods to allow easier middleware usage.
type MiddlewareMux struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func NewMiddlewareMux() *MiddlewareMux {
	return &MiddlewareMux{
		mux:         http.NewServeMux(),
		middlewares: []func(http.Handler) http.Handler{},
	}
}

// NewMiddlewareMuxWithDefaults sets up a basic middleware mux with default middleware of logger, metrics, monitoring, request id and status wrapping
func NewMiddlewareMuxWithDefaults(logger log.Logger, recorder metrics.MetricsRecorder, monitor monitoring.Monitor) *MiddlewareMux {
	mux := NewMiddlewareMux()
	mux.Use(WithStatusWrappingResponseWriter)
	mux.Use(WithRequestID)
	if logger == nil {
		logger = log.NoopLogger()
	}
	mux.Use(WithLogger(logger))

	if recorder == nil {
		recorder = &metrics.NoopRecorder{}
	}
	mux.Use(WithMetrics(recorder))

	if monitor == nil {
		monitor = &monitoring.NoopMonitor{}
	}
	mux.Use(WithMonitor(monitor))
	mux.Use(Recoverer)
	return mux
}

// Use a piece of middleware
func (mm *MiddlewareMux) Use(middleware func(http.Handler) http.Handler) {
	mm.middlewares = append(mm.middlewares, middleware)
}

// Handle a route
func (mm *MiddlewareMux) Handle(pattern string, handler http.Handler) {
	mm.mux.Handle(pattern, mm.buildHandler(handler))
}

func (mm *MiddlewareMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mm.mux.ServeHTTP(w, r)
}

func (mm *MiddlewareMux) buildHandler(base http.Handler) http.Handler {
	var f http.Handler = base
	for i := len(mm.middlewares) - 1; i >= 0; i-- {
		f = mm.middlewares[i](f)
	}
	return f
}
