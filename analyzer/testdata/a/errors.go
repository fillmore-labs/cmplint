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

import (
	"errors"
	. "errors"
	"os"

	errorsx "golang.org/x/exp/errors"
	"golang.org/x/xerrors"
)

type myError1 struct{}

func (myError1) Error() string {
	return ""
}

func (myError1) Is(_, _ error) bool {
	return false
}

type myErrorEmbedded struct {
	*myError1
}

func Errors() {
	_ = errors.Is(myError1{}, &myError1{}) // want "is false or undefined"

	_ = errors.Is(&myError1{}, &myError1{}) // want "is false or undefined"

	_ = errors.Is(&myError1{}, &myErrorEmbedded{}) // want "is false or undefined"

	_ = errors.Is(&struct{ myError1 }{}, &myError1{}) // want "is false or undefined"

	_ = Is(&myError1{}, &myError1{}) // want "is false or undefined"

	var e myError1
	_ = errors.Is(nil, &e)

	_ = errors.As(myError1{}, &myError1{})

	_ = errors.Join(&myError1{}, &myError1{})

	_ = errors.Unwrap(&myError1{})

	_ = errorsx.Is(&myError1{}, &myError1{}) // want "is false or undefined"

	_ = xerrors.Is(func() error { // want "is false or undefined"
		return &myErrorWithIs{}
	}(), &myError1{})
}

func Errors2() {
	errors := myError1{}
	_ = errors.Is(&myError1{}, &myError1{})
}

type StructWithIsField struct {
	Is func(_, _ error) bool
}

func Errors3() {
	errors := StructWithIsField{Is: func(_, _ error) bool { return false }}

	_ = errors.Is(&myError1{}, &myError1{})
}

type myErrorWithIs struct{}

func (myErrorWithIs) Error() string {
	return "my error with is"
}

func (myErrorWithIs) Is(err error) bool {
	_, ok := err.(*myErrorWithIs)

	return ok
}

type myErrorWithUnwrap struct{}

func (myErrorWithUnwrap) Error() string {
	return "my error with unwrap"
}

func (myErrorWithUnwrap) Unwrap() error {
	return os.ErrProcessDone
}

type myErrorWithUnwrapArray struct{}

func (myErrorWithUnwrapArray) Error() string {
	return "my error with unwrap"
}

func (myErrorWithUnwrapArray) Unwrap() []error {
	return []error{os.ErrProcessDone}
}

func Errors4() {
	_ = errors.Is(&myErrorWithIs{}, &myError1{})

	_ = errors.Is(&struct{ *myErrorWithIs }{}, &myError1{})

	_ = errors.Is(&myErrorWithUnwrap{}, os.ErrProcessDone)

	_ = errors.Is(&myErrorWithUnwrapArray{}, os.ErrProcessDone)

	_ = errors.Is(os.ErrProcessDone, &myErrorWithUnwrap{}) // want "type \"?myErrorWithUnwrap\"?"
}
