package coreapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/intercom/gocore/log"
	"github.com/intercom/gocore/metrics"
	"github.com/intercom/gocore/monitoring"
	"github.com/pborman/uuid"
)

// ErrAuthentication returned when Authentication has failed.
var ErrAuthentication = errors.New("Authentication Error")

func WithBasicAuth(authorisedUser, authorisedPassword string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user, password, ok := r.BasicAuth()
			if !ok {
				JSONErrorResponse(http.StatusForbidden, ErrAuthentication).WriteTo(w)
				return
			}
			if user != authorisedUser || password != authorisedPassword {
				JSONErrorResponse(http.StatusForbidden, ErrAuthentication).WriteTo(w)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// LogRequest logs the start and end of a request
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r)
		if logger == nil {
			next.ServeHTTP(w, r)
			return
		}

		logger.LogInfoMessage("request started")
		next.ServeHTTP(w, r)
		switch v := w.(type) {
		case *StatusWrappingResponseWriter:
			logger.LogInfoMessage("request_ended", "status", v.Status)
		default:
			logger.LogInfoMessage("request_ended")
		}
	})
}

// WithRequestID adds a "requestID" key to the request context.
func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "requestID", uuid.New()))
		next.ServeHTTP(w, r)
	})
}

// WithLogger adds a "logger" key to the request context.
// it will use a requestID as a standard field, if available
func WithLogger(log log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("requestID")
			if requestID != nil {
				log = log.With("requestID", requestID)
			}
			log = log.With("path", r.URL)
			r = r.WithContext(context.WithValue(r.Context(), "logger", log))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// WithMetrics adds a "metrics" key to the request context.
func WithMetrics(metrics metrics.MetricsRecorder) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "metrics", metrics))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// WithMonitor adds a "monitor" key to the request context.
func WithMonitor(monitor monitoring.Monitor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "monitor", monitor))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Get the RequestID stored in request context, or empty.
func GetRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value("requestID").(string)
	if !ok {
		return ""
	}
	return requestID
}

// Get the Logger stored in request context, or nil.
func GetLogger(r *http.Request) log.Logger {
	logger, ok := r.Context().Value("logger").(log.Logger)
	if !ok {
		return nil
	}
	return logger
}

// Get the MetricsRecorder stored in request context, or nil.
func GetMetrics(r *http.Request) metrics.MetricsRecorder {
	metric, ok := r.Context().Value("metrics").(metrics.MetricsRecorder)
	if !ok {
		return nil
	}
	return metric
}

// Get the Monitor stored in request context, or nil.
func GetMonitor(r *http.Request) monitoring.Monitor {
	monitor, ok := r.Context().Value("monitor").(monitoring.Monitor)
	if !ok {
		return nil
	}
	return monitor
}

// StatusWrappingResponseWriter wraps a http.ResponseWriter, overriding WriteHeader to keep
// a record of the Status set.
type StatusWrappingResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (rw *StatusWrappingResponseWriter) WriteHeader(status int) {
	rw.Status = status
	rw.ResponseWriter.WriteHeader(status)
}

// WithStatusWrappingResponseWriter wraps the default response writer to allow tracking of previously written status.
func WithStatusWrappingResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&StatusWrappingResponseWriter{ResponseWriter: w, Status: 0}, r)
	})
}

// Recoverer recovers from panic and prints a log line to the request's logger, if available.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				logger := GetLogger(r)
				if logger != nil {
					logger.LogErrorMessage("Request Panicked", "status", 500, "requestID", r.Context().Value("requestID"), "error", rcv)
				} else {
					debug.PrintStack()
				}
				err := errors.New(fmt.Sprint(rcv))
				JSONErrorResponse(500, err).WriteTo(w)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
