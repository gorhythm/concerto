// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package transport_test

import (
	"encoding/json"
	"testing"

	"github.com/gorhythm/concerto/transport"
)

func TestTransport_IsValid(t *testing.T) {
	if got, want := transport.TransportGRPC.IsValid(), true; got != want {
		t.Errorf(`test "valid" failed. Got %v, want %v`, got, want)
	}

	if got, want := transport.Transport(100).IsValid(), false; got != want {
		t.Errorf(`test "invalid" failed. Got %v, want %v`, got, want)
	}
}

func TestTransport_String(t *testing.T) {
	if got, want := transport.TransportGRPC.String(), "grpc"; got != want {
		t.Errorf(`test "TransportGRPC" failed. Got %q, want %q`, got, want)
	}

	if got, want := transport.Transport(100).String(), "%!Transport(100)"; got != want {
		t.Errorf(`test "Transport 100" failed. Got, %q want %q`, got, want)
	}
}

func TestTransportFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    transport.Transport
		wantErr string
	}{
		{
			name: "valid",
			args: args{
				s: "thrift",
			},
			want: transport.TransportThrift,
		},
		{
			name: "invalid",
			args: args{
				s: "unkown-transport",
			},
			wantErr: "not a valid Transport string",
		},
		{
			name: "empty string",
			args: args{
				s: "",
			},
			wantErr: "not a valid Transport string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transport.TransportFromString(tt.args.s)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Errorf(
						"test %q failed. Got error %q, want %q",
						tt.name, err, tt.wantErr,
					)
				}
				return
			}
			if err != nil {
				t.Fatalf(
					"test %q failed. Didn't expect an error, got %q",
					tt.name, err,
				)
			}

			if got != tt.want {
				t.Errorf("test %q failed. Got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTransport_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		trans   transport.Transport
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			trans:   transport.TransportGRPC,
			want:    "grpc",
			wantErr: false,
		},
		{
			name:    "invalid",
			trans:   transport.Transport(1000),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.trans.MarshalText()
			if tt.wantErr {
				if err == nil {
					t.Errorf("test %q failed. Expected an error", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf(
					"test %q failed. Didn't expect an error, got %q",
					tt.name, err,
				)
			}

			if got, want := "grpc", string(got); got != want {
				t.Errorf("test %q failed. Got %v, want %v", tt.name, got, want)
			}
		})
	}
}

func TestTransport_UnmarshalText(t *testing.T) {
	type args struct {
		text []byte
	}
	tests := []struct {
		name    string
		args    args
		want    transport.Transport
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				text: []byte("grpc"),
			},
			want:    transport.TransportGRPC,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				text: []byte("invalid"),
			},
			want:    transport.TransportUnspecified,
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				text: []byte(""),
			},
			want:    transport.TransportUnspecified,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var trans transport.Transport
			err := trans.UnmarshalText(tt.args.text)
			if tt.wantErr {
				if err == nil {
					t.Errorf("test %q failed. Expected an error", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf(
					"test %q failed. Didn't expect an error, got %q",
					tt.name, err,
				)
			}

			if got := trans; got != tt.want {
				t.Errorf("test %q failed. Got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTransport_MarshalJSON(t *testing.T) {
	type args struct {
		data map[string]transport.Transport
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				data: map[string]transport.Transport{
					"transport": transport.TransportGRPC,
				},
			},
			want:    `{"transport":"grpc"}`,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				data: map[string]transport.Transport{
					"transport": transport.Transport(100),
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args.data)
			if tt.wantErr {
				if err == nil {
					t.Errorf("test %q failed. Expected an error", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf(
					"test %q failed. Didn't expect an error, got %q",
					tt.name, err,
				)
			}

			if got := string(got); got != tt.want {
				t.Errorf("test %q failed. Got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTransport_UnmarshalJSON(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    transport.Transport
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				bytes: []byte(`{"transport": "grpc"}`),
			},
			want:    transport.TransportGRPC,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				bytes: []byte(`{"transport": "invalid"}`),
			},
			want:    transport.TransportUnspecified,
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				bytes: []byte(`{"transport": ""}`),
			},
			want:    transport.TransportUnspecified,
			wantErr: true,
		},
		{
			name: "null",
			args: args{
				bytes: []byte(`{"transport": null}`),
			},
			want:    transport.TransportUnspecified,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMap map[string]transport.Transport
			err := json.Unmarshal([]byte(tt.args.bytes), &gotMap)
			if tt.wantErr {
				if err == nil {
					t.Errorf("test %q failed. Expected an error", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf(
					"test %q failed. Didn't expect an error, got %q",
					tt.name, err,
				)
			}

			if got := gotMap["transport"]; got != tt.want {
				t.Errorf("test %q failed. Got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
