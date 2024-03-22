// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package client

import (
	icontext "context"

	ithrift "github.com/apache/thrift/lib/go/thrift"
	ikitendpoint "github.com/go-kit/kit/endpoint"

	iconcerto "github.com/gorhythm/concerto"
	iconcertocallmeta "github.com/gorhythm/concerto/endpoint/middleware/callmeta"
	iconcertotransport "github.com/gorhythm/concerto/transport"
	iconcertothrift "github.com/gorhythm/concerto/transport/thrift"

	icalc "github.com/gorhythm/concerto/sample/calc"
	imessage "github.com/gorhythm/concerto/sample/calc/concerto/message"
	icalculator "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	icalculatorendpoint "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	icalculatorthrift "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/thrift"
	ithriftgen "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

// config is a group of options for a Client.
type config struct {
	registry    *icalculatorthrift.Registry
	middlewares []ikitendpoint.Middleware
	before      []iconcertothrift.ClientRequestFunc
	after       []iconcertothrift.ClientResponseFunc
	finalizer   []iconcertothrift.ClientFinalizerFunc
}

// newConfig applies all the options to a returned config.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.registry == nil {
		cfg.registry = icalculatorthrift.DefaultRegistry
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

func WithRegistry(registry *icalculatorthrift.Registry) Option {
	return optionFunc(func(c config) config {
		c.registry = registry
		return c
	})
}

func WithMiddlewares(
	middlewares ...ikitendpoint.Middleware,
) Option {
	return optionFunc(func(c config) config {
		c.middlewares = append(c.middlewares, middlewares...)
		return c
	})
}

func WithBefore(
	before ...iconcertothrift.ClientRequestFunc,
) Option {
	return optionFunc(func(c config) config {
		c.before = append(c.before, before...)
		return c
	})
}

func WithAfter(
	after ...iconcertothrift.ClientResponseFunc,
) Option {
	return optionFunc(func(c config) config {
		c.after = append(c.after, after...)
		return c
	})
}

func WithFinalizer(fn ...iconcertothrift.ClientFinalizerFunc) Option {
	return optionFunc(func(c config) config {
		c.finalizer = append(c.finalizer, fn...)
		return c
	})
}

// New returns an CalculatorService backed by a Thrift client at the other
// end of the conn.
func New(tclient ithrift.TClient, opts ...Option) icalc.CalculatorService {
	var (
		client = ithriftgen.NewCalculatorServiceClient(
			iconcertothrift.ExtractResponseMeta(tclient),
		)
		cfg = newConfig(opts...)
	)

	calculateEndpoint := func(ctx icontext.Context, aReq any) (_ any, err error) {
		ctx, cancel := icontext.WithCancel(ctx)
		defer cancel()

		if cfg.finalizer != nil {
			defer func() {
				for _, _f := range cfg.finalizer {
					_f(ctx, err)
				}
			}()
		}

		_op, _num1, _num2, err := cfg.registry.EncodeCalculateRequest(
			aReq.(*imessage.CalculateRequest),
		)
		if err != nil {
			return nil, err
		}

		reqMD := iconcertothrift.Metadata{
			Header: iconcertothrift.HeaderMap{},
		}
		for _, _f := range cfg.before {
			_f(ctx, &reqMD)
		}

		var respMD iconcertothrift.Metadata
		r, err := client.Calculate(
			iconcertothrift.ContextWithResponseMeta(ctx, &respMD),
			_op, _num1, _num2,
		)
		if err != nil {
			return nil, err
		}

		for _, fn := range cfg.after {
			ctx = fn(ctx, respMD)
		}

		return cfg.registry.DecodeCalculateResponse(r)
	}

	return &icalculatorendpoint.Set{
		CalculateEndpoint: ikitendpoint.Chain(iconcertocallmeta.Middleware(
			iconcerto.CallMeta{
				Service:   icalculator.ServiceDesc.Service,
				Method:    icalculator.ServiceDesc.Methods.Calculate.Name,
				Transport: iconcertotransport.TransportThrift,
			},
		), cfg.middlewares...)(calculateEndpoint),
	}
}
