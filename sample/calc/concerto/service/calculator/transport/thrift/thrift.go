// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package thrift

import (
	"github.com/gorhythm/concerto/sample/calc/concerto/message"
	thrift "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

type Registry struct {
	EncodeCalculateRequest  func(*message.CalculateRequest) (thrift.Op, int64, int64, error)
	DecodeCalculateRequest  func(thrift.Op, int64, int64) (*message.CalculateRequest, error)
	EncodeCalculateResponse func(*message.CalculateResponse) (int64, error)
	DecodeCalculateResponse func(int64) (*message.CalculateResponse, error)
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

func encodeCalculateRequest(_obj *message.CalculateRequest) (op thrift.Op, num1 int64, num2 int64, err error) {
	return thrift.EncodeCalculateRequest(_obj)
}

func decodeCalculateRequest(op thrift.Op, num1 int64, num2 int64) (*message.CalculateRequest, error) {
	return thrift.DecodeCalculateRequest(op, num1, num2)
}

func encodeCalculateResponse(_obj *message.CalculateResponse) (int64, error) {
	return thrift.EncodeCalculateResponse(_obj)
}

func decodeCalculateResponse(result int64) (*message.CalculateResponse, error) {
	return thrift.DecodeCalculateResponse(result)
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
		*message.CalculateRequest,
	) (thrift.Op, int64, int64, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateRequest = fn
		return reg
	})
}

func WithDecodeCalculateRequest(
	fn func(
		thrift.Op, int64, int64,
	) (*message.CalculateRequest, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateRequest = fn
		return reg
	})
}

func WithEncodeCalculateResponse(
	fn func(
		*message.CalculateResponse,
	) (int64, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateResponse = fn
		return reg
	})
}

func WithDecodeCalculateResponse(
	fn func(
		int64,
	) (*message.CalculateResponse, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateResponse = fn
		return reg
	})
}
