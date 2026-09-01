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
	"testing"

	"github.com/stretchr/testify/assert"
)

// The names are the ones plc4j's ModbusByteOrder enum uses, because that is how they are written in
// a connection string or a tag config.
func TestByteOrderByName(t *testing.T) {
	for _, test := range []struct {
		name      string
		byteOrder ByteOrder
	}{
		{"BIG_ENDIAN", BigEndianOrder},
		{"LITTLE_ENDIAN", LittleEndianOrder},
		{"BIG_ENDIAN_BYTE_SWAP", BigEndianByteSwapOrder},
		{"LITTLE_ENDIAN_BYTE_SWAP", LittleEndianByteSwapOrder},
	} {
		t.Run(test.name, func(t *testing.T) {
			byteOrder, ok := ByteOrderByName(test.name)
			assert.True(t, ok)
			assert.Equal(t, test.byteOrder, byteOrder)
			assert.Equal(t, test.name, byteOrder.String())
		})
	}
}

func TestByteOrderByName_rejectsAnUnknownName(t *testing.T) {
	for _, name := range []string{"", "big_endian", "BIGENDIAN", "MIDDLE_ENDIAN"} {
		_, ok := ByteOrderByName(name)
		assert.False(t, ok, "%q must not be a known byte order", name)
	}
}

// A byte-swap mode keeps the endianness of the mode it is named after; the swap happens on top of
// it.
func TestByteOrderEndianness(t *testing.T) {
	assert.True(t, BigEndianOrder.isBigEndian())
	assert.True(t, BigEndianByteSwapOrder.isBigEndian())
	assert.False(t, LittleEndianOrder.isBigEndian())
	assert.False(t, LittleEndianByteSwapOrder.isBigEndian())

	assert.Equal(t, binary.BigEndian, BigEndianByteSwapOrder.bufferByteOrder())
	assert.Equal(t, binary.LittleEndian, LittleEndianByteSwapOrder.bufferByteOrder())

	assert.False(t, BigEndianOrder.swapsBytes())
	assert.False(t, LittleEndianOrder.swapsBytes())
	assert.True(t, BigEndianByteSwapOrder.swapsBytes())
	assert.True(t, LittleEndianByteSwapOrder.swapsBytes())
}
