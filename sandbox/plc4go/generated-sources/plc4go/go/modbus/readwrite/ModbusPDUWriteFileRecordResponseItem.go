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
    "math"
    "plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
    log "github.com/sirupsen/logrus"
)

type ModbusPDUWriteFileRecordResponseItem struct {
    referenceType uint8
    fileNumber uint16
    recordNumber uint16
    recordData []int8

}


func NewModbusPDUWriteFileRecordResponseItem(referenceType uint8, fileNumber uint16, recordNumber uint16, recordData []int8) spi.Message {
    return &ModbusPDUWriteFileRecordResponseItem{referenceType: referenceType, fileNumber: fileNumber, recordNumber: recordNumber, recordData: recordData}
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

func ModbusPDUWriteFileRecordResponseItemParse(io spi.ReadBuffer) spi.Message {
    var startPos = io.GetPos()
    var curPos uint16

    // Simple Field (referenceType)
    var referenceType uint8 = io.ReadUint8(8)

    // Simple Field (fileNumber)
    var fileNumber uint16 = io.ReadUint16(16)

    // Simple Field (recordNumber)
    var recordNumber uint16 = io.ReadUint16(16)

    // Implicit Field (recordLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var recordLength uint16 = io.ReadUint16(16)

    // Array field (recordData)
    // Length array
    var _recordDataLength := recordLength
    List<int8> _recordDataList = new LinkedList<>();
    recordDataEndPos := io.GetPos() + _recordDataLength;
    for ;io.getPos() < recordDataEndPos; {
        _recordDataList.add(io.ReadInt8(8));
    }
    int8[] recordData = new int8[_recordDataList.size()]
    for i := 0; i < _recordDataList.size(); i++ {
        recordData[i] = (int8) _recordDataList.get(i)
    }

    // Create the instance
    return NewModbusPDUWriteFileRecordResponseItem(referenceType, fileNumber, recordNumber, recordData)
}

