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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
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
