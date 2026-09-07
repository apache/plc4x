/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package umas

import (
	"regexp"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// SymbolQuery selects symbols of the PLC project by name. See TagHandler.ParseQuery for why this
// exists at all - plc4j has no query type and browses everything.
type SymbolQuery interface {
	apiModel.PlcQuery

	// Matches says whether a symbol name is selected by this query. The name is expected to be
	// folded to lower case already, which is how the driver keys its symbol table.
	Matches(lowerCaseSymbolName string) bool
}

type symbolQuery struct {
	query   string
	pattern *regexp.Regexp
}

var _ SymbolQuery = symbolQuery{}

// NewSymbolQuery compiles a browse query into a matcher. '*' stands for any run of characters, '?'
// for exactly one, and everything else is matched literally and case insensitively.
func NewSymbolQuery(query string) (SymbolQuery, error) {
	var pattern strings.Builder
	pattern.WriteString("^")
	for _, queryRune := range query {
		switch queryRune {
		case '*':
			pattern.WriteString(".*")
		case '?':
			pattern.WriteString(".")
		default:
			pattern.WriteString(regexp.QuoteMeta(string(queryRune)))
		}
	}
	pattern.WriteString("$")
	compiled, err := regexp.Compile("(?i)" + pattern.String())
	if err != nil {
		return nil, errors.Wrapf(err, "Error compiling the browse query %s", query)
	}
	return symbolQuery{query: query, pattern: compiled}, nil
}

func (q symbolQuery) GetQueryString() string {
	return q.query
}

func (q symbolQuery) Matches(lowerCaseSymbolName string) bool {
	return q.pattern.MatchString(lowerCaseSymbolName)
}

func (q symbolQuery) String() string {
	return "umas.SymbolQuery{" + q.query + "}"
}
