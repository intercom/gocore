package log

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestLogInfo(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("foo", "bar")
	checkLogFormatMatches(t, "level=info foo=bar\n", buf)
}

func TestLogInfoWithOneValueBecomesMessage(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("foo")
	checkLogFormatMatches(t, "level=info msg=foo\n", buf)
}

func TestLogInfoMessage(t *testing.T) {
	buf := logWithBuffer()
	LogInfoMessage("my message")
	checkLogFormatMatches(t, "level=info msg=\"my message\"\n", buf)
}

func TestLogInfoMessageWithExtra(t *testing.T) {
	buf := logWithBuffer()
	LogInfoMessage("my message", "foo", 7)
	checkLogFormatMatches(t, "level=info foo=7 msg=\"my message\"\n", buf)
}

func TestLogErrorMessageWithExtra(t *testing.T) {
	buf := logWithBuffer()
	LogErrorMessage("my message", "bar", 7.6)
	checkLogFormatMatches(t, "level=error bar=7.6 msg=\"my message\"\n", buf)
}

func TestLogInfoWithCompoundTypeArray(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("key", []string{"foo", "bar"})
	checkLogFormatMatches(t, "level=info key=\"[foo bar]\"\n", buf)
}

func TestLogInfoWithCompoundTypeMap(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("key", map[string]interface{}{"another": 12})
	checkLogFormatMatches(t, "level=info key=map[another:12]\n", buf)
}

func TestLogInfoWithCompoundTypeStruct(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("key", testTypeNotStringer{"foo"}, "bar", testTypeStringer{"bar"})
	checkLogFormatMatches(t, "level=info key={Foo:foo} bar=bar\n", buf)
}

func TestLogWithStandardFields(t *testing.T) {
	buf := logWithBuffer()

	SetStandardFields("foo", "bar")
	LogErrorMessage("uh oh")
	checkLogFormatMatches(t, "level=error foo=bar msg=\"uh oh\"\n", buf)

	LogInfoMessage("something", "key", 4)
	checkLogFormatMatches(t, "level=info foo=bar key=4 msg=something\n", buf)
}

func TestJSONLog(t *testing.T) {
	buf := bytes.Buffer{}

	SetupGlobalLoggerTo(jsonLoggerTo(&buf))
	LogInfoMessage("something", "key", 4)
	checkLogFormatMatches(t, "{\"key\":4,\"level\":\"info\",\"msg\":\"something\"}\n", &buf)
}

func TestJSONLogWithTimestamp(t *testing.T) {
	buf := bytes.Buffer{}

	SetupJSONLoggerTo(&buf)
	LogInfoMessage("something", "key", 4)
	tsl := timestampedLog{}
	json.Unmarshal(buf.Bytes(), &tsl)
	have := tsl.Timestamp.Unix()
	want := time.Now().Unix()
	difference := have - want
	// fail unless we're within a second of...
	if difference < -1 || difference > 1 {
		t.Errorf("timestamp didnt match, have %v, want %v", have, want)
	}
}

func logWithBuffer() *bytes.Buffer {
	buf := bytes.Buffer{}
	SetupGlobalLoggerTo(logfmtLoggerTo(&buf))
	return &buf
}

func checkLogFormatMatches(t *testing.T, want string, buf *bytes.Buffer) {
	have := buf.String()
	if want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
	buf.Reset()
}

type timestampedLog struct {
	Timestamp time.Time `json:"aaats"`
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
