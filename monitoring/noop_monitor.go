package monitoring

type NoopMonitor struct{}

func (*NoopMonitor) CaptureException(error) {}
