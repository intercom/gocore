package events

import libhoney "github.com/honeycombio/libhoney-go"

// HoneycombSink will send the event to HoneyComb.io, with fields added.
type HoneycombSink struct {
	builder *libhoney.Builder
}

func NewHoneycombSink(builder *libhoney.Builder) *HoneycombSink {
	return &HoneycombSink{builder: builder}
}

func (sink *HoneycombSink) SendEvent(eventName string, fields map[string]interface{}) {
	honeyEvent := sink.builder.NewEvent()
	honeyEvent.AddField("name", eventName)
	for key, value := range fields {
		honeyEvent.AddField(key, value)
	}
	honeyEvent.Send()
}
