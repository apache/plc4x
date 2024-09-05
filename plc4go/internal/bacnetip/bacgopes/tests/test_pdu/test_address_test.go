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

package test_pdu

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// Assert that the type, network, length, and address are what
//
//	they should be.  Note that the address parameter is a hex string
//	that will be converted to bytes for comparison.
//
//	:param addr: the address to match
//	:param t: the address type
//	:param n: the network number
//	:param l: the address length
//	:param a: the address expressed as hex bytes
func matchAddress(_t *testing.T, addr *bacgopes.Address, t bacgopes.AddressType, n *uint16, l *uint8, a string) {
	_t.Helper()
	assert.Equal(_t, addr.AddrType, t)
	assert.Equal(_t, addr.AddrNet, n)
	assert.Equal(_t, addr.AddrLen, l)
	if a == "" {
		assert.Nil(_t, addr.AddrAddress)
	} else {
		decodeString, err := hex.DecodeString(a)
		require.NoError(_t, err)
		assert.Equal(_t, addr.AddrAddress, decodeString)
	}
}

func init() { // TODO: maybe put in a setupsuite
	bacgopes.Settings.RouteAware = true
}

func TestAddress(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// null address
	testAddr, err := bacgopes.NewAddress(testingLogger)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.NULL_ADDRESS, nil, nil, "")
}

func TestAddressInt(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test integer local station
	testAddr, err := bacgopes.NewAddress(testingLogger, 1)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "01")
	assert.Equal(t, "1", testAddr.String())

	testAddr, err = bacgopes.NewAddress(testingLogger, 254)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "fe")
	assert.Equal(t, "254", testAddr.String())

	// Test bad integer
	_, err = bacgopes.NewAddress(testingLogger, -1)
	assert.Error(t, err)

	_, err = bacgopes.NewAddress(testingLogger, 256)
	assert.Error(t, err)
}

func TestAddressIpv4Str(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test IPv4 local station address
	testAddr, err := bacgopes.NewAddress(testingLogger, "1.2.3.4")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(6), "01020304BAC0")
	assert.Equal(t, "1.2.3.4", testAddr.String())

	// test IPv4 local station address with non-standard port
	testAddr, err = bacgopes.NewAddress(testingLogger, "1.2.3.4:47809")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(6), "01020304BAC1")
	assert.Equal(t, "1.2.3.4:47809", testAddr.String())

	// test IPv4 local station address with unrecognized port
	testAddr, err = bacgopes.NewAddress(testingLogger, "1.2.3.4:47999")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(6), "01020304bb7f")
	assert.Equal(t, "0x01020304bb7f", testAddr.String())
}

func TestAddressEthStr(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test IPv4 local station address
	testAddr, err := bacgopes.NewAddress(testingLogger, "01:02:03:04:05:06")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(6), "010203040506")
	assert.Equal(t, "0x010203040506", testAddr.String())
}

func TestAddressLocalStationStr(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test integer local station
	testAddr, err := bacgopes.NewAddress(testingLogger, "1")
	require.NoError(t, err)
	assert.Equal(t, "1", testAddr.String())

	testAddr, err = bacgopes.NewAddress(testingLogger, "254")
	require.NoError(t, err)
	assert.Equal(t, "254", testAddr.String())

	// Test bad integer
	_, err = bacgopes.NewAddress(testingLogger, 256)
	assert.Error(t, err)

	// test modern hex string
	testAddr, err = bacgopes.NewAddress(testingLogger, "0x01")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "01")
	assert.Equal(t, "1", testAddr.String())

	testAddr, err = bacgopes.NewAddress(testingLogger, "0x0102")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(2), "0102")
	assert.Equal(t, "0x0102", testAddr.String())

	// test old school hex string
	testAddr, err = bacgopes.NewAddress(testingLogger, "X'01'")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "01")
	assert.Equal(t, "1", testAddr.String())

	testAddr, err = bacgopes.NewAddress(testingLogger, "X'0102'")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(2), "0102")
	assert.Equal(t, "0x0102", testAddr.String())
}

func TestAddressLocalBroadcastStr(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test IPv4 local station address
	testAddr, err := bacgopes.NewAddress(testingLogger, "*")
	require.NoError(t, err)
	assert.Equal(t, "*", testAddr.String())
}

func TestAddressRemoteBroadcastStr(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test IPv4 local station address
	testAddr, err := bacgopes.NewAddress(testingLogger, "1:*")
	require.NoError(t, err)
	assert.Equal(t, "1:*", testAddr.String())
}

func TestAddressRemoteStationStr(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test IPv4 local station address
	testAddr, err := bacgopes.NewAddress(testingLogger, "1:2")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "02")
	assert.Equal(t, "1:2", testAddr.String())

	testAddr, err = bacgopes.NewAddress(testingLogger, "1:254")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "fe")
	assert.Equal(t, "1:254", testAddr.String())

	// test bad network and mode
	_, err = bacgopes.NewAddress(testingLogger, "65536:2")
	assert.Error(t, err)
	_, err = bacgopes.NewAddress(testingLogger, "1:256")
	assert.Error(t, err)

	// test moder hex string
	testAddr, err = bacgopes.NewAddress(testingLogger, "1:0x02")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "02")
	assert.Equal(t, "1:2", testAddr.String())

	// test bad network
	_, err = bacgopes.NewAddress(testingLogger, "65536:0x02")
	assert.Error(t, err)

	testAddr, err = bacgopes.NewAddress(testingLogger, "1:0x0203")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(2), "0203")
	assert.Equal(t, "1:0x0203", testAddr.String())

	// test old school hex
	testAddr, err = bacgopes.NewAddress(testingLogger, "1:X'02'")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "02")
	assert.Equal(t, "1:2", testAddr.String())

	testAddr, err = bacgopes.NewAddress(testingLogger, "1:X'0203'")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(2), "0203")
	assert.Equal(t, "1:0x0203", testAddr.String())

	_, err = bacgopes.NewAddress(testingLogger, "65536:X'02'")
	assert.Error(t, err)
}

func TestAddressGlobalBroadcastStr(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// test IPv4 local station address
	testAddr, err := bacgopes.NewAddress(testingLogger, "*:*")
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.GLOBAL_BROADCAST_ADDRESS, nil, nil, "")
	assert.Equal(t, "*:*", testAddr.String())
}

func TestLocalStation(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// one Parameter
	_, err := bacgopes.NewLocalStation(testingLogger, nil, nil)
	require.Error(t, err)

	// test integer
	testAddr, err := bacgopes.NewLocalStation(testingLogger, 1, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "01")
	assert.Equal(t, "1", testAddr.String())

	testAddr, err = bacgopes.NewLocalStation(testingLogger, 254, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "fe")
	assert.Equal(t, "254", testAddr.String())

	// Test bad integer
	_, err = bacgopes.NewLocalStation(testingLogger, -1, nil)
	require.Error(t, err)
	_, err = bacgopes.NewLocalStation(testingLogger, 256, nil)
	require.Error(t, err)

	// Test bytes
	xtob, err := bacgopes.Xtob("01")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewLocalStation(testingLogger, xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "01")
	assert.Equal(t, "1", testAddr.String())
	xtob, err = bacgopes.Xtob("fe")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewLocalStation(testingLogger, xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(1), "fe")
	assert.Equal(t, "254", testAddr.String())

	// multi-byte strings are hex encoded
	xtob, err = bacgopes.Xtob("0102")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewLocalStation(testingLogger, xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(2), "0102")
	assert.Equal(t, "0x0102", testAddr.String())

	xtob, err = bacgopes.Xtob("010203")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewLocalStation(testingLogger, xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(3), "010203")
	assert.Equal(t, "0x010203", testAddr.String())

	// match and IP address
	xtob, err = bacgopes.Xtob("01020304bac0")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewLocalStation(testingLogger, xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.LOCAL_STATION_ADDRESS, nil, l(6), "01020304bac0")
	assert.Equal(t, "1.2.3.4", testAddr.String())
}

func TestRemoteStation(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	// two Parameters, correct types
	_, err := bacgopes.NewRemoteStation(testingLogger, nil, nil, nil)
	require.Error(t, err)

	// test bad network
	_, err = bacgopes.NewRemoteStation(testingLogger, nil, -11, nil)
	require.Error(t, err)
}

func TestRemoteStationInts(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	net := func(i uint16) *uint16 {
		return &i
	}

	// testInteger
	testAddr, err := bacgopes.NewRemoteStation(testingLogger, net(1), 1, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "01")
	assert.Equal(t, "1:1", testAddr.String())

	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), 254, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "fe")
	assert.Equal(t, "1:254", testAddr.String())

	// test station address
	_, err = bacgopes.NewRemoteStation(testingLogger, nil, -1, nil)
	require.Error(t, err)
	_, err = bacgopes.NewRemoteStation(testingLogger, nil, 256, nil)
	require.Error(t, err)
}

func TestRemoteStationBytes(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	net := func(i uint16) *uint16 {
		return &i
	}

	// multi-byte strings are hex encoded
	xtob, err := bacgopes.Xtob("0102")
	require.NoError(t, err)
	testAddr, err := bacgopes.NewRemoteStation(testingLogger, net(1), xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(2), "0102")
	assert.Equal(t, "1:0x0102", testAddr.String())

	xtob, err = bacgopes.Xtob("010203")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(3), "010203")
	assert.Equal(t, "1:0x010203", testAddr.String())

	// match with IPv4 address
	xtob, err = bacgopes.Xtob("01020304bac0")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(6), "01020304bac0")
	assert.Equal(t, "1:1.2.3.4", testAddr.String())

	xtob, err = bacgopes.Xtob("01020304bac1")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, nil)
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(6), "01020304bac1")
	assert.Equal(t, "1:1.2.3.4:47809", testAddr.String())
}

func TestRemoteStationIntsRouted(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	net := func(i uint16) *uint16 {
		return &i
	}

	Address := func(a string) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}

	// testInteger
	testAddr, err := bacgopes.NewRemoteStation(testingLogger, net(1), 1, Address("1.2.3.4"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "01")
	assert.Equal(t, "1:1@1.2.3.4", testAddr.String())

	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), 254, Address("1.2.3.4"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "fe")
	assert.Equal(t, "1:254@1.2.3.4", testAddr.String())

	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), 254, Address("1.2.3.4:47809"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "fe")
	assert.Equal(t, "1:254@1.2.3.4:47809", testAddr.String())

	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), 254, Address("0x01020304BAC0"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "fe")
	assert.Equal(t, "1:254@1.2.3.4", testAddr.String())

	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), 254, Address("0x01020304BAC1"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(1), "fe")
	assert.Equal(t, "1:254@1.2.3.4:47809", testAddr.String())

	// test station address
	_, err = bacgopes.NewRemoteStation(testingLogger, nil, -1, nil)
	require.Error(t, err)
	_, err = bacgopes.NewRemoteStation(testingLogger, nil, 256, nil)
	require.Error(t, err)
}

func TestRemoteStationBytesRouted(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)

	net := func(i uint16) *uint16 {
		return &i
	}

	Address := func(a string) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}

	// multi-byte strings are hex encoded
	xtob, err := bacgopes.Xtob("0102")
	require.NoError(t, err)
	testAddr, err := bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("1.2.3.4"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(2), "0102")
	assert.Equal(t, "1:0x0102@1.2.3.4", testAddr.String())

	xtob, err = bacgopes.Xtob("010203")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("1.2.3.4"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(3), "010203")
	assert.Equal(t, "1:0x010203@1.2.3.4", testAddr.String())

	xtob, err = bacgopes.Xtob("010203")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("1.2.3.4:47809"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(3), "010203")
	assert.Equal(t, "1:0x010203@1.2.3.4:47809", testAddr.String())

	xtob, err = bacgopes.Xtob("010203")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("0x01020304BAC0"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(3), "010203")
	assert.Equal(t, "1:0x010203@1.2.3.4", testAddr.String())

	xtob, err = bacgopes.Xtob("010203")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("0x01020304BAC1"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(3), "010203")
	assert.Equal(t, "1:0x010203@1.2.3.4:47809", testAddr.String())

	// match with an IPv4 address
	xtob, err = bacgopes.Xtob("01020304bac0")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("1.2.3.4"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(6), "01020304bac0")
	assert.Equal(t, "1:1.2.3.4@1.2.3.4", testAddr.String())

	xtob, err = bacgopes.Xtob("01020304bac0")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("1.2.3.4:47809"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(6), "01020304bac0")
	assert.Equal(t, "1:1.2.3.4@1.2.3.4:47809", testAddr.String())

	xtob, err = bacgopes.Xtob("01020304bac0")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("0x01020304BAC0"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(6), "01020304bac0")
	assert.Equal(t, "1:1.2.3.4@1.2.3.4", testAddr.String())

	xtob, err = bacgopes.Xtob("01020304bac0")
	require.NoError(t, err)
	testAddr, err = bacgopes.NewRemoteStation(testingLogger, net(1), xtob, Address("0x01020304BAC1"))
	require.NoError(t, err)
	matchAddress(t, testAddr, bacgopes.REMOTE_STATION_ADDRESS, n(1), l(6), "01020304bac0")
	assert.Equal(t, "1:1.2.3.4@1.2.3.4:47809", testAddr.String())
}

func TestLocalBroadcast(t *testing.T) {
	testAddr := bacgopes.NewLocalBroadcast(nil)
	matchAddress(t, testAddr, bacgopes.LOCAL_BROADCAST_ADDRESS, nil, nil, "")
	assert.Equal(t, "*", testAddr.String())
}

func TestLocalBroadcastRouted(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	Address := func(a string) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}

	testAddr := bacgopes.NewLocalBroadcast(Address("1.2.3.4"))
	matchAddress(t, testAddr, bacgopes.LOCAL_BROADCAST_ADDRESS, nil, nil, "")
	assert.Equal(t, "*@1.2.3.4", testAddr.String())
}

func TestRemoteBroadcast(t *testing.T) {
	testAddr := bacgopes.NewRemoteBroadcast(1, nil)
	matchAddress(t, testAddr, bacgopes.REMOTE_BROADCAST_ADDRESS, n(1), nil, "")
	assert.Equal(t, "1:*", testAddr.String())
}

func TestRemoteBroadcastRouted(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	Address := func(a string) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}

	testAddr := bacgopes.NewRemoteBroadcast(1, Address("1.2.3.4"))
	matchAddress(t, testAddr, bacgopes.REMOTE_BROADCAST_ADDRESS, n(1), nil, "")
	assert.Equal(t, "1:*@1.2.3.4", testAddr.String())
}

func TestGlobalBroadcast(t *testing.T) {
	testAddr := bacgopes.NewGlobalBroadcast(nil)
	matchAddress(t, testAddr, bacgopes.GLOBAL_BROADCAST_ADDRESS, nil, nil, "")
	assert.Equal(t, "*:*", testAddr.String())
}

func TestGlobalBroadcastRouted(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	Address := func(a string) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}

	testAddr := bacgopes.NewGlobalBroadcast(Address("1.2.3.4"))
	matchAddress(t, testAddr, bacgopes.GLOBAL_BROADCAST_ADDRESS, nil, nil, "")
	assert.Equal(t, "*:*@1.2.3.4", testAddr.String())

	testAddr = bacgopes.NewGlobalBroadcast(Address("1.2.3.4:47809"))
	matchAddress(t, testAddr, bacgopes.GLOBAL_BROADCAST_ADDRESS, nil, nil, "")
	assert.Equal(t, "*:*@1.2.3.4:47809", testAddr.String())
}

func TestAddressEquality(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	Address := func(a any) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}
	LocalStation := func(addr any) *bacgopes.Address {
		station, err := bacgopes.NewLocalStation(testingLogger, addr, nil)
		require.NoError(t, err)
		return station
	}
	RemoteStation := func(net uint16, addr any) *bacgopes.Address {
		station, err := bacgopes.NewRemoteStation(testingLogger, &net, addr, nil)
		require.NoError(t, err)
		return station
	}
	LocalBroadcast := func() *bacgopes.Address {
		broadcast := bacgopes.NewLocalBroadcast(nil)
		return broadcast
	}
	RemoteBroadcast := func(net uint16) *bacgopes.Address {
		broadcast := bacgopes.NewRemoteBroadcast(net, nil)
		return broadcast
	}
	GlobalBroadcast := func() *bacgopes.Address {
		broadcast := bacgopes.NewGlobalBroadcast(nil)
		return broadcast
	}

	assert.True(t, Address(1).Equals(LocalStation(1)))
	assert.True(t, Address("2").Equals(LocalStation(2)))
	assert.True(t, Address("*").Equals(LocalBroadcast()))
	assert.True(t, Address("3:4").Equals(RemoteStation(3, 4)))
	assert.True(t, Address("5:*").Equals(RemoteBroadcast(5)))
	assert.True(t, Address("*:*").Equals(GlobalBroadcast()))

}

func TestAddressEqualityRouted(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	Address := func(a any) *bacgopes.Address {
		address, err := bacgopes.NewAddress(testingLogger, a)
		require.NoError(t, err)
		return address
	}
	RemoteStation := func(net uint16, addr any, route *bacgopes.Address) *bacgopes.Address {
		station, err := bacgopes.NewRemoteStation(testingLogger, &net, addr, route)
		require.NoError(t, err)
		return station
	}
	RemoteBroadcast := func(net uint16, route *bacgopes.Address) *bacgopes.Address {
		broadcast := bacgopes.NewRemoteBroadcast(net, route)
		return broadcast
	}
	GlobalBroadcast := func(route *bacgopes.Address) *bacgopes.Address {
		broadcast := bacgopes.NewGlobalBroadcast(route)
		return broadcast
	}

	assert.True(t, Address("3:4@6.7.8.9").Equals(RemoteStation(3, 4, Address("6.7.8.9"))))
	assert.True(t, Address("3:4@0x06070809BAC0").Equals(RemoteStation(3, 4, Address("6.7.8.9"))))

	assert.True(t, Address("3:4@6.7.8.9:47809").Equals(RemoteStation(3, 4, Address("6.7.8.9:47809"))))
	assert.True(t, Address("3:4@0x06070809BAC1").Equals(RemoteStation(3, 4, Address("6.7.8.9:47809"))))

	assert.True(t, Address("5:*@6.7.8.9").Equals(RemoteBroadcast(5, Address("6.7.8.9"))))
	assert.True(t, Address("5:*@0x06070809BAC0").Equals(RemoteBroadcast(5, Address("6.7.8.9"))))

	assert.True(t, Address("5:*@6.7.8.9:47809").Equals(RemoteBroadcast(5, Address("6.7.8.9:47809"))))
	assert.True(t, Address("5:*@0x06070809BAC1").Equals(RemoteBroadcast(5, Address("6.7.8.9:47809"))))

	assert.True(t, Address("*:*@6.7.8.9").Equals(GlobalBroadcast(Address("6.7.8.9"))))
	assert.True(t, Address("*:*@0x06070809BAC0").Equals(GlobalBroadcast(Address("6.7.8.9"))))

	assert.True(t, Address("*:*@6.7.8.9:47809").Equals(GlobalBroadcast(Address("6.7.8.9:47809"))))
	assert.True(t, Address("*:*@0x06070809BAC1").Equals(GlobalBroadcast(Address("6.7.8.9:47809"))))
}

func n(n uint16) *uint16 {
	return &n
}

func l(l uint8) *uint8 {
	return &l
}
