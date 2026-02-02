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

import "fillmore-labs.com/cmplint/internal/typeutil"

type funcType int

const (
	funcNone funcType = iota
	funcErr0
	funcErr1
	funcCmp0
	funcCmp1
)

// Since we have a lot of hardcoded libraries here, a check by signature might be a better heuristic.
var functions = map[typeutil.FuncName]funcType{ //nolint:gochecknoglobals
	{Path: "errors", Name: "Is"}:                                                               funcErr0,
	{Path: "golang.org/x/exp/errors", Name: "Is"}:                                              funcErr0,
	{Path: "golang.org/x/xerrors", Name: "Is"}:                                                 funcErr0,
	{Path: "github.com/pkg/errors", Name: "Is"}:                                                funcErr0,
	{Path: "github.com/friendsofgo/errors", Name: "Is"}:                                        funcErr0,
	{Path: "github.com/go-errors/errors", Name: "Is"}:                                          funcErr0,
	{Path: "github.com/go-faster/errors", Name: "Is"}:                                          funcErr0,
	{Path: "github.com/cockroachdb/errors", Name: "Is"}:                                        funcErr0,
	{Path: "github.com/cockroachdb/errors/markers", Name: "Is"}:                                funcErr0,
	{Path: "github.com/juju/errors", Name: "Is"}:                                               funcErr0,
	{Path: "gotest.tools/v3/assert", Name: "Equal"}:                                            funcCmp1,
	{Path: "gotest.tools/v3/assert", Name: "ErrorIs"}:                                          funcErr1,
	{Path: "gotest.tools/v3/assert/cmp", Name: "Equal"}:                                        funcCmp0,
	{Path: "github.com/stretchr/testify/assert", Name: "ErrorIs"}:                              funcErr1,
	{Path: "github.com/stretchr/testify/assert", Name: "ErrorIsf"}:                             funcErr1,
	{Path: "github.com/stretchr/testify/assert", Name: "NotErrorIs"}:                           funcErr1,
	{Path: "github.com/stretchr/testify/assert", Name: "NotErrorIsf"}:                          funcErr1,
	{Path: "github.com/stretchr/testify/require", Name: "ErrorIs"}:                             funcErr1,
	{Path: "github.com/stretchr/testify/require", Name: "ErrorIsf"}:                            funcErr1,
	{Path: "github.com/stretchr/testify/require", Name: "NotErrorIs"}:                          funcErr1,
	{Path: "github.com/stretchr/testify/require", Name: "NotErrorIsf"}:                         funcErr1,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "ErrorIs"}:      funcErr0,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "ErrorIsf"}:     funcErr0,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "NotErrorIs"}:   funcErr0,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "NotErrorIsf"}:  funcErr0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "ErrorIs"}:     funcErr0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "ErrorIsf"}:    funcErr0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "NotErrorIs"}:  funcErr0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "NotErrorIsf"}: funcErr0,
}
