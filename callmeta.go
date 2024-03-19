// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Copyright (c) The go-grpc-middleware Authors.
// Licensed under the Apache License 2.0.

package concerto

import (
	"context"

	"github.com/gorhythm/concerto/transport"
)

type contextKeyType int

const (
	callMetaKey contextKeyType = iota
)

var (
	// NilCallMeta is empty CallMeta, all zero values.
	NilCallMeta CallMeta
)

// ContextWithCallMeta returns a copy of ctx with v set as the call meta value.
func ContextWithCallMeta(ctx context.Context, v CallMeta) context.Context {
	return context.WithValue(ctx, callMetaKey, v)
}

// CallMetaFromContext returns the call meta from ctx.
func CallMetaFromContext(ctx context.Context) CallMeta {
	if ctx == nil {
		return NilCallMeta
	}

	if v, ok := ctx.Value(callMetaKey).(CallMeta); ok {
		return v
	}

	return NilCallMeta
}

// MethodInfo represents information about a specific method.
type MethodInfo struct {
	// Name is the name of the method.
	Name string
	// FullName is the fully qualified name of the method.
	FullName string
}

// CallMeta contains metadata related to a method call.
type CallMeta struct {
	// Service is the name of the service containing the method.
	Service string
	// Method is the name of the specific method.
	Method string
	// Transport is the transport used for the call.
	Transport transport.Transport
}

// FullMethod returns the fully qualified method name (service/method).
// Example: "concerto.test.v1.TestService/Ping".
func (c CallMeta) FullMethod() string {
	return c.Service + "/" + c.Method
}
