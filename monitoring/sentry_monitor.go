package monitoring

import (
	"fmt"

	"github.com/getsentry/raven-go"
)

type SentryMonitor struct {
	ravenClient *raven.Client
}

func NewSentryMonitor(ravenDSN string) (*SentryMonitor, error) {
	client, err := raven.NewClient(ravenDSN, nil)
	if err != nil {
		return nil, err
	}
	return &SentryMonitor{ravenClient: client}, nil
}

func (sm *SentryMonitor) CaptureException(err error) {
	sm.ravenClient.CaptureErrorAndWait(err, map[string]string{}, raven.NewStacktrace(2, 3, nil))
}

func (sm *SentryMonitor) CaptureExceptionWithTags(err error, tags ...interface{}) {
	sm.ravenClient.CaptureErrorAndWait(err, sm.convertTagListToMap(tags...), raven.NewStacktrace(2, 3, nil))
}

func (sm *SentryMonitor) convertTagListToMap(tags ...interface{}) map[string]string {
	tagMap := map[string]string{}

	if len(tags)%2 == 1 {
		tags = append(tags, nil)
	}

	for i := 0; i < len(tags); i += 2 {
		k, v := tags[i], tags[i+1]
		tagMap[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
	}
	return tagMap
}
