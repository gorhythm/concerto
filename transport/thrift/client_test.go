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
	"go.uber.org/mock/gomock"

	concertothrift "github.com/gorhythm/concerto/transport/thrift"
	"github.com/gorhythm/concerto/transport/thrift/mock"
)

func TestExtractResponseMeta(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name   string
		ctx    context.Context
		respMD thrift.ResponseMeta
		want   *concertothrift.Metadata
	}{
		{
			name: "header with value",
			ctx: concertothrift.ContextWithResponseMeta(
				context.Background(), &concertothrift.Metadata{},
			),
			respMD: thrift.ResponseMeta{
				Headers: thrift.THeaderMap{
					"key": "value",
				},
			},
			want: &concertothrift.Metadata{
				Header: concertothrift.HeaderMap{
					"key": "value",
				},
			},
		},
		{
			name: "missing preset response metadata",
			ctx:  context.Background(),
			respMD: thrift.ResponseMeta{
				Headers: thrift.THeaderMap{
					"key": "value",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTClient := mock.NewMockTClient(ctrl)
			mockTClient.
				EXPECT().
				Call(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).
				Return(
					thrift.ResponseMeta{
						Headers: thrift.THeaderMap{
							"key": "value",
						},
					},
					nil,
				).
				AnyTimes()
			tclient := concertothrift.ExtractResponseMeta(mockTClient)
			_, err := tclient.(thrift.WrappedTClient).Wrapped(tt.ctx, "", nil, nil)
			if err != nil {
				t.Fatalf("test %q failed. Didn't expect an error, got %q", tt.name, err)
			}

			got, want := concertothrift.ResponseMetaFromContext(tt.ctx), tt.want
			if !reflect.DeepEqual(got, want) {
				t.Errorf("test %q failed. Got %+v, want %+v", tt.name, got, want)
			}
		})
	}
}
