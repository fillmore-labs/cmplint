// Copyright 2025-2026 Oliver Eikemeier. All Rights Reserved.
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

package gclplugin

import cmplint "fillmore-labs.com/cmplint/analyzer"

// Settings represents the configuration options for an instance of the [Plugin].
type Settings struct {
	CheckIs *bool `json:"check-is,omitzero"`
}

// Options converts [Settings] into a list of [cmplint.Option] for the cmplint analyzer.
// It processes settings and applies them only when explicitly set (non-nil).
func (s Settings) Options() []cmplint.Option {
	var opts []cmplint.Option

	opts = appendOption(opts, s.CheckIs, cmplint.WithCheckIs)

	return opts
}

// appendOption appends a non-nil setting to an option list.
func appendOption[T, O any](opts []O, value *T, constructor func(T) O) []O {
	if value == nil {
		return opts
	}

	return append(opts, constructor(*value))
}
