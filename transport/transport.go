// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package transport

import (
	"errors"
	"strconv"
)

// Transport represents the type of transport.
type Transport int

// Enum of transport types.
const (
	TransportUnspecified Transport = iota
	TransportHTTP
	TransportJSONRPCOverHTTP
	TransportGRPC
	TransportThrift
)

// Enum value maps for Transport.
var (
	// transportText maps Transport values to their string representations.
	transportText = map[Transport]string{
		TransportHTTP:            "http",
		TransportJSONRPCOverHTTP: "jsonrpc",
		TransportGRPC:            "grpc",
		TransportThrift:          "thrift",
	}
	// transportValue maps string representations to their Transport values.
	transportValue = map[string]Transport{
		"http":    TransportHTTP,
		"jsonrpc": TransportJSONRPCOverHTTP,
		"grpc":    TransportGRPC,
		"thrift":  TransportThrift,
	}
)

// String returns the string representation of the Transport value.
func (t Transport) String() string {
	if str, ok := transportText[t]; ok {
		return str
	}
	return "%!Transport(" + strconv.FormatInt(int64(t), 10) + ")"
}

// IsValid reports whether t is a valid [Transport].
func (t Transport) IsValid() bool {
	_, found := transportText[t]
	return found
}

// FromString converts a string to a Transport value.
func FromString(s string) (Transport, error) {
	if v, ok := transportValue[s]; ok {
		return v, nil
	}

	return Transport(0), errors.New("not a valid Transport string")
}

// MarshalText converts a Transport value to text.
func (t Transport) MarshalText() ([]byte, error) {
	if str, ok := transportText[t]; ok {
		return []byte(str), nil
	}

	return nil, errors.New("not a valid Transport")
}

// UnmarshalText converts text to a Transport value.
func (t *Transport) UnmarshalText(text []byte) error {
	v, err := FromString(string(text))
	if err != nil {
		return err
	}
	*t = v
	return nil
}

// MarshalJSON converts a Transport value to JSON.
func (t Transport) MarshalJSON() ([]byte, error) {
	if str, ok := transportText[t]; ok {
		return []byte(strconv.Quote(str)), nil
	}

	return nil, errors.New("not a valid Transport")
}

// UnmarshalJSON converts JSON to a Transport value.
func (t *Transport) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	v, err := FromString(str)
	if err != nil {
		return err
	}
	*t = v

	return nil
}
