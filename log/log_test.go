package log_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/intercom/gocore/log"
)

var testTime = time.Now().UTC().Format(time.RFC3339)

func TestLogInfo(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo(testTime, "foo", "bar")
	checkLogFormatMatches(t, fmt.Sprintf("level=info time=%s foo=bar\n", testTime), buf)
}

func TestLogInfoWithOneValueBecomesMessage(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo(testTime, "foo")
	checkLogFormatMatches(t, fmt.Sprintf("level=info time=%s msg=foo\n", testTime), buf)
}

func TestLogInfoMessage(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfoMessage("my message")
	checkLogFormatMatches(t, "level=info msg=\"my message\"\n", buf)
}

func TestLogInfoMessageWithExtra(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfoMessage("my message", "foo", 7)
	checkLogFormatMatches(t, "level=info foo=7 msg=\"my message\"\n", buf)
}

func TestLogErrorMessageWithExtra(t *testing.T) {
	buf := logWithBuffer()
	log.LogErrorMessage("my message", "bar", 7.6)
	checkLogFormatMatches(t, "level=error bar=7.6 msg=\"my message\"\n", buf)
}

func TestLogInfoWithCompoundTypeArray(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo(testTime, "key", []string{"foo", "bar"})
	checkLogFormatMatches(t, fmt.Sprintf("level=info time=%s key=\"[foo bar]\"\n", testTime), buf)
}

func TestLogInfoWithCompoundTypeMap(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo(testTime, "key", map[string]interface{}{"another": 12})
	checkLogFormatMatches(t, fmt.Sprintf("level=info time=%s key=map[another:12]\n", testTime), buf)
}

func TestLogInfoWithCompoundTypeStruct(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo(testTime, "key", testTypeNotStringer{"foo"}, "bar", testTypeStringer{"bar"})
	checkLogFormatMatches(t, fmt.Sprintf("level=info time=%s key={Foo:foo} bar=bar\n", testTime), buf)
}

func TestLogWithStandardFields(t *testing.T) {
	buf := logWithBuffer()

	log.SetStandardFields("foo", "bar")
	log.LogErrorMessage("uh oh")
	checkLogFormatMatches(t, "level=error foo=bar msg=\"uh oh\"\n", buf)

	log.LogInfoMessage("something", "key", 4)
	checkLogFormatMatches(t, "level=info foo=bar key=4 msg=something\n", buf)
}

func TestJSONLog(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetupJSONLoggerTo(&buf)

	log.LogInfoMessage("something", "key", 4)
	checkLogFormatMatches(t, "{\"key\":4,\"level\":\"info\",\"msg\":\"something\"}\n", &buf)
}

func logWithBuffer() *bytes.Buffer {
	buf := bytes.Buffer{}
	log.SetupLogFmtLoggerTo(&buf)
	return &buf
}

func checkLogFormatMatches(t *testing.T, want string, buf *bytes.Buffer) {
	have := buf.String()
	if want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
	buf.Reset()
}

type testTypeNotStringer struct {
	Foo string
}

type testTypeStringer struct {
	Bar string
}

func (t testTypeStringer) String() string {
	return t.Bar
}
