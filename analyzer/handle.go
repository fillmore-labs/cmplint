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
	"go/token"

	"fillmore-labs.com/cmplint/internal/typeutil"
)

// handleBinaryExpr checks binary expressions for equality or inequality
// comparisons involving addresses of composite literals or new() calls.
func (p pass) handleBinaryExpr(n *ast.BinaryExpr) {
	switch n.Op { //nolint:exhaustive
	case token.EQL, token.NEQ:
		// Delegate to comparison for further analysis of the comparison.
		p.comparison(n, n.X, n.Y, false)
	}
}

// handleCallExpr processes function calls by identifier, specifically looking
// for `errors.Is` (from standard library or x/exp) and assertion functions
// from `github.com/stretchr/testify` that perform error comparisons.
//
// It checks if the function is one of the targeted comparison functions
// and delegates the analysis of its arguments to comparison.
func (p pass) handleCallExpr(n *ast.CallExpr, functions map[typeutil.FuncName]funcType) {
	if len(n.Args) < 2 { // Other function or multi-valued argument
		return
	}

	// Retrieve the definition of the called function.
	fun, methodExpr, ok := typeutil.FuncOf(p.TypesInfo, n.Fun)
	if !ok {
		return
	}

	funcName := typeutil.NewFuncName(fun)

	ftyp, ok := functions[funcName]
	if !ok {
		return
	}

	baseArg := 0
	if methodExpr {
		baseArg = 1
	}

	switch ftyp {
	case funcErr0:
		// Delegate analysis of errors.Is(..., ...) to comparison.
		p.comparison(n, n.Args[baseArg], n.Args[baseArg+1], true)

	case funcErr1:
		if len(n.Args) < 3+baseArg { // should not happen
			p.LogErrorf(n, "Got only %d arguments for %s, expected at least %d", len(n.Args), funcName, 3+baseArg)

			return
		}

		// Delegate analysis of assert.ErrorIs(t, ..., ...) to comparison.
		p.comparison(n, n.Args[baseArg+1], n.Args[baseArg+2], true)

	case funcCmp0:
		// Delegate analysis of cmp(..., ...) to comparison.
		p.comparison(n, n.Args[baseArg], n.Args[baseArg+1], true)

	case funcCmp1:
		if len(n.Args) < 3+baseArg { // should not happen
			p.LogErrorf(n, "Got only %d arguments for %s, expected at least %d", len(n.Args), funcName, 3+baseArg)

			return
		}

		// Delegate analysis of assert.Equal(t, ..., ...) to comparison.
		p.comparison(n, n.Args[baseArg+1], n.Args[baseArg+2], false)

	case funcNone: // should not happen
		p.LogErrorf(n, "Unconfigured function %s", funcName)

		return
	}
}
