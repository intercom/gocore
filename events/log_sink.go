package events

import "github.com/intercom/gocore/log"

// LogEventSink will log the event with level "info" to the logger passed in, using the event fields.
type LogEventSink struct {
	logger log.Logger
}

func NewLogEventSink(logger log.Logger) *LogEventSink {
	return &LogEventSink{logger: logger}
}

func (sink *LogEventSink) SendEvent(eventName string, fields map[string]interface{}) {
	keyVals := []interface{}{}
	for key, value := range fields {
		keyVals = append(keyVals, []interface{}{key, value})
	}
	sink.logger.LogInfoMessage(eventName, keyVals...)
}
