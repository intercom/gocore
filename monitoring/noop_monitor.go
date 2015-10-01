package monitoring

type NoopMonitor struct{}

func (*NoopMonitor) CaptureException(error)                                  {}
func (*NoopMonitor) CaptureExceptionWithTags(err error, tags ...interface{}) {}
