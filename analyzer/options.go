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

import (
	"log/slog"

	"fillmore-labs.com/cmplint/internal/options"
)

// Option configures specific behavior of the [analysis.Analyzer].
type Option interface {
	Apply(opts *option)
	LogAttr() slog.Attr
}

// Join creates a new Option joining the provided Option values.
//
// The result implements [slog.LogValuer], so the following is evaluated lazily:
//
//	slog.LogAttrs(ctx, slog.LevelInfo, "settings", Join(opts...).LogAttr())
func Join(opts ...Option) Option {
	return options.Join(opts)
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

// Apply sets the name field in the provided [options] struct.
func (o nameOption) Apply(opts *option) {
	opts.name = o.name
}

func (o nameOption) LogAttr() slog.Attr {
	return slog.String("name", o.name)
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

// Apply sets the doc field in the provided [options] struct.
func (o docOption) Apply(opts *option) {
	opts.doc = o.doc
}

func (o docOption) LogAttr() slog.Attr {
	return slog.String("doc", o.doc)
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

// Apply sets the checkis field in the provided [options] struct.
func (o checkisOption) Apply(opts *option) {
	opts.checkis = o.checkis
}

// LogAttr implements [Option].
func (o checkisOption) LogAttr() slog.Attr {
	return slog.Bool("check-is", o.checkis)
}
