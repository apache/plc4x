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

package bacnetip

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// newTestConnection wires a Connection with no codec/transport — enough for
// the lifecycle-flag and metadata assertions but not for sending real APDUs.
func newTestConnection(t *testing.T, options map[string][]string) *Connection {
	t.Helper()
	tm := transactions.NewRequestTransactionManager(math.MaxInt)
	return NewConnection(nil, NewTagHandler(), tm, options)
}

func TestConnection_GetMetadata_DeclaresCapabilities(t *testing.T) {
	conn := newTestConnection(t, map[string][]string{})
	md := conn.GetMetadata()
	require.NotNil(t, md)
	assert.True(t, md.CanRead(), "read must be supported")
	assert.True(t, md.CanWrite(), "write must be supported")
	assert.True(t, md.CanSubscribe(), "subscribe must be supported")
	assert.False(t, md.CanBrowse(), "browse not yet supported")
}

func TestConnection_ConfigurationParsedFromOptions(t *testing.T) {
	conn := newTestConnection(t, map[string][]string{
		"LocalDeviceId":         {"42"},
		"MaxApduLengthAccepted": {"480"},
		"SegmentationSupported": {"no-segmentation"},
	})
	assert.Equal(t, uint32(42), conn.configuration.LocalDeviceId)
	assert.Equal(t, uint16(480), conn.configuration.MaxApduLengthAccepted)
	assert.Equal(t, "no-segmentation", conn.configuration.SegmentationSupported)
}

func TestConnection_BuildersReturnNonNil(t *testing.T) {
	conn := newTestConnection(t, map[string][]string{})
	assert.NotNil(t, conn.ReadRequestBuilder())
	assert.NotNil(t, conn.WriteRequestBuilder())
	assert.NotNil(t, conn.SubscriptionRequestBuilder())
}

func TestConnection_AddSubscriber_IsIdempotent(t *testing.T) {
	conn := newTestConnection(t, map[string][]string{})
	sub := NewSubscriber(conn)
	conn.addSubscriber(sub)
	conn.addSubscriber(sub) // second add should be a no-op
	assert.Len(t, conn.subscribers, 1)
}

func TestConnection_InvokeIdGenerator_Monotonic(t *testing.T) {
	conn := newTestConnection(t, map[string][]string{})
	a := conn.invokeIdGenerator.getAndIncrement()
	b := conn.invokeIdGenerator.getAndIncrement()
	c := conn.invokeIdGenerator.getAndIncrement()
	assert.Equal(t, a+1, b)
	assert.Equal(t, b+1, c)
}

func TestConnection_StringIncludesDriverName(t *testing.T) {
	conn := newTestConnection(t, map[string][]string{})
	assert.Contains(t, conn.String(), "bacnetip")
}
