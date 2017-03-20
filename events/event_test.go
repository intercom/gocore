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

func TestDynamicField(t *testing.T) {
	sink := &testSink{t: t, events: make(chan string)}
	event := NewEvent("name", sink)

	i := 0
	fn := func() interface{} {
		i += 1
		return i
	}
	event.AddDynamicField("foo", fn)

	event.Send()
	sent := <-sink.events
	if sent != "name - map[foo:1]" {
		t.Errorf("did not send event properly to sink %s", sent)
	}

	event.Send()
	sent = <-sink.events
	if sent != "name - map[foo:2]" { // different result as different sends
		t.Errorf("did not send event properly to sink %s", sent)
	}
}

func TestDynamicFieldTwoSinks(t *testing.T) {
	sink := &testSink{t: t, events: make(chan string)}
	event := NewEvent("name", sink, sink)

	i := 0
	fn := func() interface{} {
		i += 1
		return i
	}
	event.AddDynamicField("foo", fn)

	event.Send()
	sent := <-sink.events
	if sent != "name - map[foo:1]" {
		t.Errorf("did not send event properly to sink %s", sent)
	}

	sent = <-sink.events
	if sent != "name - map[foo:1]" { // same result as multiple sends at once
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
