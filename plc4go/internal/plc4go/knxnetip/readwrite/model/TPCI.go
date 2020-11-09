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

type TPCI uint8

type ITPCI interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    TPCI_UNNUMBERED_DATA_PACKET TPCI = 0x0
    TPCI_UNNUMBERED TPCI = 0x1
    TPCI_NUMBERED_DATA_PACKET TPCI = 0x2
    TPCI_NUMBERED_CONTROL_DATA TPCI = 0x3
)

func TPCIValueOf(value uint8) TPCI {
    switch value {
        case 0x0:
            return TPCI_UNNUMBERED_DATA_PACKET
        case 0x1:
            return TPCI_UNNUMBERED
        case 0x2:
            return TPCI_NUMBERED_DATA_PACKET
        case 0x3:
            return TPCI_NUMBERED_CONTROL_DATA
    }
    return 0
}

func CastTPCI(structType interface{}) TPCI {
    castFunc := func(typ interface{}) TPCI {
        if sTPCI, ok := typ.(TPCI); ok {
            return sTPCI
        }
        return 0
    }
    return castFunc(structType)
}

func (m TPCI) LengthInBits() uint16 {
    return 2
}

func (m TPCI) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func TPCIParse(io *utils.ReadBuffer) (TPCI, error) {
    val, err := io.ReadUint8(2)
    if err != nil {
        return 0, nil
    }
    return TPCIValueOf(val), nil
}

func (e TPCI) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(2, uint8(e))
    return err
}
