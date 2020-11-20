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

type KnxMedium uint8

type IKnxMedium interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxMedium_MEDIUM_RESERVED_1 KnxMedium = 0x01
    KnxMedium_MEDIUM_TP1 KnxMedium = 0x02
    KnxMedium_MEDIUM_PL110 KnxMedium = 0x04
    KnxMedium_MEDIUM_RESERVED_2 KnxMedium = 0x08
    KnxMedium_MEDIUM_RF KnxMedium = 0x10
    KnxMedium_MEDIUM_KNX_IP KnxMedium = 0x20
)

func KnxMediumValueOf(value uint8) KnxMedium {
    switch value {
        case 0x01:
            return KnxMedium_MEDIUM_RESERVED_1
        case 0x02:
            return KnxMedium_MEDIUM_TP1
        case 0x04:
            return KnxMedium_MEDIUM_PL110
        case 0x08:
            return KnxMedium_MEDIUM_RESERVED_2
        case 0x10:
            return KnxMedium_MEDIUM_RF
        case 0x20:
            return KnxMedium_MEDIUM_KNX_IP
    }
    return 0
}

func CastKnxMedium(structType interface{}) KnxMedium {
    castFunc := func(typ interface{}) KnxMedium {
        if sKnxMedium, ok := typ.(KnxMedium); ok {
            return sKnxMedium
        }
        return 0
    }
    return castFunc(structType)
}

func (m KnxMedium) LengthInBits() uint16 {
    return 8
}

func (m KnxMedium) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxMediumParse(io *utils.ReadBuffer) (KnxMedium, error) {
    val, err := io.ReadUint8(8)
    if err != nil {
        return 0, nil
    }
    return KnxMediumValueOf(val), nil
}

func (e KnxMedium) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(8, uint8(e))
    return err
}
