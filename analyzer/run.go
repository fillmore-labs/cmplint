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
	"errors"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ErrNoInspector is returned when the ast inspector is missing.
var ErrNoInspector = errors.New("cmplint: inspector missing")

// pass wraps the analysis.Pass to provide helper methods.
type pass struct {
	*analysis.Pass
}

// run is the entry point for the cmplint analyzer.
// It traverses the AST of the package's files and reports diagnostics
// for problematic comparisons.
func run(a *analysis.Pass) (any, error) {
	in, ok := a.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspector
	}

	p := pass{a}

	for n := range in.PreorderSeq((*ast.BinaryExpr)(nil), (*ast.CallExpr)(nil)) {
		switch n := n.(type) {
		case *ast.BinaryExpr: // Process equality and inequality operations.
			p.handleBinaryExpr(n)

		case *ast.CallExpr: // Check for errors.Is(x, y) and testify functions.
			p.handleCallExpr(n)
		}
	}

	return any(nil), nil
}
