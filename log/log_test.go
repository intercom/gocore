package log_test

import (
	"bytes"
	"testing"

	"github.com/intercom/gocore/log"
)

func TestLogInfoMessage(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetupLoggerTo(&buf)
	log.LogInfoMessage("my message")
	if want, have := "level=info msg=\"my message\"\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func TestLogInfoMessageWithExtra(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetupLoggerTo(&buf)
	log.LogInfoMessage("my message", "foo", 7)
	if want, have := "level=info foo=7 msg=\"my message\"\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func TestLogErrorMessageWithExtra(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetupLoggerTo(&buf)
	log.LogErrorMessage("my message", "bar", 7.6)
	if want, have := "level=error bar=7.6 msg=\"my message\"\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func TestLogWithStandardFields(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetupLoggerTo(&buf)

	log.SetStandardFields("foo", "bar")
	log.LogErrorMessage("uh oh")
	if want, have := "level=error foo=bar msg=\"uh oh\"\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	log.LogInfoMessage("something", "key", 4)
	if want, have := "level=info foo=bar key=4 msg=something\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}
