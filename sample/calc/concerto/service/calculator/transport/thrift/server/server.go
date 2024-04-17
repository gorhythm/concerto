// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package server

import (
	icontext "context"

	ithrift "github.com/apache/thrift/lib/go/thrift"
	ikittransport "github.com/go-kit/kit/transport"

	iconcerto "github.com/gorhythm/concerto"
	iconcertotransport "github.com/gorhythm/concerto/transport"
	iconcertothrift "github.com/gorhythm/concerto/transport/thrift"

	imessage "github.com/gorhythm/concerto/sample/calc/concerto/message"
	icalculator "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator"
	icalculatorendpoint "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/endpoint"
	icalculatorthrift "github.com/gorhythm/concerto/sample/calc/concerto/service/calculator/transport/thrift"
	ithriftgen "github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/concerto/sample/calculator/v1"
)

// config is a group of options for a Server.
type config struct {
	registry     *icalculatorthrift.Registry
	before       []iconcertothrift.ServerRequestFunc
	after        []iconcertothrift.ServerResponseFunc
	finalizer    []iconcertothrift.ServerFinalizerFunc
	errorHandler ikittransport.ErrorHandler
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

// An Option sets options to config.
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

func WithBefore(
	before ...iconcertothrift.ServerRequestFunc,
) Option {
	return optionFunc(func(c config) config {
		c.before = append(c.before, before...)
		return c
	})
}

func WithAfter(
	after ...iconcertothrift.ServerResponseFunc,
) Option {
	return optionFunc(func(c config) config {
		c.after = append(c.after, after...)
		return c
	})
}

func WithErrorHandler(errorHandler ikittransport.ErrorHandler) Option {
	return optionFunc(func(c config) config {
		c.errorHandler = errorHandler
		return c
	})
}

func WithFinalizer(fn ...iconcertothrift.ServerFinalizerFunc) Option {
	return optionFunc(func(c config) config {
		c.finalizer = append(c.finalizer, fn...)
		return c
	})
}

type server struct {
	registry     *icalculatorthrift.Registry
	before       []iconcertothrift.ServerRequestFunc
	after        []iconcertothrift.ServerResponseFunc
	finalizer    []iconcertothrift.ServerFinalizerFunc
	errorHandler ikittransport.ErrorHandler
	endpoints    *icalculatorendpoint.Set
}

// New makes a set of endpoints available as a Thrift service.
func New(endpoints *icalculatorendpoint.Set, opts ...Option) ithriftgen.CalculatorService {
	cfg := newConfig(opts...)
	return &server{
		registry:     cfg.registry,
		before:       cfg.before,
		after:        cfg.after,
		finalizer:    cfg.finalizer,
		errorHandler: cfg.errorHandler,
		endpoints:    endpoints,
	}
}

func (s *server) Calculate(
	ctx icontext.Context, _op ithriftgen.Op, _num1 int64, _num2 int64,
) (_result int64, err error) {
	req, err := s.registry.DecodeCalculateRequest(_op, _num1, _num2)
	if err != nil {
		return 0, err
	}

	if len(s.finalizer) > 0 {
		defer func() {
			for _, _f := range s.finalizer {
				_f(ctx, err)
			}
		}()
	}

	reqMD := iconcertothrift.Metadata{
		Header: iconcertothrift.NewHeaderContext(ctx),
	}
	for _, _f := range s.before {
		ctx = _f(ctx, reqMD)
	}

	aResp, err := s.endpoints.CalculateEndpoint(
		iconcerto.ContextWithCallMeta(
			ctx,
			iconcerto.CallMeta{
				Service:   icalculator.ServiceDesc.Service,
				Method:    icalculator.ServiceDesc.Methods.Calculate.Name,
				Transport: iconcertotransport.TransportThrift,
			},
		),
		req,
	)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return 0, err
	}

	respMD := iconcertothrift.Metadata{
		Header: iconcertothrift.HeaderMap{},
	}

	for _, fn := range s.after {
		ctx = fn(ctx, &respMD)
	}

	respHeader := respMD.Header.(iconcertothrift.HeaderMap)
	if len(respHeader) > 0 {
		if helper, ok := ithrift.GetResponseHelper(ctx); ok {
			for k, v := range respHeader {
				helper.SetHeader(k, v)
			}
		}
	}

	resp, err := s.registry.EncodeCalculateResponse(
		aResp.(*imessage.CalculateResponse),
	)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return 0, err
	}

	return resp, nil
}
