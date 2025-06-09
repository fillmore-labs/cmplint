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
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	// Name is the default name.
	Name = "cmplint"

	// Doc is the default documentation.
	Doc = `cmplint is a Go linter (static analysis tool) that detects comparisons against
the address of newly created values, such as ptr == &MyStruct{} or ptr == new(MyStruct).
These comparisons are almost always incorrect, as each expression creates a unique
allocation at runtime, usually yielding false or undefined results.

Example of code flagged by cmplint:

	err := json.Unmarshal(msg, &es)
	if errors.Is(err, &json.UnmarshalTypeError{}) { // flagged
		//...
	}`

	// URL is the home page of the [analysis.Analyzer].
	URL = "https://pkg.go.dev/fillmore-labs.com/cmplint"
)

// Analyzer is the [analysis.Analyzer] for the cmplint linter.
//
// It checks for comparisons directly against the address of a composite literal
// or a newly allocated zero value using `new()`.
var Analyzer = New() //nolint:gochecknoglobals

// New creates a new instance of the [analysis.Analyzer] for the cmplint linter.
//
// It accepts a variadic number of [Option] arguments to customize its behavior,
// such as its name, documentation, or specific checking rules.
func New(opts ...Option) *analysis.Analyzer {
	o := makeOptions(opts)

	return &analysis.Analyzer{
		Name: o.name,
		Doc:  o.doc,
		URL:  URL,

		Flags: o.flags(),
		Run:   o.run,

		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}
