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
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func TestSerializePlcValue(t *testing.T) {
	anyController := readWriteModel.ControllerType_ANY

	t.Run("STRING honors the tag string length", func(t *testing.T) {
		tag := NewStringTag(readWriteModel.MemoryArea_DATA_BLOCKS, 1, 0, 0, 1, 10, readWriteModel.TransportSize_STRING)
		item, err := serializePlcValue(tag, spiValues.NewPlcSTRING("ABC"), anyController)
		require.NoError(t, err)
		data := item.GetData()
		// maxLen(1) + curLen(1) + 10 padded chars = 12 bytes, NOT the 256 of the old
		// silently-default-254 behavior.
		require.Len(t, data, 12)
		assert.Equal(t, uint8(10), data[0])
		assert.Equal(t, uint8(3), data[1])
		assert.Equal(t, []byte("ABC"), data[2:5])
	})
	t.Run("BOOL array is bit packed", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 10, readWriteModel.TransportSize_BOOL)
		list := make([]apiValues.PlcValue, 10)
		for i := range list {
			list[i] = spiValues.NewPlcBOOL(i == 0 || i == 1 || i == 9)
		}
		item, err := serializePlcValue(tag, spiValues.NewPlcList(list), anyController)
		require.NoError(t, err)
		assert.Equal(t, []byte{0b00000011, 0b00000010}, item.GetData())
		assert.Equal(t, readWriteModel.DataTransportSize_BYTE_WORD_DWORD, item.GetTransportSize())
	})
	t.Run("S5TIME serializes to two bytes", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 1, readWriteModel.TransportSize_S5TIME)
		item, err := serializePlcValue(tag, spiValues.NewPlcTIME(2*time.Second), anyController)
		require.NoError(t, err)
		// 2000ms -> finest fitting time base 10ms (0x0), counter 200 BCD
		assert.Equal(t, []byte{0x02, 0x00}, item.GetData())
	})
	t.Run("INT array concatenates elements", func(t *testing.T) {
		tag := NewTag(readWriteModel.MemoryArea_FLAGS_MARKERS, 0, 0, 0, 2, readWriteModel.TransportSize_INT)
		list := []apiValues.PlcValue{spiValues.NewPlcINT(1), spiValues.NewPlcINT(2)}
		item, err := serializePlcValue(tag, spiValues.NewPlcList(list), anyController)
		require.NoError(t, err)
		assert.Equal(t, []byte{0x00, 0x01, 0x00, 0x02}, item.GetData())
	})
}

// A write of whole bytes overwrites every bit in them, so only the bytes a BOOL array covers
// completely are written packed; each bit it shares a byte with goes out on its own.
func TestWriteItemsFor(t *testing.T) {
	handler := NewTagHandler()
	describe := func(address string) []string {
		tag, err := handler.ParseTag(address)
		require.NoError(t, err)
		var pieces []string
		for _, item := range writeItemsFor("t", tag) {
			piece := item.tag.(PlcTag)
			pieces = append(pieces, fmt.Sprintf("%d.%d x%d @%d", piece.GetByteOffset(), piece.GetBitOffset(),
				piece.GetNumElements(), item.firstElement))
		}
		return pieces
	}

	assert.Equal(t, []string{"16.0 x16 @0"}, describe("%DB1:16.0[0..15]:BOOL"))
	assert.Equal(t, []string{"16.0 x8 @0", "17.0 x1 @8", "17.1 x1 @9"}, describe("%DB1:16.0[0..9]:BOOL"))
	assert.Equal(t, []string{"16.6 x1 @0", "16.7 x1 @1", "17.0 x8 @2", "18.0 x1 @10", "18.1 x1 @11"},
		describe("%DB1:16.6[0..11]:BOOL"))
	assert.Equal(t, []string{"16.3 x1 @0", "16.4 x1 @1", "16.5 x1 @2"}, describe("%DB1:16.3[0..2]:BOOL"))
	assert.Equal(t, []string{"0.0 x4 @0"}, describe("%DB1.DBW0[0..3]:INT"), "any other tag is one item")
	assert.Equal(t, []string{"16.3 x1 @0"}, describe("%DB1:16.3:BOOL"), "and so is a single BOOL")
}

func TestWriteItemValueOf(t *testing.T) {
	list := spiValues.NewPlcList([]apiValues.PlcValue{
		spiValues.NewPlcBOOL(true), spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(true),
	})
	single := writeItem{tag: NewTag(readWriteModel.MemoryArea_DATA_BLOCKS, 1, 16, 7, 1, readWriteModel.TransportSize_BOOL),
		firstElement: 2, split: true}
	value, err := single.valueOf(list)
	require.NoError(t, err)
	assert.True(t, value.GetBool())
	assert.False(t, value.IsList(), "a single bit is written as a scalar")

	packed := writeItem{tag: NewTag(readWriteModel.MemoryArea_DATA_BLOCKS, 1, 17, 0, 8, readWriteModel.TransportSize_BOOL),
		firstElement: 1, split: true}
	_, err = packed.valueOf(list)
	assert.Error(t, err, "the list has no elements 1 to 8")

	_, err = single.valueOf(spiValues.NewPlcBOOL(true))
	assert.Error(t, err, "a BOOL array takes a list")
}

// The PLC answers per item, and a tag written in several items is only OK if all of them were.
func TestApplyWriteResponseMergesTheItemsOfATag(t *testing.T) {
	handler := NewTagHandler()
	bools, err := handler.ParseTag("%DB1:16.0[0..9]:BOOL")
	require.NoError(t, err)
	word, err := handler.ParseTag("%DB1.DBW20:INT")
	require.NoError(t, err)
	request := spiModel.NewDefaultPlcWriteRequest(
		map[string]apiModel.PlcTag{"bools": bools, "word": word},
		[]string{"bools", "word"},
		map[string]apiValues.PlcValue{},
		nil, nil)

	// bools is three items - byte 16, 17.0 and 17.1 - and the second one fails.
	statuses := []readWriteModel.DataTransportErrorCode{
		readWriteModel.DataTransportErrorCode_OK,
		readWriteModel.DataTransportErrorCode_ACCESS_DENIED,
		readWriteModel.DataTransportErrorCode_OK,
		readWriteModel.DataTransportErrorCode_OK,
	}
	var items []readWriteModel.S7VarPayloadStatusItem
	for _, status := range statuses {
		items = append(items, readWriteModel.NewS7VarPayloadStatusItem(status))
	}
	response := readWriteModel.NewS7MessageResponseData(1, readWriteModel.NewS7ParameterWriteVarResponse(uint8(len(items))),
		readWriteModel.NewS7PayloadWriteVarResponse(items), 0, 0)

	var sent []writeChunkItem
	for _, item := range writeItemsOf(request) {
		sent = append(sent, writeChunkItem{item: item})
	}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	require.NoError(t, (&Writer{log: zerolog.Nop()}).applyWriteResponse(response, sent, responseCodes))
	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, responseCodes["bools"])
	assert.Equal(t, apiModel.PlcResponseCode_OK, responseCodes["word"])
}
