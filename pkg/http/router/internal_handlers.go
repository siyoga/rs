package router

import (
	"context"
	"fmt"
	"github.com/siyoga/rollstory/pkg/http/router/response"
	"net/http"
	"runtime/debug"
)

type DefaultErr struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type BadRequestErrHandler func(err error, resp response.Builder)

type InternalErrHandler func(err interface{}, resp response.Builder)

type PanicHandler func(ctx context.Context, e interface{})

var (
	DefaultInternalErrHandler = func(err interface{}, resp response.Builder) {
		resp.
			SetHeader("Content-Type", "application/json").
			SetJsonBody(&DefaultErr{
				Message: "internal error",
				Code:    http.StatusInternalServerError,
			}).
			SetStatusCode(http.StatusInternalServerError)
	}

	DefaultBadRequestErrHandler = func(err error, resp response.Builder) {
		resp.
			SetHeader("Content-Type", "application/json").
			SetJsonBody(&DefaultErr{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			}).
			SetStatusCode(http.StatusBadRequest)
	}

	DefaultPanicHandler = func(logger Logger) func(ctx context.Context, e interface{}) {
		return func(ctx context.Context, e interface{}) {
			logger.Error(ctx, fmt.Sprintf("Recover after panic: %+v\n%s", e, debug.Stack()))
		}
	}
)
