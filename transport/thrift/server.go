// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// The MIT License (MIT)
// Copyright (c) 2015 Peter Bourgon.

package thrift

import "context"

// ServerFinalizerFunc can be used to perform work at the end of an gRPC
// request, after the response has been written to the client.
type ServerFinalizerFunc func(ctx context.Context, err error)
