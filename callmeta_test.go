// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Copyright (c) The go-grpc-middleware Authors.
// Licensed under the Apache License 2.0.

package concerto_test

import (
	"context"
	"testing"

	"github.com/gorhythm/concerto"
	"github.com/gorhythm/concerto/transport"
)

func TestCallMeta_FullMethod(t *testing.T) {
	callMeta := concerto.CallMeta{
		Service:   "concerto.test.v1.TestService",
		Method:    "Ping",
		Transport: transport.TransportGRPC,
	}
	want := "concerto.test.v1.TestService/Ping"
	got := callMeta.FullMethod()

	if got != want {
		t.Errorf("test failed. Got %q, want %q", got, want)
	}
}

func TestCallMetaFromContext(t *testing.T) {
	callMeta := concerto.CallMeta{
		Service:   "concerto.test.v1.TestService",
		Method:    "Ping",
		Transport: transport.TransportGRPC,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want concerto.CallMeta
	}{
		{
			name: "found",
			args: args{
				ctx: concerto.ContextWithCallMeta(context.Background(), callMeta),
			},
			want: callMeta,
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
			},
			want: concerto.NilCallMeta,
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
			},
			want: concerto.NilCallMeta,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := concerto.CallMetaFromContext(tt.args.ctx); got != tt.want {
				t.Errorf(
					"test %q failed. Got %+v, want %+v",
					tt.name, got, tt.want,
				)
			}
		})
	}
}

func TestContextWithCallMeta(t *testing.T) {
	callMeta := concerto.CallMeta{
		Service:   "concerto.test.v1.TestService",
		Method:    "Ping",
		Transport: transport.TransportGRPC,
	}
	want := callMeta
	got := concerto.CallMetaFromContext(
		concerto.ContextWithCallMeta(context.Background(), callMeta),
	)

	if got != want {
		t.Errorf("test failed. Got %+v, want %+v", got, want)
	}
}
