// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package calc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gorhythm/concerto"
)

//concerto:method grpc:"Calculator"
//concerto:method thrift:"Calculator"
//concerto:service http:"/v1/calculator"
type CalculatorService interface {
	//concerto:method grpc:"Calculate"
	//concerto:method thrift:"Calculate,exceptions=1:concerto.rpc.status.Error"
	//concerto:method http:"POST /calculate,additional_routes=POST calculate/{op}/{num1}/{num2}"
	Calculate(
		ctx context.Context,
		//concerto:field proto:"1" thrift:"1" json:"op"
		op Op,
		//concerto:field proto:"2,jstype = JS_STRING" thrift:"2" json:"num1"
		num1 int64,
		//concerto:field proto:"3,jstype = JS_STRING" thrift:"3" json:"num2"
		num2 int64,
	) (
		//concerto:field proto:"1,jstype = JS_STRING" thrift:"1" json:"result"
		result int64,
		err error,
	)
}

// Op represents the type of arithmetic operation.
//
//concerto:enum textcase:"camel"
type Op int32

const (
	//concerto:value text:"Add"
	OpAdd      Op = iota // Add represents addition operation.
	OpSubtract           // Subtract represents subtraction operation.
	OpMultiply           // Multiply represents multiplication operation.
	OpDivide             // Divide represents division operation.
)

type calculatorService struct {
	logger *slog.Logger
}

var _ CalculatorService = (*calculatorService)(nil)

// NewCalculatorService constructs and returns an implementation of
// CalculatorService.
func NewCalculatorService(
	logger *slog.Logger,
) CalculatorService {
	return &calculatorService{
		logger: logger,
	}
}

// Calculate implements CalculatorService.
func (s *calculatorService) Calculate(
	ctx context.Context, op Op, num1 int64, num2 int64,
) (int64, error) {
	s.logger.Info(
		fmt.Sprintf("calculate(%s, %d, %d)", op, num1, num2),
		"transport", concerto.CallMetaFromContext(ctx).Transport,
	)
	switch op {
	case OpAdd:
		return num1 + num2, nil
	case OpSubtract:
		return num1 - num2, nil
	case OpMultiply:
		return num1 * num2, nil
	case OpDivide:
		if num2 == 0 {
			return 0, errors.New("division by zero is undefined")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", op)
	}
}
