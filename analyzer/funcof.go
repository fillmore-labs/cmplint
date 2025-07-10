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
	"go/types"
)

// funcOf resolves an expression that identifies a function to its underlying [*types.Func] object.
// It handles simple function identifiers, selector expressions (e.g., `pkg.Func`), and
// generic function instantiations (e.g., `Func[T]` or `Func[T, U]`).
// It returns the resolved function and true if successful, or nil and false otherwise.
func (p pass) funcOf(e ast.Expr) (fun *types.Func, ok bool) {
	return funcOf(p.TypesInfo, e)
}

func funcOf(info *types.Info, ex ast.Expr) (fun *types.Func, ok bool) {
	for {
		switch e := ex.(type) {
		case *ast.Ident:
			fun, ok = info.Uses[e].(*types.Func)

			return fun, ok

		case *ast.SelectorExpr:
			fun, ok = info.Uses[e.Sel].(*types.Func) // types.Checker calls recordUse from recordSelection

			return fun, ok

		case *ast.IndexExpr:
			if !info.Types[e.Index].IsType() {
				return nil, false
			}

			ex = e.X // Unwrap generic expression

		case *ast.IndexListExpr:
			if len(e.Indices) == 0 || !info.Types[e.Indices[0]].IsType() { // should not happen
				return nil, false
			}

			ex = e.X // Unwrap generic expression

		case *ast.ParenExpr:
			ex = e.X

		default:
			return nil, false
		}
	}
}
