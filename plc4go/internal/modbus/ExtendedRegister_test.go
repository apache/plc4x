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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// The flat extended-register address names a file and an offset within it, and a request that runs
// past the end of a file continues at offset 0 of the next one (plc4j
// ModbusTcpConnection.getReadRequestPdu).
func TestSplitExtendedRegister(t *testing.T) {
	for _, test := range []struct {
		name        string
		address     uint16
		lengthWords uint16
		expected    []extendedRegisterGroup
	}{
		{
			"the first register of the first file",
			0, 1,
			[]extendedRegisterGroup{{fileNumber: 1, recordNumber: 0, lengthWords: 1}},
		},
		{
			"an offset inside the first file",
			1234, 2,
			[]extendedRegisterGroup{{fileNumber: 1, recordNumber: 1234, lengthWords: 2}},
		},
		{
			"the last register that still fits into the first file",
			9999, 1,
			[]extendedRegisterGroup{{fileNumber: 1, recordNumber: 9999, lengthWords: 1}},
		},
		{
			"a request filling the first file exactly is not split",
			0, 10000,
			[]extendedRegisterGroup{{fileNumber: 1, recordNumber: 0, lengthWords: 10000}},
		},
		{
			"the first register of the second file",
			10000, 1,
			[]extendedRegisterGroup{{fileNumber: 2, recordNumber: 0, lengthWords: 1}},
		},
		{
			"a request spanning a file boundary",
			9999, 2,
			[]extendedRegisterGroup{
				{fileNumber: 1, recordNumber: 9999, lengthWords: 1},
				{fileNumber: 2, recordNumber: 0, lengthWords: 1},
			},
		},
		{
			"a request spanning a boundary in the middle of the area",
			25000, 3,
			[]extendedRegisterGroup{
				{fileNumber: 3, recordNumber: 5000, lengthWords: 3},
			},
		},
		{
			"a request running off the end of the third file",
			29998, 5,
			[]extendedRegisterGroup{
				{fileNumber: 3, recordNumber: 9998, lengthWords: 2},
				{fileNumber: 4, recordNumber: 0, lengthWords: 3},
			},
		},
		{
			// plc4j stops after two groups and would hand the second one 20001 registers, more
			// than a file holds.
			"a request spanning more than two files",
			9999, 20002,
			[]extendedRegisterGroup{
				{fileNumber: 1, recordNumber: 9999, lengthWords: 1},
				{fileNumber: 2, recordNumber: 0, lengthWords: 10000},
				{fileNumber: 3, recordNumber: 0, lengthWords: 10000},
				{fileNumber: 4, recordNumber: 0, lengthWords: 1},
			},
		},
		{
			"the very end of the addressable area",
			65535, 1,
			[]extendedRegisterGroup{{fileNumber: 7, recordNumber: 5535, lengthWords: 1}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, splitExtendedRegister(test.address, test.lengthWords))
		})
	}
}

// Nothing to read means no items at all rather than one empty one, which wouldn't serialize into a
// valid request.
func TestSplitExtendedRegister_zeroLengthYieldsNoGroups(t *testing.T) {
	assert.Empty(t, splitExtendedRegister(1234, 0))
}

// An extended register is read with FC 0x14, addressed by file and offset within that file.
func TestReader_extendedRegisterUsesReadFileRecord(t *testing.T) {
	adu := capturedReadRequest(t, DefaultConfiguration(), parseTag(t, "extended-register:1:UDINT"))

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUReadFileRecordRequest)
	require.True(t, ok, "expected a read-file-record request, got %T", adu.GetPdu())
	items := pdu.GetItems()
	require.Len(t, items, 1)
	assert.Equal(t, extendedRegisterReferenceType, items[0].GetReferenceType())
	assert.Equal(t, uint16(1), items[0].GetFileNumber())
	// The extended register area is addressed starting at zero, so the logical address 1 is
	// register 1 of the first file and not register 0 (plc4j
	// ModbusTagExtendedRegister.getLogicalAddress).
	assert.Equal(t, uint16(1), items[0].GetRecordNumber())
	// A UDINT occupies two registers.
	assert.Equal(t, uint16(2), items[0].GetRecordLength())
}

// A read that runs past the end of a file is sent as one item per file it touches.
func TestReader_extendedRegisterSplitsAtTheFileBoundary(t *testing.T) {
	// Logical address 9999 is register 9999 of the first file, so the second of the two
	// registers a UDINT occupies falls into the second file.
	adu := capturedReadRequest(t, DefaultConfiguration(), parseTag(t, "extended-register:9999:UDINT"))

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUReadFileRecordRequest)
	require.True(t, ok, "expected a read-file-record request, got %T", adu.GetPdu())
	items := pdu.GetItems()
	require.Len(t, items, 2)
	assert.Equal(t, uint16(1), items[0].GetFileNumber())
	assert.Equal(t, uint16(9999), items[0].GetRecordNumber())
	assert.Equal(t, uint16(1), items[0].GetRecordLength())
	assert.Equal(t, uint16(2), items[1].GetFileNumber())
	assert.Equal(t, uint16(0), items[1].GetRecordNumber())
	assert.Equal(t, uint16(1), items[1].GetRecordLength())
}

// The registers of every item of the response together are the value; a value split across a file
// boundary would otherwise be decoded from the first half alone.
func TestReader_extendedRegisterResponseIsAssembledFromAllItems(t *testing.T) {
	reader := NewReader(DefaultConfiguration(), newCaptureCodec(nil))
	request := readRequestFor(t, reader, parseTag(t, "extended-register:9999:UDINT"))
	responseAdu := readWriteModel.NewModbusTcpADU(1, 1, readWriteModel.NewModbusPDUReadFileRecordResponse(
		[]readWriteModel.ModbusPDUReadFileRecordResponseItem{
			readWriteModel.NewModbusPDUReadFileRecordResponseItem(extendedRegisterReferenceType, []byte{0x01, 0x02}),
			readWriteModel.NewModbusPDUReadFileRecordResponseItem(extendedRegisterReferenceType, []byte{0x03, 0x04}),
		}))

	response, err := reader.ToPlc4xReadResponse(responseAdu, request)

	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("tag"))
	assert.Equal(t, uint32(0x01020304), response.GetValue("tag").GetUint32())
}

// An extended register is written with FC 0x15, addressed the same way the read is.
func TestWriter_extendedRegisterUsesWriteFileRecord(t *testing.T) {
	adu, _ := capturedRequestPdu(t, parseTag(t, "extended-register:1:UDINT"), spiValues.NewPlcUDINT(0x01020304))

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteFileRecordRequest)
	require.True(t, ok, "expected a write-file-record request, got %T", adu.GetPdu())
	items := pdu.GetItems()
	require.Len(t, items, 1)
	assert.Equal(t, extendedRegisterReferenceType, items[0].GetReferenceType())
	assert.Equal(t, uint16(1), items[0].GetFileNumber())
	assert.Equal(t, uint16(1), items[0].GetRecordNumber())
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, items[0].GetRecordData())
}

// A write that runs past the end of a file is cut up along that boundary, and each item carries
// exactly the registers that fall into its file.
func TestWriter_extendedRegisterSplitsPayloadAtTheFileBoundary(t *testing.T) {
	adu, _ := capturedRequestPdu(t, parseTag(t, "extended-register:9999:UDINT"), spiValues.NewPlcUDINT(0x01020304))

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteFileRecordRequest)
	require.True(t, ok, "expected a write-file-record request, got %T", adu.GetPdu())
	items := pdu.GetItems()
	require.Len(t, items, 2)
	assert.Equal(t, uint16(1), items[0].GetFileNumber())
	assert.Equal(t, uint16(9999), items[0].GetRecordNumber())
	assert.Equal(t, []byte{0x01, 0x02}, items[0].GetRecordData())
	assert.Equal(t, uint16(2), items[1].GetFileNumber())
	assert.Equal(t, uint16(0), items[1].GetRecordNumber())
	assert.Equal(t, []byte{0x03, 0x04}, items[1].GetRecordData())
}

// A tag that names its own unit and byte order is honoured on the extended-register path too.
func TestWriter_extendedRegisterHonoursTheTagSettings(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.unitIdentifier = 17
	adu, _ := capturedRequestPduWith(t, configuration,
		parseTag(t, "extended-register:1:UDINT{unit-id: 5, byte-order: 'LITTLE_ENDIAN'}"),
		spiValues.NewPlcUDINT(0x01020304))

	assert.Equal(t, uint8(5), adu.GetUnitIdentifier())
	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteFileRecordRequest)
	require.True(t, ok, "expected a write-file-record request, got %T", adu.GetPdu())
	require.Len(t, pdu.GetItems(), 1)
	assert.Equal(t, []byte{0x04, 0x03, 0x02, 0x01}, pdu.GetItems()[0].GetRecordData())
}

// A write-file-record response that isn't an exception is a successful write.
func TestWriter_extendedRegisterWriteResponse(t *testing.T) {
	writer := NewWriter(DefaultConfiguration(), newCaptureCodec(nil))
	tag := parseTag(t, "extended-register:1:UDINT")
	request := writeRequestFor(t, writer, tag, spiValues.NewPlcUDINT(0x01020304))
	requestPdu := readWriteModel.NewModbusPDUWriteFileRecordRequest(
		[]readWriteModel.ModbusPDUWriteFileRecordRequestItem{
			readWriteModel.NewModbusPDUWriteFileRecordRequestItem(extendedRegisterReferenceType, 1, 1, []byte{0x01, 0x02, 0x03, 0x04}),
		})
	responsePdu := readWriteModel.NewModbusPDUWriteFileRecordResponse(
		[]readWriteModel.ModbusPDUWriteFileRecordResponseItem{
			readWriteModel.NewModbusPDUWriteFileRecordResponseItem(extendedRegisterReferenceType, 1, 1, []byte{0x01, 0x02, 0x03, 0x04}),
		})

	t.Run("the echo of the request we sent", func(t *testing.T) {
		response, err := writer.ToPlc4xWriteResponse(
			readWriteModel.NewModbusTcpADU(1, 1, requestPdu),
			readWriteModel.NewModbusTcpADU(1, 1, responsePdu),
			request,
		)
		require.NoError(t, err)
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("tag"))
	})

	// A response that doesn't belong to the request we sent must be an error rather than a panic
	// on an unchecked type assertion.
	t.Run("an echo for a request we never sent", func(t *testing.T) {
		response, err := writer.ToPlc4xWriteResponse(
			readWriteModel.NewModbusTcpADU(1, 1, readWriteModel.NewModbusPDUWriteSingleCoilRequest(2, 0xFF00)),
			readWriteModel.NewModbusTcpADU(1, 1, responsePdu),
			request,
		)
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	// The exception codes are mapped the same way they are for every other function code.
	t.Run("an exception", func(t *testing.T) {
		response, err := writer.ToPlc4xWriteResponse(
			readWriteModel.NewModbusTcpADU(1, 1, requestPdu),
			readWriteModel.NewModbusTcpADU(1, 1, readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS)),
			request,
		)
		require.NoError(t, err)
		assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, response.GetResponseCode("tag"))
	})
}
