//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type ModbusPDUWriteFileRecordResponseItem struct {
	referenceType uint8
	fileNumber    uint16
	recordNumber  uint16
	recordData    []int8
}

// The corresponding interface
type IModbusPDUWriteFileRecordResponseItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewModbusPDUWriteFileRecordResponseItem(referenceType uint8, fileNumber uint16, recordNumber uint16, recordData []int8) spi.Message {
	return &ModbusPDUWriteFileRecordResponseItem{referenceType: referenceType, fileNumber: fileNumber, recordNumber: recordNumber, recordData: recordData}
}

func CastIModbusPDUWriteFileRecordResponseItem(structType interface{}) IModbusPDUWriteFileRecordResponseItem {
	castFunc := func(typ interface{}) IModbusPDUWriteFileRecordResponseItem {
		if iModbusPDUWriteFileRecordResponseItem, ok := typ.(IModbusPDUWriteFileRecordResponseItem); ok {
			return iModbusPDUWriteFileRecordResponseItem
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUWriteFileRecordResponseItem(structType interface{}) ModbusPDUWriteFileRecordResponseItem {
	castFunc := func(typ interface{}) ModbusPDUWriteFileRecordResponseItem {
		if sModbusPDUWriteFileRecordResponseItem, ok := typ.(ModbusPDUWriteFileRecordResponseItem); ok {
			return sModbusPDUWriteFileRecordResponseItem
		}
		return ModbusPDUWriteFileRecordResponseItem{}
	}
	return castFunc(structType)
}

func (m ModbusPDUWriteFileRecordResponseItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (referenceType)
	lengthInBits += 8

	// Simple field (fileNumber)
	lengthInBits += 16

	// Simple field (recordNumber)
	lengthInBits += 16

	// Implicit Field (recordLength)
	lengthInBits += 16

	// Array field
	if len(m.recordData) > 0 {
		lengthInBits += 8 * uint16(len(m.recordData))
	}

	return lengthInBits
}

func (m ModbusPDUWriteFileRecordResponseItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUWriteFileRecordResponseItemParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (referenceType)
	referenceType, _referenceTypeErr := io.ReadUint8(8)
	if _referenceTypeErr != nil {
		return nil, errors.New("Error parsing 'referenceType' field " + _referenceTypeErr.Error())
	}

	// Simple Field (fileNumber)
	fileNumber, _fileNumberErr := io.ReadUint16(16)
	if _fileNumberErr != nil {
		return nil, errors.New("Error parsing 'fileNumber' field " + _fileNumberErr.Error())
	}

	// Simple Field (recordNumber)
	recordNumber, _recordNumberErr := io.ReadUint16(16)
	if _recordNumberErr != nil {
		return nil, errors.New("Error parsing 'recordNumber' field " + _recordNumberErr.Error())
	}

	// Implicit Field (recordLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	recordLength, _recordLengthErr := io.ReadUint16(16)
	if _recordLengthErr != nil {
		return nil, errors.New("Error parsing 'recordLength' field " + _recordLengthErr.Error())
	}

	// Array field (recordData)
	var recordData []int8
	// Length array
	_recordDataLength := recordLength
	_recordDataEndPos := io.GetPos() + uint16(_recordDataLength)
	for io.GetPos() < _recordDataEndPos {
		_recordDataVal, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'recordData' field " + _err.Error())
		}
		recordData = append(recordData, _recordDataVal)
	}

	// Create the instance
	return NewModbusPDUWriteFileRecordResponseItem(referenceType, fileNumber, recordNumber, recordData), nil
}

func (m ModbusPDUWriteFileRecordResponseItem) Serialize(io spi.WriteBuffer) {

	// Simple Field (referenceType)
	referenceType := uint8(m.referenceType)
	io.WriteUint8(8, (referenceType))

	// Simple Field (fileNumber)
	fileNumber := uint16(m.fileNumber)
	io.WriteUint16(16, (fileNumber))

	// Simple Field (recordNumber)
	recordNumber := uint16(m.recordNumber)
	io.WriteUint16(16, (recordNumber))

	// Implicit Field (recordLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	recordLength := uint16(uint16(uint16(len(m.recordData))) / uint16(uint16(2)))
	io.WriteUint16(16, (recordLength))

	// Array Field (recordData)
	if m.recordData != nil {
		for _, _element := range m.recordData {
			io.WriteInt8(8, _element)
		}
	}

}
