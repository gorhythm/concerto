// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package callmeta

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/gorhythm/concerto"
)

// Middleware returns a [github.com/go-kit/kit/endpoint.Middleware] that adds
// the call meta to endpoint context.
func Middleware(m concerto.CallMeta) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			return next(concerto.ContextWithCallMeta(ctx, m), req)
		}
	}
}
