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

type Status uint8

type IStatus interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    Status_NO_ERROR Status = 0x00
    Status_PROTOCOL_TYPE_NOT_SUPPORTED Status = 0x01
    Status_UNSUPPORTED_PROTOCOL_VERSION Status = 0x02
    Status_OUT_OF_ORDER_SEQUENCE_NUMBER Status = 0x04
    Status_INVALID_CONNECTION_ID Status = 0x21
    Status_CONNECTION_TYPE_NOT_SUPPORTED Status = 0x22
    Status_CONNECTION_OPTION_NOT_SUPPORTED Status = 0x23
    Status_NO_MORE_CONNECTIONS Status = 0x24
    Status_NO_MORE_UNIQUE_CONNECTIONS Status = 0x25
    Status_DATA_CONNECTION Status = 0x26
    Status_KNX_CONNECTION Status = 0x27
    Status_TUNNELLING_LAYER_NOT_SUPPORTED Status = 0x29
)

func StatusValueOf(value uint8) Status {
    switch value {
        case 0x00:
            return Status_NO_ERROR
        case 0x01:
            return Status_PROTOCOL_TYPE_NOT_SUPPORTED
        case 0x02:
            return Status_UNSUPPORTED_PROTOCOL_VERSION
        case 0x04:
            return Status_OUT_OF_ORDER_SEQUENCE_NUMBER
        case 0x21:
            return Status_INVALID_CONNECTION_ID
        case 0x22:
            return Status_CONNECTION_TYPE_NOT_SUPPORTED
        case 0x23:
            return Status_CONNECTION_OPTION_NOT_SUPPORTED
        case 0x24:
            return Status_NO_MORE_CONNECTIONS
        case 0x25:
            return Status_NO_MORE_UNIQUE_CONNECTIONS
        case 0x26:
            return Status_DATA_CONNECTION
        case 0x27:
            return Status_KNX_CONNECTION
        case 0x29:
            return Status_TUNNELLING_LAYER_NOT_SUPPORTED
    }
    return 0
}

func CastStatus(structType interface{}) Status {
    castFunc := func(typ interface{}) Status {
        if sStatus, ok := typ.(Status); ok {
            return sStatus
        }
        return 0
    }
    return castFunc(structType)
}

func (m Status) LengthInBits() uint16 {
    return 8
}

func (m Status) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func StatusParse(io *utils.ReadBuffer) (Status, error) {
    val, err := io.ReadUint8(8)
    if err != nil {
        return 0, nil
    }
    return StatusValueOf(val), nil
}

func (e Status) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(8, uint8(e))
    return err
}
