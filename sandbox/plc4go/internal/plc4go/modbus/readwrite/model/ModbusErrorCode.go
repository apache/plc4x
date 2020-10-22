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

type ModbusErrorCode uint8

type IModbusErrorCode interface {
    spi.Message
    Serialize(io spi.WriteBuffer) error
}

const(
    ModbusErrorCode_ILLEGAL_FUNCTION ModbusErrorCode = 1
    ModbusErrorCode_ILLEGAL_DATA_ADDRESS ModbusErrorCode = 2
    ModbusErrorCode_ILLEGAL_DATA_VALUE ModbusErrorCode = 3
    ModbusErrorCode_SLAVE_DEVICE_FAILURE ModbusErrorCode = 4
    ModbusErrorCode_ACKNOWLEDGE ModbusErrorCode = 5
    ModbusErrorCode_SLAVE_DEVICE_BUSY ModbusErrorCode = 6
    ModbusErrorCode_NEGATIVE_ACKNOWLEDGE ModbusErrorCode = 7
    ModbusErrorCode_MEMORY_PARITY_ERROR ModbusErrorCode = 8
    ModbusErrorCode_GATEWAY_PATH_UNAVAILABLE ModbusErrorCode = 10
    ModbusErrorCode_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND ModbusErrorCode = 11
)

func ModbusErrorCodeValueOf(value uint8) ModbusErrorCode {
    switch value {
        case 1:
            return ModbusErrorCode_ILLEGAL_FUNCTION
        case 10:
            return ModbusErrorCode_GATEWAY_PATH_UNAVAILABLE
        case 11:
            return ModbusErrorCode_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND
        case 2:
            return ModbusErrorCode_ILLEGAL_DATA_ADDRESS
        case 3:
            return ModbusErrorCode_ILLEGAL_DATA_VALUE
        case 4:
            return ModbusErrorCode_SLAVE_DEVICE_FAILURE
        case 5:
            return ModbusErrorCode_ACKNOWLEDGE
        case 6:
            return ModbusErrorCode_SLAVE_DEVICE_BUSY
        case 7:
            return ModbusErrorCode_NEGATIVE_ACKNOWLEDGE
        case 8:
            return ModbusErrorCode_MEMORY_PARITY_ERROR
    }
    return 0
}

func CastModbusErrorCode(structType interface{}) ModbusErrorCode {
    castFunc := func(typ interface{}) ModbusErrorCode {
        if sModbusErrorCode, ok := typ.(ModbusErrorCode); ok {
            return sModbusErrorCode
        }
        return 0
    }
    return castFunc(structType)
}

func (m ModbusErrorCode) LengthInBits() uint16 {
    return 8
}

func (m ModbusErrorCode) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusErrorCodeParse(io *spi.ReadBuffer) (ModbusErrorCode, error) {
    // TODO: Implement ...
    return 0, nil
}

func (e ModbusErrorCode) Serialize(io spi.WriteBuffer) error {
    // TODO: Implement ...
    return nil
}
