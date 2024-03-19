// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package grpc

import (
	"context"

	"github.com/gorhythm/concerto/sample/calc/concerto/message"
	"github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
)

type Registry struct {
	EncodeCalculateRequest  func(context.Context, any) (any, error)
	DecodeCalculateRequest  func(context.Context, any) (any, error)
	EncodeCalculateResponse func(context.Context, any) (any, error)
	DecodeCalculateResponse func(context.Context, any) (any, error)
}

// DefaultRegistry is the default registry.
var DefaultRegistry = NewRegistry()

// NewRegistry initializes and return a new instance of Registry with default
// encoder and decoder functions.
func NewRegistry(opts ...RegistryOption) *Registry {
	var reg = Registry{}
	for _, opt := range opts {
		reg = opt.apply(reg)
	}
	if reg.EncodeCalculateRequest == nil {
		reg.EncodeCalculateRequest = encodeCalculateRequest
	}
	if reg.DecodeCalculateRequest == nil {
		reg.DecodeCalculateRequest = decodeCalculateRequest
	}
	if reg.EncodeCalculateResponse == nil {
		reg.EncodeCalculateResponse = encodeCalculateResponse
	}
	if reg.DecodeCalculateResponse == nil {
		reg.DecodeCalculateResponse = decodeCalculateResponse
	}
	return &reg
}

func encodeCalculateRequest(ctx context.Context, obj any) (any, error) {
	return calc.EncodeCalculateRequest(obj.(*message.CalculateRequest))
}

func decodeCalculateRequest(ctx context.Context, obj any) (any, error) {
	return calc.DecodeCalculateRequest(obj.(*calc.CalculateRequest))
}

func encodeCalculateResponse(ctx context.Context, obj any) (any, error) {
	return calc.EncodeCalculateResponse(obj.(*message.CalculateResponse))
}

func decodeCalculateResponse(ctx context.Context, obj any) (any, error) {
	return calc.DecodeCalculateResponse(obj.(*calc.CalculateResponse))
}

// A RegistryOption sets options to Registry.
type RegistryOption interface {
	apply(Registry) Registry
}

type registryOptionFunc func(Registry) Registry

func (fn registryOptionFunc) apply(reg Registry) Registry {
	return fn(reg)
}

func WithEncodeCalculateRequest(
	fn func(
		context.Context, *message.CalculateRequest,
	) (*calc.CalculateRequest, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateRequest = func(
			ctx context.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*message.CalculateRequest))
		}
		return reg
	})
}

func WithDecodeCalculateRequest(
	fn func(
		context.Context, *calc.CalculateRequest,
	) (*message.CalculateRequest, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateRequest = func(
			ctx context.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*calc.CalculateRequest))
		}
		return reg
	})
}

func WithEncodeCalculateResponse(
	fn func(
		context.Context, *message.CalculateResponse,
	) (*calc.CalculateResponse, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateResponse = func(
			ctx context.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*message.CalculateResponse))
		}
		return reg
	})
}

func WithDecodeCalculateResponse(
	fn func(
		context.Context, *calc.CalculateResponse,
	) (*message.CalculateResponse, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateResponse = func(
			ctx context.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*calc.CalculateResponse))
		}
		return reg
	})
}
