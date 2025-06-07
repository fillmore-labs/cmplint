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
	"flag"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Documentation constants.
const (
	Name = "cmplint"
	Doc  = `cmplint is a Go linter (static analysis tool) that detects comparisons against
the address of newly created values, such as ptr == &MyStruct{} or ptr == new(MyStruct).
These comparisons are almost always incorrect, as each expression creates a unique
allocation at runtime, usually yielding false or undefined results.

Example of code flagged by cmplint:

	err := json.Unmarshal(msg, &es)
	if errors.Is(err, &json.UnmarshalTypeError{}) { // flagged
		//...
	}`

	URL = "https://pkg.go.dev/fillmore-labs.com/cmplint"
)

//nolint:gochecknoglobals
var (
	// Analyzer is the [analysis.Analyzer] for the cmplint linter.
	// It checks for comparisons directly against the address of a composite literal
	// or a newly allocated zero value using `new()`.
	Analyzer = &analysis.Analyzer{
		Name:  Name,
		Doc:   Doc,
		URL:   URL,
		Flags: flags(),
		Run:   run,

		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	// CheckIs indicates whether to check for `Is(error) bool` and unwrap methods.
	CheckIs = true
)

func flags() flag.FlagSet {
	var flags flag.FlagSet

	flags.BoolVar(&CheckIs, "check-is", true, `suppresses diagnostic on errors if the type has an "Is" method`)

	return flags
}
