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
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// TagHandler parses UMAS tag addresses, which are symbol names of the PLC project. Ported from
// plc4j's UmasTagHandler.
type TagHandler struct {
}

func NewTagHandler() TagHandler {
	return TagHandler{}
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	if !symbolicAddressPattern.MatchString(tagAddress) {
		return nil, errors.Errorf("Unable to parse %s as a UMAS symbol name", tagAddress)
	}
	return NewTag(tagAddress), nil
}

// ParseQuery turns a browse query into a symbol-name filter.
//
// Deliberate deviation from plc4j, whose UmasTagHandler.parseQuery returns null and whose browse
// then answers every query with the whole symbol table. A null query isn't an option here - the
// plc4go browse-request builder stores what this returns and the browser dispatches on it - and
// while at it a query which actually selects something makes a browse request with more than one
// query mean something. The syntax is the smallest thing that covers "everything" and "this
// prefix": '*' matches any run of characters, '?' exactly one, matching is case insensitive
// (symbol names are folded to lower case throughout the driver), and a query without a wildcard
// matches that one symbol.
func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("A UMAS browse query can't be empty, use '*' to browse every symbol")
	}
	return NewSymbolQuery(query)
}
