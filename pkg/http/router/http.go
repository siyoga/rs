package router

import "net/http"

type router struct {
	router *http.ServeMux
	logger Logger

	badReqHandler      BadRequestErrHandler
	internalErrHandler InternalErrHandler

	middlewares []func(string, http.Handler) http.Handler
}

type Router interface {
	Use(middleware func(string, http.Handler) http.Handler)
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func NewRouter(
	logger Logger,
	badReqHandler BadRequestErrHandler,
	internalErrHandler InternalErrHandler,
	panicHandler PanicHandler,
) Router {
	r := &router{
		router:             http.NewServeMux(),
		logger:             logger,
		badReqHandler:      badReqHandler,
		internalErrHandler: internalErrHandler,
		middlewares:        make([]func(string, http.Handler) http.Handler, 0),
	}

	r.Use(panicMiddleware(internalErrHandler, panicHandler))

	return r
}

// Use добавляет middleware в цепочку
func (r *router) Use(middleware func(string, http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

// Handle регистрирует обработчик с применением всех middlewares
func (r *router) Handle(pattern string, handler http.Handler) {
	finalHandler := handler

	// Применяем middlewares в обратном порядке (последний добавленный выполнится первым)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](pattern, finalHandler)
	}

	r.router.Handle(pattern, finalHandler)
}

// HandleFunc регистрирует функцию-обработчик с применением всех middlewares
func (r *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Handle(pattern, http.HandlerFunc(handler))
}

// ServeHTTP реализует интерфейс http.Handler
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
