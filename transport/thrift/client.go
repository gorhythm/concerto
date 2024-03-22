// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// The MIT License (MIT)
// Copyright (c) 2015 Peter Bourgon.

package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

// ClientFinalizerFunc can be used to perform work at the end of a client Thrift
// request, after the response is returned. The principal intended use is for
// error logging. Additional response parameters are provided in the context
// under keys with the ContextKeyResponse prefix.
// NOTE: err may be nil. There maybe also no additional response parameters
// depending on when an error occurs.
type ClientFinalizerFunc func(ctx context.Context, err error)

// ExtractResponseMeta is a Thrift middleware that retrieves the Thrift response metadata
// and injects it into the context. This is useful for passing along metadata from the
// Thrift response to other parts of the application that have access to the context.
func ExtractResponseMeta(next thrift.TClient) thrift.TClient {
	return thrift.WrappedTClient{
		Wrapped: func(
			ctx context.Context, method string, args, result thrift.TStruct,
		) (thrift.ResponseMeta, error) {
			tRespMeta, err := next.Call(ctx, method, args, result)
			if err != nil {
				return tRespMeta, err
			}

			// Check if the context has a Metadata value and if so, update its Header field
			// with the headers from the Thrift response metadata.
			if cRespMeta, ok := ctx.Value(contextKeyResponseMeta).(*Metadata); ok && cRespMeta != nil {
				cRespMeta.Header = HeaderMap(tRespMeta.Headers)
			}

			return tRespMeta, nil
		},
	}
}
