//go:build linux

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

package serialport

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPorts(t *testing.T) {
	ports, err := ListPorts()
	require.NoError(t, err)
	// CI machines may expose zero real UARTs; pin the invariants that hold
	// regardless: device paths, no duplicates, sorted output.
	for _, p := range ports {
		assert.True(t, strings.HasPrefix(p, "/dev/"), "port %q must be a /dev path", p)
	}
	assert.True(t, sort.StringsAreSorted(ports), "ports must be sorted")
	seen := map[string]bool{}
	for _, p := range ports {
		assert.False(t, seen[p], "duplicate %q", p)
		seen[p] = true
	}
}
