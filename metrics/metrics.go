package metrics

import "time"

// MetricsRecorder global instance.
var globalMetrics MetricsRecorder

// Public interface for recording metrics.
type MetricsRecorder interface {
	IncrementCount(metricName string)
	IncrementCountBy(metricName string, amount int)
	MeasureSince(metricName string, since time.Time)
	MeasureDurationMS(metricName string, durationMS float32)
	SetGauge(metricName string, val float32)
	SetPrefix(prefix string)
	WithTag(key, value string) MetricsRecorder
}

// Package-level default initialization of the Metrics global.
// Initializes it to a no-op implementation;
// later calls can replace it by calling SetMetricsGlobal.
func init() {
	globalMetrics = &NoopRecorder{}
}

// Public initialization function to initialize the Metrics global.
// If you're using metrics, this should be called before any goroutines
// using them are started.
//
// (If you don't care about metrics, you don't need to call this function;
// nothing will break, since a no-op metrics sink is used by default.)
func SetMetricsGlobal(recorder MetricsRecorder) {
	globalMetrics = recorder
}

// Increment Count by 1 for Metric by name
func IncrementCount(metricName string) {
	globalMetrics.IncrementCount(metricName)
}

// Increment Count by amount for Metric by name
func IncrementCountBy(metricName string, amount int) {
	globalMetrics.IncrementCountBy(metricName, amount)
}

// Measure Time since given for Metric by name
func MeasureSince(metricName string, since time.Time) {
	globalMetrics.MeasureSince(metricName, since)
}

// Gauge value for Metric by name
func SetGauge(metricName string, val float32) {
	globalMetrics.SetGauge(metricName, val)
}

// Set Prefix for all Metrics collected
func SetPrefix(prefix string) {
	globalMetrics.SetPrefix(prefix)
}

// WithTag returns a new MetricsRecorder that has the tags added to it.
func WithTag(key, value string) MetricsRecorder {
	return globalMetrics.WithTag(key, value)
}
