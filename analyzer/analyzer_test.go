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

package analyzer_test

import (
	"errors"
	"go/version"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	. "fillmore-labs.com/cmplint/analyzer"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	const minVersion = "go1.25.0"
	if version.Compare(runtime.Version(), minVersion) < 0 {
		t.Skipf("needs minimal version %s", minVersion)
	}

	dir := analysistest.TestData()

	tests := []struct {
		name    string
		options Options
		flags   map[string]string
		pkg     string
	}{
		{
			name:    "default",
			options: nil,
			pkg:     "./a",
		},
		{
			name: "check-is=false",
			options: Options{
				WithCheckIs(false),
			},
			pkg: "./b",
		},
		{
			name: "check-is=false via flags",
			options: Options{
				WithName("ptrequality"),
				WithDoc("Documentation"),
			},
			flags: map[string]string{
				"check-is": "false",
			},
			pkg: "./b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := New(tt.options)
			for f, v := range tt.flags {
				if err := a.Flags.Set(f, v); err != nil {
					t.Fatalf("Can't set flag %s=%s: %v", f, v, err)
				}
			}

			analysistest.Run(t, dir, a, tt.pkg)
		})
	}
}

func TestMissingInspector(t *testing.T) {
	t.Parallel()

	if _, err := New().Run(&analysis.Pass{}); !errors.Is(err, ErrNoInspector) {
		t.Errorf("Expected ErrNoInspector, got %v", err)
	}
}
