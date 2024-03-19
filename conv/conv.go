// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package conv

// Ptr returns a pointer to a new memory location that holds the value of v.
func Ptr[T any](v T) *T {
	return &v
}
