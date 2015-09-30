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

func TestMonitoringSetupSentryReturnsNilOnFail(t *testing.T) {
	sm := monitoring.NewSentryMonitor("abc123daef")
	if sm != nil {
		t.Errorf("Expected nil, got %v", sm)
	}
}

type TestMonitor struct {
	latestError error
}

func (tm *TestMonitor) CaptureException(err error) {
	tm.latestError = err
}
