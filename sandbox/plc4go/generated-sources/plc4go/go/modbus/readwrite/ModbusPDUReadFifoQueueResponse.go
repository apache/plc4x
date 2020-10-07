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

type ModbusPDUReadFifoQueueResponse struct {
    fifoValue []uint16
    ModbusPDU
}

func (m ModbusPDUReadFifoQueueResponse) initialize() ModbusPDU {
    return m.ModbusPDU
}

func NewModbusPDUReadFifoQueueResponse(fifoValue []uint16) ModbusPDUInitializer {
    return &ModbusPDUReadFifoQueueResponse{fifoValue: fifoValue}
}

func (m ModbusPDUReadFifoQueueResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Implicit Field (byteCount)
    lengthInBits += 16

    // Implicit Field (fifoCount)
    lengthInBits += 16

    // Array field
    if len(m.fifoValue) > 0 {
        lengthInBits += 16 * uint16(len(m.fifoValue))
    }

    return lengthInBits
}

func (m ModbusPDUReadFifoQueueResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadFifoQueueResponseParse(io spi.ReadBuffer) ModbusPDUInitializer {
    var startPos = io.GetPos()
    var curPos uint16

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var byteCount uint16 = io.ReadUint16(16)

    // Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var fifoCount uint16 = io.ReadUint16(16)

    // Array field (fifoValue)
    // Count array
    if fifoCount > math.MaxUint8 {
        throw new ParseException("Array count of " + (fifoCount) + " exceeds the maximum allowed count of " + math.MaxUint8);
    }
    uint16[] fifoValue;
    {
        var itemCount := fifoCount
        fifoValue = new uint16[itemCount]
        for curItem := 0; curItem < itemCount; curItem++ {
            
            fifoValue[curItem] = io.ReadUint16(16)
        }
    }

    // Create the instance
    return NewModbusPDUReadFifoQueueResponse(fifoValue)
}

