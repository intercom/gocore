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

func (dd *DatadogStatsdRecorder) IncrementCount(metricName string) {
	dd.IncrementCountBy(metricName, 1)
}

func (dd *DatadogStatsdRecorder) IncrementCountBy(metricName string, amount int) {
	dd.sink.IncrCounterWithTags(
		dd.withPrefixAndServiceName(metricName, "counter"),
		float32(amount),
		dd.tags,
	)
}

func (dd *DatadogStatsdRecorder) MeasureSince(metricName string, since time.Time) {
	now := time.Now()
	elapsed := now.Sub(since)
	msec := float32(elapsed.Nanoseconds()) / float32(time.Millisecond)

	dd.sink.AddSampleWithTags(
		dd.withPrefixAndServiceName(metricName, "timer"),
		msec,
		dd.tags,
	)
}

func (dd *DatadogStatsdRecorder) SetGauge(metricName string, val float32) {
	dd.sink.SetGaugeWithTags(
		dd.withPrefixAndServiceName(metricName, "gauge"),
		val,
		dd.tags,
	)
}

// WithTag returns a new DatadogStatsdRecorder that has the tags added to it.
func (dd *DatadogStatsdRecorder) WithTag(key, value string) MetricsRecorder {
	dd.tags = append(dd.tags, flatKey(key, value))
	return &DatadogStatsdRecorder{StatsdRecorder: dd.StatsdRecorder, sink: dd.sink, tags: dd.tags}
}

func (dd *DatadogStatsdRecorder) GetTags() []string {
	return dd.tags
}

// adds prefix, service name prefix, and type prefix
func (dd *DatadogStatsdRecorder) withPrefixAndServiceName(metricName, typeStr string) []string {
	key := dd.prefixedMetricName(metricName)
	if dd.StatsdRecorder.Metrics.EnableTypePrefix {
		key = insert(0, typeStr, key)
	}
	if dd.StatsdRecorder.Metrics.ServiceName != "" {
		key = insert(0, dd.StatsdRecorder.Metrics.ServiceName, key)
	}
	return key
}

// Inserts a string value at an index into the slice
func insert(i int, v string, s []string) []string {
	s = append(s, "")
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}
