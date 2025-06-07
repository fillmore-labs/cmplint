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

/*
cmplint is a Go linter (static analysis tool) that detects comparisons against
the address of newly created values, such as ptr == &MyStruct{} or ptr == new(MyStruct).
These comparisons are almost always incorrect, as each expression creates a unique
allocation at runtime, usually yielding false or undefined results.

Usage:

	cmplint [flags] [package ...]

The flags are:

	-c int
		display offending line with this many lines of context (default -1)
	-V
		print version and exit

# Examples

To check the current package:

	cmplint .

To check all packages in the current module:

	cmplint ./...

Example of code flagged by cmplint:

	err := json.Unmarshal(msg, &es)
	if errors.Is(err, &json.UnmarshalTypeError{}) { // flagged
		//...
	}

# Links

- Equality of Pointers to Zero-Sized Types, https://blog.fillmore-labs.com/posts/zerosized-1/
*/
package main
