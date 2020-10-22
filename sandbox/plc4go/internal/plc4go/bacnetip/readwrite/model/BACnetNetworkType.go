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

import "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"

type BACnetNetworkType uint8

type IBACnetNetworkType interface {
    spi.Message
    Serialize(io spi.WriteBuffer) error
}

const(
    BACnetNetworkType_ETHERNET BACnetNetworkType = 0x0
    BACnetNetworkType_ARCNET BACnetNetworkType = 0x1
    BACnetNetworkType_MSTP BACnetNetworkType = 0x2
    BACnetNetworkType_PTP BACnetNetworkType = 0x3
    BACnetNetworkType_LONTALK BACnetNetworkType = 0x4
    BACnetNetworkType_IPV4 BACnetNetworkType = 0x5
    BACnetNetworkType_ZIGBEE BACnetNetworkType = 0x6
    BACnetNetworkType_VIRTUAL BACnetNetworkType = 0x7
    BACnetNetworkType_REMOVED_NON_BACNET BACnetNetworkType = 0x8
    BACnetNetworkType_IPV6 BACnetNetworkType = 0x9
    BACnetNetworkType_SERIAL BACnetNetworkType = 0xA
)


func CastBACnetNetworkType(structType interface{}) BACnetNetworkType {
    castFunc := func(typ interface{}) BACnetNetworkType {
        if sBACnetNetworkType, ok := typ.(BACnetNetworkType); ok {
            return sBACnetNetworkType
        }
        return 0
    }
    return castFunc(structType)
}

func (m BACnetNetworkType) LengthInBits() uint16 {
    return 4
}

func (m BACnetNetworkType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetNetworkTypeParse(io *spi.ReadBuffer) (BACnetNetworkType, error) {
    // TODO: Implement ...
    return 0, nil
}

func (e BACnetNetworkType) Serialize(io spi.WriteBuffer) error {
    // TODO: Implement ...
    return nil
}
