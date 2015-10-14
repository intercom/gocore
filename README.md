### GoCore

A library with a set of standardized functions for applications written in Go at Intercom.

To use:

Checkout into your gopath.

Vendor it into your project, making sure dependencies are satisfied (GoCore does not vendor it's own dependencies)

#### Logs

Structured logs in standard LogFmt or JSON format, with levels.

```go


// don't have to set a namespace, can just import and reference via "log" if you don't need the default logger too.
import corelog "github.com/intercom/gocore/log"

func main() {
  // setup logger before using
  corelog.SetupLogFmtLoggerTo(os.Stderr)

  // or JSON format:
  corelog.SetupJSONLoggerTo(os.Stderr)

  // log messages
  corelog.LogInfoMessage("reading items")

  // structured logging
  corelog.LogInfo("read_item_count", 4, "status", "finished")

  // both
  corelog.LogInfoMessage("reading items", "read_item_count", 4, "status", "finished")

  // setting standard fields
  corelog.SetStandardFields("instance_id", "67daf")
}
```

#### Metrics

Standardised Metrics options, for Global setup or individual.

```go

import "github.com/intercom/gocore/metrics"

func main() {
  // set the global metrics recorder to a new statsd recorder with namespace
  // without this, global metrics are no-opped.
  statsdRecorder, err := metrics.NewStatsdRecorder("127.0.0.1:8888", "namespace")
  if err == nil {
    metrics.SetMetricsGlobal(statsdRecorder)
  }

  // use global metrics
  metrics.IncrementCount("metricName")
  metrics.MeasureSince("metricName", startTime)

  // set prefix for all global metrics
  metrics.SetPrefix("prefixName")

  // create a new metric instance for separate collections:
  perAppMetrics, _ := metrics.NewStatsdRecorder("127.0.0.1:8889", "per-app-namespace")

  // use same recording methods
  perAppMetrics.IncrementCount("metricName")
}
```

To add a new recorder, implement the MetricsRecorder interface.

#### Monitoring

Standardised Monitoring options, for Global setup or individual. Currently, monitoring to Sentry is implemented.

```go
import "github.com/intercom/gocore/monitoring"

func main() {
  // set the global monitoring recorder to a new sentry monitor
  // without this, global monitoring is no-opped.
  sentryMonitor, err := monitoring.NewSentryMonitor("sentryDSN")
  if err == nil {
    monitoring.SetMonitoringGlobal(sentryMonitor)
  }

  // use global monitoring
  monitoring.CaptureException(errors.New("NewError"))

  // create a new monitoring instance:
  sentryMonitoring, _ := monitoring.NewSentryMonitor("sentryDSN")

  // use same capture method
  sentryMonitoring.CaptureException(errors.New("NewError"))
}
```

#### Dependencies

GoKit Log, Levels:

```
"github.com/go-kit/kit/log"
"github.com/go-kit/kit/log/levels"
```

Armon/Go-Metrics:

```
"github.com/armon/go-metrics"
```

Sentry/Raven:

```
"github.com/getsentry/raven-go"
```
