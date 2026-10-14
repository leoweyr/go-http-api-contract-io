package server

import (
	"net/http"
)

// MiddlewareChain sequences a set of middlewares applied in declaration order around a final handler.
type MiddlewareChain struct {
	middlewares []Middleware
}

// NewMiddlewareChain creates a MiddlewareChain from the given middlewares so the first declared runs outermost.
func NewMiddlewareChain(middlewares ...Middleware) MiddlewareChain {
	var chain MiddlewareChain = MiddlewareChain{middlewares: middlewares}

	return chain
}

// Use returns a new MiddlewareChain with the given middleware appended to the sequence.
func (chain MiddlewareChain) Use(middleware Middleware) MiddlewareChain {
	var extended []Middleware = make([]Middleware, 0, len(chain.middlewares)+1)

	extended = append(extended, chain.middlewares...)
	extended = append(extended, middleware)

	return MiddlewareChain{middlewares: extended}
}

// Then wraps the final handler with every middleware in the chain, applying them so the first declared middleware runs outermost, and returns the composed handler.
func (chain MiddlewareChain) Then(handler http.Handler) http.Handler {
	var composed http.Handler = handler

	var index int

	for index = len(chain.middlewares) - 1; index >= 0; index = index - 1 {
		composed = chain.middlewares[index].Handle(composed)
	}

	return composed
}
