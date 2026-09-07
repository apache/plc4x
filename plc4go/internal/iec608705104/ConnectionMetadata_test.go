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

package iec608705104

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestConnection_MetadataMatchesBuilders(t *testing.T) {
	c := NewConnection(DefaultConfiguration(), nil, map[string][]string{}, NewTagHandler())
	metadata := c.GetMetadata()

	// IEC 60870-5-104 as this driver speaks it is push only: the controlled station reports what its
	// own configuration says it should, and there is no request for the current value of a point. No
	// command path either - a command needs an I-format frame, a send sequence number and the
	// activation/confirmation/termination dance around it, none of which this driver does yet.
	// plc4j's Iec60870514PlcDriver says the same by overriding canSubscribe and nothing else.
	assert.False(t, metadata.CanRead(), "the station is never asked for a value")
	assert.False(t, metadata.CanWrite(), "no command path")
	assert.True(t, metadata.CanSubscribe())
	assert.False(t, metadata.CanBrowse(), "a general interrogation answers with ASDUs, not a directory")

	testutils.AssertBuilderMatchesCapability(t, "reading", metadata.CanRead(), func() any { return c.ReadRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "writing", metadata.CanWrite(), func() any { return c.WriteRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return c.SubscriptionRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "subscribing", metadata.CanSubscribe(), func() any { return c.UnsubscriptionRequestBuilder() })
	testutils.AssertBuilderMatchesCapability(t, "browsing", metadata.CanBrowse(), func() any { return c.BrowseRequestBuilder() })
}
