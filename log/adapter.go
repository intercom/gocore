package log

import (
	"io"
	"regexp"
	"time"
)

func NewWriterToLog(log Logger) io.Writer {
	return &StdlibAdapter{Logger: log}
}

type StdlibAdapter struct {
	Logger Logger
}

func (a *StdlibAdapter) Write(p []byte) (int, error) {
	result := a.subexps(p)
	keyvals := []interface{}{}
	var timestamp string
	if date, ok := result["date"]; ok && date != "" {
		timestamp = date
	}
	if time, ok := result["time"]; ok && time != "" {
		if timestamp != "" {
			timestamp += " "
		}
		timestamp += time
	}
	if timestamp != "" {
		t, _ := time.Parse("2006/01/02 15:04:05", timestamp)
		keyvals = append(keyvals, "timestamp", t)
	}
	if msg, ok := result["msg"]; ok {
		keyvals = append(keyvals, "msg", msg)
	}
	a.Logger.LogError(keyvals...)
	return len(p), nil
}

const (
	logRegexpDate = `(?P<date>[0-9]{4}/[0-9]{2}/[0-9]{2})?[ ]?`
	logRegexpTime = `(?P<time>[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?)?[ ]?`
	logRegexpFile = `(?P<file>.+?:[0-9]+)?`
	logRegexpMsg  = `(: )?(?P<msg>.*)`
)

var (
	logRegexp = regexp.MustCompile(logRegexpDate + logRegexpTime + logRegexpFile + logRegexpMsg)
)

func (a *StdlibAdapter) subexps(line []byte) map[string]string {
	m := logRegexp.FindSubmatch(line)
	if len(m) < len(logRegexp.SubexpNames()) {
		return map[string]string{}
	}
	result := map[string]string{}
	for i, name := range logRegexp.SubexpNames() {
		result[name] = string(m[i])
	}
	return result
}
