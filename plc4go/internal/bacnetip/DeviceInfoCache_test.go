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
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceInfoCache_PutAndGet(t *testing.T) {
	c := NewDeviceInfoCache()
	addr := &net.UDPAddr{IP: net.IPv4(10, 0, 0, 5), Port: 47808}
	c.Put(DeviceInfo{DeviceId: 1234, Address: addr, SourceNetwork: 0, MaxApdu: 1476, VendorId: 5})

	got := c.Get(1234)
	require.NotNil(t, got)
	assert.Equal(t, uint32(1234), got.DeviceId)
	assert.Equal(t, addr.String(), got.Address.String())
	assert.False(t, got.LastSeen.IsZero(), "Put should populate LastSeen")
}

func TestDeviceInfoCache_GetUnknown_ReturnsNil(t *testing.T) {
	c := NewDeviceInfoCache()
	assert.Nil(t, c.Get(999))
}

func TestDeviceInfoCache_Touch_BumpsLastSeen(t *testing.T) {
	c := NewDeviceInfoCache()
	c.Put(DeviceInfo{DeviceId: 7})
	original := c.Get(7).LastSeen

	// Sleep just enough to make sure the monotonic clock moves.
	time.Sleep(time.Millisecond)
	c.Touch(7)
	updated := c.Get(7).LastSeen
	assert.True(t, updated.After(original), "Touch should advance LastSeen")
}

func TestDeviceInfoCache_Touch_UnknownIsNoOp(t *testing.T) {
	c := NewDeviceInfoCache()
	// Just don't panic.
	c.Touch(404)
	assert.Empty(t, c.All())
}

func TestDeviceInfoCache_All_ReturnsSnapshot(t *testing.T) {
	c := NewDeviceInfoCache()
	c.Put(DeviceInfo{DeviceId: 1})
	c.Put(DeviceInfo{DeviceId: 2})
	all := c.All()
	assert.Len(t, all, 2)

	// Mutating the snapshot must not affect the cache.
	all[0].DeviceId = 9999
	again := c.All()
	for _, e := range again {
		assert.NotEqual(t, uint32(9999), e.DeviceId)
	}
}

func TestDeviceInfoCache_LoadStatic_ValidEntries(t *testing.T) {
	c := NewDeviceInfoCache()
	err := c.LoadStatic("1234@0:10.0.0.5:47808, 42@7:10.0.0.6:47808")
	require.NoError(t, err)

	d1 := c.Get(1234)
	require.NotNil(t, d1)
	assert.Equal(t, uint16(0), d1.SourceNetwork)
	assert.Equal(t, "10.0.0.5:47808", d1.Address.String())

	d2 := c.Get(42)
	require.NotNil(t, d2)
	assert.Equal(t, uint16(7), d2.SourceNetwork)
}

func TestDeviceInfoCache_LoadStatic_Empty(t *testing.T) {
	c := NewDeviceInfoCache()
	require.NoError(t, c.LoadStatic(""))
	assert.Empty(t, c.All())
}

func TestDeviceInfoCache_LoadStatic_PartialFailureReturnsErr(t *testing.T) {
	c := NewDeviceInfoCache()
	err := c.LoadStatic("1234@0:10.0.0.5:47808,bogus,5678@1:10.0.0.6:47808")
	require.Error(t, err)
	// Valid entries should still be in the cache.
	assert.NotNil(t, c.Get(1234))
	assert.NotNil(t, c.Get(5678))
}

func TestParseStaticDevice_BadInputs(t *testing.T) {
	cases := []string{
		"",                  // empty entry
		"42",                // missing '@'
		"42@",               // missing network
		"42@1",              // missing host:port
		"42@nan:10.0.0.1:1", // non-numeric network
		"42@0:10.0.0.1",     // missing port
	}
	for _, c := range cases {
		_, err := parseStaticDevice(c)
		assert.Error(t, err, "input %q should fail", c)
	}
}
