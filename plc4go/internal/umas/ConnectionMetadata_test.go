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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestConnection_MetadataMatchesBuilders(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	metadata := connection.GetMetadata()

	// Reading, writing and browsing are what plc4j's UmasDriver advertises. Subscribing is there
	// because it is emulated by polling the read path, the way plc4j's
	// PollingSubscriptionConnectionBase does it for this driver.
	assert.True(t, metadata.CanRead())
	assert.True(t, metadata.CanWrite())
	assert.True(t, metadata.CanSubscribe())
	assert.True(t, metadata.CanBrowse())

	testutils.AssertBuilderMatchesCapability(t, "reading", metadata.CanRead(), func() any { return connection.ReadRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "writing", metadata.CanWrite(), func() any { return connection.WriteRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return connection.SubscriptionRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return connection.UnsubscriptionRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "browsing", metadata.CanBrowse(), func() any { return connection.BrowseRequestBuilder() })
}
