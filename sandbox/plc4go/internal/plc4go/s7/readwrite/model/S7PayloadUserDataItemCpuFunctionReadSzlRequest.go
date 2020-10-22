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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
)

// The data-structure of this message
type S7PayloadUserDataItemCpuFunctionReadSzlRequest struct {
    S7PayloadUserDataItem
}

// The corresponding interface
type IS7PayloadUserDataItemCpuFunctionReadSzlRequest interface {
    IS7PayloadUserDataItem
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) CpuFunctionType() uint8 {
    return 0x04
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) initialize(returnCode IDataTransportErrorCode, transportSize IDataTransportSize, szlId ISzlId, szlIndex uint16) spi.Message {
    m.ReturnCode = returnCode
    m.TransportSize = transportSize
    m.SzlId = szlId
    m.SzlIndex = szlIndex
    return m
}

func NewS7PayloadUserDataItemCpuFunctionReadSzlRequest() S7PayloadUserDataItemInitializer {
    return &S7PayloadUserDataItemCpuFunctionReadSzlRequest{}
}

func CastIS7PayloadUserDataItemCpuFunctionReadSzlRequest(structType interface{}) IS7PayloadUserDataItemCpuFunctionReadSzlRequest {
    castFunc := func(typ interface{}) IS7PayloadUserDataItemCpuFunctionReadSzlRequest {
        if iS7PayloadUserDataItemCpuFunctionReadSzlRequest, ok := typ.(IS7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
            return iS7PayloadUserDataItemCpuFunctionReadSzlRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7PayloadUserDataItemCpuFunctionReadSzlRequest(structType interface{}) S7PayloadUserDataItemCpuFunctionReadSzlRequest {
    castFunc := func(typ interface{}) S7PayloadUserDataItemCpuFunctionReadSzlRequest {
        if sS7PayloadUserDataItemCpuFunctionReadSzlRequest, ok := typ.(S7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
            return sS7PayloadUserDataItemCpuFunctionReadSzlRequest
        }
        if sS7PayloadUserDataItemCpuFunctionReadSzlRequest, ok := typ.(*S7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
            return *sS7PayloadUserDataItemCpuFunctionReadSzlRequest
        }
        return S7PayloadUserDataItemCpuFunctionReadSzlRequest{}
    }
    return castFunc(structType)
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.S7PayloadUserDataItem.LengthInBits()

    return lengthInBits
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadUserDataItemCpuFunctionReadSzlRequestParse(io *spi.ReadBuffer) (S7PayloadUserDataItemInitializer, error) {

    // Create the instance
    return NewS7PayloadUserDataItemCpuFunctionReadSzlRequest(), nil
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return S7PayloadUserDataItemSerialize(io, m.S7PayloadUserDataItem, CastIS7PayloadUserDataItem(m), ser)
}
