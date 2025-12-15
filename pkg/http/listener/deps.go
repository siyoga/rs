package listener

import "context"

type Logger interface {
	Error(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	WithField(ctx context.Context, k string, v interface{}) context.Context
}
