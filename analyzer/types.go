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

import (
	"go/token"
	"go/types"
)

// errorIsInterface holds a reference to the `interface{ Is(error) bool }` type.
// This is used to check if a type implements the optional error comparison interface.
//
//nolint:gochecknoglobals
var (
	// `interface{ Is(error) bool }`.
	errorIsInterface = newErrorIsInterface()

	// `interface{ Unwrap() error }`.
	errorUnwrapInterface = newErrorUnwrapInterface()

	// `interface{ Unwrap() []error }`.
	errorUnwrapArrayInterface = newErrorUnwrapArrayInterface()
)

// newErrorIsInterface constructs and returns a new [types.Interface] representing
// the `interface{ Is(error) bool }` type.
func newErrorIsInterface() *types.Interface {
	const isMethodName = "Is"

	var noPkg *types.Package

	params := singleVar(errorType())
	results := singleVar(types.Typ[types.Bool])
	sig := types.NewSignatureType(nil, nil, nil, params, results, false)
	isFunc := types.NewFunc(token.NoPos, noPkg, isMethodName, sig)

	return types.NewInterfaceType([]*types.Func{isFunc}, nil).Complete()
}

// newErrorUnwrapInterface constructs and returns a new [types.Interface] representing
// the `interface{ Unwrap() error }` type.
func newErrorUnwrapInterface() *types.Interface {
	const unwrapMethodName = "Unwrap"

	var noPkg *types.Package

	results := singleVar(errorType())
	sig := types.NewSignatureType(nil, nil, nil, nil, results, false)
	unwrapFunc := types.NewFunc(token.NoPos, noPkg, unwrapMethodName, sig)

	return types.NewInterfaceType([]*types.Func{unwrapFunc}, nil).Complete()
}

// newErrorUnwrapArrayInterface constructs and returns a new [types.Interface] representing
// the `interface{ Unwrap() []error }` type.
func newErrorUnwrapArrayInterface() *types.Interface {
	const unwrapMethodName = "Unwrap"

	var noPkg *types.Package

	results := singleVar(types.NewSlice(errorType()))
	sig := types.NewSignatureType(nil, nil, nil, nil, results, false)
	unwrapFunc := types.NewFunc(token.NoPos, noPkg, unwrapMethodName, sig)

	return types.NewInterfaceType([]*types.Func{unwrapFunc}, nil).Complete()
}

func errorType() types.Type {
	const errorTypName = "error"

	return types.Universe.Lookup(errorTypName).Type()
}

// singleVar constructs a [types.Tuple] containing a single unnamed variable
// of the specified [types.Type]. This is a helper for creating function signatures.
func singleVar(t types.Type) *types.Tuple {
	var noPkg *types.Package

	return types.NewTuple(types.NewVar(token.NoPos, noPkg, "", t))
}
