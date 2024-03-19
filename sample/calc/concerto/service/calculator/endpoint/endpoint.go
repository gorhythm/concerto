// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package endpoint

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"

	"github.com/gorhythm/concerto/endpoint"

	"github.com/gorhythm/concerto/sample/calc"
	"github.com/gorhythm/concerto/sample/calc/concerto/message"
)

// Set collects all of the endpoints that compose a CalculatorService. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Set struct {
	CalculateEndpoint kitendpoint.Endpoint
}

var _ calc.CalculatorService = &Set{}

// New constructs and returns a new instance of [Set].
func New(
	svc calc.CalculatorService, opts ...endpoint.Option,
) *Set {
	var (
		cfg       = endpoint.NewConfig(opts...)
		endpoints = Set{}
	)
	{
		e := func(
			ctx context.Context, req any,
		) (any, error) {
			_req := req.(*message.CalculateRequest)
			_result, err := svc.Calculate(ctx, _req.Op, _req.Num1, _req.Num2)
			if err != nil {
				return nil, err
			}

			return &message.CalculateResponse{Result: _result}, err
		}
		endpoints.CalculateEndpoint = cfg.ApplyMiddlewares(e)
	}

	return &endpoints
}

// Calculate implements the CalculatorService interface, so Set may be used as
// a service. This is primarily useful in the context of a client library.
func (e *Set) Calculate(
	ctx context.Context, op calc.Op, num1 int64, num2 int64,
) (int64, error) {
	resp, err := e.CalculateEndpoint(
		ctx,
		&message.CalculateRequest{
			Op:   op,
			Num1: num1,
			Num2: num2,
		},
	)

	if err != nil {
		return 0, err
	}

	return resp.(*message.CalculateResponse).Result, err
}
