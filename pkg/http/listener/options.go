package listener

import (
	"net"
	"net/http"
	"time"
)

type Option func(*options)

type options struct {
	netListener     net.Listener
	server          []func(server *http.Server)
	mw              []func(handler http.Handler) http.Handler
	shutdownTimeout time.Duration
	logger          Logger

	recoverPanics bool
}

func With(logger Logger) Option {
	return func(o *options) {
		if logger != nil {
			WithLogger(logger)(o)
		}
	}
}

func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.server = append(
			o.server,
			func(s *http.Server) {
				s.IdleTimeout = timeout
			},
		)
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.server = append(
			o.server,
			func(s *http.Server) {
				s.ReadTimeout = timeout
			})
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.server = append(
			o.server,
			func(s *http.Server) {
				s.WriteTimeout = timeout
			})
	}
}
