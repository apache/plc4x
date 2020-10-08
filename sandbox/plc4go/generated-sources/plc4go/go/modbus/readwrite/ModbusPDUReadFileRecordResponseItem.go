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

type ModbusPDUReadFileRecordResponseItem struct {
    referenceType uint8
    data []int8

}


func NewModbusPDUReadFileRecordResponseItem(referenceType uint8, data []int8) spi.Message {
    return &ModbusPDUReadFileRecordResponseItem{referenceType: referenceType, data: data}
}

func (m ModbusPDUReadFileRecordResponseItem) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Implicit Field (dataLength)
    lengthInBits += 8

    // Simple field (referenceType)
    lengthInBits += 8

    // Array field
    if len(m.data) > 0 {
        lengthInBits += 8 * uint16(len(m.data))
    }

    return lengthInBits
}

func (m ModbusPDUReadFileRecordResponseItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordResponseItemParse(io spi.ReadBuffer) spi.Message {
    var startPos = io.GetPos()
    var curPos uint16

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var dataLength uint8 = io.ReadUint8(8)

    // Simple Field (referenceType)
    var referenceType uint8 = io.ReadUint8(8)

    // Array field (data)
    // Length array
    var _dataLength := (dataLength) - (1)
    List<int8> _dataList = new LinkedList<>();
    dataEndPos := io.GetPos() + _dataLength;
    for ;io.getPos() < dataEndPos; {
        _dataList.add(io.ReadInt8(8));
    }
    int8[] data = new int8[_dataList.size()]
    for i := 0; i < _dataList.size(); i++ {
        data[i] = (int8) _dataList.get(i)
    }

    // Create the instance
    return NewModbusPDUReadFileRecordResponseItem(referenceType, data)
}

