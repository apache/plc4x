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
    "reflect"
)

// The data-structure of this message
type S7PayloadUserDataItem struct {
    ReturnCode IDataTransportErrorCode
    TransportSize IDataTransportSize
    SzlId ISzlId
    SzlIndex uint16

}

// The corresponding interface
type IS7PayloadUserDataItem interface {
    spi.Message
    CpuFunctionType() uint8
    Serialize(io utils.WriteBuffer) error
}

type S7PayloadUserDataItemInitializer interface {
    initialize(returnCode IDataTransportErrorCode, transportSize IDataTransportSize, szlId ISzlId, szlIndex uint16) spi.Message
}

func S7PayloadUserDataItemCpuFunctionType(m IS7PayloadUserDataItem) uint8 {
    return m.CpuFunctionType()
}


func CastIS7PayloadUserDataItem(structType interface{}) IS7PayloadUserDataItem {
    castFunc := func(typ interface{}) IS7PayloadUserDataItem {
        if iS7PayloadUserDataItem, ok := typ.(IS7PayloadUserDataItem); ok {
            return iS7PayloadUserDataItem
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7PayloadUserDataItem(structType interface{}) S7PayloadUserDataItem {
    castFunc := func(typ interface{}) S7PayloadUserDataItem {
        if sS7PayloadUserDataItem, ok := typ.(S7PayloadUserDataItem); ok {
            return sS7PayloadUserDataItem
        }
        if sS7PayloadUserDataItem, ok := typ.(*S7PayloadUserDataItem); ok {
            return *sS7PayloadUserDataItem
        }
        return S7PayloadUserDataItem{}
    }
    return castFunc(structType)
}

func (m S7PayloadUserDataItem) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Enum Field (returnCode)
    lengthInBits += 8

    // Enum Field (transportSize)
    lengthInBits += 8

    // Implicit Field (dataLength)
    lengthInBits += 16

    // Simple field (szlId)
    lengthInBits += m.SzlId.LengthInBits()

    // Simple field (szlIndex)
    lengthInBits += 16

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits
}

func (m S7PayloadUserDataItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadUserDataItemParse(io *utils.ReadBuffer, cpuFunctionType uint8) (spi.Message, error) {

    // Enum field (returnCode)
    returnCode, _returnCodeErr := DataTransportErrorCodeParse(io)
    if _returnCodeErr != nil {
        return nil, errors.New("Error parsing 'returnCode' field " + _returnCodeErr.Error())
    }

    // Enum field (transportSize)
    transportSize, _transportSizeErr := DataTransportSizeParse(io)
    if _transportSizeErr != nil {
        return nil, errors.New("Error parsing 'transportSize' field " + _transportSizeErr.Error())
    }

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _dataLengthErr := io.ReadUint16(16)
    if _dataLengthErr != nil {
        return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Simple Field (szlId)
    _szlIdMessage, _err := SzlIdParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'szlId'. " + _err.Error())
    }
    var szlId ISzlId
    szlId, _szlIdOk := _szlIdMessage.(ISzlId)
    if !_szlIdOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_szlIdMessage).Name() + " to ISzlId")
    }

    // Simple Field (szlIndex)
    szlIndex, _szlIndexErr := io.ReadUint16(16)
    if _szlIndexErr != nil {
        return nil, errors.New("Error parsing 'szlIndex' field " + _szlIndexErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer S7PayloadUserDataItemInitializer
    var typeSwitchError error
    switch {
    case cpuFunctionType == 0x04:
        initializer, typeSwitchError = S7PayloadUserDataItemCpuFunctionReadSzlRequestParse(io)
    case cpuFunctionType == 0x08:
        initializer, typeSwitchError = S7PayloadUserDataItemCpuFunctionReadSzlResponseParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Create the instance
    return initializer.initialize(returnCode, transportSize, szlId, szlIndex), nil
}

func S7PayloadUserDataItemSerialize(io utils.WriteBuffer, m S7PayloadUserDataItem, i IS7PayloadUserDataItem, childSerialize func() error) error {

    // Enum field (returnCode)
    returnCode := CastDataTransportErrorCode(m.ReturnCode)
    _returnCodeErr := returnCode.Serialize(io)
    if _returnCodeErr != nil {
        return errors.New("Error serializing 'returnCode' field " + _returnCodeErr.Error())
    }

    // Enum field (transportSize)
    transportSize := CastDataTransportSize(m.TransportSize)
    _transportSizeErr := transportSize.Serialize(io)
    if _transportSizeErr != nil {
        return errors.New("Error serializing 'transportSize' field " + _transportSizeErr.Error())
    }

    // Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    dataLength := uint16(uint16(uint16(m.LengthInBytes())) - uint16(uint16(4)))
    _dataLengthErr := io.WriteUint16(16, (dataLength))
    if _dataLengthErr != nil {
        return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
    }

    // Simple Field (szlId)
    szlId := CastISzlId(m.SzlId)
    _szlIdErr := szlId.Serialize(io)
    if _szlIdErr != nil {
        return errors.New("Error serializing 'szlId' field " + _szlIdErr.Error())
    }

    // Simple Field (szlIndex)
    szlIndex := uint16(m.SzlIndex)
    _szlIndexErr := io.WriteUint16(16, (szlIndex))
    if _szlIndexErr != nil {
        return errors.New("Error serializing 'szlIndex' field " + _szlIndexErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}
