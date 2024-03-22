// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// The MIT License (MIT)
// Copyright (c) 2015 Peter Bourgon.

package thrift_test

import (
	"context"
	"reflect"
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
	want := args.val

	md := concertothrift.Metadata{
		Header: concertothrift.HeaderMap{},
	}
	concertothrift.SetRequestHeader(args.key, args.val)(context.Background(), &md)
	got, _ := md.Header.Get(args.key)
	if got != want {
		t.Errorf("test failed. Got %q, want %q", got, want)
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
		md   concertothrift.Metadata
		want any
	}{
		{
			name: "found",
			args: args{
				ctxKey:    "ctxKey",
				headerKey: "headerKey",
			},
			md: concertothrift.Metadata{
				Header: concertothrift.HeaderMap{
					"headerKey": "value",
				},
			},
			want: "value",
		},
		{
			name: "not found",
			args: args{
				ctxKey:    "ctxKey",
				headerKey: "headerKey",
			},
			md: concertothrift.Metadata{
				Header: concertothrift.HeaderMap{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := concertothrift.SetRequestContext(tt.args.ctxKey, tt.args.headerKey)(
				context.Background(), tt.md,
			)
			if got := ctx.Value(tt.args.ctxKey); got != tt.want {
				t.Errorf("test %q failed. Got %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestHeaderMap_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		hm        concertothrift.HeaderMap
		args      args
		wantValue string
		wantOk    bool
	}{
		{
			hm: concertothrift.HeaderMap{
				"key": "value",
			},
			args: args{
				key: "key",
			},
			wantValue: "value",
			wantOk:    true,
		},
		{
			hm: concertothrift.HeaderMap{
				"key": "",
			},
			args: args{
				key: "key",
			},
			wantValue: "",
			wantOk:    true,
		},
		{
			hm: concertothrift.HeaderMap{
				"key": "value",
			},
			args: args{
				key: "other-key",
			},
			wantValue: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		gotValue, gotOk := tt.hm.Get(tt.args.key)
		if gotValue != tt.wantValue {
			t.Errorf("test failed. GotValue = %v, want %v", gotValue, tt.wantValue)
		}
		if gotOk != tt.wantOk {
			t.Errorf("test failed. GotOk = %v, want %v", gotOk, tt.wantOk)
		}
	}
}

func TestHeaderContext_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		kvs       map[string]string
		args      args
		wantValue string
		wantOk    bool
	}{
		{
			kvs: map[string]string{
				"key": "value",
			},
			args: args{
				key: "key",
			},
			wantValue: "value",
			wantOk:    true,
		},
		{
			kvs: map[string]string{
				"key": "",
			},
			args: args{
				key: "key",
			},
			wantValue: "",
			wantOk:    true,
		},
		{
			kvs: map[string]string{
				"key": "value",
			},
			args: args{
				key: "other-key",
			},
			wantValue: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		var (
			ctx = context.Background()
			ks  []string
		)
		for k, v := range tt.kvs {
			ctx = thrift.SetHeader(ctx, k, v)
			ks = append(ks, k)
		}
		ctx = thrift.SetReadHeaderList(ctx, ks)

		hc := concertothrift.NewHeaderContext(ctx)

		gotValue, gotOk := hc.Get(tt.args.key)
		if gotValue != tt.wantValue {
			t.Errorf("test failed. GotValue = %v, want %v", gotValue, tt.wantValue)
		}
		if gotOk != tt.wantOk {
			t.Errorf("test failed. GotOk = %v, want %v", gotOk, tt.wantOk)
		}
	}
}

func TestHeaderMap_Keys(t *testing.T) {
	hm := concertothrift.HeaderMap{
		"key1": "value1",
		"key2": "values",
	}
	want := []string{"key1", "key2"}

	if got := hm.Keys(); !reflect.DeepEqual(got, want) {
		t.Errorf("test failed. Got %v, want %v", got, want)
	}
}

func TestHeaderContext_Keys(t *testing.T) {
	ctx := thrift.SetReadHeaderList(context.Background(), []string{"key1", "key2"})
	want := []string{"key1", "key2"}

	hc := concertothrift.NewHeaderContext(ctx)

	if got := hc.Keys(); !reflect.DeepEqual(got, want) {
		t.Errorf("test failed. Got %v, want %v", got, want)
	}
}

func TestHeaderMap_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		hm   concertothrift.HeaderMap
		args args
		want concertothrift.HeaderMap
	}{
		{
			hm: concertothrift.HeaderMap{},
			args: args{
				key:   "key",
				value: "value",
			},
			want: concertothrift.HeaderMap{
				"key": "value",
			},
		},
		{
			hm: concertothrift.HeaderMap{
				"key": "value",
			},
			args: args{
				key:   "other-key",
				value: "value",
			},
			want: concertothrift.HeaderMap{
				"key":       "value",
				"other-key": "value",
			},
		},
		{
			hm: concertothrift.HeaderMap{
				"key": "value",
			},
			args: args{
				key:   "key",
				value: "override-value",
			},
			want: concertothrift.HeaderMap{
				"key": "override-value",
			},
		},
	}
	for _, tt := range tests {
		tt.hm.Set(tt.args.key, tt.args.value)
		if got, want := tt.hm, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("test failed. Got %v, want %v", got, want)
		}
	}
}
