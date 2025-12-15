package listener

import (
	"os"
	"syscall"
	"time"
)

const (
	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 30 * time.Second
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}
