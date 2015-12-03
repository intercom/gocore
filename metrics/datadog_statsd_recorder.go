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
	sink, _ := datadog.NewDogStatsdSink(statsiteEndpoint, hostname)
	config := metrics.DefaultConfig(namespace)
	config.EnableHostname = false
	m, _ := metrics.New(config, sink)
	return &DatadogStatsdRecorder{StatsdRecorder: &StatsdRecorder{m, ""}, sink: sink, tags: []string{}}, nil
}

func (t *DatadogStatsdRecorder) IncrementCount(metricName string) {
	// TODO(JO): use with tag calls
	t.sink.IncrCounter(t.prefixedMetricName(metricName), 1)
}

func (t *DatadogStatsdRecorder) IncrementCountBy(metricName string, amount int) {
	// TODO(JO): use with tag calls
	t.sink.IncrCounter(t.prefixedMetricName(metricName), float32(amount))
}

func (t *DatadogStatsdRecorder) MeasureSince(metricName string, since time.Time) {
	now := time.Now()
	elapsed := now.Sub(since)
	msec := float32(elapsed.Nanoseconds()) / float32(time.Millisecond)
	// TODO(JO): use with tag tag calls
	t.sink.AddSample(t.prefixedMetricName(metricName), msec)
}

func (t *DatadogStatsdRecorder) SetGauge(metricName string, val float32) {
	// TODO(JO): use with tag calls
	t.sink.SetGauge(t.prefixedMetricName(metricName), val)
}

// WithTag returns a new DatadogStatsdRecorder that has the tags added to it.
func (t *DatadogStatsdRecorder) WithTag(key, value string) MetricsRecorder {
	t.tags = append(t.tags, flatKey(key, value))
	return &DatadogStatsdRecorder{t.StatsdRecorder, t.sink, t.tags}
}

func (t *DatadogStatsdRecorder) GetTags() []string {
	return t.tags
}

func test() {
	t, _ := NewDatadogStatsdRecorder("endpoint", "namespace", "hostname")
	t.WithTag("foo", "bar").IncrementCount("damnit")
	tagged := t.WithTag("doo", "you")
	tagged.MeasureSince("foo", time.Now())

	b := NewTeedMetricsRecorder(t)
	b.IncrementCount("metricName string")
	b.WithTag("key", "value").MeasureSince("foo", time.Now())
}
