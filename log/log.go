package log

import (
	"io"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/levels"
)

var (
	logger levels.Levels // level based logger
)

// Public initialization function to initialize the logger global.
// This should be called before any goroutines using the logger are started
func SetupLoggerTo(writer io.Writer) {
	logger = levels.New(kitlog.NewLogfmtLogger(writer))
}

// Convenience setup for Stderr
func SetupLoggerToStderr() {
	SetupLoggerTo(os.Stderr)
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
	logger.Info(keyvals...)
}

// Log a series of key, values to Error
func LogError(keyvals ...interface{}) {
	logger.Error(keyvals...)
}

// Sets standard fields on the logger, for all calls
func SetStandardFields(keyvals ...interface{}) {
	logger = logger.With(keyvals...)
}

// Package-level default initialization of the logger.
// Initializes it to a no-op implementation;
// later calls can replace it by calling SetupLogger.
func init() {
	logger = levels.New(kitlog.NewNopLogger())
}
