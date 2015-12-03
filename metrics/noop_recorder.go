package metrics

import "time"

// A nop-op MetricsRecorder implementation.
// For use when real metrics systems are unavailable.
type NoopRecorder struct{}

func (*NoopRecorder) IncrementCount(string)                       {}
func (*NoopRecorder) IncrementCountBy(string, int)                {}
func (*NoopRecorder) MeasureSince(string, time.Time)              {}
func (*NoopRecorder) SetGauge(string, float32)                    {}
func (*NoopRecorder) SetPrefix(string)                            {}
func (n *NoopRecorder) WithTag(key, value string) MetricsRecorder { return n }
