// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package thrift_test

import (
	"context"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"

	concertothrift "github.com/gorhythm/concerto/transport/thrift"
)

func TestSetRequestHeader(t *testing.T) {
	args := struct {
		key string
		val string
	}{
		key: "key",
		val: "value",
	}
	hm := thrift.THeaderMap{}
	concertothrift.SetRequestHeader(args.key, args.val)(context.Background(), hm)
	if got, want := hm[args.key], args.val; got != want {
		t.Errorf("test failed. Got %v, want %v", got, want)
	}
}

func TestSetRequestContext(t *testing.T) {
	type args struct {
		ctxKey    any
		headerKey string
	}

	tests := []struct {
		name string
		args args
		ctx  context.Context
		want any
	}{
		{
			name: "found",
			args: args{
				ctxKey:    "ctxKey",
				headerKey: "headerKey",
			},
			ctx: thrift.SetHeader(
				context.Background(),
				"headerKey", "value",
			),
			want: "value",
		},
		{
			name: "not found",
			args: args{
				ctxKey:    "ctxKey",
				headerKey: "headerKey",
			},
			ctx:  context.Background(),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := concertothrift.SetRequestContext(tt.args.ctxKey, tt.args.headerKey)(tt.ctx)
			if got := ctx.Value(tt.args.ctxKey); got != tt.want {
				t.Errorf("test %q failed. Got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
