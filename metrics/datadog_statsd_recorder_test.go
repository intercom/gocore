package metrics_test

import (
	"net"
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

func TestDatadogStatsdTagsMakeNewInstance(t *testing.T) {
	recorder, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8125", "namespace", "hostname")
	tagged := recorder.WithTag("tagkey", "tagvalue")

	if len(recorder.GetTags()) != 0 {
		t.Error("modified old recorder")
	}

	if len(tagged.(*metrics.DatadogStatsdRecorder).GetTags()) != 1 {
		t.Error("did not modify new recorder")
	}
}

func TestDatadogStatsdMultiTags(t *testing.T) {
	recorder, _ := metrics.NewDatadogStatsdRecorder("127.0.0.1:8125", "namespace", "hostname")
	tagged := recorder.WithTag("tagkey", "tagvalue")
	tagged = tagged.WithTag("anotherkey", "anothervalue")
	tags := tagged.(*metrics.DatadogStatsdRecorder).GetTags()

	if want, have := "tagkey:tagvalue", tags[0]; want != have {
		t.Errorf("want %#v tag, have %#v tag", want, have)
	}

	if want, have := "anotherkey:anothervalue", tags[1]; want != have {
		t.Errorf("want %#v tag, have %#v tag", want, have)
	}
}

const (
	DogStatsdAddr = "127.0.0.1:7254"
)

func mockNewDogStatsdSink(addr string, tags []string, tagWithHostname bool) *metrics.DatadogStatsdRecorder {
	dog, _ := metrics.NewDatadogStatsdRecorder(addr, "namespace", "hostname")
	return dog
}

func setupTestServerAndBuffer(t *testing.T) (*net.UDPConn, []byte) {
	udpAddr, err := net.ResolveUDPAddr("udp", DogStatsdAddr)
	if err != nil {
		t.Fatal(err)
	}
	server, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	return server, make([]byte, 1024)
}

func TestDogStatsdSink(t *testing.T) {
	server, buf := setupTestServerAndBuffer(t)
	defer server.Close()

	dog := mockNewDogStatsdSink(DogStatsdAddr, []string{}, false)
	dog.IncrementCountBy("counter", 4)
	assertServerMatchesExpected(t, server, buf, "namespace.counter:4|c")
}

func TestDogStatsdSinkWithTag(t *testing.T) {
	server, buf := setupTestServerAndBuffer(t)
	defer server.Close()

	dog := mockNewDogStatsdSink(DogStatsdAddr, []string{}, false)
	tagged := dog.WithTag("tagkey", "tagvalue")
	tagged.IncrementCountBy("counter", 4)
	assertServerMatchesExpected(t, server, buf, "namespace.counter:4|c|#tagkey:tagvalue")
}

func assertServerMatchesExpected(t *testing.T, server *net.UDPConn, buf []byte, expected string) {
	n, _ := server.Read(buf)
	msg := buf[:n]
	if string(msg) != expected {
		t.Fatalf("Line %s does not match expected: %s", string(msg), expected)
	}
}
