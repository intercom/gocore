package log

import (
	"fmt"
	"io"
	"reflect"

	kitlog "github.com/go-kit/kit/log"
)

var (
	GlobalLogger Logger
)

// Logger represents a logger with inheritable context
type Logger interface {
	LogInfoMessage(message string, keyvalues ...interface{})
	LogErrorMessage(message string, keyvalues ...interface{})
	LogInfo(keyvals ...interface{})
	LogError(keyvals ...interface{})
	With(keyvals ...interface{}) Logger
}

// Public initialization function to initialize the logger global.
// Logfmt format.
// This should be called before any goroutines using the logger are started
func SetupLogFmtLoggerTo(writer io.Writer) {
	l := LogfmtLoggerTo(writer)
	GlobalLogger = l
}

// Public initialization function to initialize the logger global.
// JSON format.
// This should be called before any goroutines using the logger are started
func SetupJSONLoggerTo(writer io.Writer) {
	l := JSONLoggerTo(writer)
	GlobalLogger = l
}

// Log a message to Info, with optional keyvalues
func LogInfoMessage(message string, keyvalues ...interface{}) {
	GlobalLogger.LogInfoMessage(message, keyvalues...)
}

// Log a message to Error, with optional keyvalues
func LogErrorMessage(message string, keyvalues ...interface{}) {
	GlobalLogger.LogErrorMessage(message, keyvalues...)
}

// Log a series of key, values to Info
func LogInfo(keyvals ...interface{}) {
	GlobalLogger.LogInfo(keyvals...)
}

// Log a series of key, values to Error
func LogError(keyvals ...interface{}) {
	GlobalLogger.LogError(keyvals...)
}

// Sets standard fields on the logger, for all calls
func SetStandardFields(keyvals ...interface{}) {
	GlobalLogger = GlobalLogger.With(keyvals...)
}

func JSONLoggerTo(writer io.Writer) *CoreLogger {
	return NewCoreLogger(kitlog.NewJSONLogger(writer))
}

func LogfmtLoggerTo(writer io.Writer) *CoreLogger {
	return NewCoreLogger(kitlog.NewLogfmtLogger(writer))
}

func NoopLogger() *CoreLogger {
	return NewCoreLogger(kitlog.NewNopLogger())
}

// Encode compound values using %+v. To use a custom encoding, use a type that implements fmt.Stringer
func encodeCompoundValues(keyvals ...interface{}) []interface{} {
	if len(keyvals)%2 == 1 {
		keyvals = append(keyvals, nil) // missing a value
	}

	for i := 0; i < len(keyvals); i += 2 {
		_, v := keyvals[i], keyvals[i+1]

		rvalue := reflect.ValueOf(v)
		switch rvalue.Kind() {
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
			keyvals[i+1] = fmt.Sprintf("%+v", v)
		}
	}
	return keyvals
}

// Package-level default initialization of the logger.
// Initializes it to a no-op implementation;
// later calls can replace it by calling SetupLogger.
func init() {
	GlobalLogger = NoopLogger()
}
