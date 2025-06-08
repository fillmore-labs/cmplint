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

//nolint:ireturn
package analyzer

import "flag"

// options holds the configurable parameters for the analyzer.
type options struct {
	name    string
	doc     string
	checkis bool
}

// newOptions returns a new [options] struct initialized with default values.
// These defaults are used if no overriding options or flags are provided.
func newOptions(opts Options) *options {
	o := &options{ // Defaults
		name:    Name,
		doc:     Doc,
		checkis: true,
	}
	o.apply(opts)

	return o
}

// apply iterates through a list of [Option] values and applies each one
// to the [options] struct `o`, modifying it in place.
func (o *options) apply(opts Options) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// flags returns a [flag.FlagSet] containing command-line flags that can
// configure the analyzer's behavior. These flags correspond to the fields
// in the [options] struct.
// The returned FlagSet is used by the analysis driver to parse command-line
// arguments for the analyzer.
func (o *options) flags() flag.FlagSet {
	var flags flag.FlagSet

	flags.BoolVar(&o.checkis, "check-is", o.checkis,
		`suppress diagnostic on errors.Is if the compared type has an "Is(error) bool" method`)

	return flags
}

// Option configures specific behavior of the [analysis.Analyzer].
type Option interface {
	apply(opts *options)
}

// Options is a slice of [Option] values. It also implements the [Option] interface,
// allowing multiple options to be applied sequentially as a single Option.
type Options []Option

// apply applies each Option in the Options slice to the given [options] struct.
func (o Options) apply(opts *options) {
	for _, opt := range o {
		opt.apply(opts)
	}
}

// WithName returns an [Option] that sets a custom name for the analyzer.
// This name is used in analyzer reports and documentation.
func WithName(name string) Option {
	return nameOption{name: name}
}

// nameOption implements the [Option] interface to set the analyzer's name.
type nameOption struct {
	name string
}

// apply sets the name field in the provided [options] struct.
func (o nameOption) apply(opts *options) {
	opts.name = o.name
}

// WithDoc returns an [Option] that sets custom documentation for the analyzer.
// This documentation is used by analysis drivers.
func WithDoc(doc string) Option {
	return docOption{doc: doc}
}

// docOption implements the [Option] interface to set the analyzer's documentation.
type docOption struct {
	doc string
}

// apply sets the doc field in the provided [options] struct.
func (o docOption) apply(opts *options) {
	opts.doc = o.doc
}

// WithCheckIs returns an [Option] that configures the diagnostic suppression behavior
// related to the `Is(error) bool` method.
// If `checkis` is true (the default), diagnostics for `errors.Is(err, &MyError{})`
// may be suppressed if `*MyError` implements `Is(error) bool`.
// If false, this specific suppression heuristic is disabled.
func WithCheckIs(checkis bool) Option {
	return checkisOption{checkis: checkis}
}

// checkisOption implements the [Option] interface to configure the `checkis` behavior.
type checkisOption struct {
	checkis bool
}

// apply sets the checkis field in the provided [options] struct.
func (o checkisOption) apply(opts *options) {
	opts.checkis = o.checkis
}
