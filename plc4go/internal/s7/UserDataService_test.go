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

func TestDecodeControllerType(t *testing.T) {
	tests := []struct {
		descriptor string
		expected   readWriteModel.ControllerType
	}{
		{"6ES7 212-1AE40-0XB0", readWriteModel.ControllerType_S7_1200},
		{"6ES7 512-1DK01-0AB0", readWriteModel.ControllerType_S7_1500},
		{"6ES7 315-2EH14-0AB0", readWriteModel.ControllerType_S7_300},
		{"6ES7 412-2XJ05-0AB0", readWriteModel.ControllerType_S7_400},
		{"CPU 315-2 PN/DP", readWriteModel.ControllerType_S7_300},
		{"CPU 412-2", readWriteModel.ControllerType_S7_400},
		{"CPU 1212C", readWriteModel.ControllerType_S7_1200},
		{"CPU 1516-3", readWriteModel.ControllerType_S7_1500},
		{"", readWriteModel.ControllerType_ANY},
		{"garbage", readWriteModel.ControllerType_ANY},
		{"6ES7", readWriteModel.ControllerType_ANY},
	}
	for _, test := range tests {
		t.Run(test.descriptor, func(t *testing.T) {
			assert.Equal(t, test.expected, decodeControllerType(test.descriptor))
		})
	}
}

func TestFindArticleNumber(t *testing.T) {
	t.Run("embedded order number", func(t *testing.T) {
		data := append(append([]byte{0x00, 0x01, 0xFF}, []byte("6ES7 212-1AE40-0XB0")...), 0x00, 0x00)
		assert.Equal(t, "6ES7 212-1AE40-0XB0", findArticleNumber(data))
	})
	t.Run("no order number", func(t *testing.T) {
		assert.Equal(t, "", findArticleNumber([]byte("hello world, no siemens here")))
	})
	t.Run("model name fallback", func(t *testing.T) {
		data := append(append([]byte{0x00, 0x01}, []byte("CPU 315-2 PN/DP")...), 0x00)
		assert.Equal(t, "CPU 315-2 PN/DP", findCpuModelName(data))
	})
}

func TestParseSzlProbeResponse(t *testing.T) {
	makeResponse := func(items []byte, errorCode *uint16) readWriteModel.S7Message {
		return readWriteModel.NewS7MessageUserData(
			11,
			readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{
				readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
					0x12, 0x08, 0x04, 0x01, 0x00,
					ptr(uint8(0)), ptr(uint8(1)), errorCode,
				),
			}),
			readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{
				readWriteModel.NewS7PayloadUserDataItemCpuFunctionReadSzlResponse(
					readWriteModel.DataTransportErrorCode_OK,
					readWriteModel.DataTransportSize_OCTET_STRING,
					uint16(len(items)),
					items,
				),
			}),
		)
	}

	t.Run("successful identification", func(t *testing.T) {
		// SZL header (8 bytes) followed by an item starting with a 2 byte index and the MLFB
		items := append([]byte{0x00, 0x11, 0x00, 0x01, 0x00, 0x1C, 0x00, 0x01, 0x00, 0x01}, []byte("6ES7 512-1DK01-0AB0 ")...)
		article, controllerType, err := parseSzlProbeResponse(makeResponse(items, ptr(uint16(0))))
		require.NoError(t, err)
		assert.Equal(t, "6ES7 512-1DK01-0AB0", article)
		assert.Equal(t, readWriteModel.ControllerType_S7_1500, controllerType)
	})
	t.Run("plc rejects the szl id", func(t *testing.T) {
		_, _, err := parseSzlProbeResponse(makeResponse([]byte{0x00, 0x00, 0x00, 0x00}, ptr(uint16(0xD401))))
		assert.Error(t, err)
	})
	t.Run("too short payload", func(t *testing.T) {
		_, _, err := parseSzlProbeResponse(makeResponse([]byte{0x00}, ptr(uint16(0))))
		assert.Error(t, err)
	})
}

func TestSupportsUserDataServices(t *testing.T) {
	assert.True(t, supportsUserDataServices(readWriteModel.ControllerType_S7_300))
	assert.True(t, supportsUserDataServices(readWriteModel.ControllerType_S7_1500))
	assert.False(t, supportsUserDataServices(readWriteModel.ControllerType_LOGO))
	assert.False(t, supportsUserDataServices(readWriteModel.ControllerType_S7_200))
	assert.False(t, supportsUserDataServices(readWriteModel.ControllerType_ANY))
}

func TestParseListBlocksOfTypeResponse(t *testing.T) {
	makeResponse := func(items []byte, errorCode *uint16) readWriteModel.S7Message {
		return readWriteModel.NewS7MessageUserData(
			12,
			readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{
				readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
					0x12, 0x08, 0x03, 0x02, 0x00,
					ptr(uint8(0)), ptr(uint8(1)), errorCode,
				),
			}),
			readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{
				readWriteModel.NewS7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse(
					readWriteModel.DataTransportErrorCode_OK,
					readWriteModel.DataTransportSize_OCTET_STRING,
					uint16(len(items)),
					items,
				),
			}),
		)
	}

	t.Run("two data blocks", func(t *testing.T) {
		blockNumbers, err := parseListBlocksOfTypeResponse(makeResponse([]byte{
			0x00, 0x01, 0x22, 0x01, // DB1
			0x00, 0x45, 0x22, 0x01, // DB69
		}, ptr(uint16(0))))
		require.NoError(t, err)
		assert.Equal(t, []uint16{1, 69}, blockNumbers)
	})
	t.Run("rejected request", func(t *testing.T) {
		_, err := parseListBlocksOfTypeResponse(makeResponse([]byte{}, ptr(uint16(0xD401))))
		assert.Error(t, err)
	})
}

func ptr[T any](value T) *T {
	return &value
}
