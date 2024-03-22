// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package thrift

//go:generate mockgen -destination=./mock/thift.go -copyright_file=./copyright.txt -package=mock github.com/apache/thrift/lib/go/thrift TClient
