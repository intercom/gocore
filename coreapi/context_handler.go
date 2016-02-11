package coreapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/intercom/gocore/log"
	"github.com/intercom/gocore/metrics"
	"github.com/intercom/gocore/monitoring"
	"github.com/pborman/uuid"

	"golang.org/x/net/context"
)

// ContextHandlerFunc is the signature of API handlers that use this system.
type ContextHandlerFunc func(*ContextHandler, http.ResponseWriter, *http.Request)

// A ContextHandler holds per-request state, such as loggers and metrics setup with per-request information.
type ContextHandler struct {
	context.Context
	Cancel      context.CancelFunc
	Logger      *log.CoreLogger
	Metrics     metrics.MetricsRecorder
	Monitor     monitoring.Monitor
	RequestID   string
	handlerFunc ContextHandlerFunc
}

// ServeHTTP makes ContextHandler satisifies the http.Handler interface
func (ch *ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.Context, ch.Cancel = context.WithCancel(context.Background())
	ch.RequestID = uuid.New()
	ch.Context = context.WithValue(ch.Context, "requestID", ch.RequestID)
	ch.Logger = ch.Logger.SetStandardFields("requestID", ch.RequestID)

	// measure timing info
	defer metrics.MeasureSince(fmt.Sprintf("api.%s_%s", r.URL.String(), r.Method), time.Now())

	defer func(ctx context.Context) {
		// Panic recovery
		if rcv := recover(); rcv != nil {
			log.LogErrorMessage("Request Panicked", "status", 500, "requestID", ctx.Value("requestID"), "error", rcv)
			ch.Metrics.IncrementCount(fmt.Sprintf("api.%s_%s.error", r.URL.String(), r.Method))
			err := errors.New(fmt.Sprint(rcv))
			ch.Monitor.CaptureExceptionWithTags(err, "requestID", ctx.Value("requestID"), "endpoint", r.URL.String())
			JSONErrorResponse(500, err).WriteTo(w)
		}
		ch.Cancel()
	}(ch.Context)

	ch.handlerFunc(ch, w, r)
}
