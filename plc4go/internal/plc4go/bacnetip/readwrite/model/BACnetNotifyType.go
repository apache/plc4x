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

type BACnetNotifyType uint8

type IBACnetNotifyType interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    BACnetNotifyType_ALARM BACnetNotifyType = 0x0
    BACnetNotifyType_EVENT BACnetNotifyType = 0x1
    BACnetNotifyType_ACK_NOTIFICATION BACnetNotifyType = 0x2
)

func BACnetNotifyTypeValueOf(value uint8) BACnetNotifyType {
    switch value {
        case 0x0:
            return BACnetNotifyType_ALARM
        case 0x1:
            return BACnetNotifyType_EVENT
        case 0x2:
            return BACnetNotifyType_ACK_NOTIFICATION
    }
    return 0
}

func CastBACnetNotifyType(structType interface{}) BACnetNotifyType {
    castFunc := func(typ interface{}) BACnetNotifyType {
        if sBACnetNotifyType, ok := typ.(BACnetNotifyType); ok {
            return sBACnetNotifyType
        }
        return 0
    }
    return castFunc(structType)
}

func (m BACnetNotifyType) LengthInBits() uint16 {
    return 4
}

func (m BACnetNotifyType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetNotifyTypeParse(io *utils.ReadBuffer) (BACnetNotifyType, error) {
    val, err := io.ReadUint8(4)
    if err != nil {
        return 0, nil
    }
    return BACnetNotifyTypeValueOf(val), nil
}

func (e BACnetNotifyType) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(4, uint8(e))
    return err
}
