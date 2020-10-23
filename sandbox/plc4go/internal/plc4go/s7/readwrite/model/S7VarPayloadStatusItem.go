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
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type S7VarPayloadStatusItem struct {
    ReturnCode IDataTransportErrorCode

}

// The corresponding interface
type IS7VarPayloadStatusItem interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewS7VarPayloadStatusItem(returnCode IDataTransportErrorCode) spi.Message {
    return &S7VarPayloadStatusItem{ReturnCode: returnCode}
}

func CastIS7VarPayloadStatusItem(structType interface{}) IS7VarPayloadStatusItem {
    castFunc := func(typ interface{}) IS7VarPayloadStatusItem {
        if iS7VarPayloadStatusItem, ok := typ.(IS7VarPayloadStatusItem); ok {
            return iS7VarPayloadStatusItem
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7VarPayloadStatusItem(structType interface{}) S7VarPayloadStatusItem {
    castFunc := func(typ interface{}) S7VarPayloadStatusItem {
        if sS7VarPayloadStatusItem, ok := typ.(S7VarPayloadStatusItem); ok {
            return sS7VarPayloadStatusItem
        }
        if sS7VarPayloadStatusItem, ok := typ.(*S7VarPayloadStatusItem); ok {
            return *sS7VarPayloadStatusItem
        }
        return S7VarPayloadStatusItem{}
    }
    return castFunc(structType)
}

func (m S7VarPayloadStatusItem) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Enum Field (returnCode)
    lengthInBits += 8

    return lengthInBits
}

func (m S7VarPayloadStatusItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7VarPayloadStatusItemParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Enum field (returnCode)
    returnCode, _returnCodeErr := DataTransportErrorCodeParse(io)
    if _returnCodeErr != nil {
        return nil, errors.New("Error parsing 'returnCode' field " + _returnCodeErr.Error())
    }

    // Create the instance
    return NewS7VarPayloadStatusItem(returnCode), nil
}

func (m S7VarPayloadStatusItem) Serialize(io utils.WriteBuffer) error {

    // Enum field (returnCode)
    returnCode := CastDataTransportErrorCode(m.ReturnCode)
    _returnCodeErr := returnCode.Serialize(io)
    if _returnCodeErr != nil {
        return errors.New("Error serializing 'returnCode' field " + _returnCodeErr.Error())
    }

    return nil
}
