package metrics_test

import (
	"testing"
	"time"

	"github.com/intercom/gocore/metrics"
)

func TestDatadogStatsdTagKeyValue(t *testing.T) {
	recorder, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8125", "namespace", "hostname")
	tagged := recorder.WithTag("tagkey", "tagvalue")
	tagged.MeasureSince("foo", time.Now())
	tags := tagged.(*metrics.DatadogStatsdRecorder).GetTags()

	if want, have := "tagkey:tagvalue", tags[0]; want != have {
		t.Errorf("want %#v tag, have %#v tag", want, have)
	}
}

func TestDatadogStatsdMultiTags(t *testing.T) {
	recorder, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8125", "namespace", "hostname")
	tagged := recorder.WithTag("tagkey", "tagvalue")
	tagged = tagged.WithTag("anotherkey", "anothervalue")
	tags := tagged.(*metrics.DatadogStatsdRecorder).GetTags()

	if want, have := "anotherkey:anothervalue", tags[1]; want != have {
		t.Errorf("want %#v tag, have %#v tag", want, have)
	}

}
