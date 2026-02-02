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

package analyzer //nolint:testpackage

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestExprToString(t *testing.T) { //nolint:funlen
	t.Parallel()

	p := pass{
		Pass: &analysis.Pass{
			Fset: token.NewFileSet(),
		},
	}

	testCases := []struct {
		name string
		expr string
		want string
	}{
		{
			name: "simple identifier",
			expr: "myVar",
			want: "myVar",
		},
		{
			name: "binary expression",
			expr: "a + b",
			want: "a + b",
		},
		{
			name: "call expression",
			expr: "myFunc(1, \"hello\")",
			want: "myFunc(1, \"hello\")",
		},
		{
			name: "composite literal",
			expr: "MyType{Field: 42}",
			want: "MyType{Field: 42}",
		},
		{
			name: "address of composite literal",
			expr: "&MyType{}",
			want: "&MyType{}",
		},
		{
			name: "selector expression",
			expr: "pkg.Var",
			want: "pkg.Var",
		},
		{
			name: "nil expression should be invalid",
			expr: "",
			want: "invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var e ast.Expr

			if tc.expr != "" {
				expr, err := parser.ParseExprFrom(p.Fset, "test.go", tc.expr, 0)
				if err != nil {
					t.Fatalf("failed to parse expression %q: %v", tc.expr, err)
				}
				e = expr
			}

			if got := p.exprToString(e); got != tc.want {
				t.Errorf("exprToString() = %q, want %q", got, tc.want)
			}
		})
	}
}
