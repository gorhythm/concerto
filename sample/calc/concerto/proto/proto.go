// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package proto

//go:generate protoc --proto_path=idl --go_out=. --go-grpc_out=. --go_opt=module=github.com/gorhythm/concerto/sample/calc/concerto/proto --go-grpc_opt=module=github.com/gorhythm/concerto/sample/calc/concerto/proto concerto/sample/calc/v1/calc_v1.proto
