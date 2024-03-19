// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package calc

import (
	"fmt"
	"strconv"
)

// Enum value maps for Op.
var (
	_opName = map[Op]string{
		OpAdd:      "Add",
		OpSubtract: "SUBTRACT",
		OpMultiply: "MULTIPLY",
		OpDivide:   "DIVIDE",
	}
	_opValue = map[string]Op{
		"Add":      OpAdd,
		"SUBTRACT": OpSubtract,
		"MULTIPLY": OpMultiply,
		"DIVIDE":   OpDivide,
	}
)

func (x Op) String() string {
	if str, ok := _opName[x]; ok {
		return str
	}

	return "%!Op(" + strconv.FormatInt(int64(x), 10) + ")"
}

func OpFromString(s string) (Op, error) {
	if v, ok := _opValue[s]; ok {
		return v, nil
	}

	return Op(0), fmt.Errorf("not a valid Op string")
}

func (x Op) MarshalText() ([]byte, error) {
	if str, ok := _opName[x]; ok {
		return []byte(str), nil
	}

	return nil, fmt.Errorf("not a valid Op")
}

func (x *Op) UnmarshalText(text []byte) error {
	v, err := OpFromString(string(text))
	if err != nil {
		return err
	}
	*x = v
	return nil
}

func (x Op) MarshalJSON() ([]byte, error) {
	if str, ok := _opName[x]; ok {
		return []byte(strconv.Quote(str)), nil
	}

	return nil, fmt.Errorf("not a valid Op")
}

func (x *Op) UnmarshalJSON(text []byte) error {
	str, err := strconv.Unquote(string(text))
	if err != nil {
		return err
	}

	v, err := OpFromString(str)
	if err != nil {
		return err
	}
	*x = v
	return nil
}
