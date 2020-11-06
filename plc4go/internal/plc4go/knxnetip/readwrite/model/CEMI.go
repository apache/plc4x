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
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMI struct {
    Child ICEMIChild
    ICEMI
    ICEMIParent
}

// The corresponding interface
type ICEMI interface {
    MessageCode() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type ICEMIParent interface {
    SerializeParent(io utils.WriteBuffer, child ICEMI, serializeChildFunction func() error) error
}

type ICEMIChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *CEMI)
    ICEMI
}

func NewCEMI() *CEMI {
    return &CEMI{}
}

func CastCEMI(structType interface{}) CEMI {
    castFunc := func(typ interface{}) CEMI {
        if casted, ok := typ.(CEMI); ok {
            return casted
        }
        if casted, ok := typ.(*CEMI); ok {
            return *casted
        }
        return CEMI{}
    }
    return castFunc(structType)
}

func (m *CEMI) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (messageCode)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *CEMI) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIParse(io *utils.ReadBuffer, size uint8) (*CEMI, error) {

    // Discriminator Field (messageCode) (Used as input to a switch field)
    messageCode, _messageCodeErr := io.ReadUint8(8)
    if _messageCodeErr != nil {
        return nil, errors.New("Error parsing 'messageCode' field " + _messageCodeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *CEMI
    var typeSwitchError error
    switch {
    case messageCode == 0x11:
        _parent, typeSwitchError = CEMIDataReqParse(io)
    case messageCode == 0x2E:
        _parent, typeSwitchError = CEMIDataConParse(io)
    case messageCode == 0x29:
        _parent, typeSwitchError = CEMIDataIndParse(io)
    case messageCode == 0x10:
        _parent, typeSwitchError = CEMIRawReqParse(io)
    case messageCode == 0x2F:
        _parent, typeSwitchError = CEMIRawConParse(io)
    case messageCode == 0x2D:
        _parent, typeSwitchError = CEMIRawIndParse(io)
    case messageCode == 0x13:
        _parent, typeSwitchError = CEMIPollDataReqParse(io)
    case messageCode == 0x25:
        _parent, typeSwitchError = CEMIPollDataConParse(io)
    case messageCode == 0x2B:
        _parent, typeSwitchError = CEMIBusmonIndParse(io)
    case messageCode == 0xFC:
        _parent, typeSwitchError = CEMIMPropReadReqParse(io)
    case messageCode == 0xFB:
        _parent, typeSwitchError = CEMIMPropReadConParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *CEMI) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *CEMI) SerializeParent(io utils.WriteBuffer, child ICEMI, serializeChildFunction func() error) error {

    // Discriminator Field (messageCode) (Used as input to a switch field)
    messageCode := uint8(child.MessageCode())
    _messageCodeErr := io.WriteUint8(8, (messageCode))
    if _messageCodeErr != nil {
        return errors.New("Error serializing 'messageCode' field " + _messageCodeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *CEMI) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            }
        }
    }
}

func (m *CEMI) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.CEMI"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

