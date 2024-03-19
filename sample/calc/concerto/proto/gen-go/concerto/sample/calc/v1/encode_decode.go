// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package calc

import (
	"github.com/gorhythm/concerto/sample/calc"
	"github.com/gorhythm/concerto/sample/calc/concerto/message"
)

func EncodeCalculateRequest(in *message.CalculateRequest) (*CalculateRequest, error) {
	return &CalculateRequest{
		Op:   Op(in.Op),
		Num1: in.Num1,
		Num2: in.Num2,
	}, nil
}

func DecodeCalculateRequest(in *CalculateRequest) (*message.CalculateRequest, error) {
	return &message.CalculateRequest{
		Op:   calc.Op(in.Op),
		Num1: in.Num1,
		Num2: in.Num2,
	}, nil
}

func EncodeCalculateResponse(in *message.CalculateResponse) (*CalculateResponse, error) {
	return &CalculateResponse{
		Result: in.Result,
	}, nil
}

func DecodeCalculateResponse(in *CalculateResponse) (*message.CalculateResponse, error) {
	return &message.CalculateResponse{
		Result: in.Result,
	}, nil
}
