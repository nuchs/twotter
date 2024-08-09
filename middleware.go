package main

import "net/http"

// Middleware is represented as a function that take an http.Handler and returns
// a new http.Handler. The new handler composes the functionality of the
// original handler with the middleware.
type MiddleWare func(h http.Handler) http.Handler

// bind this middleware to the next one, producing a new middleware which is
// comprised of the functionality of both. When handling request, the current
// middleware will get the opporunity to manage the request before passing it to
// the next middleware.
func (mw MiddleWare) bind(next MiddleWare) MiddleWare {
	return func(h http.Handler) http.Handler {
		return mw(next(h))
	}
}

// Call is a convenience function that allows a middleware chain to terminate in
// a http.HandlerFunc rather than a http.Handler.
func (mw MiddleWare) Call(h http.HandlerFunc) http.Handler {
	return mw(http.HandlerFunc(h))
}

// ChainMiddleWare is used to bind a sequence of middlewares into a single one. The chain
// can be reused with multiple terminating http.Handlers.
func ChainMiddleWare(mws ...MiddleWare) MiddleWare {
	if len(mws) == 0 {
		return nil
	}

	mw := mws[0]
	for _, m := range mws[1:] {
		mw = mw.bind(m)
	}

	return mw
}
