package log

import (
	"fmt"
	"io"
	"reflect"
	"time"

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
	logger = levels.New(kitlog.NewLogfmtLogger(writer))
	// I feel bad doing this, but the json logger sorts the fields
	// in alphabetical order, so this ensures its first
	logger = logger.With("ats", kitlog.DefaultTimestampUTC)
}

// Public initialization function to initialize the logger global.
// JSON format.
// This should be called before any goroutines using the logger are started
func SetupJSONLoggerTo(writer io.Writer) {
	logger = levels.New(kitlog.NewJSONLogger(writer))
	// This is to keep it consistent with the Json Logger
	logger = logger.With("ats", kitlog.DefaultTimestampUTC)
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
