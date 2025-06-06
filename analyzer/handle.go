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
	"go/types"
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

// handleCallExpr checks call expressions to identify specific functions
// like errors.Is or testify assertion functions that perform comparisons
// and might involve addresses of composite literals or new() calls.
func (p pass) handleCallExpr(n *ast.CallExpr) {
	switch fun := ast.Unparen(n.Fun).(type) {
	case *ast.SelectorExpr:
		// Optimize to ignore method calls and function pointer fields.
		// Should be benchmarked to verify it is really faster.
		if _, ok := p.TypesInfo.Selections[fun]; ok {
			return
		}

		p.handleCallIdent(n, fun.Sel)

	case *ast.Ident:
		p.handleCallIdent(n, fun)
	}
}

// handleCallIdent processes function calls by identifier, specifically looking
// for `errors.Is` (from standard library or x/exp) and assertion functions
// from `github.com/stretchr/testify` that perform error comparisons.
//
// It checks if the function is one of the targeted comparison functions
// and delegates the analysis of its arguments to comparison.
func (p pass) handleCallIdent(n *ast.CallExpr, fun *ast.Ident) {
	path, name, ok := p.pathAndName(fun)
	if !ok {
		return
	}

	switch path {
	case "errors", "golang.org/x/exp/errors", "golang.org/x/xerrors":
		if len(n.Args) != 2 {
			return // errors.Is takes exactly two arguments
		}

		switch name { //nolint:gocritic
		case "Is":
			// Delegate analysis of errors.Is(..., ...) to comparison.
			p.comparison(n, n.Args[0], n.Args[1], true)
		}

	case "github.com/stretchr/testify/assert", "github.com/stretchr/testify/require":
		if len(n.Args) < 3 {
			return // Testify comparison functions typically take at least t, expected, actual
		}

		switch name {
		// assert.Equal does not compare object identity, but uses [reflect.DeepEqual].
		case "ErrorIs", "ErrorIsf", "NotErrorIs", "NotErrorIsf":
			// Delegate analysis of assert.ErrorIs(t, ..., ...) to comparison.
			p.comparison(n, n.Args[1], n.Args[2], true)
		}
	}
}

// pathAndName retrieves the  package path and name of the object used by an identifier.
// It returns the path, name, and a boolean indicating success.
func (p pass) pathAndName(id *ast.Ident) (path, name string, ok bool) {
	var pkg *types.Package

	obj, ok := p.TypesInfo.Uses[id]
	if ok {
		pkg = obj.Pkg()
	}

	if pkg == nil { // In Universe scope.
		return "", "", false
	}

	return pkg.Path(), obj.Name(), true
}
