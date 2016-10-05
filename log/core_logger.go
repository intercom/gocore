package log

import (
	"time"

	"github.com/go-kit/kit/log/levels"
)

type CoreLogger struct {
	levels.Levels
	useTimestamp bool
	level        LogLevel
}

type LogLevel int

const (
	INFO_LEVEL LogLevel = iota
	ERROR_LEVEL
)

func NewCoreLogger(l levels.Levels) *CoreLogger {
	return &CoreLogger{Levels: l, level: INFO_LEVEL}
}

func (cl *CoreLogger) LogInfoMessage(message string, keyvalues ...interface{}) {
	cl.LogInfo(append(keyvalues, "msg", message)...)
}

func (cl *CoreLogger) LogErrorMessage(message string, keyvalues ...interface{}) {
	cl.LogError(append(keyvalues, "msg", message)...)
}

func (cl *CoreLogger) LogInfo(keyvals ...interface{}) {
	if cl.level > INFO_LEVEL {
		return
	}
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	cl.Levels.Info().Log(encodeCompoundValues(cl.logTimestamp(keyvals)...)...)
}

func (cl *CoreLogger) LogError(keyvals ...interface{}) {
	if cl.level > ERROR_LEVEL {
		return
	}
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	cl.Levels.Error().Log(encodeCompoundValues(cl.logTimestamp(keyvals)...)...)
}

func (cl *CoreLogger) SetStandardFields(keyvals ...interface{}) *CoreLogger {
	encoded := encodeCompoundValues(keyvals...)
	newLogger := NewCoreLogger(cl.Levels.With(encoded...))
	newLogger.UseTimestamp(cl.useTimestamp)
	newLogger.SetLevel(cl.level)
	return newLogger
}

func (cl *CoreLogger) SetLevel(level LogLevel) {
	cl.level = level
}

func (cl *CoreLogger) UseTimestamp(shouldUse bool) {
	cl.useTimestamp = shouldUse
}

func (cl *CoreLogger) logTimestamp(keyvals []interface{}) []interface{} {
	if cl.useTimestamp {
		return append(keyvals, "timestamp", defaultTimeUTC())
	}
	return keyvals
}

func defaultTimeUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}
