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

package eip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
)

// TestSliceServiceHappyPath exercises a valid 2-entry MultipleServiceResponse
// offset table, so the reversed-table test below isn't the only case covered.
func TestSliceServiceHappyPath(t *testing.T) {
	service1, err := readWriteModel.NewCipReadResponse(0, 0,
		readWriteModel.NewCIPData(readWriteModel.CIPDataTypeCode_DINT, []byte{0x01, 0x00, 0x00, 0x00})).Serialize()
	require.NoError(t, err)
	service2, err := readWriteModel.NewCipReadResponse(0, 0,
		readWriteModel.NewCIPData(readWriteModel.CIPDataTypeCode_DINT, []byte{0x02, 0x00, 0x00, 0x00})).Serialize()
	require.NoError(t, err)

	servicesData := append(append([]byte{}, service1...), service2...)
	offsets := []uint16{0, uint16(len(service1))}

	first, ok := sliceService(servicesData, offsets, 0, 2)
	require.True(t, ok)
	firstResponse, isReadResponse := first.(readWriteModel.CipReadResponse)
	require.True(t, isReadResponse)
	assert.Equal(t, []byte{0x01, 0x00, 0x00, 0x00}, firstResponse.GetData().GetData())

	second, ok := sliceService(servicesData, offsets, 1, 2)
	require.True(t, ok)
	secondResponse, isReadResponse := second.(readWriteModel.CipReadResponse)
	require.True(t, isReadResponse)
	assert.Equal(t, []byte{0x02, 0x00, 0x00, 0x00}, secondResponse.GetData().GetData())
}

// TestSliceServiceReversedOffsetsIsRejected guards against a lying/reversed
// offset table wrapping around in unsigned 16-bit arithmetic before the
// bounds guards run. offsets[0] sits near 0xFFFF while offsets[1] sits near
// 0, so a naive `uint16(offsets[i] - offsets[0])` subtraction would wrap to a
// small positive number instead of tripping the "offset < 0" guard.
func TestSliceServiceReversedOffsetsIsRejected(t *testing.T) {
	servicesData := make([]byte, 8)
	offsets := []uint16{0xFFF0, 2}

	service, ok := sliceService(servicesData, offsets, 1, 2)
	assert.False(t, ok)
	assert.Nil(t, service)
}
