package monitoring_test

import (
	"errors"
	"testing"

	"github.com/intercom/gocore/monitoring"
)

func TestMonitoringGlobal(t *testing.T) {
	monitoring.CaptureException(errors.New("Foo")) // doesn't blow up when no global set

	tm := &TestMonitor{}
	monitoring.SetMonitoringGlobal(tm)
	monitoring.CaptureException(errors.New("Bar"))

	if tm.latestError.Error() != "Bar" {
		t.Errorf("Expected Bar, got %v", tm.latestError.Error())
	}
}

func TestMonitoringIndividual(t *testing.T) {
	tm := &TestMonitor{}
	tm.CaptureException(errors.New("Bar"))

	if tm.latestError.Error() != "Bar" {
		t.Errorf("Expected Bar, got %v", tm.latestError.Error())
	}
}

func TestMonitoringWithTags(t *testing.T) {
	tm := &TestMonitor{}
	tm.CaptureExceptionWithTags(errors.New("Bar"), "tag", "value")

	if tm.latestError.Error() != "Bar" {
		t.Errorf("Expected Bar, got %v", tm.latestError.Error())
	}
	if tm.latestTags[0] != "tag" {
		t.Errorf("Expected tag, got %v", tm.latestTags[0])
	}
	if tm.latestTags[1] != "value" {
		t.Errorf("Expected value, got %v", tm.latestTags[1])
	}
}

func TestMonitoringSetupSentryReturnsNilOnFail(t *testing.T) {
	sm := monitoring.NewSentryMonitor("abc123daef")
	if sm != nil {
		t.Errorf("Expected nil, got %v", sm)
	}
}

type TestMonitor struct {
	latestError error
	latestTags  []interface{}
}

func (tm *TestMonitor) CaptureException(err error) {
	tm.latestError = err
}

func (tm *TestMonitor) CaptureExceptionWithTags(err error, tags ...interface{}) {
	tm.latestError = err
	tm.latestTags = tags
}
