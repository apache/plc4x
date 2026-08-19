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

package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// handOutBuilder calls builder, reporting what it handed out and whether it panicked instead.
func handOutBuilder(builder func() any) (handedOut any, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	return builder(), false
}

// AssertAdvertisedCapabilityHasBuilder checks the direction that always holds: a connection that
// advertises a capability has to hand out a builder for it. The converse is not asserted, because a
// driver may hand out a builder whose operations report "not implemented" while honestly advertising
// the capability as absent. Use AssertBuilderMatchesCapability for connections where both
// directions are supposed to hold.
func AssertAdvertisedCapabilityHasBuilder(t *testing.T, capability string, advertised bool, builder func() any) {
	t.Helper()
	if !advertised {
		return
	}
	handedOut, panicked := handOutBuilder(builder)
	assert.False(t, panicked, "%s is advertised but the builder panicked", capability)
	assert.NotNil(t, handedOut, "%s is advertised but no builder was handed out", capability)
}

// AssertBuilderMatchesCapability checks that a connection hands out a request builder exactly when
// its metadata advertises the matching capability. An advertised capability has to produce a
// builder, and an unadvertised one has to refuse rather than hand out a builder that cannot work.
func AssertBuilderMatchesCapability(t *testing.T, capability string, advertised bool, builder func() any) {
	t.Helper()
	if advertised {
		AssertAdvertisedCapabilityHasBuilder(t, capability, advertised, builder)
		return
	}
	_, panicked := handOutBuilder(builder)
	assert.True(t, panicked, "%s is not advertised but a builder was handed out", capability)
}
