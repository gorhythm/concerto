// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package client

import (
	ikitendpoint "github.com/go-kit/kit/endpoint"
	ikitgrpc "github.com/go-kit/kit/transport/grpc"
	igooglegrpc "google.golang.org/grpc"

	iconcerto "github.com/gorhythm/concerto"
	iconcertocallmeta "github.com/gorhythm/concerto/endpoint/middleware/callmeta"
	iconcertotransport "github.com/gorhythm/concerto/transport"

	icalc "github.com/gorhythm/concerto/sample/calc"
	icalcpb "github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
	icalculator "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	icalculatorendpoint "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	icalculatorgrpc "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/grpc"
)

// config is a group of options for a Client.
type config struct {
	registry    *icalculatorgrpc.Registry
	middlewares []ikitendpoint.Middleware
	options     []ikitgrpc.ClientOption
}

// newConfig applies all the options to a returned config.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.registry == nil {
		cfg.registry = icalculatorgrpc.DefaultRegistry
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

func WithRegistry(registry *icalculatorgrpc.Registry) Option {
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

func WithTransportOptions(
	opts ...ikitgrpc.ClientOption,
) Option {
	return optionFunc(func(c config) config {
		c.options = append(c.options, opts...)
		return c
	})
}

// New returns an CalculatorService backed by a gRPC server at the other
// end of the conn.
func New(conn *igooglegrpc.ClientConn, opts ...Option) icalc.CalculatorService {
	cfg := newConfig(opts...)
	calculateEndpoint := ikitgrpc.NewClient(
		conn,
		"concerto.sample.calc.v1.CalculatorService",
		"Calculate",
		cfg.registry.EncodeCalculateRequest,
		cfg.registry.DecodeCalculateResponse,
		icalcpb.CalculateResponse{},
		cfg.options...,
	).Endpoint()

	return &icalculatorendpoint.Set{
		CalculateEndpoint: ikitendpoint.Chain(iconcertocallmeta.Middleware(
			iconcerto.CallMeta{
				Service:   icalculator.ServiceDesc.Service,
				Method:    icalculator.ServiceDesc.Methods.Calculate.Name,
				Transport: iconcertotransport.TransportGRPC,
			},
		), cfg.middlewares...)(calculateEndpoint),
	}
}
