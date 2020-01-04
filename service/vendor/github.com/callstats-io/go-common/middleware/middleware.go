package middleware

import "net/http"

// Middleware is the function signature all middlewares are expected to conform to
// Each middleware is expected to either write a final response using the ResponseWriter passed as argument
// OR call the next handler with the writer and the (perhaps modified) request
// Handlers are called in the order they registered:
// For example:
// b := Middleware{}
// b.Use(MiddlewareA)
// b.Use(MiddlewareB)
// b.Use(MiddlewareC)
// b.Handle(Handler)
// is equivalent to A(B(C(D)))
type Middleware func(h http.HandlerFunc) http.HandlerFunc

// Configurator implements the container for middleware
type Configurator struct {
	middlewares []Middleware
}

// NewConfigurator returns a new configurator with initial capacity of 10 middlewares.
// The space is automatically increased if required when a new middleware is added.
func NewConfigurator() *Configurator {
	return NewConfiguratorWithCapacity(10)
}

// NewConfiguratorWithCapacity returns a new Configurator with initial space for N middlewares.
// The space is automatically increased if required when a new middleware is added.
func NewConfiguratorWithCapacity(cap int) *Configurator {
	return &Configurator{
		middlewares: make([]Middleware, 0, cap),
	}
}

// Use registers a Middleware to be run around the actual handler.
func (c *Configurator) Use(h Middleware) {
	c.middlewares = append(c.middlewares, h)
}

// Append registers a Middleware to be run around the actual handler after all already registered middlewares.
func (c *Configurator) Append(h Middleware) *Configurator {
	c.middlewares = append(c.middlewares, h)
	return c
}

// Prepend registers a Middleware to be run around the actual handler before all already registered middlewares.
func (c *Configurator) Prepend(h Middleware) *Configurator {
	if len(c.middlewares) == 0 {
		return c.Append(h)
	}
	last := c.middlewares[len(c.middlewares)-1]
	for i := len(c.middlewares) - 1; i > 0; i-- {
		c.middlewares[i] = c.middlewares[i-1]
	}
	c.middlewares[0] = h
	c.middlewares = append(c.middlewares, last)
	return c
}

// Sink wraps the http.HandlerFunc with the currently registered middlewares.
// Any middlewares registered after calling this function won't be used for the handler.
func (c *Configurator) Sink(h http.HandlerFunc) http.HandlerFunc {
	return c.chain(h, c.middlewares)
}

// ToHandler creates the middleware chain with final call to the given handler func.
// Optionally you can pass an arbitrary list of middlewares that should be run around the HandlerFunc.
func (c *Configurator) ToHandler(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	cc := c.Clone()
	for _, mw := range middlewares {
		cc.Use(mw)
	}
	return cc.chain(h, cc.middlewares)
}

// Clone creates the middleware chain with final call to the given handler func
func (c *Configurator) Clone() *Configurator {
	cNew := NewConfiguratorWithCapacity(len(c.middlewares))
	// doing a copy(dst,src) would require setting the len of new slice so just copy manually
	for _, m := range c.middlewares {
		cNew.middlewares = append(cNew.middlewares, m)
	}
	return cNew
}

func (c *Configurator) chain(handler http.HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	if len(middlewares) == 0 {
		// return final handler func
		return handler
	}

	// call first middleware with the next one as arg etc. up to the handler
	first := middlewares[0]
	next := c.chain(handler, middlewares[1:])
	return first(next)
}
