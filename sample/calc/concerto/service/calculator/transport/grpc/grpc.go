// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package grpc

import (
	icontext "context"

	imessage "github.com/gorhythm/concerto/sample/calc/concerto/message"
	icalc "github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
)

type Registry struct {
	EncodeCalculateRequest  func(icontext.Context, any) (any, error)
	DecodeCalculateRequest  func(icontext.Context, any) (any, error)
	EncodeCalculateResponse func(icontext.Context, any) (any, error)
	DecodeCalculateResponse func(icontext.Context, any) (any, error)
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

func encodeCalculateRequest(_ icontext.Context, aReq any) (any, error) {
	return icalc.EncodeCalculateRequest(aReq.(*imessage.CalculateRequest))
}

func decodeCalculateRequest(_ icontext.Context, aReq any) (any, error) {
	return icalc.DecodeCalculateRequest(aReq.(*icalc.CalculateRequest))
}

func encodeCalculateResponse(_ icontext.Context, aResp any) (any, error) {
	return icalc.EncodeCalculateResponse(aResp.(*imessage.CalculateResponse))
}

func decodeCalculateResponse(_ icontext.Context, aResp any) (any, error) {
	return icalc.DecodeCalculateResponse(aResp.(*icalc.CalculateResponse))
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
		icontext.Context, *imessage.CalculateRequest,
	) (*icalc.CalculateRequest, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateRequest = func(
			ctx icontext.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*imessage.CalculateRequest))
		}
		return reg
	})
}

func WithDecodeCalculateRequest(
	fn func(
		icontext.Context, *icalc.CalculateRequest,
	) (*imessage.CalculateRequest, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateRequest = func(
			ctx icontext.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*icalc.CalculateRequest))
		}
		return reg
	})
}

func WithEncodeCalculateResponse(
	fn func(
		icontext.Context, *imessage.CalculateResponse,
	) (*icalc.CalculateResponse, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateResponse = func(
			ctx icontext.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*imessage.CalculateResponse))
		}
		return reg
	})
}

func WithDecodeCalculateResponse(
	fn func(
		icontext.Context, *icalc.CalculateResponse,
	) (*imessage.CalculateResponse, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateResponse = func(
			ctx icontext.Context, obj any,
		) (any, error) {
			return fn(ctx, obj.(*icalc.CalculateResponse))
		}
		return reg
	})
}
