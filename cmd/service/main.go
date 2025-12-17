package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/siyoga/rollstory/internal/api/ping"
	"github.com/siyoga/rollstory/internal/generated/api"
	bootstrap "github.com/siyoga/rollstory/internal/init"
	"github.com/siyoga/rollstory/pkg/http/listener"
	"github.com/siyoga/rollstory/pkg/http/router"
	"github.com/siyoga/rollstory/pkg/logger"
)

const (
	success = 0
	fail    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	log, err := logger.New()
	if err != nil {
		fmt.Println("while init logger: ", err.Error())
		return fail
	}
	defer log.Flush(time.Second)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	container := bootstrap.NewContainer()
	rt := router.NewRouter(
		log,
		router.DefaultBadRequestErrHandler,
		router.DefaultInternalErrHandler,
		router.DefaultPanicHandler(log),
	)

	// register handlers
	// Each RPC handler implements api.StrictServerInterface partially
	// For now, we have only ping handler, in future we'll combine multiple handlers
	if err := container.Invoke(func(pingHandler *ping.Handler) {
		// Use ping handler directly as it implements StrictServerInterface
		strictHandler := api.NewStrictHandler(pingHandler, nil)
		api.HandlerWithOptions(strictHandler, api.StdHTTPServerOptions{
			BaseRouter: rt,
		})
	}); err != nil {
		log.Error(ctx, "failed to invoke handlers", err)
		return fail
	}

	listenerReadTimeout, err := strconv.Atoi(os.Getenv("LISTENER_READ_TIMEOUT"))
	if err != nil {
		log.Warning(ctx, "can not parse listenerReadTimeout, using default value = 5")
		listenerReadTimeout = 5 // Default value
	}

	listenerWriteTimeout, err := strconv.Atoi(os.Getenv("LISTENER_WRITE_TIMEOUT"))
	if err != nil {
		log.Warning(ctx, "can not parse listenerWriteTimeout, using default value = 5")
		listenerWriteTimeout = 5 // Default value
	}

	listenerIdleTimeout, err := strconv.Atoi(os.Getenv("LISTENER_IDLE_TIMEOUT"))
	if err != nil {
		log.Warning(ctx, "can not parse listenerIdleTimeout, using default value = 5")
		listenerIdleTimeout = 5 // Default value
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Warning(ctx, "can not parse port, using default value = 8080")
		port = 8080
	}

	ln := listener.New(
		listener.With(log),
		listener.WithIdleTimeout(time.Duration(listenerIdleTimeout)*time.Second),
		listener.WithReadTimeout(time.Duration(listenerReadTimeout)*time.Second),
		listener.WithWriteTimeout(time.Duration(listenerWriteTimeout)*time.Second),
	)

	log.Info(ctx, fmt.Sprintf("listening on port %d", port))
	if err := ln.Listen(ctx, port, rt); err != nil {
		log.Error(ctx, "server error", err)
		return fail
	}

	return success
}
