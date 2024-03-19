// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package calculator

import (
	concerto "github.com/gorhythm/concerto"
)

// ServiceDesc is the description for CalculatorService service.
var ServiceDesc = struct {
	Service string
	Methods struct {
		Calculate concerto.MethodInfo
	}
}{
	Service: "concerto.sample.calc.v1.CalculatorService",
	Methods: struct {
		Calculate concerto.MethodInfo
	}{
		Calculate: concerto.MethodInfo{
			Name:     "Calculate",
			FullName: "concerto.sample.calc.v1.CalculatorService/Calculate",
		},
	},
}
