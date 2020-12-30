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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

// The data-structure of this message
type LDataInd struct {
    AdditionalInformationLength uint8
    AdditionalInformation []*CEMIAdditionalInformation
    DataFrame *LDataFrame
    Parent *CEMI
    ILDataInd
}

// The corresponding interface
type ILDataInd interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *LDataInd) MessageCode() uint8 {
    return 0x29
}


func (m *LDataInd) InitializeParent(parent *CEMI) {
}

func NewLDataInd(additionalInformationLength uint8, additionalInformation []*CEMIAdditionalInformation, dataFrame *LDataFrame, ) *CEMI {
    child := &LDataInd{
        AdditionalInformationLength: additionalInformationLength,
        AdditionalInformation: additionalInformation,
        DataFrame: dataFrame,
        Parent: NewCEMI(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastLDataInd(structType interface{}) *LDataInd {
    castFunc := func(typ interface{}) *LDataInd {
        if casted, ok := typ.(LDataInd); ok {
            return &casted
        }
        if casted, ok := typ.(*LDataInd); ok {
            return casted
        }
        if casted, ok := typ.(CEMI); ok {
            return CastLDataInd(casted.Child)
        }
        if casted, ok := typ.(*CEMI); ok {
            return CastLDataInd(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *LDataInd) GetTypeName() string {
    return "LDataInd"
}

func (m *LDataInd) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (additionalInformationLength)
    lengthInBits += 8

    // Array field
    if len(m.AdditionalInformation) > 0 {
        for _, element := range m.AdditionalInformation {
            lengthInBits += element.LengthInBits()
        }
    }

    // Simple field (dataFrame)
    lengthInBits += m.DataFrame.LengthInBits()

    return lengthInBits
}

func (m *LDataInd) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func LDataIndParse(io *utils.ReadBuffer) (*CEMI, error) {

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

    // Simple Field (dataFrame)
    dataFrame, _dataFrameErr := LDataFrameParse(io)
    if _dataFrameErr != nil {
        return nil, errors.New("Error parsing 'dataFrame' field " + _dataFrameErr.Error())
    }

    // Create a partially initialized instance
    _child := &LDataInd{
        AdditionalInformationLength: additionalInformationLength,
        AdditionalInformation: additionalInformation,
        DataFrame: dataFrame,
        Parent: &CEMI{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *LDataInd) Serialize(io utils.WriteBuffer) error {
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

    // Simple Field (dataFrame)
    _dataFrameErr := m.DataFrame.Serialize(io)
    if _dataFrameErr != nil {
        return errors.New("Error serializing 'dataFrame' field " + _dataFrameErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *LDataInd) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "dataFrame":
                var dt *LDataFrame
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.DataFrame = dt
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

func (m *LDataInd) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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
    if err := e.EncodeElement(m.DataFrame, xml.StartElement{Name: xml.Name{Local: "dataFrame"}}); err != nil {
        return err
    }
    return nil
}

