// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package v1

import (
	"github.com/gorhythm/concerto/sample/calc"
	"github.com/gorhythm/concerto/sample/calc/concerto/message"
)

func EncodeCalculateRequest(_obj *message.CalculateRequest) (op Op, num1 int64, num2 int64, err error) {
	op = Op(_obj.Op)
	num1 = _obj.Num1
	num2 = _obj.Num2

	return
}

func DecodeCalculateRequest(op Op, num1 int64, num2 int64) (*message.CalculateRequest, error) {
	return &message.CalculateRequest{
		Op:   calc.Op(op),
		Num1: num1,
		Num2: num2,
	}, nil
}

func EncodeCalculateResponse(_obj *message.CalculateResponse) (int64, error) {
	return _obj.Result, nil
}

func DecodeCalculateResponse(result int64) (*message.CalculateResponse, error) {
	return &message.CalculateResponse{
		Result: result,
	}, nil
}
