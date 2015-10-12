package metrics

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
)

// A MetricsRecorder implementation which is just a wrapper around
// the go-metrics library, to record to Statsd.
type MultiStatsdRecorder struct {
	Metrics []*metrics.Metrics
	prefix  string
}

// Takes a slice of host:port strings of the statsite endpoints to write to.
func NewMultiStatsdRecorder(statsiteEndpoints []string, namespace string) (*MultiStatsdRecorder, error) {
	recorder := &MultiStatsdRecorder{}
	for i, endpoint := range statsiteEndpoints {
		if endpoint == "" {
			return nil, fmt.Errorf("Uninitialized MultiStatsdRecorder %d", i)
		}
		// statsdsink can be used to send to statssite over UDP
		// https://github.com/armon/go-metrics/blob/master/statsd.go#L21
		sink, _ := metrics.NewStatsdSink(endpoint)
		config := metrics.DefaultConfig(namespace)
		config.EnableHostname = false
		m, _ := metrics.New(config, sink)
		recorder.Metrics = append(recorder.Metrics, m)
	}
	return recorder, nil
}

func (m *MultiStatsdRecorder) IncrementCount(metricName string) {
	for _, recorder := range m.Metrics {
		recorder.IncrCounter(m.prefixedMetricName(metricName), 1)
	}
}

func (m *MultiStatsdRecorder) IncrementCountBy(metricName string, amount int) {
	for _, recorder := range m.Metrics {
		recorder.IncrCounter(m.prefixedMetricName(metricName), float32(amount))
	}
}

func (m *MultiStatsdRecorder) MeasureSince(metricName string, since time.Time) {
	for _, recorder := range m.Metrics {
		recorder.MeasureSince(m.prefixedMetricName(metricName), since)
	}
}

func (m *MultiStatsdRecorder) SetGauge(metricName string, val float32) {
	for _, recorder := range m.Metrics {
		recorder.SetGauge(m.prefixedMetricName(metricName), val)
	}
}

func (m *MultiStatsdRecorder) SetPrefix(prefix string) {
	m.prefix = prefix
}

func (m *MultiStatsdRecorder) prefixedMetricName(metricName string) []string {
	if m.prefix == "" {
		return []string{metricName}
	}
	return []string{m.prefix, metricName}
}
