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

package s7

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

func TestParsePlcValue(t *testing.T) {
	ctx := t.Context()
	anyController := readWriteModel.ControllerType_ANY

	t.Run("single BOOL", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 1, readWriteModel.TransportSize_BOOL)
		value, err := parsePlcValue(ctx, tag, []byte{0x01}, anyController)
		require.NoError(t, err)
		assert.True(t, value.GetBool())
	})
	t.Run("single INT", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 1, readWriteModel.TransportSize_INT)
		value, err := parsePlcValue(ctx, tag, []byte{0x00, 0x04}, anyController)
		require.NoError(t, err)
		assert.Equal(t, int16(4), value.GetInt16())
	})
	t.Run("single REAL", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 1, readWriteModel.TransportSize_REAL)
		value, err := parsePlcValue(ctx, tag, []byte{0x41, 0x28, 0x00, 0x00}, anyController)
		require.NoError(t, err)
		assert.Equal(t, float32(10.5), value.GetFloat32())
	})
	t.Run("single STRING with fixed length", func(t *testing.T) {
		tag := NewStringTag(readWriteModel.MemoryArea_DATA_BLOCKS, 1, 0, 0, 1, 5, readWriteModel.TransportSize_STRING)
		// maxLen=5, curLen=3, "ABC" + 2 bytes padding
		value, err := parsePlcValue(ctx, tag, []byte{0x05, 0x03, 'A', 'B', 'C', 0x00, 0x00}, anyController)
		require.NoError(t, err)
		assert.Equal(t, "ABC", value.GetString())
	})
	t.Run("INT array", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 2, readWriteModel.TransportSize_INT)
		value, err := parsePlcValue(ctx, tag, []byte{0x00, 0x01, 0x00, 0x02}, anyController)
		require.NoError(t, err)
		require.True(t, value.IsList())
		list := value.GetList()
		require.Len(t, list, 2)
		assert.Equal(t, int16(1), list[0].GetInt16())
		assert.Equal(t, int16(2), list[1].GetInt16())
	})
	t.Run("BYTE array stays raw", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 4, readWriteModel.TransportSize_BYTE)
		value, err := parsePlcValue(ctx, tag, []byte{0xDE, 0xAD, 0xBE, 0xEF}, anyController)
		require.NoError(t, err)
		assert.Equal(t, []byte{0xDE, 0xAD, 0xBE, 0xEF}, value.GetRaw())
	})
	t.Run("BOOL array unpacks bits LSB first", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 10, readWriteModel.TransportSize_BOOL)
		value, err := parsePlcValue(ctx, tag, []byte{0b00000011, 0b00000010}, anyController)
		require.NoError(t, err)
		require.True(t, value.IsList())
		list := value.GetList()
		require.Len(t, list, 10)
		assert.True(t, list[0].GetBool())
		assert.True(t, list[1].GetBool())
		assert.False(t, list[2].GetBool())
		assert.False(t, list[8].GetBool())
		assert.True(t, list[9].GetBool())
	})
}

func TestEncodeS7Address(t *testing.T) {
	t.Run("S5TIME is requested as byte array", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 1, readWriteModel.TransportSize_S5TIME)
		address, err := encodeS7Address(tag)
		require.NoError(t, err)
		addressAny := address.(readWriteModel.S7AddressAny)
		assert.Equal(t, readWriteModel.TransportSize_BYTE, addressAny.GetTransportSize())
		assert.Equal(t, uint16(2), addressAny.GetNumberOfElements())
	})
	t.Run("TIME is requested as byte array", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 1, readWriteModel.TransportSize_TIME)
		address, err := encodeS7Address(tag)
		require.NoError(t, err)
		addressAny := address.(readWriteModel.S7AddressAny)
		assert.Equal(t, readWriteModel.TransportSize_BYTE, addressAny.GetTransportSize())
		assert.Equal(t, uint16(4), addressAny.GetNumberOfElements())
	})
	t.Run("BOOL array is requested packed", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 10, readWriteModel.TransportSize_BOOL)
		address, err := encodeS7Address(tag)
		require.NoError(t, err)
		addressAny := address.(readWriteModel.S7AddressAny)
		assert.Equal(t, readWriteModel.TransportSize_BYTE, addressAny.GetTransportSize())
		assert.Equal(t, uint16(2), addressAny.GetNumberOfElements())
	})
	t.Run("STRING converts to CHAR elements", func(t *testing.T) {
		tag := NewStringTag(readWriteModel.MemoryArea_DATA_BLOCKS, 1, 0, 0, 1, 20, readWriteModel.TransportSize_STRING)
		address, err := encodeS7Address(tag)
		require.NoError(t, err)
		addressAny := address.(readWriteModel.S7AddressAny)
		assert.Equal(t, readWriteModel.TransportSize_CHAR, addressAny.GetTransportSize())
		assert.Equal(t, uint16(22), addressAny.GetNumberOfElements())
	})
}
