package coreapi_test

import (
	"testing"

	"github.com/intercom/gocore/coreapi"
)

func TestBasicAuth(t *testing.T) {
	auth := &coreapi.BasicAuth{User: "user", Pass: "pass"}
	if auth.CheckBasicAuth("user", "bar") {
		t.Errorf("Expected false, got true")
	}
	if auth.CheckBasicAuth("foo", "pass") {
		t.Errorf("Expected false, got true")
	}
	if !auth.CheckBasicAuth("user", "pass") {
		t.Errorf("Expected true, got false")
	}
}
