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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/options"
)

func TestNewDiscoverer_DefaultTimeout(t *testing.T) {
	d := NewDiscoverer()
	assert.Equal(t, 5*time.Second, d.discoveryTimeout)
}

func TestSetDiscoveryTimeout_Override(t *testing.T) {
	d := NewDiscoverer()
	d.SetDiscoveryTimeout(2 * time.Second)
	assert.Equal(t, 2*time.Second, d.discoveryTimeout)
}

func TestSetDiscoveryTimeout_ZeroFallsBackToDefault(t *testing.T) {
	d := NewDiscoverer()
	d.SetDiscoveryTimeout(0)
	assert.Equal(t, 5*time.Second, d.discoveryTimeout)
}

func TestSetDiscoveryTimeout_NegativeFallsBackToDefault(t *testing.T) {
	d := NewDiscoverer()
	d.SetDiscoveryTimeout(-1 * time.Second)
	assert.Equal(t, 5*time.Second, d.discoveryTimeout)
}

func TestResolveBacnetUDPAddr(t *testing.T) {
	// Host only — default port applies.
	addr, err := resolveBacnetUDPAddr("192.168.1.50", 47808)
	require.NoError(t, err)
	assert.Equal(t, "192.168.1.50", addr.IP.String())
	assert.Equal(t, 47808, addr.Port)

	// Host:port — explicit port wins.
	addr, err = resolveBacnetUDPAddr("10.0.0.5:47809", 47808)
	require.NoError(t, err)
	assert.Equal(t, "10.0.0.5", addr.IP.String())
	assert.Equal(t, 47809, addr.Port)

	// Invalid host.
	_, err = resolveBacnetUDPAddr("not-an-ip", 47808)
	require.Error(t, err)
}

func TestExtractProtocolSpecificOptions_RemoteAddress(t *testing.T) {
	opts := []options.WithDiscoveryOption{
		options.WithDiscoveryOptionProtocolSpecific("remote-address", "192.168.1.50"),
		options.WithDiscoveryOptionProtocolSpecific("bacnet-port", 47808),
	}
	specific, err := extractProtocolSpecificOptions(opts)
	require.NoError(t, err)
	assert.Equal(t, "192.168.1.50", specific.remoteAddress)
}
