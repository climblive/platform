package rest

import (
	"net/http"
	"slices"
)

type Middleware = func(http.Handler) http.Handler

type Mux struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewMux() *Mux {
	return &Mux{
		mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
	}
}

func (m *Mux) RegisterMiddleware(mw Middleware) {
	m.middlewares = append(m.middlewares, mw)
}

func (m *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.middlewares[0](m.middlewares[1](http.HandlerFunc(handler)))

	var chain http.Handler = http.HandlerFunc(handler)

	for _, mw := range slices.Backward(m.middlewares) {
		chain = mw(chain)
	}

	m.mux.Handle(pattern, chain)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
