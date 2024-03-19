// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package client

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	kitendpoint "github.com/go-kit/kit/endpoint"

	"github.com/gorhythm/concerto"
	"github.com/gorhythm/concerto/endpoint/middleware/callmeta"
	"github.com/gorhythm/concerto/transport"
	concertothrift "github.com/gorhythm/concerto/transport/thrift"

	"github.com/gorhythm/concerto/sample/calc"
	"github.com/gorhythm/concerto/sample/calc/concerto/message"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	calculatorthrift "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/thrift"
	thriftgen "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

// config is a group of options for a Client.
type config struct {
	registry    *calculatorthrift.Registry
	middlewares []kitendpoint.Middleware
	before      []concertothrift.ClientRequestFunc
}

// newConfig applies all the options to a returned config.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.registry == nil {
		cfg.registry = calculatorthrift.DefaultRegistry
	}
	return cfg
}

// A Option sets options to config.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(cfg config) config {
	return fn(cfg)
}

func WithRegistry(registry *calculatorthrift.Registry) Option {
	return optionFunc(func(c config) config {
		c.registry = registry
		return c
	})
}

func WithMiddlewares(
	middlewares ...kitendpoint.Middleware,
) Option {
	return optionFunc(func(c config) config {
		c.middlewares = append(c.middlewares, middlewares...)
		return c
	})
}

func WithClientBefore(
	before ...concertothrift.ClientRequestFunc,
) Option {
	return optionFunc(func(c config) config {
		c.before = append(c.before, before...)
		return c
	})
}

// New returns an CalculatorService backed by a gRPC server at the other
// end of the conn.
func New(client *thriftgen.CalculatorServiceClient, opts ...Option) calc.CalculatorService {
	cfg := newConfig(opts...)

	calculateEndpoint := func(_ctx context.Context, _req any) (any, error) {
		op, num1, num2, err := cfg.registry.EncodeCalculateRequest(
			_req.(*message.CalculateRequest),
		)
		if err != nil {
			return nil, err
		}

		_hm := thrift.THeaderMap{}
		for _, fn := range cfg.before {
			fn(_ctx, _hm)
		}

		_hmKeys := make([]string, 0, len(_hm))
		for k, v := range _hm {
			thrift.SetHeader(_ctx, k, v)
			_hmKeys = append(_hmKeys, k)
		}
		thrift.SetWriteHeaderList(_ctx, _hmKeys)

		_resp, err := client.Calculate(_ctx, op, num1, num2)
		if err != nil {
			return nil, err
		}

		return cfg.registry.DecodeCalculateResponse(_resp)
	}

	return &endpoint.Set{
		CalculateEndpoint: kitendpoint.Chain(callmeta.Middleware(
			concerto.CallMeta{
				Service:   calculator.ServiceDesc.Service,
				Method:    calculator.ServiceDesc.Methods.Calculate.Name,
				Transport: transport.TransportThrift,
			},
		), cfg.middlewares...)(calculateEndpoint),
	}
}
