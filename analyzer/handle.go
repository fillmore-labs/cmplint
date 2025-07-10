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
func (p pass) handleCallExpr(n *ast.CallExpr) {
	// Retrieve the definition of the called function.
	fun, ok := p.funcOf(n.Fun)
	if !ok {
		return
	}

	var path string
	if pkg := fun.Pkg(); pkg != nil {
		path = pkg.Path()
	}

	name := fun.Name()

	switch path {
	case "errors", "golang.org/x/exp/errors", "golang.org/x/xerrors", "github.com/pkg/errors":
		if len(n.Args) != 2 {
			return // errors.Is takes exactly two arguments.
		}

		switch name {
		case "Is":
			// Delegate analysis of errors.Is(..., ...) to comparison.
			p.comparison(n, n.Args[0], n.Args[1], true)

		default:
		}

	case "github.com/stretchr/testify/assert", "github.com/stretchr/testify/require":
		if len(n.Args) < 3 {
			return // Testify comparison functions typically take at least t, expected, actual.
		}

		switch name {
		// assert.Equal does not compare object identity, but uses [reflect.DeepEqual].
		case "ErrorIs", "ErrorIsf", "NotErrorIs", "NotErrorIsf":
			// Delegate analysis of assert.ErrorIs(t, ..., ...) to comparison.
			p.comparison(n, n.Args[1], n.Args[2], true)
		}

	case "gotest.tools/v3/assert":
		if len(n.Args) < 3 {
			return // gotest.tools comparison functions typically take at least t, expected, actual.
		}

		switch name {
		case "Equal":
			// Delegate analysis of assert.Equal(t, ..., ...) to comparison.
			p.comparison(n, n.Args[1], n.Args[2], false)

		case "ErrorIs":
			// Delegate analysis of assert.ErrorIs(t, ..., ...) to comparison.
			p.comparison(n, n.Args[1], n.Args[2], true)
		}
	}
}
