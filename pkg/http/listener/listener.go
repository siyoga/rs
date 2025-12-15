package listener

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type HTTPListener struct {
	options
}

func New(opts ...Option) *HTTPListener {
	o := options{
		shutdownTimeout: defaultShutdownTimeout,
		recoverPanics:   true,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &HTTPListener{
		options: o,
	}
}

func (l *HTTPListener) Listen(ctx context.Context, port int, handler http.Handler) error {
	for i := len(l.mw) - 1; i >= 0; i-- {
		handler = l.mw[i](handler)
	}

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		ReadTimeout:    defaultReadTimeout,
		IdleTimeout:    defaultIdleTimeout,
		WriteTimeout:   defaultWriteTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ErrorLog:       errorLog(l.logger),
	}

	for _, opt := range l.server {
		opt(server)
	}

	errChan := make(chan error)
	signalsChan := make(chan os.Signal, 2)
	signal.Notify(signalsChan, shutdownSignals...)

	go l.listen(errChan, server, l.netListener)

	var err error
	select {
	case err = <-errChan:
		_ = l.shutdown(server, l.shutdownTimeout)

	case <-signalsChan:
		_ = l.shutdown(server, l.shutdownTimeout)
	case <-ctx.Done():
		err = ctx.Err()
		if e := l.shutdown(server, l.shutdownTimeout); e != nil {
			err = e
		}
	}

	return err
}

func (l *HTTPListener) listen(ch chan error, server *http.Server, netListener net.Listener) {
	var err error

	if netListener != nil {
		err = server.Serve(netListener)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		ch <- err
	}
}

func (l *HTTPListener) shutdown(server *http.Server, timeout time.Duration) error {
	var cancel func()
	ctx := context.Background()

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	server.SetKeepAlivesEnabled(false)
	return server.Shutdown(ctx)
}

func errorLog(logger Logger) *log.Logger {
	if logger == nil {
		return log.New(os.Stderr, "", log.LstdFlags)
	}

	return log.New(writerFunc(func(p []byte) (n int, err error) {
		l := len(p)
		if bytes.HasPrefix(p, []byte("http: panic serving ")) {
			// Skip logging of panic, handled by server itself. This kind of errors would be logged
			// by middlewares
			return l, nil
		}
		p = bytes.TrimRight(p, "\n")
		logger.Error(context.Background(), string(p))

		return l, nil
	}), "", 0)
}
