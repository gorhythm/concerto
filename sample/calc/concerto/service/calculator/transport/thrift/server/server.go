// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package server

import (
	"context"

	"github.com/gorhythm/concerto"
	"github.com/gorhythm/concerto/transport"
	concertothrift "github.com/gorhythm/concerto/transport/thrift"

	"github.com/gorhythm/concerto/sample/calc/concerto/message"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	"github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/thrift"
	thriftgen "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

// config is a group of options for a Server.
type config struct {
	registry *thrift.Registry
	before   []concertothrift.ServerRequestFunc
}

// newConfig applies all the options to a returned config.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.registry == nil {
		cfg.registry = thrift.DefaultRegistry
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

func WithRegistry(registry *thrift.Registry) Option {
	return optionFunc(func(c config) config {
		c.registry = registry
		return c
	})
}

func WithServerBefore(
	before ...concertothrift.ServerRequestFunc,
) Option {
	return optionFunc(func(c config) config {
		c.before = append(c.before, before...)
		return c
	})
}

type server struct {
	cfg       config
	endpoints *endpoint.Set
}

// New makes a set of endpoints available as a Thrift service.
func New(endpoints *endpoint.Set, opts ...Option) thriftgen.CalculatorService {
	cfg := newConfig(opts...)
	return &server{
		cfg:       cfg,
		endpoints: endpoints,
	}
}

func (s *server) Calculate(
	ctx context.Context, op thriftgen.Op, num1 int64, num2 int64,
) (result int64, err error) {
	req, err := s.cfg.registry.DecodeCalculateRequest(op, num1, num2)
	if err != nil {
		return 0, err
	}

	for _, f := range s.cfg.before {
		ctx = f(ctx)
	}

	resp, err := s.endpoints.CalculateEndpoint(
		concerto.ContextWithCallMeta(
			ctx,
			concerto.CallMeta{
				Service:   calculator.ServiceDesc.Service,
				Method:    calculator.ServiceDesc.Methods.Calculate.Name,
				Transport: transport.TransportThrift,
			},
		),
		req,
	)
	if err != nil {
		return 0, err
	}

	return s.cfg.registry.EncodeCalculateResponse(
		resp.(*message.CalculateResponse),
	)
}
