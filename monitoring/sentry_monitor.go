package monitoring

import "github.com/getsentry/raven-go"

type SentryMonitor struct {
	ravenClient *raven.Client
}

func NewSentryMonitor(ravenDSN string) *SentryMonitor {
	client, err := raven.NewClient(ravenDSN, nil)
	if err != nil {
		return nil
	}
	return &SentryMonitor{ravenClient: client}
}

func (rm *SentryMonitor) CaptureException(err error) {
	rm.ravenClient.CaptureErrorAndWait(err, map[string]string{})
}
