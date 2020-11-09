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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

type APCI uint8

type IAPCI interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    APCI_GROUP_VALUE_READ_PDU APCI = 0x0
    APCI_GROUP_VALUE_RESPONSE_PDU APCI = 0x1
    APCI_GROUP_VALUE_WRITE_PDU APCI = 0x2
    APCI_INDIVIDUAL_ADDRESS_WRITE_PDU APCI = 0x3
    APCI_INDIVIDUAL_ADDRESS_READ_PDU APCI = 0x4
    APCI_INDIVIDUAL_ADDRESS_RESPONSE_PDU APCI = 0x5
    APCI_ADC_READ_PDU APCI = 0x6
    APCI_ADC_RESPONSE_PDU APCI = 0x7
    APCI_MEMORY_READ_PDU APCI = 0x8
    APCI_MEMORY_RESPONSE_PDU APCI = 0x9
    APCI_MEMORY_WRITE_PDU APCI = 0xA
    APCI_USER_MESSAGE_PDU APCI = 0xB
    APCI_DEVICE_DESCRIPTOR_READ_PDU APCI = 0xC
    APCI_DEVICE_DESCRIPTOR_RESPONSE_PDU APCI = 0xD
    APCI_RESTART_PDU APCI = 0xE
    APCI_OTHER_PDU APCI = 0xF
)

func APCIValueOf(value uint8) APCI {
    switch value {
        case 0x0:
            return APCI_GROUP_VALUE_READ_PDU
        case 0x1:
            return APCI_GROUP_VALUE_RESPONSE_PDU
        case 0x2:
            return APCI_GROUP_VALUE_WRITE_PDU
        case 0x3:
            return APCI_INDIVIDUAL_ADDRESS_WRITE_PDU
        case 0x4:
            return APCI_INDIVIDUAL_ADDRESS_READ_PDU
        case 0x5:
            return APCI_INDIVIDUAL_ADDRESS_RESPONSE_PDU
        case 0x6:
            return APCI_ADC_READ_PDU
        case 0x7:
            return APCI_ADC_RESPONSE_PDU
        case 0x8:
            return APCI_MEMORY_READ_PDU
        case 0x9:
            return APCI_MEMORY_RESPONSE_PDU
        case 0xA:
            return APCI_MEMORY_WRITE_PDU
        case 0xB:
            return APCI_USER_MESSAGE_PDU
        case 0xC:
            return APCI_DEVICE_DESCRIPTOR_READ_PDU
        case 0xD:
            return APCI_DEVICE_DESCRIPTOR_RESPONSE_PDU
        case 0xE:
            return APCI_RESTART_PDU
        case 0xF:
            return APCI_OTHER_PDU
    }
    return 0
}

func CastAPCI(structType interface{}) APCI {
    castFunc := func(typ interface{}) APCI {
        if sAPCI, ok := typ.(APCI); ok {
            return sAPCI
        }
        return 0
    }
    return castFunc(structType)
}

func (m APCI) LengthInBits() uint16 {
    return 4
}

func (m APCI) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func APCIParse(io *utils.ReadBuffer) (APCI, error) {
    val, err := io.ReadUint8(4)
    if err != nil {
        return 0, nil
    }
    return APCIValueOf(val), nil
}

func (e APCI) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(4, uint8(e))
    return err
}
