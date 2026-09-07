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

package modbus

import (
	"encoding/binary"
)

// ByteOrder is the layout a device uses for a value that spans more than one byte. Ported from
// plc4j's ModbusByteOrder: nothing in the protocol announces it, so it has to be configured, either
// for the whole connection (the default-payload-byte-order option) or for a single tag (the
// byte-order tag config).
type ByteOrder uint8

const (
	// BigEndianOrder is the order the specification asks for: [1, 2, 3, 4] and
	// [1, 2, 3, 4, 5, 6, 7, 8].
	BigEndianOrder ByteOrder = iota
	// LittleEndianOrder reverses the whole value: [4, 3, 2, 1] and [8, 7, 6, 5, 4, 3, 2, 1].
	LittleEndianOrder
	// BigEndianByteSwapOrder swaps the two bytes inside every register, leaving the registers
	// themselves in order: [2, 1, 4, 3] and [2, 1, 4, 3, 6, 5, 8, 7].
	BigEndianByteSwapOrder
	// LittleEndianByteSwapOrder reverses the registers but not the bytes inside them:
	// [3, 4, 1, 2] and [7, 8, 5, 6, 3, 4, 1, 2].
	LittleEndianByteSwapOrder
)

// byteOrderNames spells the modes the way plc4j's ModbusByteOrder enum does, which is also how they
// are written in a connection string or a tag config.
var byteOrderNames = map[ByteOrder]string{
	BigEndianOrder:            "BIG_ENDIAN",
	LittleEndianOrder:         "LITTLE_ENDIAN",
	BigEndianByteSwapOrder:    "BIG_ENDIAN_BYTE_SWAP",
	LittleEndianByteSwapOrder: "LITTLE_ENDIAN_BYTE_SWAP",
}

func (b ByteOrder) String() string {
	if name, ok := byteOrderNames[b]; ok {
		return name
	}
	return "UNKNOWN"
}

// ByteOrderByName resolves one of the four names above, and reports whether it knew it.
func ByteOrderByName(name string) (ByteOrder, bool) {
	for byteOrder, candidate := range byteOrderNames {
		if candidate == name {
			return byteOrder, true
		}
	}
	return BigEndianOrder, false
}

// isBigEndian says whether the bytes of one value run from the most significant to the least
// significant one. The two byte-swap modes only shuffle bytes around afterwards, so they keep the
// endianness of the mode they are named after.
func (b ByteOrder) isBigEndian() bool {
	return b == BigEndianOrder || b == BigEndianByteSwapOrder
}

// bufferByteOrder is what the read and write buffers have to be told to make a multi-byte value
// come out in this order.
func (b ByteOrder) bufferByteOrder() binary.ByteOrder {
	if b.isBigEndian() {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

// swapsBytes says whether the two bytes of every register still have to be exchanged.
func (b ByteOrder) swapsBytes() bool {
	return b == BigEndianByteSwapOrder || b == LittleEndianByteSwapOrder
}

// byteSwap exchanges the two bytes of every register. A trailing odd byte is not part of a whole
// register and is left where it is (plc4j ModbusTcpConnection.byteSwap).
func byteSwap(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	for i := 0; i+1 < len(out); i += 2 {
		out[i], out[i+1] = out[i+1], out[i]
	}
	return out
}
