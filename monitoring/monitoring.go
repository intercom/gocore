package monitoring

var (
	globalMonitor Monitor
)

type Monitor interface {
	CaptureException(err error)
}

// Package-level default initialization of the Monitoring global.
// Initializes it to a no-op implementation;
// later calls can replace it by calling SetMonitoringGlobal.
func init() {
	globalMonitor = &NoopMonitor{}
}

// setup the monitoring global
func SetMonitoringGlobal(monitor Monitor) {
	if monitor != nil {
		globalMonitor = monitor
	}
}

// Capture an exception and send to sentry
func CaptureException(err error) {
	globalMonitor.CaptureException(err)
}
