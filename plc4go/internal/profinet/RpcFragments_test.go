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

package profinet

import (
	"context"
	"encoding/binary"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/profinet/readwrite/model"
)

// manyBlocks is a request body far longer than one datagram: n control blocks of 32 bytes.
func manyBlocks(n int) []readWriteModel.PnIoCm_Block {
	var blocks []readWriteModel.PnIoCm_Block
	for i := range n {
		uuid := make([]byte, 16)
		uuid[15] = byte(i)
		blocks = append(blocks, readWriteModel.NewPnIoCm_Control_Request_ParameterEnd(blockVersionHigh, blockVersionLow, readWriteModel.NewUuid(uuid), uint16(i), controlCommandPrmEnd))
	}
	return blocks
}

func serializedRequest(t *testing.T, encoding readWriteModel.IntegerEncoding, activity readWriteModel.DceRpc_ActivityUuid, blocks []readWriteModel.PnIoCm_Block) []byte {
	packet := readWriteModel.NewDceRpc_Packet(readWriteModel.DceRpc_PacketType_REQUEST, true, false, false,
		encoding, readWriteModel.CharacterEncoding_ASCII, readWriteModel.FloatingPointEncoding_IEEE,
		readWriteModel.NewDceRpc_ObjectUuid(0, 1, 0x0001, 0x0002), readWriteModel.NewDceRpc_InterfaceUuid_DeviceInterface(),
		activity, 0, 7, readWriteModel.DceRpc_Operation_CONNECT, 0, readWriteModel.NewPnIoCm_Packet_Req(argsMaximum, argsMaximum, 0, blocks))
	data, err := packet.Serialize()
	require.NoError(t, err)
	return data
}

func TestFragmentRoundTrip(t *testing.T) {
	activity := newActivity()
	data := serializedRequest(t, readWriteModel.IntegerEncoding_BIG_ENDIAN, activity, manyBlocks(60))
	require.Greater(t, len(data)-rpcHeaderLength, defaultRPCFragmentBody, "the request needs fragments")

	fragments := fragmentPacket(data, 256, true)
	require.Len(t, fragments, (len(data)-rpcHeaderLength+255)/256)
	r := &reassembly{}
	var whole []byte
	for i, fragment := range fragments {
		assert.LessOrEqual(t, len(fragment), rpcHeaderLength+256)
		h, ok := parseRPCHeader(fragment)
		require.True(t, ok)
		assert.True(t, h.isFragment())
		assert.Equal(t, uint16(i), h.fragment)
		assert.Equal(t, byte(i), fragment[rpcOffsetSerialLow])
		assert.Equal(t, callKey{activityKey(activity), 7}, h.key())
		last := i == len(fragments)-1
		assert.Equal(t, last, h.flags&rpcFlagLastFragment != 0)
		assert.Equal(t, last, h.flags&rpcFlagNoFack != 0, "a FACK is wanted for every fragment but the last")
		var accepted, done bool
		whole, accepted, done = r.add(h, fragment)
		assert.True(t, accepted)
		assert.Equal(t, last, done)
	}
	assert.Equal(t, data, whole, "reassembly gives back the packet byte for byte")
	_, err := readWriteModel.DceRpc_PacketParse(context.Background(), whole)
	require.NoError(t, err)

	// Out of order: dropped until the expected fragment comes.
	r = &reassembly{}
	h0, _ := parseRPCHeader(fragments[0])
	h2, _ := parseRPCHeader(fragments[2])
	_, _, _ = r.add(h0, fragments[0])
	_, accepted, _ := r.add(h2, fragments[2])
	assert.False(t, accepted)

	assert.Equal(t, [][]byte{data}, fragmentPacket(data, len(data), true), "a packet that fits is not cut")

	fack := fackFor(fragments[1])
	h, ok := parseRPCHeader(fack)
	require.True(t, ok)
	assert.Equal(t, byte(rpcPacketTypeFack), h.packetType)
	assert.Equal(t, uint16(1), h.fragment)
	assert.Equal(t, 0, h.bodyLength)
}

func TestRPCHeaderLittleEndian(t *testing.T) {
	activity := newActivity()
	data := serializedRequest(t, readWriteModel.IntegerEncoding_LITTLE_ENDIAN, activity, manyBlocks(2))
	h, ok := parseRPCHeader(data)
	require.True(t, ok)
	assert.Equal(t, binary.LittleEndian, h.order)
	assert.Equal(t, callKey{activityKey(activity), 7}, h.key())
	assert.Equal(t, len(data)-rpcHeaderLength, h.bodyLength)
}

// TestFragmentedCall sends a request needing fragments to another endpoint standing in for the
// device, whose response needs fragments too.
func TestFragmentedCall(t *testing.T) {
	device := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 2), Port: pnioCMPort}
	var served atomic.Int32
	server, err := listenRPC(device, func(operation readWriteModel.DceRpc_Operation, request readWriteModel.PnIoCm_Packet_Req) (readWriteModel.PnIoCm_Packet_Res, error) {
		served.Store(int32(len(request.GetBlocks())))
		var blocks []readWriteModel.PnIoCm_Block
		for _, block := range request.GetBlocks() {
			b := block.(readWriteModel.PnIoCm_Control_Request_ParameterEnd)
			blocks = append(blocks, readWriteModel.NewPnIoCm_Control_Response_ParameterEnd(blockVersionHigh, blockVersionLow, b.GetArUuid(), b.GetSessionKey(), controlCommandDone))
		}
		return readWriteModel.NewPnIoCm_Packet_Res(0, 0, 0, 0, argsMaximum, 0, blocks), nil
	}, 300, zerolog.Nop())
	if err != nil {
		t.Skipf("no second loopback address here: %v", err)
	}
	defer server.close()
	client, err := listenRPC(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: pnioCMPort}, nil, 256, zerolog.Nop())
	require.NoError(t, err)
	defer client.close()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	started := time.Now()
	response, err := client.call(ctx, device, readWriteModel.NewDceRpc_ObjectUuid(0, 1, 1, 2), newActivity(), 0,
		readWriteModel.DceRpc_Operation_CONNECT, readWriteModel.NewPnIoCm_Packet_Req(argsMaximum, argsMaximum, 0, manyBlocks(60)))
	require.NoError(t, err)
	assert.Equal(t, int32(60), served.Load(), "the device got every block of the fragmented request")
	assert.Len(t, response.GetBlocks(), 60, "the fragmented response came back whole")
	assert.Less(t, time.Since(started), 2*time.Second, "the device's FACKs pace the fragments, no waiting for timeouts")
}
