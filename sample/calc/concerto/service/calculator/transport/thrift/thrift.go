// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package thrift

import (
	imessage "github.com/gorhythm/concerto/sample/calc/concerto/message"
	iv1 "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

type Registry struct {
	EncodeCalculateRequest  func(*imessage.CalculateRequest) (iv1.Op, int64, int64, error)
	DecodeCalculateRequest  func(iv1.Op, int64, int64) (*imessage.CalculateRequest, error)
	EncodeCalculateResponse func(*imessage.CalculateResponse) (int64, error)
	DecodeCalculateResponse func(int64) (*imessage.CalculateResponse, error)
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

func encodeCalculateRequest(req *imessage.CalculateRequest) (_op iv1.Op, _num1 int64, _num2 int64, err error) {
	return iv1.EncodeCalculateRequest(req)
}

func decodeCalculateRequest(_op iv1.Op, _num1 int64, _num2 int64) (req *imessage.CalculateRequest, err error) {
	return iv1.DecodeCalculateRequest(_op, _num1, _num2)
}

func encodeCalculateResponse(resp *imessage.CalculateResponse) (_result int64, err error) {
	return iv1.EncodeCalculateResponse(resp)
}

func decodeCalculateResponse(_result int64) (resp *imessage.CalculateResponse, err error) {
	return iv1.DecodeCalculateResponse(_result)
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
	fn func(*imessage.CalculateRequest) (op iv1.Op, num1 int64, num2 int64, _ error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateRequest = fn
		return reg
	})
}

func WithDecodeCalculateRequest(
	fn func(op iv1.Op, num1 int64, num2 int64) (*imessage.CalculateRequest, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateRequest = fn
		return reg
	})
}

func WithEncodeCalculateResponse(
	fn func(*imessage.CalculateResponse) (result int64, _ error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.EncodeCalculateResponse = fn
		return reg
	})
}

func WithDecodeCalculateResponse(
	fn func(result int64) (*imessage.CalculateResponse, error),
) RegistryOption {
	return registryOptionFunc(func(reg Registry) Registry {
		reg.DecodeCalculateResponse = fn
		return reg
	})
}
