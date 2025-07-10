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
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strings"
)

// comparison analyzes a comparison operation (either binary like `==` or
// a function call like `errors.Is`) to determine if one of the operands is
// the address of a composite literal or a new() call.
//
// It reports a diagnostic if such a comparison is found, providing additional context
// if the comparison involves zero-sized types.
func (p pass) comparison(n ast.Node, left, right ast.Expr, isError bool) {
	var (
		t      types.Type // The type of T in a &T{} or new(T) operand
		isLeft bool       // operand detected is on the left side of the comparison
		other  ast.Expr   // isLeft ? right : left
	)

	// Determine if one of the operands is a new literal (&T{} or new(T)), check the left first.
	if tl, ok := p.isAddrOfCompLitOrNew(left); ok {
		t, other, isLeft = tl, right, true
	} else if tr, ok := p.isAddrOfCompLitOrNew(right); ok {
		t, other, isLeft = tr, left, false
	} else {
		return
	}

	// The `isLeft` flag is used by `shouldSuppressDiagnostic` to consider `Unwrap` methods
	// if the new literal is the first argument in an error comparison (`errors.Is(&T{}, target)`).
	if t != nil && isError && p.checkis && shouldSuppressDiagnostic(t, isLeft) {
		return
	}

	// Determine if the comparison is with a zero-sized type and the other operand is not nil.
	// In this case, the result is undefined.
	isUndefined := false

	if t != nil && IsZeroSized(t) {
		otherType, ok := p.TypesInfo.Types[other]
		isUndefined = !ok || !otherType.IsNil()
	}

	// Report diagnostic
	typeName := "invalid type"
	if t != nil {
		typeName = types.TypeString(t, types.RelativeTo(p.Pkg))
	}

	otherStr := p.exprToString(other)

	if isUndefined {
		p.ReportRangef(n,
			"Result of comparison of %q with address of new zero-sized variable of type %q is false or undefined",
			otherStr, typeName)
	} else {
		p.ReportRangef(n,
			"Result of comparison of %q with address of new variable of type %q is always false",
			otherStr, typeName)
	}
}

// exprToString converts an AST expression to its string representation.
func (p pass) exprToString(e ast.Expr) string {
	var s strings.Builder
	if err := format.Node(&s, p.Fset, e); err != nil {
		return "invalid"
	}

	return s.String()
}

// isAddrOfCompLitOrNew checks if the given AST expression `x` represents
// the address of a composite literal (`&T{...}`) or a call to the built-in
// `new()` function (`new(T)`).
// It returns the element type `T` of the resulting pointer.
func (p pass) isAddrOfCompLitOrNew(x ast.Expr) (typ types.Type, ok bool) {
	switch e := ast.Unparen(x).(type) {
	case *ast.UnaryExpr:
		if e.Op != token.AND {
			return nil, false // not &...
		}

		cl, ok := ast.Unparen(e.X).(*ast.CompositeLit)
		if !ok {
			return nil, false // not &...{}
		}

		typ = p.TypesInfo.Types[cl].Type

		return typ, true

	case *ast.CallExpr:
		if len(e.Args) != 1 {
			return nil, false // some function
		}

		funType := p.TypesInfo.Types[e.Fun]
		if !funType.IsBuiltin() {
			return nil, false // no built-in
		}

		fun, ok := ast.Unparen(e.Fun).(*ast.Ident)
		if !ok || fun.Name != "new" {
			return nil, false // not new(...)
		}

		typ = p.TypesInfo.TypeOf(e.Args[0])

		return typ, true

	default:
		return nil, false
	}
}

// shouldSuppressDiagnostic determines whether a diagnostic should be suppressed.
// This is primarily relevant for `errors.Is` calls, where certain patterns involving
// `Is` or `Unwrap` methods might make the comparison legitimate despite involving a new address.
func shouldSuppressDiagnostic(t types.Type, isLeft bool) bool {
	ptr := types.NewPointer(t) // ptr is *T, the type of the newly created literal

	// The standard library `errors.Is(err, target)` function checks if `err` (or an error
	// in its `Unwrap` tree) matches `target`. This matching can occur in several ways.
	// For this linter, which flags `errors.Is(err, &T{})` (where `&T{}` is the `target`),
	// we are concerned with two scenarios for suppression:

	// 1. If `target` could be matched by an `Is(error) bool` method of `err`, assume it would
	//    be the `Is(error) bool` method of `target` and suppress the diagnostic in this case.
	if types.Implements(ptr, errorIsInterface) {
		return true
	}

	// 2. `err.Unwrap()`: If `err` is the newly created literal (`isLeft` is true),
	//    and its type `*T` implements an `Unwrap` method, `errors.Is` will traverse
	//    the unwrapped errors. The comparison might then be valid against an unwrapped error.
	//    Thus, we suppress the diagnostic in this case.
	if isLeft &&
		(types.Implements(ptr, errorUnwrapInterface) ||
			types.Implements(ptr, errorUnwrapArrayInterface)) {
		return true
	}

	// We do not have dynamic runtime types, these heuristics rely on static type information
	// and seem to work well in practice.

	return false
}
