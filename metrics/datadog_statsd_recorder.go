package metrics

import (
	"errors"
	"fmt"
	"time"
)
import (
	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/datadog"
)

// DatadogStatsdRecorder wraps a StatsdRecorder and allows tagging of metrics
type DatadogStatsdRecorder struct {
	*StatsdRecorder
	sink *datadog.DogStatsdSink // we need direct access to the strongly typed underlying sink
	tags []string
}

func flatKey(key, value string) string {
	if value != "" {
		return fmt.Sprintf("%s:%s", key, value)
	}
	return key
}

func NewDatadogStatsdRecorder(statsiteEndpoint, namespace, hostname string) (*DatadogStatsdRecorder, error) {
	if statsiteEndpoint == "" {
		return nil, errors.New("Uninitialized DatadogStatsdRecorder")
	}
	sink, err := datadog.NewDogStatsdSink(statsiteEndpoint, hostname)
	if err != nil {
		return nil, err
	}
	config := metrics.DefaultConfig(namespace)
	config.EnableHostname = false
	m, _ := metrics.New(config, sink)
	return &DatadogStatsdRecorder{StatsdRecorder: &StatsdRecorder{m, ""}, sink: sink, tags: []string{}}, nil
}

func (t *DatadogStatsdRecorder) IncrementCount(metricName string) {
	t.sink.IncrCounterWithTags(t.prefixedMetricName(metricName), 1, t.tags)
}

func (t *DatadogStatsdRecorder) IncrementCountBy(metricName string, amount int) {
	t.sink.IncrCounterWithTags(t.prefixedMetricName(metricName), float32(amount), t.tags)
}

func (t *DatadogStatsdRecorder) MeasureSince(metricName string, since time.Time) {
	now := time.Now()
	elapsed := now.Sub(since)
	msec := float32(elapsed.Nanoseconds()) / float32(time.Millisecond)
	t.sink.AddSampleWithTags(t.prefixedMetricName(metricName), msec, t.tags)
}

func (t *DatadogStatsdRecorder) SetGauge(metricName string, val float32) {
	t.sink.SetGaugeWithTags(t.prefixedMetricName(metricName), val, t.tags)
}

// WithTag returns a new DatadogStatsdRecorder that has the tags added to it.
func (t *DatadogStatsdRecorder) WithTag(key, value string) MetricsRecorder {
	t.tags = append(t.tags, flatKey(key, value))
	return &DatadogStatsdRecorder{StatsdRecorder: t.StatsdRecorder, sink: t.sink, tags: t.tags}
}

func (t *DatadogStatsdRecorder) GetTags() []string {
	return t.tags
}
