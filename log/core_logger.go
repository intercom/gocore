package log

import (
	"time"

	"github.com/go-kit/kit/log"
)

type CoreLogger struct {
	log.Logger
	hideTimestamp bool
}

func NewCoreLogger(l log.Logger) *CoreLogger {
	return &CoreLogger{Logger: l}
}

func (cl *CoreLogger) LogInfoMessage(message string, keyvalues ...interface{}) {
	cl.LogInfo(append(keyvalues, "msg", message)...)
}

func (cl *CoreLogger) LogErrorMessage(message string, keyvalues ...interface{}) {
	cl.LogError(append(keyvalues, "msg", message)...)
}

func (cl *CoreLogger) LogInfo(keyvals ...interface{}) {
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	cl.Logger.Log(encodeCompoundValues(append(cl.logTimestamp(keyvals), "level", "info")...)...)
}

func (cl *CoreLogger) LogError(keyvals ...interface{}) {
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	cl.Logger.Log(encodeCompoundValues(append(cl.logTimestamp(keyvals), "level", "error")...)...)
}

func (cl *CoreLogger) SetStandardFields(keyvals ...interface{}) Logger {
	kitLogger := log.With(cl.Logger, keyvals...)
	newLogger := NewCoreLogger(kitLogger)
	newLogger.hideTimestamp = cl.hideTimestamp
	return newLogger
}

func (cl *CoreLogger) With(keyvals ...interface{}) Logger {
	return cl.SetStandardFields(keyvals...)
}

func (cl *CoreLogger) logTimestamp(keyvals []interface{}) []interface{} {
	if !cl.hideTimestamp {
		return append(keyvals, "timestamp", defaultTimeUTC())
	}
	return keyvals
}

func defaultTimeUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}
