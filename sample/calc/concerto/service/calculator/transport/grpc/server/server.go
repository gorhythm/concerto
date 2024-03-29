// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package server

import (
	icontext "context"

	ikitgrpc "github.com/go-kit/kit/transport/grpc"

	iconcerto "github.com/gorhythm/concerto"
	iconcertotransport "github.com/gorhythm/concerto/transport"

	icalc "github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
	icalculator "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	icalculatorendpoint "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	icalculatorgrpc "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/grpc"
)

// config is a group of options for a Server.
type config struct {
	registry *icalculatorgrpc.Registry
	options  []ikitgrpc.ServerOption
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

func WithTransportOptions(
	opts ...ikitgrpc.ServerOption,
) Option {
	return optionFunc(func(c config) config {
		c.options = append(c.options, opts...)
		return c
	})
}

type server struct {
	icalc.UnimplementedCalculatorServiceServer
	calculatorHandler ikitgrpc.Handler
}

// New makes a set of endpoints available as a gRPC
// CalculatorServiceServer.
func New(endpoints *icalculatorendpoint.Set, opts ...Option) icalc.CalculatorServiceServer {
	cfg := newConfig(opts...)
	return &server{
		calculatorHandler: ikitgrpc.NewServer(
			endpoints.CalculateEndpoint,
			cfg.registry.DecodeCalculateRequest,
			cfg.registry.EncodeCalculateResponse,
			cfg.options...,
		),
	}
}

func (s *server) Calculate(
	ctx icontext.Context, req *icalc.CalculateRequest,
) (*icalc.CalculateResponse, error) {
	_, aResp, err := s.calculatorHandler.ServeGRPC(
		iconcerto.ContextWithCallMeta(
			ctx,
			iconcerto.CallMeta{
				Service:   icalculator.ServiceDesc.Service,
				Method:    icalculator.ServiceDesc.Methods.Calculate.Name,
				Transport: iconcertotransport.TransportGRPC,
			},
		),
		req,
	)
	if err != nil {
		return nil, err
	}

	return aResp.(*icalc.CalculateResponse), nil
}
