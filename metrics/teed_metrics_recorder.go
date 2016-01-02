package metrics

import "time"

type TeedMetricsRecorder struct {
	metrics []MetricsRecorder
	prefix  string
}

func NewTeedMetricsRecorder(metrics ...MetricsRecorder) *TeedMetricsRecorder {
	return &TeedMetricsRecorder{metrics: metrics}
}

func (t *TeedMetricsRecorder) IncrementCount(metricName string) {
	for _, m := range t.metrics {
		m.IncrementCount(metricName)
	}
}

func (t *TeedMetricsRecorder) IncrementCountBy(metricName string, amount int) {
	for _, m := range t.metrics {
		m.IncrementCountBy(metricName, amount)
	}
}

func (t *TeedMetricsRecorder) MeasureSince(metricName string, since time.Time) {
	for _, m := range t.metrics {
		m.MeasureSince(metricName, since)
	}
}

func (t *TeedMetricsRecorder) SetGauge(metricName string, val float32) {
	for _, m := range t.metrics {
		m.SetGauge(metricName, val)
	}
}

func (t *TeedMetricsRecorder) SetPrefix(prefix string) {
	for _, m := range t.metrics {
		m.SetPrefix(prefix)
	}
}

func (t *TeedMetricsRecorder) WithTag(key, value string) MetricsRecorder {
	newRecorder := TeedMetricsRecorder{prefix: t.prefix, metrics: []MetricsRecorder{}}
	for _, m := range t.metrics {
		newRecorder.metrics = append(newRecorder.metrics, m.WithTag(key, value))
	}
	return &newRecorder
}

func (t *TeedMetricsRecorder) GetMetrics() []MetricsRecorder {
	return t.metrics
}
