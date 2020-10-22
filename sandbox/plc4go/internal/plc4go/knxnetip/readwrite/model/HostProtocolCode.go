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

type HostProtocolCode uint8

type IHostProtocolCode interface {
    spi.Message
    Serialize(io spi.WriteBuffer) error
}

const(
    HostProtocolCode_IPV4_UDP HostProtocolCode = 0x01
    HostProtocolCode_IPV4_TCP HostProtocolCode = 0x02
)

func HostProtocolCodeValueOf(value uint8) HostProtocolCode {
    switch value {
        case 0x01:
            return HostProtocolCode_IPV4_UDP
        case 0x02:
            return HostProtocolCode_IPV4_TCP
    }
    return 0
}

func CastHostProtocolCode(structType interface{}) HostProtocolCode {
    castFunc := func(typ interface{}) HostProtocolCode {
        if sHostProtocolCode, ok := typ.(HostProtocolCode); ok {
            return sHostProtocolCode
        }
        return 0
    }
    return castFunc(structType)
}

func (m HostProtocolCode) LengthInBits() uint16 {
    return 8
}

func (m HostProtocolCode) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func HostProtocolCodeParse(io *spi.ReadBuffer) (HostProtocolCode, error) {
    // TODO: Implement ...
    return 0, nil
}

func (e HostProtocolCode) Serialize(io spi.WriteBuffer) error {
    // TODO: Implement ...
    return nil
}
