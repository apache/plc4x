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

type ModbusPDUReadWriteMultipleHoldingRegistersRequest struct {
    readStartingAddress uint16
    readQuantity uint16
    writeStartingAddress uint16
    writeQuantity uint16
    value []int8
    ModbusPDU
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) initialize() spi.Message {
    return spi.Message(m)
}

func NewModbusPDUReadWriteMultipleHoldingRegistersRequest(readStartingAddress uint16, readQuantity uint16, writeStartingAddress uint16, writeQuantity uint16, value []int8) ModbusPDUInitializer {
    return &ModbusPDUReadWriteMultipleHoldingRegistersRequest{readStartingAddress: readStartingAddress, readQuantity: readQuantity, writeStartingAddress: writeStartingAddress, writeQuantity: writeQuantity, value: value}
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (readStartingAddress)
    lengthInBits += 16

    // Simple field (readQuantity)
    lengthInBits += 16

    // Simple field (writeStartingAddress)
    lengthInBits += 16

    // Simple field (writeQuantity)
    lengthInBits += 16

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.value) > 0 {
        lengthInBits += 8 * uint16(len(m.value))
    }

    return lengthInBits
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadWriteMultipleHoldingRegistersRequestParse(io spi.ReadBuffer) ModbusPDUInitializer {
    var startPos = io.GetPos()
    var curPos uint16

    // Simple Field (readStartingAddress)
    var readStartingAddress uint16 = io.ReadUint16(16)

    // Simple Field (readQuantity)
    var readQuantity uint16 = io.ReadUint16(16)

    // Simple Field (writeStartingAddress)
    var writeStartingAddress uint16 = io.ReadUint16(16)

    // Simple Field (writeQuantity)
    var writeQuantity uint16 = io.ReadUint16(16)

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var byteCount uint8 = io.ReadUint8(8)

    // Array field (value)
    // Count array
    if byteCount > math.MaxUint8 {
        throw new ParseException("Array count of " + (byteCount) + " exceeds the maximum allowed count of " + math.MaxUint8);
    }
    int8[] value;
    {
        var itemCount := byteCount
        value = new int8[itemCount]
        for curItem := 0; curItem < itemCount; curItem++ {
            
            value[curItem] = io.ReadInt8(8)
        }
    }

    // Create the instance
    return NewModbusPDUReadWriteMultipleHoldingRegistersRequest(readStartingAddress, readQuantity, writeStartingAddress, writeQuantity, value)
}

