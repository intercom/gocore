package events

import "sync"

// Event wraps some data and some sinks to send it to.
type Event struct {
	*sync.RWMutex
	Name   string
	fields map[string]interface{}
	sinks  []EventSink
}

// EventSink defines an interface for accepting data for an event
type EventSink interface {
	SendEvent(eventName string, fields map[string]interface{})
}

func NewEvent(name string, sinks ...EventSink) *Event {
	return &Event{RWMutex: &sync.RWMutex{}, Name: name, fields: map[string]interface{}{}, sinks: sinks}
}

// Send event to all sinks asynchronously
func (event *Event) Send() {
	for _, sink := range event.sinks {
		go sink.SendEvent(event.Name, event.getFieldsCopy())
	}
}

// Send to a particular event sink asynchronously, won't send to other sinks defined on the event prior
func (event *Event) SendTo(sink EventSink) {
	go sink.SendEvent(event.Name, event.getFieldsCopy())
}

// Add a field to this event
func (event *Event) AddField(key string, value interface{}) {
	event.Lock()
	defer event.Unlock()
	event.fields[key] = value
}

// Get a copy of the fields for this event
func (event *Event) getFieldsCopy() map[string]interface{} {
	event.RLock()
	defer event.RUnlock()

	copyFields := map[string]interface{}{}
	for key, value := range event.fields {
		copyFields[key] = value
	}
	return copyFields
}
