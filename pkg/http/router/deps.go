package router

import (
	"context"
)

type Logger interface {
	Error(ctx context.Context, args ...interface{})
}
