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
)

// The data-structure of this message
type CEMIBusmonInd struct {
    AdditionalInformationLength uint8
    AdditionalInformation []*CEMIAdditionalInformation
    CemiFrame *CEMIFrame
    Parent *CEMI
    ICEMIBusmonInd
}

// The corresponding interface
type ICEMIBusmonInd interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *CEMIBusmonInd) MessageCode() uint8 {
    return 0x2B
}


func (m *CEMIBusmonInd) InitializeParent(parent *CEMI) {
}

func NewCEMIBusmonInd(additionalInformationLength uint8, additionalInformation []*CEMIAdditionalInformation, cemiFrame *CEMIFrame, ) *CEMI {
    child := &CEMIBusmonInd{
        AdditionalInformationLength: additionalInformationLength,
        AdditionalInformation: additionalInformation,
        CemiFrame: cemiFrame,
        Parent: NewCEMI(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCEMIBusmonInd(structType interface{}) *CEMIBusmonInd {
    castFunc := func(typ interface{}) *CEMIBusmonInd {
        if casted, ok := typ.(CEMIBusmonInd); ok {
            return &casted
        }
        if casted, ok := typ.(*CEMIBusmonInd); ok {
            return casted
        }
        if casted, ok := typ.(CEMI); ok {
            return CastCEMIBusmonInd(casted.Child)
        }
        if casted, ok := typ.(*CEMI); ok {
            return CastCEMIBusmonInd(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *CEMIBusmonInd) GetTypeName() string {
    return "CEMIBusmonInd"
}

func (m *CEMIBusmonInd) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (additionalInformationLength)
    lengthInBits += 8

    // Array field
    if len(m.AdditionalInformation) > 0 {
        for _, element := range m.AdditionalInformation {
            lengthInBits += element.LengthInBits()
        }
    }

    // Simple field (cemiFrame)
    lengthInBits += m.CemiFrame.LengthInBits()

    return lengthInBits
}

func (m *CEMIBusmonInd) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIBusmonIndParse(io *utils.ReadBuffer) (*CEMI, error) {

    // Simple Field (additionalInformationLength)
    additionalInformationLength, _additionalInformationLengthErr := io.ReadUint8(8)
    if _additionalInformationLengthErr != nil {
        return nil, errors.New("Error parsing 'additionalInformationLength' field " + _additionalInformationLengthErr.Error())
    }

    // Array field (additionalInformation)
    // Length array
    additionalInformation := make([]*CEMIAdditionalInformation, 0)
    _additionalInformationLength := additionalInformationLength
    _additionalInformationEndPos := io.GetPos() + uint16(_additionalInformationLength)
    for ;io.GetPos() < _additionalInformationEndPos; {
        _item, _err := CEMIAdditionalInformationParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'additionalInformation' field " + _err.Error())
        }
        additionalInformation = append(additionalInformation, _item)
    }

    // Simple Field (cemiFrame)
    cemiFrame, _cemiFrameErr := CEMIFrameParse(io)
    if _cemiFrameErr != nil {
        return nil, errors.New("Error parsing 'cemiFrame' field " + _cemiFrameErr.Error())
    }

    // Create a partially initialized instance
    _child := &CEMIBusmonInd{
        AdditionalInformationLength: additionalInformationLength,
        AdditionalInformation: additionalInformation,
        CemiFrame: cemiFrame,
        Parent: &CEMI{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *CEMIBusmonInd) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (additionalInformationLength)
    additionalInformationLength := uint8(m.AdditionalInformationLength)
    _additionalInformationLengthErr := io.WriteUint8(8, (additionalInformationLength))
    if _additionalInformationLengthErr != nil {
        return errors.New("Error serializing 'additionalInformationLength' field " + _additionalInformationLengthErr.Error())
    }

    // Array Field (additionalInformation)
    if m.AdditionalInformation != nil {
        for _, _element := range m.AdditionalInformation {
            _elementErr := _element.Serialize(io)
            if _elementErr != nil {
                return errors.New("Error serializing 'additionalInformation' field " + _elementErr.Error())
            }
        }
    }

    // Simple Field (cemiFrame)
    _cemiFrameErr := m.CemiFrame.Serialize(io)
    if _cemiFrameErr != nil {
        return errors.New("Error serializing 'cemiFrame' field " + _cemiFrameErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *CEMIBusmonInd) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "additionalInformationLength":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.AdditionalInformationLength = data
            case "additionalInformation":
                var _values []*CEMIAdditionalInformation
                var dt *CEMIAdditionalInformation
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                _values = append(_values, dt)
                m.AdditionalInformation = _values
            case "cemiFrame":
                var dt *CEMIFrame
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.CemiFrame = dt
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

func (m *CEMIBusmonInd) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.AdditionalInformationLength, xml.StartElement{Name: xml.Name{Local: "additionalInformationLength"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "additionalInformation"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.AdditionalInformation, xml.StartElement{Name: xml.Name{Local: "additionalInformation"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "additionalInformation"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.CemiFrame, xml.StartElement{Name: xml.Name{Local: "cemiFrame"}}); err != nil {
        return err
    }
    return nil
}

