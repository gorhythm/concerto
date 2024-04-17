package http

import (
	"context"
	"net/http"
)

// ServerRequestFunc may take information from an HTTP request and put it into a request context.
// In Servers, ServerRequestFunc are executed prior to invoking the endpoint.
type ServerRequestFunc func(context.Context, *http.Request) context.Context

// SetRequestHeader returns a RequestFunc that sets the given header.
func SetRequestHeader(key, val string) ServerRequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		r.Header.Set(key, val)
		return ctx
	}
}

// SetRequestContext returns a ServerRequestFunc that associates a value from
// the request context with a new user-defined key.
func SetRequestContext(ctxKey any, val string) ServerRequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		ctx = context.WithValue(ctx, ctxKey, val)
		return ctx
	}
}

// ServerResponseFunc may take information from a request context and use it to manipulate a ResponseWriter.
// ServerResponseFuncs are only executed in servers, after invoking the endpoint but prior to writing a response.
type ServerResponseFunc func(context.Context, http.ResponseWriter) context.Context

// SetResponseHeader returns a ServerResponseFunc that sets the given header.
func SetResponseHeader(key, val string) ServerResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter) context.Context {
		w.Header().Set(key, val)
		return ctx
	}
}

// ClientRequestFunc creates an outgoing HTTP request based on the passed request object.
// It's designed to be used in HTTP clients, for client-side endpoints.
// It's a more powerful version of EncodeRequestFunc,
// and can be used if more fine-grained control of the HTTP request is required.
type ClientRequestFunc func(context.Context, interface{}) (*http.Request, error)

// ClientResponseFunc may take information from an HTTP request and make the response available for consumption.
// ClientResponseFuncs are only executed in clients, after a request has been made, but prior to it being decoded.
type ClientResponseFunc func(context.Context, *http.Response) context.Context
