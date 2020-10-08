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
    "reflect"
)

type ModbusSerialADU struct {
    transactionId uint16
    length uint16
    address uint8
    pdu ModbusPDU

}


func NewModbusSerialADU(transactionId uint16, length uint16, address uint8, pdu ModbusPDU) spi.Message {
    return &ModbusSerialADU{transactionId: transactionId, length: length, address: address, pdu: pdu}
}

func (m ModbusSerialADU) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Simple field (transactionId)
    lengthInBits += 16

    // Reserved Field (reserved)
    lengthInBits += 16

    // Simple field (length)
    lengthInBits += 16

    // Simple field (address)
    lengthInBits += 8

    // Simple field (pdu)
    lengthInBits += m.pdu.LengthInBits()

    return lengthInBits
}

func (m ModbusSerialADU) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusSerialADUParse(io spi.ReadBuffer, response bool) spi.Message {
    var startPos = io.GetPos()
    var curPos uint16

    // Simple Field (transactionId)
    var transactionId uint16 = io.ReadUint16(16)

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        var reserved uint16 = io.ReadUint16(16)
        if reserved != uint16(0x0000) {
            log.WithFields(log.Fields{
                "expected value": uint16(0x0000),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (length)
    var length uint16 = io.ReadUint16(16)

    // Simple Field (address)
    var address uint8 = io.ReadUint8(8)

    // Simple Field (pdu)
    var _pduMessage spi.Message = ModbusPDUParse(io, bool(response))
    var pdu ModbusPDU
    pdu, _pduOk := _pduMessage.(ModbusPDU)
    if !_pduOk {
        log.WithFields(log.Fields{
            "expected type": "ModbusPDU",
            "got type": reflect.TypeOf(_pduMessage),
        }).Error("Couldn't cast message")
        throw new RuntimeException()
    }

    // Create the instance
    return NewModbusSerialADU(transactionId, length, address, pdu)
}

