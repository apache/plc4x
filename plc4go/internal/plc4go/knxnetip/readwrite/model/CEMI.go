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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    "reflect"
    "strings"
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
    xml.Marshaler
}

type ICEMIParent interface {
    SerializeParent(io utils.WriteBuffer, child ICEMI, serializeChildFunction func() error) error
    GetTypeName() string
}

type ICEMIChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *CEMI)
    GetTypeName() string
    ICEMI
}

func NewCEMI() *CEMI {
    return &CEMI{}
}

func CastCEMI(structType interface{}) *CEMI {
    castFunc := func(typ interface{}) *CEMI {
        if casted, ok := typ.(CEMI); ok {
            return &casted
        }
        if casted, ok := typ.(*CEMI); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *CEMI) GetTypeName() string {
    return "CEMI"
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
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
                default:
                    switch start.Attr[0].Value {
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataReq":
                            var dt *CEMIDataReq
                            if m.Child != nil {
                                dt = m.Child.(*CEMIDataReq)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataCon":
                            var dt *CEMIDataCon
                            if m.Child != nil {
                                dt = m.Child.(*CEMIDataCon)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataInd":
                            var dt *CEMIDataInd
                            if m.Child != nil {
                                dt = m.Child.(*CEMIDataInd)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIRawReq":
                            var dt *CEMIRawReq
                            if m.Child != nil {
                                dt = m.Child.(*CEMIRawReq)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIRawCon":
                            var dt *CEMIRawCon
                            if m.Child != nil {
                                dt = m.Child.(*CEMIRawCon)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIRawInd":
                            var dt *CEMIRawInd
                            if m.Child != nil {
                                dt = m.Child.(*CEMIRawInd)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIPollDataReq":
                            var dt *CEMIPollDataReq
                            if m.Child != nil {
                                dt = m.Child.(*CEMIPollDataReq)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIPollDataCon":
                            var dt *CEMIPollDataCon
                            if m.Child != nil {
                                dt = m.Child.(*CEMIPollDataCon)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIBusmonInd":
                            var dt *CEMIBusmonInd
                            if m.Child != nil {
                                dt = m.Child.(*CEMIBusmonInd)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIMPropReadReq":
                            var dt *CEMIMPropReadReq
                            if m.Child != nil {
                                dt = m.Child.(*CEMIMPropReadReq)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIMPropReadCon":
                            var dt *CEMIMPropReadCon
                            if m.Child != nil {
                                dt = m.Child.(*CEMIMPropReadCon)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                    }
            }
        }
    }
}

func (m *CEMI) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.knxnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

