package events

import (
	"fmt"
	"testing"
)

func TestNewEvent(t *testing.T) {
	sink := &testSink{t: t, events: make(chan string)}
	event := NewEvent("name", sink)
	if event.Name != "name" {
		t.Errorf("name not set")
	}

	event.AddField("foo", "bar")
	event.Send()
	sent := <-sink.events
	if sent != "name - map[foo:bar]" {
		t.Errorf("did not send event properly to sink %s", sent)
	}
}

func TestEventSendTo(t *testing.T) {
	sink := &testSink{t: t, events: make(chan string)}
	event := NewEvent("name")
	event.AddField("foo", "bar")
	event.SendTo(sink)
	sent := <-sink.events
	if sent != "name - map[foo:bar]" {
		t.Errorf("did not send event properly to sink %s", sent)
	}
}

func TestEventTwoSinks(t *testing.T) {
	sink := &testSink{t: t, events: make(chan string)}
	event := NewEvent("name", sink, sink)
	event.AddField("foo", "bar")
	event.Send()
	sent := <-sink.events
	if sent != "name - map[foo:bar]" {
		t.Errorf("did not send event properly to sink %s", sent)
	}

	sent = <-sink.events
	if sent != "name - map[foo:bar]" {
		t.Errorf("did not send event properly to sink %s", sent)
	}
}

type testSink struct {
	t      *testing.T
	events chan string
}

func (sink *testSink) SendEvent(eventName string, fields map[string]interface{}) {
	sink.events <- fmt.Sprintf("%s - %v", eventName, fields)
}
