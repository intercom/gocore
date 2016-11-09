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
	sd, _ := metrics.NewStatsdRecorder("127.0.0.1:8888", "namespace")
	sd.IncrementCount("countMetric") // doesn't panic
	metrics.SetMetricsGlobal(sd)
	metrics.IncrementCount("countMetric")
}

func TestGetDatadogStatsdMetric(t *testing.T) {
	dd, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8888", "namespace", "hostname")
	dd.IncrementCount("countMetric")
	metrics.SetMetricsGlobal(dd)
	metrics.IncrementCount("countMetric")
}

func TestDatadogStatsdMetricTags(t *testing.T) {
	var dd metrics.MetricsRecorder
	dd, _ = metrics.NewDatadogStatsdRecorder("127.0.0.1:8888", "namespace", "hostname")
	dd = dd.WithTag("foo", "1")
	dd = dd.WithTag("bar", "2")

	tags := dd.(*metrics.DatadogStatsdRecorder).GetTags()

	if want, have := 2, len(tags); want != have {
		t.Errorf("want %#v tags, have %#v tags", want, have)
	}

	if want, have := "foo:1", tags[0]; want != have {
		t.Errorf("want first tag %#v, have %#v", want, have)
	}

	if want, have := "bar:2", tags[1]; want != have {
		t.Errorf("want second tag %#v, have %#v", want, have)
	}
}

func TestGetTeedMetricsRecorder(t *testing.T) {
	dd, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8125", "namespace", "hostname")
	teed := metrics.NewTeedMetricsRecorder(dd)
	tagged := teed.WithTag("key", "")
	tags := (tagged.(*metrics.TeedMetricsRecorder).GetMetrics()[0]).(*metrics.DatadogStatsdRecorder).GetTags()

	if want, have := "key", tags[0]; want != have {
		t.Errorf("want %#v tag, have %#v tag", want, have)
	}
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
func (tr *TestRecorder) MeasureSince(string, time.Time)                    {}
func (tr *TestRecorder) SetGauge(string, float32)                          {}
func (tr *TestRecorder) SetPrefix(string)                                  {}
func (tr *TestRecorder) WithTag(key, value string) metrics.MetricsRecorder { return tr }
