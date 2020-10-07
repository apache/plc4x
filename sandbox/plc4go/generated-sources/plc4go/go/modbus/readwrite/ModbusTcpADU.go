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

type ModbusTcpADU struct {
    transactionIdentifier uint16
    unitIdentifier uint8
    pdu ModbusPDU

}


func NewModbusTcpADU(transactionIdentifier uint16, unitIdentifier uint8, pdu ModbusPDU) ModbusTcpADU {
    return &ModbusTcpADU{transactionIdentifier: transactionIdentifier, unitIdentifier: unitIdentifier, pdu: pdu}
}

func (m ModbusTcpADU) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Simple field (transactionIdentifier)
    lengthInBits += 16

    // Const Field (protocolIdentifier)
    lengthInBits += 16

    // Implicit Field (length)
    lengthInBits += 16

    // Simple field (unitIdentifier)
    lengthInBits += 8

    // Simple field (pdu)
    lengthInBits += m.pdu.LengthInBits()

    return lengthInBits
}

func (m ModbusTcpADU) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusTcpADUParse(io spi.ReadBuffer, response bool) ModbusTcpADU {
    var startPos = io.GetPos()
    var curPos uint16

    // Simple Field (transactionIdentifier)
    var transactionIdentifier uint16 = io.ReadUint16(16)

    // Const Field (protocolIdentifier)
    uint16 protocolIdentifier = io.ReadUint16(16)
    if protocolIdentifier != ModbusTcpADU.PROTOCOLIDENTIFIER {
        throw new ParseException("Expected constant value " + ModbusTcpADU.PROTOCOLIDENTIFIER + " but got " + protocolIdentifier)
    }

    // Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var length uint16 = io.ReadUint16(16)

    // Simple Field (unitIdentifier)
    var unitIdentifier uint8 = io.ReadUint8(8)

    // Simple Field (pdu)
    var pdu ModbusPDU = ModbusPDUIO.staticParse(io, (bool) (response))

    // Create the instance
    return NewModbusTcpADU(transactionIdentifier, unitIdentifier, pdu)
}

