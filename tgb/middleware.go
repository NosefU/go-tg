package tgb

import "context"

// Middleware is used to modify or process the Update before it is passed to the handler.
// All middleware should be registered before the handlers registration.
//
// Deprecated: Causes performance issues and breaks routing logic. Use GlobalMiddlewareFunc
type Middleware interface {
	Wrap(Handler) Handler
}

type GlobalMiddlewareFunc func(context.Context, *Update) (context.Context, *Update, error)

// MiddlewareFunc is used to modify or process the Update before it is passed to the handler.
// All middleware should be registered before the handlers registration.
//
// Deprecated: Causes performance issues and breaks routing logic. Use GlobalMiddlewareFunc
type MiddlewareFunc func(Handler) Handler

func (m MiddlewareFunc) Wrap(h Handler) Handler {
	return m(h)
}

type globalChain []GlobalMiddlewareFunc

// Append extends a chain, adding the specified middleware
// as the last ones in the request flow.
func (c globalChain) Append(mws ...GlobalMiddlewareFunc) globalChain {
	result := make(globalChain, 0, len(c)+len(mws))
	result = append(result, c...)
	result = append(result, mws...)
	return result
}

// chain - is middleware chain, that stores registered middlewares
//
// Deprecated: Causes performance issues and breaks routing logic. Use globalChain
type chain []Middleware

// Append extends a chain, adding the specified middleware
// as the last ones in the request flow.
func (c chain) Append(mws ...Middleware) chain {
	result := make(chain, 0, len(c)+len(mws))
	result = append(result, c...)
	result = append(result, mws...)
	return result
}

// Then wraps handler with middleware chain.
func (c chain) Then(handler Handler) Handler {
	for i := range c {
		handler = c[len(c)-1-i].Wrap(handler)
	}

	return handler
}
