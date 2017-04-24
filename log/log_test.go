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
	checkLogFormatMatches(t, "foo=bar level=info\n", buf)
}

func TestLogInfoWithOneValueBecomesMessage(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("foo")
	checkLogFormatMatches(t, "msg=foo level=info\n", buf)
}

func TestLogInfoMessage(t *testing.T) {
	buf := logWithBuffer()
	LogInfoMessage("my message")
	checkLogFormatMatches(t, "msg=\"my message\" level=info\n", buf)
}

func TestLogInfoMessageWithExtra(t *testing.T) {
	buf := logWithBuffer()
	LogInfoMessage("my message", "foo", 7)
	checkLogFormatMatches(t, "foo=7 msg=\"my message\" level=info\n", buf)
}

func TestLogErrorMessageWithExtra(t *testing.T) {
	buf := logWithBuffer()
	LogErrorMessage("my message", "bar", 7.6)
	checkLogFormatMatches(t, "bar=7.6 msg=\"my message\" level=error\n", buf)
}

func TestLogInfoWithCompoundTypeArray(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("key", []string{"foo", "bar"})
	checkLogFormatMatches(t, "key=\"[foo bar]\" level=info\n", buf)
}

func TestLogInfoWithCompoundTypeMap(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("key", map[string]interface{}{"another": 12})
	checkLogFormatMatches(t, "key=map[another:12] level=info\n", buf)
}

func TestLogInfoWithCompoundTypeStruct(t *testing.T) {
	buf := logWithBuffer()
	LogInfo("key", testTypeNotStringer{"foo"}, "bar", testTypeStringer{"bar"})
	checkLogFormatMatches(t, "key={Foo:foo} bar=bar level=info\n", buf)
}

func TestLogWithStandardFields(t *testing.T) {
	buf := logWithBuffer()

	SetStandardFields("foo", "bar")
	LogErrorMessage("uh oh")
	checkLogFormatMatches(t, "foo=bar msg=\"uh oh\" level=error\n", buf)

	SetStandardFields("zap", "zam")
	LogInfoMessage("something", "key", 4)
	checkLogFormatMatches(t, "foo=bar zap=zam key=4 msg=something level=info\n", buf)
}

func TestLogWithStandardFieldsAndTimestamp(t *testing.T) {
	buf := bytes.Buffer{}
	log := JSONLoggerTo(&buf)

	log.SetStandardFields("foo", "bar")
	log.LogErrorMessage("uh oh")
	assertTimestampWithin(t, &buf)
}

func TestLogWithStandardFieldsMakesNewLogger(t *testing.T) {
	buf := bytes.Buffer{}
	logger := LogfmtLoggerTo(&buf)
	logger.hideTimestamp = true
	l2 := logger.SetStandardFields("foo", "bar")

	logger.LogErrorMessage("uh oh")
	checkLogFormatMatches(t, "msg=\"uh oh\" level=error\n", &buf)

	l2.LogErrorMessage("uh oh")
	checkLogFormatMatches(t, "foo=bar msg=\"uh oh\" level=error\n", &buf)
}

func TestJSONLog(t *testing.T) {
	buf := bytes.Buffer{}
	logger := JSONLoggerTo(&buf)
	logger.hideTimestamp = true
	logger.LogInfoMessage("something", "key", 4)
	checkLogFormatMatches(t, "{\"key\":4,\"level\":\"info\",\"msg\":\"something\"}\n", &buf)
}

func TestJSONLogWithTimestamp(t *testing.T) {
	buf := bytes.Buffer{}
	log := JSONLoggerTo(&buf)
	log.LogInfoMessage("something", "key", 4)
	assertTimestampWithin(t, &buf)
}

func assertTimestampWithin(t *testing.T, buf *bytes.Buffer) {
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
	log := LogfmtLoggerTo(&buf)
	log.hideTimestamp = true
	GlobalLogger = log
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
	Timestamp time.Time `json:"timestamp"`
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
