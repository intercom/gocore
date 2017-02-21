package events

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/intercom/gocore/log"
)

func TestLogSink(t *testing.T) {
	buf := bytes.Buffer{}
	logger := log.JSONLoggerTo(&buf)

	sink := &testLogEventSink{t: t, done: make(chan bool), logSink: NewLogEventSink(logger)}
	event := NewEvent("eventname", sink)

	event.AddField("foo", "bar")
	event.Send()

	<-sink.done // wait for event to be sent

	el := eventLog{}
	json.Unmarshal(buf.Bytes(), &el)
	if el.Name != "eventname" {
		t.Errorf("did not send json format log for msg=eventname")
	}
	if el.Foo != "bar" {
		t.Errorf("did not send json format log for foo=bar")
	}
}

type testLogEventSink struct {
	t       *testing.T
	done    chan bool
	logSink *LogEventSink
}

func (sink *testLogEventSink) SendEvent(eventName string, fields map[string]interface{}) {
	sink.logSink.SendEvent(eventName, fields)
	sink.done <- true
}

type eventLog struct {
	Name string `json:"msg"`
	Foo  string `json:"foo"`
}
