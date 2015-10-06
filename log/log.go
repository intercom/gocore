package log

import (
	"fmt"
	"io"
	"reflect"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/levels"
)

var (
	logger levels.Levels // level based logger
)

// Public initialization function to initialize the logger global.
// Logfmt format.
// This should be called before any goroutines using the logger are started
func SetupLogFmtLoggerTo(writer io.Writer) {
	l := logfmtLoggerTo(writer)
	l = l.With("aaats", kitlog.DefaultTimestampUTC)
	SetupGlobalLoggerTo(l)
}

// Public initialization function to initialize the logger global.
// JSON format.
// This should be called before any goroutines using the logger are started
func SetupJSONLoggerTo(writer io.Writer) {
	l := jsonLoggerTo(writer)
	// cloudwatch takes the first instance of a timestamp matching the format.
	// the json logger sorts alphabetically, so this ensures its first
	l = l.With("aaats", kitlog.DefaultTimestampUTC)
	SetupGlobalLoggerTo(l)
}

// Public initialization function to set a logger global.
func SetupGlobalLoggerTo(inLogger levels.Levels) {
	logger = inLogger
}

// Log a message to Info, with optional keyvalues
func LogInfoMessage(message string, keyvalues ...interface{}) {
	LogInfo(append(keyvalues, "msg", message)...)
}

// Log a message to Error, with optional keyvalues
func LogErrorMessage(message string, keyvalues ...interface{}) {
	LogError(append(keyvalues, "msg", message)...)
}

// Log a series of key, values to Info
func LogInfo(keyvals ...interface{}) {
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	logger.Info(encodeCompoundValues(keyvals...)...)
}

// Log a series of key, values to Error
func LogError(keyvals ...interface{}) {
	if len(keyvals) == 1 {
		keyvals = []interface{}{"msg", keyvals[0]}
	}
	logger.Error(encodeCompoundValues(keyvals...)...)
}

// Sets standard fields on the logger, for all calls
func SetStandardFields(keyvals ...interface{}) {
	logger = logger.With(encodeCompoundValues(keyvals...)...)
}

func jsonLoggerTo(writer io.Writer) levels.Levels {
	return levels.New(kitlog.NewJSONLogger(writer))
}

func logfmtLoggerTo(writer io.Writer) levels.Levels {
	return levels.New(kitlog.NewLogfmtLogger(writer))
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
	logger = levels.New(kitlog.NewNopLogger())
}
