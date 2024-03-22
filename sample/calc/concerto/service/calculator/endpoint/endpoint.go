// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package endpoint

import (
	icontext "context"

	ikitendpoint "github.com/go-kit/kit/endpoint"

	iconcertoendpoint "github.com/gorhythm/concerto/endpoint"

	icalc "github.com/gorhythm/concerto/sample/calc"
	imessage "github.com/gorhythm/concerto/sample/calc/concerto/message"
)

// Set collects all of the endpoints that compose a CalculatorService. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Set struct {
	CalculateEndpoint ikitendpoint.Endpoint
}

var _ icalc.CalculatorService = &Set{}

// New constructs and returns a new instance of [Set].
func New(
	svc icalc.CalculatorService, opts ...iconcertoendpoint.Option,
) *Set {
	var (
		cfg = iconcertoendpoint.NewConfig(opts...)
		set = Set{}
	)
	{
		e := func(
			ctx icontext.Context, aReq any,
		) (any, error) {
			req := aReq.(*imessage.CalculateRequest)
			result, err := svc.Calculate(ctx, req.Op, req.Num1, req.Num2)
			if err != nil {
				return nil, err
			}

			return &imessage.CalculateResponse{Result: result}, err
		}
		if middleware := cfg.Middleware(); middleware != nil {
			e = middleware(e)
		}
		set.CalculateEndpoint = e
	}

	return &set
}

// Calculate implements the CalculatorService interface, so Set may be used as
// a service. This is primarily useful in the icontext of a client library.
func (s *Set) Calculate(
	ctx icontext.Context, _op icalc.Op, _num1 int64, _num2 int64,
) (int64, error) {
	aResp, err := s.CalculateEndpoint(
		ctx,
		&imessage.CalculateRequest{
			Op:   _op,
			Num1: _num1,
			Num2: _num2,
		},
	)

	if err != nil {
		return 0, err
	}

	return aResp.(*imessage.CalculateResponse).Result, err
}
