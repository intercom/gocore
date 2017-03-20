package events

import "sync"

// Event wraps some data and some sinks to send it to.
type Event struct {
	*sync.RWMutex
	Name          string
	fields        map[string]interface{}
	dynamicFields map[string]func() interface{}
	sinks         []EventSink
}

// EventSink defines an interface for accepting data for an event
type EventSink interface {
	SendEvent(eventName string, fields map[string]interface{})
}

func NewEvent(name string, sinks ...EventSink) *Event {
	return &Event{RWMutex: &sync.RWMutex{}, Name: name, fields: map[string]interface{}{}, dynamicFields: map[string]func() interface{}{}, sinks: sinks}
}

// Send event to all sinks asynchronously
func (event *Event) Send() {
	dynamicFields := event.evaluateDynamicFields()
	for _, sink := range event.sinks {
		go sink.SendEvent(event.Name, event.getFieldsCopy(dynamicFields))
	}
}

// Send to a particular event sink asynchronously, won't send to other sinks defined on the event prior
func (event *Event) SendTo(sink EventSink) {
	go sink.SendEvent(event.Name, event.getFieldsCopy(event.evaluateDynamicFields()))
}

// Add a field to this event
func (event *Event) AddField(key string, value interface{}) {
	event.Lock()
	defer event.Unlock()
	event.fields[key] = value
}

// Add a dynamic field to this event; this will be evaluated each time Send() or SendTo() is called.
func (event *Event) AddDynamicField(key string, fn func() interface{}) {
	event.Lock()
	defer event.Unlock()
	event.dynamicFields[key] = fn
}

// Get a copy of the fields for this event
func (event *Event) getFieldsCopy(extraFieldsToCopy map[string]interface{}) map[string]interface{} {
	event.RLock()
	defer event.RUnlock()

	copyFields := map[string]interface{}{}
	for key, value := range event.fields {
		copyFields[key] = value
	}

	for key, value := range extraFieldsToCopy {
		copyFields[key] = value
	}
	return copyFields
}

// Evaluate dynamic fields
func (event *Event) evaluateDynamicFields() map[string]interface{} {
	event.RLock()
	defer event.RUnlock()

	copyFields := map[string]interface{}{}

	for key, fn := range event.dynamicFields {
		copyFields[key] = fn()
	}
	return copyFields
}
