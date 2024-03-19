// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

// ClientRequestFunc may take information from context and use it to construct
// headers to be transported to the server. ClientRequestFuncs are executed
// after creating the request but prior to sending the Thrift request to the
// server.
type ClientRequestFunc func(ctx context.Context, hm thrift.THeaderMap) context.Context

// ServerRequestFunc may extracts information from the request context and
// associates it with a new user-defined key. ServerRequestFunc allows for
// unified access these data in service or endpoint middleware, similar to other
// transport types (such as gRPC, HTTP, ...). ServerRequestFuncs are executed
// prior to invoking the endpoint.
type ServerRequestFunc func(ctx context.Context) context.Context

// SetRequestHeader returns a ClientRequestFunc that sets the specified header
// key-value pair.
func SetRequestHeader(key, val string) ClientRequestFunc {
	return func(ctx context.Context, hm thrift.THeaderMap) context.Context {
		hm[key] = val
		return ctx
	}
}

// SetRequestContext returns a ServerRequestFunc that associates a value from
// the request context with a new user-defined key.
func SetRequestContext(ctxKey any, headerKey string) ServerRequestFunc {
	return func(ctx context.Context) context.Context {
		val, found := thrift.GetHeader(ctx, headerKey)
		if !found {
			return ctx
		}
		return context.WithValue(ctx, ctxKey, val)
	}
}
