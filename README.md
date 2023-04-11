# go-healthcheck

### Usage

```shell
go get github.com/goforbroke1006/go-healthcheck
```

```go
package main

import (
	"context"
	
	"github.com/goforbroke1006/go-healthcheck"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	healthcheck.Panel().Start(ctx, healthcheck.DefaultAddr)
	
	// connect external resources

	healthcheck.Panel().SetHealthy()
	
	// load caches

	healthcheck.Panel().SetReady()
}

```