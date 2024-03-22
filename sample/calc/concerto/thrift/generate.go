// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package thrift

//go:generate thrift  -I "../../../../types/thrift" -r --gen "go:package_prefix=github.com/gorhythm/concerto/sample/calc/concerto/thrift/gen-go/" idl/concerto.sample.calculator.v1.thrift
