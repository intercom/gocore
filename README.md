### GoCore

A library with a set of standardised functions for applications written in Go at Intercom.

To use:

Checkout into your gopath.

Vendor it into your project, making sure dependencies are satisified (GoCore does not vendor it's own dependencies)

#### Logs

Structured logs in standard LogFmt format, with levels.

```go

import corelog "github.com/intercom/gocore/log"

func main() {
	// setup logger before using
	corelog.SetupLoggerToStderr()

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

Coming Soon

#### Sentry

Coming Soon

#### Dependencies

GoKit Log, Levels:

```
"github.com/go-kit/kit/log"
"github.com/go-kit/kit/log/levels"
```
