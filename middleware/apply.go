package middleware

import "net/http"

// Middleware is a middleware
type Middleware func(http.Handler) http.HandlerFunc

// Apply will apply all middlewares to initial handler
func Apply(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// ApplyFunc is a wrapper for Apply function
// to make it easier to call it with
// implicit HandlerFunc func(w http.ResponseWriter, r *http.Request)
func ApplyFunc(h http.HandlerFunc, mws ...Middleware) http.Handler {
	return Apply(h, mws...)
}
