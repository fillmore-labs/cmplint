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

package a

func Comparison() {
	_ = &struct{}{} == new(struct{}) // want "is false or undefined"

	_ = nil != new(struct{}) // want "is false or undefined"

	_ = nil != &[0]byte{} // want "is false or undefined"

	_ = &struct{ _ int }{} == new(struct{ _ int }) // want "is always false"

	_ = struct{}{} == struct{}{}

	_ = true == !false

	_ = nil == make(map[string]int)
}

func Comparison2() {
	new := func(_ int) *struct{} { return nil }

	_ = nil != new(0)

	new1 := func(_ int) *struct{} { return nil }

	_ = nil != new1(0)

	new2 := func(_, _ int) *struct{} { return nil }

	_ = nil != new2(0, 0)
}

func Comparison3() {
	type MyStruct struct{ _ int }

	if (nil == &MyStruct{}) { // want "is always false"
		// ...
	}

	if nil == (&MyStruct{}) { // want "is always false"
		// ...
	}

	if nil == &(MyStruct{}) { // want "is always false"
		// ...
	}

	if nil == ((new)(MyStruct)) { // want "is always false"
		// ...
	}
}
