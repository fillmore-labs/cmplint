// Copyright 2025 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package analyzer

import "go/types"

// isZeroSized determines recursively whether the underlying type of t is provably zero-sized.
// It checks arrays of length 0 and structs where all fields are zero-sized.
func isZeroSized(t types.Type) bool {
	switch u := t.Underlying().(type) {
	case *types.Array:
		// An array is zero-sized if its length is 0 or its element type is zero-sized.
		// The former (len 0) is the primary case for zero-sized arrays.
		return u.Len() == 0 || isZeroSized(u.Elem())

	case *types.Struct:
		// A struct is zero-sized if all its fields are zero-sized.
		for i := range u.NumFields() {
			if !isZeroSized(u.Field(i).Type()) {
				return false
			}
		}

		// All fields are zero-sized, so the struct is zero-sized.
		return true

	default:
		// Other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
		// are not zero-sized.
		return false
	}
}
