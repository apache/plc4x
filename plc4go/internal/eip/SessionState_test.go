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

package eip

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
)

func routingAddressBytes(t *testing.T, segments []readWriteModel.PathSegment) int {
	total := 0
	for _, segment := range segments {
		total += int(segment.GetLengthInBytes(context.Background()))
	}
	return total
}

func TestBuildRoutingAddressDefault(t *testing.T) {
	configuration := Configuration{backplane: 1, slot: 3}
	segments, pathSize := buildRoutingAddress(zerolog.Nop(), configuration)
	// port segment (backplane 1 / slot 3) + ClassID 2 + InstanceID 1
	require.Len(t, segments, 3)
	portSegment, ok := segments[0].(readWriteModel.PortSegment)
	require.True(t, ok)
	normal, ok := portSegment.GetSegmentType().(readWriteModel.PortSegmentNormal)
	require.True(t, ok)
	assert.Equal(t, uint8(1), normal.GetPort())
	assert.Equal(t, uint8(3), normal.GetLinkAddress())
	totalBytes := routingAddressBytes(t, segments)
	expected := (totalBytes + totalBytes%2) / 2
	assert.Equal(t, uint8(expected), pathSize)
}

func TestBuildRoutingAddressCommunicationPath(t *testing.T) {
	configuration := Configuration{communicationPath: "1,4,2,192.168.0.1,1,1"}
	segments, _ := buildRoutingAddress(zerolog.Nop(), configuration)
	// three routed hops + ClassID + InstanceID
	require.Len(t, segments, 5)
	hop2, ok := segments[1].(readWriteModel.PortSegment)
	require.True(t, ok)
	extended, ok := hop2.GetSegmentType().(readWriteModel.PortSegmentExtended)
	require.True(t, ok)
	assert.Equal(t, uint8(2), extended.GetPort())
	// "192.168.0.1" is 11 chars -> reported length 11, padded to "192.168.0.1\x00"
	assert.Equal(t, uint8(11), extended.GetLinkAddressSize())
	assert.Equal(t, "192.168.0.1\x00", extended.GetAddress())
}

func TestNewSessionStateSerialNumber(t *testing.T) {
	fixed := NewSessionState(zerolog.Nop(), Configuration{connectionSerialNumber: 4660})
	assert.Equal(t, uint16(4660), fixed.connectionSerialNumber)
	random := NewSessionState(zerolog.Nop(), Configuration{})
	assert.NotEqual(t, uint16(0), random.connectionSerialNumber)
}

func TestNextSequenceCount(t *testing.T) {
	state := NewSessionState(zerolog.Nop(), Configuration{})
	assert.Equal(t, uint16(1), state.nextSequenceCount())
	assert.Equal(t, uint16(2), state.nextSequenceCount())
}
