package ping

import (
	"context"

	appPing "github.com/siyoga/rollstory/internal/app/ping"
)

// PingHandler defines the interface for ping business logic
type PingHandler interface {
	Handle(ctx context.Context) (*appPing.Response, error)
}

// ErrorHandler defines the interface for error handling
type ErrorHandler interface {
	Handle(ctx context.Context, err error) error
}