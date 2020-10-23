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
type S7MessageResponseData struct {
    ErrorClass uint8
    ErrorCode uint8
    S7Message
}

// The corresponding interface
type IS7MessageResponseData interface {
    IS7Message
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7MessageResponseData) MessageType() uint8 {
    return 0x03
}

func (m S7MessageResponseData) initialize(tpduReference uint16, parameter *IS7Parameter, payload *IS7Payload) spi.Message {
    m.TpduReference = tpduReference
    m.Parameter = parameter
    m.Payload = payload
    return m
}

func NewS7MessageResponseData(errorClass uint8, errorCode uint8) S7MessageInitializer {
    return &S7MessageResponseData{ErrorClass: errorClass, ErrorCode: errorCode}
}

func CastIS7MessageResponseData(structType interface{}) IS7MessageResponseData {
    castFunc := func(typ interface{}) IS7MessageResponseData {
        if iS7MessageResponseData, ok := typ.(IS7MessageResponseData); ok {
            return iS7MessageResponseData
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7MessageResponseData(structType interface{}) S7MessageResponseData {
    castFunc := func(typ interface{}) S7MessageResponseData {
        if sS7MessageResponseData, ok := typ.(S7MessageResponseData); ok {
            return sS7MessageResponseData
        }
        if sS7MessageResponseData, ok := typ.(*S7MessageResponseData); ok {
            return *sS7MessageResponseData
        }
        return S7MessageResponseData{}
    }
    return castFunc(structType)
}

func (m S7MessageResponseData) LengthInBits() uint16 {
    var lengthInBits uint16 = m.S7Message.LengthInBits()

    // Simple field (errorClass)
    lengthInBits += 8

    // Simple field (errorCode)
    lengthInBits += 8

    return lengthInBits
}

func (m S7MessageResponseData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7MessageResponseDataParse(io *utils.ReadBuffer) (S7MessageInitializer, error) {

    // Simple Field (errorClass)
    errorClass, _errorClassErr := io.ReadUint8(8)
    if _errorClassErr != nil {
        return nil, errors.New("Error parsing 'errorClass' field " + _errorClassErr.Error())
    }

    // Simple Field (errorCode)
    errorCode, _errorCodeErr := io.ReadUint8(8)
    if _errorCodeErr != nil {
        return nil, errors.New("Error parsing 'errorCode' field " + _errorCodeErr.Error())
    }

    // Create the instance
    return NewS7MessageResponseData(errorClass, errorCode), nil
}

func (m S7MessageResponseData) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (errorClass)
    errorClass := uint8(m.ErrorClass)
    _errorClassErr := io.WriteUint8(8, (errorClass))
    if _errorClassErr != nil {
        return errors.New("Error serializing 'errorClass' field " + _errorClassErr.Error())
    }

    // Simple Field (errorCode)
    errorCode := uint8(m.ErrorCode)
    _errorCodeErr := io.WriteUint8(8, (errorCode))
    if _errorCodeErr != nil {
        return errors.New("Error serializing 'errorCode' field " + _errorCodeErr.Error())
    }

        return nil
    }
    return S7MessageSerialize(io, m.S7Message, CastIS7Message(m), ser)
}
