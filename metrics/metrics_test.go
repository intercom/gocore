package metrics_test

import (
	"testing"
	"time"

	"github.com/intercom/gocore/metrics"
)

func TestSetMetricsGlobal(t *testing.T) {
	metrics.IncrementCount("countMetric") // doesn't blow up when no global set

	tr := TestRecorder{metrics: map[string]interface{}{}}
	metrics.SetMetricsGlobal(&tr)
	metrics.IncrementCount("countMetric")
	if want, have := 1, tr.metrics["countMetric"]; want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	metrics.IncrementCountBy("countMetric", 3)
	if want, have := 4, tr.metrics["countMetric"]; want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func TestGetStatsdMetric(t *testing.T) {
	sd := metrics.NewStatsdRecorder("127.0.0.1:8888", "namespace")
	sd.IncrementCount("countMetric") // doesn't panic
	metrics.SetMetricsGlobal(sd)
	metrics.IncrementCount("countMetric")
}

type TestRecorder struct {
	metrics map[string]interface{}
}

func (tr *TestRecorder) IncrementCount(metricName string) {
	tr.IncrementCountBy(metricName, 1)
}

func (tr *TestRecorder) IncrementCountBy(metricName string, val int) {
	if _, present := tr.metrics[metricName]; present {
		tr.metrics[metricName] = tr.metrics[metricName].(int) + val
	} else {
		tr.metrics[metricName] = val
	}
}

// noops
func (tr *TestRecorder) MeasureSince(string, time.Time) {}
func (tr *TestRecorder) SetGauge(string, float32)       {}
func (tr *TestRecorder) SetPrefix(string)               {}
