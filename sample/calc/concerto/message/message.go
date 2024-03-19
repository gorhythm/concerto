// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package message

import "github.com/gorhythm/concerto/sample/calc"

type CalculateRequest struct {
	Op   calc.Op
	Num1 int64
	Num2 int64
}

type CalculateResponse struct {
	Result int64
}
