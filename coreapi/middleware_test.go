package coreapi_test

import (
	"net/http"
	"testing"

	"os"

	"net/http/httptest"

	"github.com/intercom/gocore/coreapi"
	"github.com/intercom/gocore/log"
	"github.com/intercom/gocore/metrics"
	"github.com/intercom/gocore/monitoring"
)

func TestWithStatusWrappingHandler(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {
		_, ok := w.(*coreapi.StatusWrappingResponseWriter)
		if !ok {
			t.Errorf("was not a status wrapping handler")
		}
	}

	f := coreapi.WithStatusWrappingResponseWriter(http.HandlerFunc(next))
	f.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil))
}

func TestWithRequestID(t *testing.T) {
	next := func(w http.ResponseWriter, r *http.Request) {
		id := coreapi.GetRequestID(r)
		if id == "" {
			t.Errorf("did not find a request id")
		}
	}

	f := coreapi.WithRequestID(http.HandlerFunc(next))
	f.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil))
}

func TestWithLogger(t *testing.T) {
	logger := log.JSONLoggerTo(os.Stderr)
	next := func(w http.ResponseWriter, r *http.Request) {
		found := coreapi.GetLogger(r)
		if found == nil {
			t.Errorf("did not find a logger")
		}
	}

	f := coreapi.WithLogger(logger)(http.HandlerFunc(next))
	f.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil))
}

func TestWithLoggerAndRequestID(t *testing.T) {
	logger := log.JSONLoggerTo(os.Stderr)
	next := func(w http.ResponseWriter, r *http.Request) {
		found := coreapi.GetLogger(r)
		if found == nil {
			t.Errorf("did not find a logger")
		}
	}

	f := coreapi.WithRequestID(coreapi.WithLogger(logger)(http.HandlerFunc(next)))
	f.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil))
}

func TestWithMetrics(t *testing.T) {
	metr := &metrics.NoopRecorder{}
	next := func(w http.ResponseWriter, r *http.Request) {
		found := coreapi.GetMetrics(r)
		if found != metr {
			t.Errorf("did not find a metrics object")
		}
	}

	f := coreapi.WithMetrics(metr)(http.HandlerFunc(next))
	f.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil))
}

func TestWithMonitor(t *testing.T) {
	monit := &monitoring.NoopMonitor{}
	next := func(w http.ResponseWriter, r *http.Request) {
		found := coreapi.GetMonitor(r)
		if found != monit {
			t.Errorf("did not find a monitor object")
		}
	}

	f := coreapi.WithMonitor(monit)(http.HandlerFunc(next))
	f.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/foo", nil))
}
