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
package readwrite

import (
	log "github.com/sirupsen/logrus"
	"math"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
)

type ModbusPDUReadFileRecordRequestItem struct {
	referenceType uint8
	fileNumber    uint16
	recordNumber  uint16
	recordLength  uint16
}

func NewModbusPDUReadFileRecordRequestItem(referenceType uint8, fileNumber uint16, recordNumber uint16, recordLength uint16) spi.Message {
	return &ModbusPDUReadFileRecordRequestItem{referenceType: referenceType, fileNumber: fileNumber, recordNumber: recordNumber, recordLength: recordLength}
}

func (m ModbusPDUReadFileRecordRequestItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (referenceType)
	lengthInBits += 8

	// Simple field (fileNumber)
	lengthInBits += 16

	// Simple field (recordNumber)
	lengthInBits += 16

	// Simple field (recordLength)
	lengthInBits += 16

	return lengthInBits
}

func (m ModbusPDUReadFileRecordRequestItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordRequestItemParse(io spi.ReadBuffer) spi.Message {
	var startPos = io.GetPos()
	var curPos uint16

	// Simple Field (referenceType)
	var referenceType uint8 = io.ReadUint8(8)

	// Simple Field (fileNumber)
	var fileNumber uint16 = io.ReadUint16(16)

	// Simple Field (recordNumber)
	var recordNumber uint16 = io.ReadUint16(16)

	// Simple Field (recordLength)
	var recordLength uint16 = io.ReadUint16(16)

	// Create the instance
	return NewModbusPDUReadFileRecordRequestItem(referenceType, fileNumber, recordNumber, recordLength)
}
