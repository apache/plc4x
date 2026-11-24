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

// Connectionless DCE/RPC fragmentation (IEC 61158-6-10 and DCE 1.1 RPC). A packet whose body does
// not fit one datagram goes out as fragments: copies of the 80 byte header with the fragment flag,
// a fragment number counting from 0 (also in the serial number bytes, as p-net writes it) and the
// body length of the piece, the last one flagged last fragment. The receiver acknowledges each
// fragment but the last with a header-only FACK unless the fragment says no FACK is wanted.
//
// The generated model reads and writes whole packets only (its fragment flag and number are
// constants), so fragments are cut from, and put back together into, a whole serialized packet at
// the byte level, before the model sees it.

import (
	"encoding/binary"
)

const (
	rpcHeaderLength = 80

	rpcOffsetPacketType  = 1
	rpcOffsetFlags       = 2
	rpcOffsetDrep        = 4
	rpcOffsetSerialHigh  = 7
	rpcOffsetActivity    = 40
	rpcOffsetSequence    = 64
	rpcOffsetBodyLength  = 74
	rpcOffsetFragmentNum = 76
	rpcOffsetSerialLow   = 79

	rpcFlagLastFragment = 0x02
	rpcFlagFragment     = 0x04
	rpcFlagNoFack       = 0x08

	rpcPacketTypeFack = 0x09

	// defaultRPCFragmentBody is the body each fragment carries: p-net, like other devices, takes
	// at most 1440 bytes of UDP payload per datagram, header included.
	defaultRPCFragmentBody = 1440 - rpcHeaderLength
	// maxRPCReassembly bounds a reassembled packet: its body length must fit the header's 16 bits.
	maxRPCReassembly = 0xFFFF
)

// rpcHeader is what the fragment layer reads from a packet's header.
type rpcHeader struct {
	packetType byte
	flags      byte
	order      binary.ByteOrder
	activity   [16]byte
	sequence   uint32
	fragment   uint16
	bodyLength int
}

// parseRPCHeader reads the header of a datagram; ok is false when it is too short to be one or
// announces more body than follows.
func parseRPCHeader(data []byte) (rpcHeader, bool) {
	if len(data) < rpcHeaderLength || data[0] != 4 {
		return rpcHeader{}, false
	}
	h := rpcHeader{packetType: data[rpcOffsetPacketType], flags: data[rpcOffsetFlags], order: binary.BigEndian}
	if data[rpcOffsetDrep]&0xF0 != 0 {
		h.order = binary.LittleEndian
	}
	// The activity UUID's first three fields are integers in the packet's byte order; the key is
	// big-endian, as activityKey makes it.
	copy(h.activity[:], data[rpcOffsetActivity:rpcOffsetActivity+16])
	if h.order == binary.LittleEndian {
		reverse(h.activity[0:4])
		reverse(h.activity[4:6])
		reverse(h.activity[6:8])
	}
	h.sequence = h.order.Uint32(data[rpcOffsetSequence:])
	h.fragment = h.order.Uint16(data[rpcOffsetFragmentNum:])
	h.bodyLength = int(h.order.Uint16(data[rpcOffsetBodyLength:]))
	if len(data) < rpcHeaderLength+h.bodyLength {
		return rpcHeader{}, false
	}
	return h, true
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func (h rpcHeader) key() callKey {
	return callKey{h.activity, h.sequence}
}

func (h rpcHeader) isFragment() bool {
	return h.flags&rpcFlagFragment != 0
}

// fragmentPacket cuts a serialized packet into datagrams of at most bodySize body bytes each. A
// packet that fits is returned as it is. wantFack asks the receiver to acknowledge every fragment
// but the last.
func fragmentPacket(data []byte, bodySize int, wantFack bool) [][]byte {
	body := data[rpcHeaderLength:]
	if len(body) <= bodySize {
		return [][]byte{data}
	}
	h, _ := parseRPCHeader(data)
	order := h.order
	if order == nil {
		order = binary.BigEndian
	}
	var fragments [][]byte
	for number := 0; len(body) > 0; number++ {
		n := min(bodySize, len(body))
		fragment := make([]byte, rpcHeaderLength+n)
		copy(fragment, data[:rpcHeaderLength])
		copy(fragment[rpcHeaderLength:], body[:n])
		body = body[n:]
		flags := fragment[rpcOffsetFlags]&^(rpcFlagLastFragment|rpcFlagNoFack) | rpcFlagFragment
		if len(body) == 0 {
			flags |= rpcFlagLastFragment | rpcFlagNoFack
		} else if !wantFack {
			flags |= rpcFlagNoFack
		}
		fragment[rpcOffsetFlags] = flags
		order.PutUint16(fragment[rpcOffsetBodyLength:], uint16(n))
		order.PutUint16(fragment[rpcOffsetFragmentNum:], uint16(number))
		fragment[rpcOffsetSerialHigh] = byte(number >> 8)
		fragment[rpcOffsetSerialLow] = byte(number)
		fragments = append(fragments, fragment)
	}
	return fragments
}

// fackFor is the header-only FACK acknowledging a fragment, as p-net sends it.
func fackFor(fragment []byte) []byte {
	fack := make([]byte, rpcHeaderLength)
	copy(fack, fragment[:rpcHeaderLength])
	fack[rpcOffsetPacketType] = rpcPacketTypeFack
	fack[rpcOffsetFlags] = 0
	fack[rpcOffsetBodyLength], fack[rpcOffsetBodyLength+1] = 0, 0
	return fack
}

// reassembly collects the fragments of one packet.
type reassembly struct {
	header []byte
	body   []byte
	next   uint16
}

// add takes the next fragment; done is true with the whole packet once the last has arrived. A
// fragment out of order is dropped (the sender sends it again) and a first fragment starts over.
func (r *reassembly) add(h rpcHeader, fragment []byte) (whole []byte, accepted, done bool) {
	if h.fragment == 0 {
		r.header = append([]byte(nil), fragment[:rpcHeaderLength]...)
		r.body = r.body[:0]
		r.next = 0
	}
	if r.header == nil || h.fragment != r.next || len(r.body)+h.bodyLength > maxRPCReassembly {
		return nil, false, false
	}
	r.body = append(r.body, fragment[rpcHeaderLength:rpcHeaderLength+h.bodyLength]...)
	r.next++
	if h.flags&rpcFlagLastFragment == 0 {
		return nil, true, false
	}
	whole = make([]byte, rpcHeaderLength+len(r.body))
	copy(whole, r.header)
	copy(whole[rpcHeaderLength:], r.body)
	whole[rpcOffsetFlags] &^= rpcFlagFragment | rpcFlagLastFragment | rpcFlagNoFack
	h.order.PutUint16(whole[rpcOffsetBodyLength:], uint16(len(r.body)))
	h.order.PutUint16(whole[rpcOffsetFragmentNum:], 0)
	whole[rpcOffsetSerialHigh] = 0
	whole[rpcOffsetSerialLow] = 0
	return whole, true, true
}
