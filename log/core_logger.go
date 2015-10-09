package log

import "github.com/go-kit/kit/log/levels"

type CoreLogger struct {
	levels.Levels
}

func NewCoreLogger(l levels.Levels) *CoreLogger {
	return &CoreLogger{Levels: l}
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
	cl.Levels.Info(encodeCompoundValues(keyvals...)...)
}

func (cl *CoreLogger) LogError(keyvals ...interface{}) {
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	cl.Levels.Error(encodeCompoundValues(keyvals...)...)
}

func (cl *CoreLogger) SetStandardFields(keyvals ...interface{}) {
	encoded := encodeCompoundValues(keyvals...)
	cl.Levels = cl.Levels.With(encoded...)
}
