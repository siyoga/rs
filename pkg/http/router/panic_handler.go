package router

import (
	"github.com/siyoga/rollstory/pkg/http/router/response"
	"net/http"
)

func panicMiddleware(errHandler InternalErrHandler, panicHandler PanicHandler) func(string, http.Handler) http.Handler {
	return func(pattern string, next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			handlerResp := response.NewResponse()
			defer func() {
				if err := recover(); err != nil {
					ctx := req.Context()
					panicHandler(ctx, err)
					errHandler(ctx, handlerResp)
					_ = handlerResp.Send(resp)
				}
			}()

			next.ServeHTTP(resp, req)
		})
	}
}
