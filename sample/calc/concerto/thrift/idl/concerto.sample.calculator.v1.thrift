// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

include "concerto.rpc.status.thrift"

service CalculatorService {
  i64 calculate(1: Op op, 2: i64 num1, 3: i64 num2) throws (1: concerto.rpc.status.Error e1)
}

// Op represents the type of arithmetic operation.
enum Op {
  // Add represents addition operation.
  ADD = 0,
  // Subtract represents subtraction operation.
  SUBTRACT = 1,
  // Multiply represents multiplication operation.
  MULTIPLY = 2,
  // Divide represents division operation.
  DIVIDE = 3,
}
