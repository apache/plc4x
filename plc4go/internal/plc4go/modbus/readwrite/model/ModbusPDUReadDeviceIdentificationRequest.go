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
    "encoding/xml"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDUReadDeviceIdentificationRequest struct {
    Parent *ModbusPDU
    IModbusPDUReadDeviceIdentificationRequest
}

// The corresponding interface
type IModbusPDUReadDeviceIdentificationRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ModbusPDUReadDeviceIdentificationRequest) ErrorFlag() bool {
    return false
}

func (m *ModbusPDUReadDeviceIdentificationRequest) FunctionFlag() uint8 {
    return 0x2B
}

func (m *ModbusPDUReadDeviceIdentificationRequest) Response() bool {
    return false
}


func (m *ModbusPDUReadDeviceIdentificationRequest) InitializeParent(parent *ModbusPDU) {
}

func NewModbusPDUReadDeviceIdentificationRequest() *ModbusPDU {
    child := &ModbusPDUReadDeviceIdentificationRequest{
        Parent: NewModbusPDU(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastModbusPDUReadDeviceIdentificationRequest(structType interface{}) ModbusPDUReadDeviceIdentificationRequest {
    castFunc := func(typ interface{}) ModbusPDUReadDeviceIdentificationRequest {
        if casted, ok := typ.(ModbusPDUReadDeviceIdentificationRequest); ok {
            return casted
        }
        if casted, ok := typ.(*ModbusPDUReadDeviceIdentificationRequest); ok {
            return *casted
        }
        if casted, ok := typ.(ModbusPDU); ok {
            return CastModbusPDUReadDeviceIdentificationRequest(casted.Child)
        }
        if casted, ok := typ.(*ModbusPDU); ok {
            return CastModbusPDUReadDeviceIdentificationRequest(casted.Child)
        }
        return ModbusPDUReadDeviceIdentificationRequest{}
    }
    return castFunc(structType)
}

func (m *ModbusPDUReadDeviceIdentificationRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *ModbusPDUReadDeviceIdentificationRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadDeviceIdentificationRequestParse(io *utils.ReadBuffer) (*ModbusPDU, error) {

    // Create a partially initialized instance
    _child := &ModbusPDUReadDeviceIdentificationRequest{
        Parent: &ModbusPDU{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ModbusPDUReadDeviceIdentificationRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ModbusPDUReadDeviceIdentificationRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            }
        }
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *ModbusPDUReadDeviceIdentificationRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

