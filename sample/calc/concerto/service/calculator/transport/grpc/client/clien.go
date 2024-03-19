// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package client

import (
	kitendpoint "github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	googlegrpc "google.golang.org/grpc"

	"github.com/gorhythm/concerto"
	"github.com/gorhythm/concerto/endpoint/middleware/callmeta"
	"github.com/gorhythm/concerto/transport"

	"github.com/gorhythm/concerto/sample/calc"
	pb "github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/grpc"
)

// config is a group of options for a Client.
type config struct {
	registry    *grpc.Registry
	middlewares []kitendpoint.Middleware
	options     []kitgrpc.ClientOption
}

// newConfig applies all the options to a returned config.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.registry == nil {
		cfg.registry = grpc.DefaultRegistry
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

func WithRegistry(registry *grpc.Registry) Option {
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

func WithTransportOptions(
	opts ...kitgrpc.ClientOption,
) Option {
	return optionFunc(func(c config) config {
		c.options = append(c.options, opts...)
		return c
	})
}

// New returns an CalculatorService backed by a gRPC server at the other
// end of the conn.
func New(conn *googlegrpc.ClientConn, opts ...Option) calc.CalculatorService {
	cfg := newConfig(opts...)
	calculateEndpoint := kitgrpc.NewClient(
		conn,
		"concerto.sample.calc.v1.CalculatorService",
		"Calculate",
		cfg.registry.EncodeCalculateRequest,
		cfg.registry.DecodeCalculateResponse,
		pb.CalculateResponse{},
		cfg.options...,
	).Endpoint()

	return &endpoint.Set{
		CalculateEndpoint: kitendpoint.Chain(callmeta.Middleware(
			concerto.CallMeta{
				Service:   calculator.ServiceDesc.Service,
				Method:    calculator.ServiceDesc.Methods.Calculate.Name,
				Transport: transport.TransportGRPC,
			},
		), cfg.middlewares...)(calculateEndpoint),
	}
}
