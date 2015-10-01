package metrics

import (
	"errors"
	"time"

	"github.com/armon/go-metrics"
)

// A MetricsRecorder implementation which is just a wrapper around
// the go-metrics library, to record to Statsd.
type StatsdRecorder struct {
	*metrics.Metrics
	prefix string
}

// Takes a host:port string of the statsite endpoint to write to.
func NewStatsdRecorder(statsiteEndpoint, namespace string) (*StatsdRecorder, error) {
	if statsiteEndpoint == "" {
		return nil, errors.New("Uninitialized StatsdRecorder")
	}
	// statsdsink can be used to send to statssite over UDP
	// https://github.com/armon/go-metrics/blob/master/statsd.go#L21
	sink, _ := metrics.NewStatsdSink(statsiteEndpoint)
	config := metrics.DefaultConfig(namespace)
	config.EnableHostname = false
	m, _ := metrics.New(config, sink)
	return &StatsdRecorder{Metrics: m}, nil
}

func (m *StatsdRecorder) IncrementCount(metricName string) {
	m.Metrics.IncrCounter(m.prefixedMetricName(metricName), 1)
}

func (m *StatsdRecorder) IncrementCountBy(metricName string, amount int) {
	m.Metrics.IncrCounter(m.prefixedMetricName(metricName), float32(amount))
}

func (m *StatsdRecorder) MeasureSince(metricName string, since time.Time) {
	m.Metrics.MeasureSince(m.prefixedMetricName(metricName), since)
}

func (m *StatsdRecorder) SetGauge(metricName string, val float32) {
	m.Metrics.SetGauge(m.prefixedMetricName(metricName), val)
}

func (m *StatsdRecorder) SetPrefix(prefix string) {
	m.prefix = prefix
}

func (m *StatsdRecorder) prefixedMetricName(metricName string) []string {
	if m.prefix == "" {
		return []string{metricName}
	}
	return []string{m.prefix, metricName}
}
