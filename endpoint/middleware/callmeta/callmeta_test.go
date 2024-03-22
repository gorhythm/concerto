// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package callmeta_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/gorhythm/concerto"
	"github.com/gorhythm/concerto/endpoint/middleware/callmeta"
	"github.com/gorhythm/concerto/transport"
)

func TestMiddleware(t *testing.T) {
	var (
		args = struct {
			callMeta concerto.CallMeta
		}{
			callMeta: concerto.CallMeta{
				Service:   "concerto.test.v1.TestService",
				Method:    "Ping",
				Transport: transport.TransportGRPC,
			},
		}
		want = args.callMeta
		got  concerto.CallMeta
	)

	_, err := callmeta.Middleware(args.callMeta)(
		func(ctx context.Context, _ any) (any, error) {
			got = concerto.CallMetaFromContext(ctx)
			return struct{}{}, nil
		},
	)(context.Background(), struct{}{})
	if err != nil {
		t.Fatalf("test failed. Didn't expect an error, got %q", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("test failed. Got %+v, want %+v", got, want)
	}
}
