// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package v1

import (
	icalc "github.com/gorhythm/concerto/sample/calc"
	"github.com/gorhythm/concerto/sample/calc/concerto/message"
)

func EncodeCalculateRequest(req *message.CalculateRequest) (_op Op, _num1 int64, _num2 int64, err error) {
	_op = Op(req.Op)
	_num1 = req.Num1
	_num2 = req.Num2

	return
}

func DecodeCalculateRequest(_op Op, _num1 int64, _num2 int64) (*message.CalculateRequest, error) {
	return &message.CalculateRequest{
		Op:   icalc.Op(_op),
		Num1: _num1,
		Num2: _num2,
	}, nil
}

func EncodeCalculateResponse(resp *message.CalculateResponse) (int64, error) {
	return resp.Result, nil
}

func DecodeCalculateResponse(_result int64) (*message.CalculateResponse, error) {
	return &message.CalculateResponse{
		Result: _result,
	}, nil
}
