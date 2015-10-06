package log_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/intercom/gocore/log"
)

func TestLogInfo(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo("foo", "bar")
	CheckLogFormatIgnoreTimestamp(t, `level=info ats=\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ foo=bar\n`, buf)
}

func TestLogInfoWithOneValueBecomesMessage(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo("foo")
	checkLogFormatMatches(t, "level=info msg=foo\n", buf)
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
	log.LogInfo("key", []string{"foo", "bar"})
	checkLogFormatMatches(t, "level=info key=\"[foo bar]\"\n", buf)
}

func TestLogInfoWithCompoundTypeMap(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo("key", map[string]interface{}{"another": 12})
	checkLogFormatMatches(t, "level=info key=map[another:12]\n", buf)
}

func TestLogInfoWithCompoundTypeStruct(t *testing.T) {
	buf := logWithBuffer()
	log.LogInfo("key", testTypeNotStringer{"foo"}, "bar", testTypeStringer{"bar"})
	checkLogFormatMatches(t, "level=info key={Foo:foo} bar=bar\n", buf)
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

func CheckLogFormatIgnoreTimestamp(t *testing.T, want string, buf *bytes.Buffer) {
	have := buf.String()
	var regexLine = regexp.MustCompile(want)
	fmt.Println(have)
	fmt.Println(regexLine)
	fmt.Println(regexLine.FindStringSubmatch(want))
	if regexLine.MatchString(have) {
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
