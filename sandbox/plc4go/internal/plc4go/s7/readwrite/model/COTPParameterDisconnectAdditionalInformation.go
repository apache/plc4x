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
)

// The data-structure of this message
type COTPParameterDisconnectAdditionalInformation struct {
    Data []uint8
    COTPParameter
}

// The corresponding interface
type ICOTPParameterDisconnectAdditionalInformation interface {
    ICOTPParameter
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m COTPParameterDisconnectAdditionalInformation) ParameterType() uint8 {
    return 0xE0
}

func (m COTPParameterDisconnectAdditionalInformation) initialize() spi.Message {
    return m
}

func NewCOTPParameterDisconnectAdditionalInformation(data []uint8) COTPParameterInitializer {
    return &COTPParameterDisconnectAdditionalInformation{Data: data}
}

func CastICOTPParameterDisconnectAdditionalInformation(structType interface{}) ICOTPParameterDisconnectAdditionalInformation {
    castFunc := func(typ interface{}) ICOTPParameterDisconnectAdditionalInformation {
        if iCOTPParameterDisconnectAdditionalInformation, ok := typ.(ICOTPParameterDisconnectAdditionalInformation); ok {
            return iCOTPParameterDisconnectAdditionalInformation
        }
        return nil
    }
    return castFunc(structType)
}

func CastCOTPParameterDisconnectAdditionalInformation(structType interface{}) COTPParameterDisconnectAdditionalInformation {
    castFunc := func(typ interface{}) COTPParameterDisconnectAdditionalInformation {
        if sCOTPParameterDisconnectAdditionalInformation, ok := typ.(COTPParameterDisconnectAdditionalInformation); ok {
            return sCOTPParameterDisconnectAdditionalInformation
        }
        if sCOTPParameterDisconnectAdditionalInformation, ok := typ.(*COTPParameterDisconnectAdditionalInformation); ok {
            return *sCOTPParameterDisconnectAdditionalInformation
        }
        return COTPParameterDisconnectAdditionalInformation{}
    }
    return castFunc(structType)
}

func (m COTPParameterDisconnectAdditionalInformation) LengthInBits() uint16 {
    var lengthInBits uint16 = m.COTPParameter.LengthInBits()

    // Array field
    if len(m.Data) > 0 {
        lengthInBits += 8 * uint16(len(m.Data))
    }

    return lengthInBits
}

func (m COTPParameterDisconnectAdditionalInformation) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPParameterDisconnectAdditionalInformationParse(io *spi.ReadBuffer, rest uint8) (COTPParameterInitializer, error) {

    // Array field (data)
    // Count array
    data := make([]uint8, rest)
    for curItem := uint16(0); curItem < uint16(rest); curItem++ {

        _item, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'data' field " + _err.Error())
        }
        data[curItem] = _item
    }

    // Create the instance
    return NewCOTPParameterDisconnectAdditionalInformation(data), nil
}

func (m COTPParameterDisconnectAdditionalInformation) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Array Field (data)
    if m.Data != nil {
        for _, _element := range m.Data {
            _elementErr := io.WriteUint8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'data' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return COTPParameterSerialize(io, m.COTPParameter, CastICOTPParameter(m), ser)
}
