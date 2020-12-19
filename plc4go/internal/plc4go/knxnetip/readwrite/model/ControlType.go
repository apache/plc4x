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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

type ControlType uint8

type IControlType interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    ControlType_CONNECT ControlType = 0x0
    ControlType_DISCONNECT ControlType = 0x1
    ControlType_ACK ControlType = 0x2
    ControlType_NACK ControlType = 0x3
)

func ControlTypeValueOf(value uint8) ControlType {
    switch value {
        case 0x0:
            return ControlType_CONNECT
        case 0x1:
            return ControlType_DISCONNECT
        case 0x2:
            return ControlType_ACK
        case 0x3:
            return ControlType_NACK
    }
    return 0
}

func CastControlType(structType interface{}) ControlType {
    castFunc := func(typ interface{}) ControlType {
        if sControlType, ok := typ.(ControlType); ok {
            return sControlType
        }
        return 0
    }
    return castFunc(structType)
}

func (m ControlType) LengthInBits() uint16 {
    return 2
}

func (m ControlType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ControlTypeParse(io *utils.ReadBuffer) (ControlType, error) {
    val, err := io.ReadUint8(2)
    if err != nil {
        return 0, nil
    }
    return ControlTypeValueOf(val), nil
}

func (e ControlType) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(2, uint8(e))
    return err
}

func (e ControlType) String() string {
    switch e {
    case ControlType_CONNECT:
        return "CONNECT"
    case ControlType_DISCONNECT:
        return "DISCONNECT"
    case ControlType_ACK:
        return "ACK"
    case ControlType_NACK:
        return "NACK"
    }
    return ""
}
