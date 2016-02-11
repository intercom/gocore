package coreapi_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/intercom/gocore/coreapi"
	"github.com/intercom/gocore/log"
	"github.com/intercom/gocore/metrics"
	"github.com/intercom/gocore/monitoring"
)

func TestServeMux(t *testing.T) {
	buf := bytes.Buffer{}
	logger := log.LogfmtLoggerTo(&buf)
	recorder, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8888", "namespace", "hostname")

	mux := coreapi.ServeMuxWithDefaults(logger, recorder, &monitoring.NoopMonitor{})
	handler := &TestHandler{t: t}
	endpoint := mux.EndpointFor("/test", handler.testHandlerFunc)
	ts := httptest.NewServer(endpoint)
	defer ts.Close()

	data := url.Values{}
	data.Set("test", "foo")
	res, _ := http.PostForm(ts.URL, data)

	if 200 != res.StatusCode {
		t.Errorf("should have status %#v, have status %#v", 200, res.StatusCode)
	}

	checkLogFormatMatches(t, fmt.Sprintf("level=error url=/test requestID=%s msg=\"uh oh\"\n", handler.lastRequestID), &buf)

	tags := recorder.GetTags()
	if want, have := "url:/test", tags[0]; want != have {
		t.Errorf("want first tag %#v, have %#v", want, have)
	}

}

type TestHandler struct {
	t             *testing.T
	lastRequestID string
}

func (handler *TestHandler) testHandlerFunc(ctx *coreapi.ContextHandler, w http.ResponseWriter, r *http.Request) {
	handler.lastRequestID = ctx.RequestID
	if "foo" != r.FormValue("test") {
		handler.t.Errorf("should have form data %#v, have form data %#v", "foo", r.FormValue("test"))
	}
	ctx.Logger.LogErrorMessage("uh oh")
}

func checkLogFormatMatches(t *testing.T, want string, buf *bytes.Buffer) {
	have := buf.String()
	if want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
	buf.Reset()
}
