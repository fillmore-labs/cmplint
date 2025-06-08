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

// New creates a new instance of the analyzer.
// It accepts a variadic number of [Option] arguments to customize its behavior,
// such as its name, documentation, or specific checking rules.
func New(opts ...Option) *analysis.Analyzer {
	o := newOptions(opts)

	a := &analysis.Analyzer{
		Name:  o.name,
		Doc:   o.doc,
		URL:   URL,
		Flags: o.flags(),
		Run:   o.run,

		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return a
}

// run is the main analysis function for the analyzer.
func (o *options) run(a *analysis.Pass) (any, error) {
	in, ok := a.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspector
	}

	p := pass{Pass: a, checkis: o.checkis}

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

// ErrNoInspector is returned by the analyzer's Run method if the required
// AST inspector result is not available from its dependencies.
var ErrNoInspector = errors.New("inspector missing")

// pass wraps the analysis.Pass and includes analyzer-specific state
// like configuration options. It provides helper methods for the analysis logic.
type pass struct {
	*analysis.Pass
	checkis bool
}
