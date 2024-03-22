// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// The MIT License (MIT)
// Copyright (c) 2015 Peter Bourgon.

package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	"golang.org/x/exp/maps"
)

type contextKey int

const (
	contextKeyResponseMeta contextKey = iota
)

// ContextWithResponseMeta returns a copy of ctx with v set as the response
// metadata value.
func ContextWithResponseMeta(ctx context.Context, v *Metadata) context.Context {
	return context.WithValue(ctx, contextKeyResponseMeta, v)
}

// ResponseMetaFromContext returns the response metadata from ctx.
func ResponseMetaFromContext(ctx context.Context) *Metadata {
	if ctx == nil {
		return nil
	}

	if v, ok := ctx.Value(contextKeyResponseMeta).(*Metadata); ok {
		return v
	}

	return nil
}

// Header contains methods for managing Thrift headers.
type Header interface {
	// Get returns the value associated with the given key. If the key is
	// present in the header, the value (which may be empty) is returned.
	// Otherwise, the returned value will be the empty string. The ok return
	// value reports whether the value is explicitly set in the header.
	Get(key string) (value string, ok bool)

	// Set sets the header entry associated with key to the value.
	Set(key, value string)

	// Keys returns the keys of the header.
	Keys() []string
}

// HeaderMap represents a Thrift header backed by a map.
type HeaderMap map[string]string

// Get returns the value associated with the given key. If the key is present
// in the header, the value (which may be empty) is returned. Otherwise, the
// returned value will be the empty string. The ok return value reports whether
// the value is explicitly set in the header.
func (h HeaderMap) Get(key string) (value string, ok bool) {
	v, found := h[key]
	return v, found
}

// Set sets the value associated with the given key in the header map.
func (h HeaderMap) Set(key string, value string) {
	h[key] = value
}

// Keys returns the keys of the header.
func (h HeaderMap) Keys() []string {
	return maps.Keys(h)
}

// HeaderContext represents a Thrift header backed by a context.
type HeaderContext struct {
	ctx context.Context
}

// NewHeaderContext creates and returns a new instance of HeaderContext with
// provider context.
func NewHeaderContext(ctx context.Context) *HeaderContext {
	return &HeaderContext{ctx}
}

// Get returns the value associated with the given key. If the key is present
// in the header, the value (which may be empty) is returned. Otherwise, the
// returned value will be the empty string. The ok return value reports whether
// the value is explicitly set in the header.
func (h *HeaderContext) Get(key string) (string, bool) {
	return thrift.GetHeader(h.ctx, key)
}

// Set is a no-op method for HeaderContext, implementing the Header interface.
// NOTE: Cannot set to HeaderContext.
func (h *HeaderContext) Set(string, string) {}

// Keys returns a slice containing all the keys present in the context's
// headers.
func (h *HeaderContext) Keys() []string {
	return thrift.GetReadHeaderList(h.ctx)
}

// Metadata represents the metadata attached to the request or response.
type Metadata struct {
	// The headers in the request or response.
	Header Header
}

// ClientRequestFunc may take information from context and use it to construct
// headers to be transported to the server. ClientRequestFuncs are executed
// after creating the request but prior to sending the Thrift request to the
// server.
type ClientRequestFunc func(ctx context.Context, md *Metadata) context.Context

// ServerRequestFunc may extracts information from the request context and
// associates it with a new user-defined key. ServerRequestFunc allows for
// unified access these data in service or endpoint middleware, similar to other
// transport types (such as gRPC, HTTP, ...). ServerRequestFuncs are executed
// prior to invoking the endpoint.
type ServerRequestFunc func(ctx context.Context, md Metadata) context.Context

// ServerResponseFunc may take information from a request context and use it to
// manipulate the gRPC response metadata headers and trailers. ResponseFuncs are
// only executed in servers, after invoking the endpoint but prior to writing a
// response.
type ServerResponseFunc func(ctx context.Context, md *Metadata) context.Context

// ClientResponseFunc may take information from a Thrift response metadata and
// make the responses available for consumption. ClientResponseFuncs are only
// executed in clients, after a request has been made, but prior to it being
// decoded.
type ClientResponseFunc func(ctx context.Context, md Metadata) context.Context

// SetRequestHeader returns a ClientRequestFunc that sets the specified header
// key-value pair.
func SetRequestHeader(key, val string) ClientRequestFunc {
	return func(ctx context.Context, md *Metadata) context.Context {
		md.Header.Set(key, val)
		return ctx
	}
}

// SetRequestContext returns a ServerRequestFunc that associates a value from
// the request context with a new user-defined key.
func SetRequestContext(ctxKey any, headerKey string) ServerRequestFunc {
	return func(ctx context.Context, md Metadata) context.Context {
		val, found := md.Header.Get(headerKey)
		if !found {
			return ctx
		}

		return context.WithValue(ctx, ctxKey, val)
	}
}
