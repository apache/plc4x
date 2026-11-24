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

package ads

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// strings.Split always returns at least one element for a non-empty query, so the guard that
// joins the second address segment has to test for more than one segment. Testing for more than
// zero made it index an empty slice, and a query with no dot in it - "*" above all - took the
// browse goroutine down with "index out of range [0] with length 0".
func TestFilterSymbols_queryWithoutADotDoesNotPanic(t *testing.T) {
	connection := connectionWithArraySymbol()

	for _, query := range []string{"*", "**", "MAIN", "nosuchsymbol"} {
		t.Run(query, func(t *testing.T) {
			assert.NotPanics(t, func() {
				assert.Nil(t, connection.filterSymbols(query))
			})
		})
	}
}

// A dotted query keeps working exactly as before: the symbol is still assembled from the first
// two segments.
func TestFilterSymbols_dottedQueryStillResolves(t *testing.T) {
	connection := connectionWithArraySymbol()

	assert.NotPanics(t, func() { connection.filterSymbols("MAIN.arr") })
	assert.Nil(t, connection.filterSymbols("MAIN.nosuchsymbol"))
}
