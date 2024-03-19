// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package endpoint_test

import (
	"context"
	"fmt"

	kitendpoint "github.com/go-kit/kit/endpoint"

	"github.com/gorhythm/concerto/endpoint"
)

func ExampleWithMiddlewares() {
	cfg := endpoint.NewConfig(
		endpoint.WithMiddlewares(
			annotate("first"),
			annotate("second"),
			annotate("third"),
		),
	)

	e := cfg.ApplyMiddlewares(myEndpoint)
	if _, err := e(context.Background(), struct{}{}); err != nil {
		panic(err)
	}

	// Output:
	// first pre
	// second pre
	// third pre
	// my endpoint!
	// third post
	// second post
	// first post
}

func annotate(s string) kitendpoint.Middleware {
	return func(next kitendpoint.Endpoint) kitendpoint.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			fmt.Println(s, "pre")
			defer fmt.Println(s, "post")
			return next(ctx, req)
		}
	}
}

func myEndpoint(context.Context, any) (any, error) {
	fmt.Println("my endpoint!")
	return struct{}{}, nil
}
