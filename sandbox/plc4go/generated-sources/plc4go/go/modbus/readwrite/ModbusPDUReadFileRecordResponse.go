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

type ModbusPDUReadFileRecordResponse struct {
    items []ModbusPDUReadFileRecordResponseItem
    ModbusPDU
}

func (m ModbusPDUReadFileRecordResponse) initialize() spi.Message {
    return spi.Message(m)
}

func NewModbusPDUReadFileRecordResponse(items []ModbusPDUReadFileRecordResponseItem) ModbusPDUInitializer {
    return &ModbusPDUReadFileRecordResponse{items: items}
}

func (m ModbusPDUReadFileRecordResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.items) > 0 {
        for _, element := range m.items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m ModbusPDUReadFileRecordResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordResponseParse(io spi.ReadBuffer) ModbusPDUInitializer {
    var startPos = io.GetPos()
    var curPos uint16

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var byteCount uint8 = io.ReadUint8(8)

    // Array field (items)
    // Length array
    var _itemsLength := byteCount
    List<ModbusPDUReadFileRecordResponseItem> _itemsList = new LinkedList<>();
    itemsEndPos := io.GetPos() + _itemsLength;
    for ;io.getPos() < itemsEndPos; {
        _itemsList.add(ModbusPDUReadFileRecordResponseItemIO.staticParse(io));
    }
    ModbusPDUReadFileRecordResponseItem[] items = _itemsList.toArray(new ModbusPDUReadFileRecordResponseItem[0])

    // Create the instance
    return NewModbusPDUReadFileRecordResponse(items)
}

