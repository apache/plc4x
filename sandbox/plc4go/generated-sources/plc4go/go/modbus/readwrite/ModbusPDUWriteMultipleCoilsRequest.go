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
)

type ModbusPDUWriteMultipleCoilsRequest struct {
    startingAddress uint16
    quantity uint16
    value []int8
    ModbusPDU
}

func (m ModbusPDUWriteMultipleCoilsRequest) initialize() ModbusPDU {
    return m.ModbusPDU
}

func NewModbusPDUWriteMultipleCoilsRequest(startingAddress uint16, quantity uint16, value []int8) ModbusPDUInitializer {
    return &ModbusPDUWriteMultipleCoilsRequest{startingAddress: startingAddress, quantity: quantity, value: value}
}

func (m ModbusPDUWriteMultipleCoilsRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (startingAddress)
    lengthInBits += 16

    // Simple field (quantity)
    lengthInBits += 16

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.value) > 0 {
        lengthInBits += 8 * uint16(len(m.value))
    }

    return lengthInBits
}

func (m ModbusPDUWriteMultipleCoilsRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUWriteMultipleCoilsRequestParse(io spi.ReadBuffer) ModbusPDUInitializer {
    var startPos = io.GetPos()
    var curPos uint16

    // Simple Field (startingAddress)
    var startingAddress uint16 = io.ReadUint16(16)

    // Simple Field (quantity)
    var quantity uint16 = io.ReadUint16(16)

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
    return NewModbusPDUWriteMultipleCoilsRequest(startingAddress, quantity, value)
}

