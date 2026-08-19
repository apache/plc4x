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

package modbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertBuilderMatchesCapability checks that a connection hands out a request builder exactly when
// its metadata advertises the matching capability. An advertised capability has to produce a
// builder, and an unadvertised one has to refuse rather than hand out a builder that cannot work.
func assertBuilderMatchesCapability(t *testing.T, capability string, advertised bool, builder func() any) {
	t.Helper()
	var handedOut any
	panicked := func() (panicked bool) {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		handedOut = builder()
		return false
	}()
	if advertised {
		assert.False(t, panicked, "%s is advertised but the builder panicked", capability)
		assert.NotNil(t, handedOut, "%s is advertised but no builder was handed out", capability)
		return
	}
	assert.True(t, panicked, "%s is not advertised but a builder was handed out", capability)
}

func TestConnection_MetadataMatchesBuilders(t *testing.T) {
	c := NewConnection(Configuration{}, nil, map[string][]string{}, NewTagHandler())
	metadata := c.GetMetadata()

	assert.True(t, metadata.CanRead())
	assert.True(t, metadata.CanWrite())
	assert.False(t, metadata.CanSubscribe(), "the modbus driver has no subscriber")
	assert.False(t, metadata.CanBrowse(), "the modbus driver has no browser")

	assertBuilderMatchesCapability(t, "reading", metadata.CanRead(), func() any { return c.ReadRequestBuilder() })
	assertBuilderMatchesCapability(t, "writing", metadata.CanWrite(), func() any { return c.WriteRequestBuilder() })
	assertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return c.SubscriptionRequestBuilder() })
	assertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return c.UnsubscriptionRequestBuilder() })
	assertBuilderMatchesCapability(t, "browsing", metadata.CanBrowse(), func() any { return c.BrowseRequestBuilder() })
}
