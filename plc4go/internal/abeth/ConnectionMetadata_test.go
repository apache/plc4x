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

package abeth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestConnection_MetadataMatchesBuilders(t *testing.T) {
	c, _ := newTestConnection(t, DefaultConfiguration())
	metadata := c.GetMetadata()

	// ab-eth has no write path at all: the mspec knows a single command, the DF1 "protected typed
	// logical read", and plc4j's AbEthDriver overrides canRead and nothing else. Subscribing is
	// there because it is emulated by polling the read path.
	assert.True(t, metadata.CanRead())
	assert.False(t, metadata.CanWrite(), "ab-eth has no write command")
	assert.True(t, metadata.CanSubscribe())
	assert.False(t, metadata.CanBrowse(), "the ab-eth driver has no browser")

	testutils.AssertBuilderMatchesCapability(t, "reading", metadata.CanRead(), func() any { return c.ReadRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "writing", metadata.CanWrite(), func() any { return c.WriteRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return c.SubscriptionRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return c.UnsubscriptionRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "browsing", metadata.CanBrowse(), func() any { return c.BrowseRequestBuilder() })
}
