// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package server

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"github.com/gorhythm/concerto"
	"github.com/gorhythm/concerto/transport"

	pb "github.com/gorhythm/concerto/sample/calc/concerto/proto/gen-go/concerto/sample/calc/v1"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/grpc"
)

// config is a group of options for a Server.
type config struct {
	registry *grpc.Registry
	options  []kitgrpc.ServerOption
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

func WithTransportOptions(
	opts ...kitgrpc.ServerOption,
) Option {
	return optionFunc(func(c config) config {
		c.options = append(c.options, opts...)
		return c
	})
}

type server struct {
	pb.UnimplementedCalculatorServiceServer
	calculatorHandler kitgrpc.Handler
}

// New makes a set of endpoints available as a gRPC
// CalculatorServiceServer.
func New(endpoints *endpoint.Set, opts ...Option) pb.CalculatorServiceServer {
	cfg := newConfig(opts...)
	return &server{
		calculatorHandler: kitgrpc.NewServer(
			endpoints.CalculateEndpoint,
			cfg.registry.DecodeCalculateRequest,
			cfg.registry.EncodeCalculateResponse,
			cfg.options...,
		),
	}
}

func (s *server) Calculate(
	ctx context.Context, req *pb.CalculateRequest,
) (*pb.CalculateResponse, error) {
	_, resp, err := s.calculatorHandler.ServeGRPC(
		concerto.ContextWithCallMeta(
			ctx,
			concerto.CallMeta{
				Service:   calculator.ServiceDesc.Service,
				Method:    calculator.ServiceDesc.Methods.Calculate.Name,
				Transport: transport.TransportGRPC,
			},
		),
		req,
	)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CalculateResponse), nil
}
