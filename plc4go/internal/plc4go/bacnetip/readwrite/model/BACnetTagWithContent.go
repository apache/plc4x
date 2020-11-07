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
    "strconv"
)

// Constant values.
const BACnetTagWithContent_OPENTAG uint8 = 0x2e
const BACnetTagWithContent_CLOSINGTAG uint8 = 0x2f

// The data-structure of this message
type BACnetTagWithContent struct {
    TypeOrTagNumber uint8
    ContextSpecificTag uint8
    LengthValueType uint8
    ExtTagNumber *uint8
    ExtLength *uint8
    PropertyIdentifier []uint8
    Value *BACnetTag
    IBACnetTagWithContent
}

// The corresponding interface
type IBACnetTagWithContent interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewBACnetTagWithContent(typeOrTagNumber uint8, contextSpecificTag uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8, propertyIdentifier []uint8, value *BACnetTag) *BACnetTagWithContent {
    return &BACnetTagWithContent{TypeOrTagNumber: typeOrTagNumber, ContextSpecificTag: contextSpecificTag, LengthValueType: lengthValueType, ExtTagNumber: extTagNumber, ExtLength: extLength, PropertyIdentifier: propertyIdentifier, Value: value}
}

func CastBACnetTagWithContent(structType interface{}) BACnetTagWithContent {
    castFunc := func(typ interface{}) BACnetTagWithContent {
        if casted, ok := typ.(BACnetTagWithContent); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetTagWithContent); ok {
            return *casted
        }
        return BACnetTagWithContent{}
    }
    return castFunc(structType)
}

func (m *BACnetTagWithContent) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (typeOrTagNumber)
    lengthInBits += 4

    // Simple field (contextSpecificTag)
    lengthInBits += 1

    // Simple field (lengthValueType)
    lengthInBits += 3

    // Optional Field (extTagNumber)
    if m.ExtTagNumber != nil {
        lengthInBits += 8
    }

    // Optional Field (extLength)
    if m.ExtLength != nil {
        lengthInBits += 8
    }

    // Array field
    if len(m.PropertyIdentifier) > 0 {
        lengthInBits += 8 * uint16(len(m.PropertyIdentifier))
    }

    // Const Field (openTag)
    lengthInBits += 8

    // Simple field (value)
    lengthInBits += m.Value.LengthInBits()

    // Const Field (closingTag)
    lengthInBits += 8

    return lengthInBits
}

func (m *BACnetTagWithContent) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagWithContentParse(io *utils.ReadBuffer) (*BACnetTagWithContent, error) {

    // Simple Field (typeOrTagNumber)
    typeOrTagNumber, _typeOrTagNumberErr := io.ReadUint8(4)
    if _typeOrTagNumberErr != nil {
        return nil, errors.New("Error parsing 'typeOrTagNumber' field " + _typeOrTagNumberErr.Error())
    }

    // Simple Field (contextSpecificTag)
    contextSpecificTag, _contextSpecificTagErr := io.ReadUint8(1)
    if _contextSpecificTagErr != nil {
        return nil, errors.New("Error parsing 'contextSpecificTag' field " + _contextSpecificTagErr.Error())
    }

    // Simple Field (lengthValueType)
    lengthValueType, _lengthValueTypeErr := io.ReadUint8(3)
    if _lengthValueTypeErr != nil {
        return nil, errors.New("Error parsing 'lengthValueType' field " + _lengthValueTypeErr.Error())
    }

    // Optional Field (extTagNumber) (Can be skipped, if a given expression evaluates to false)
    var extTagNumber *uint8 = nil
    if bool((typeOrTagNumber) == ((15))) {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'extTagNumber' field " + _err.Error())
        }

        extTagNumber = &_val
    }

    // Optional Field (extLength) (Can be skipped, if a given expression evaluates to false)
    var extLength *uint8 = nil
    if bool((lengthValueType) == ((5))) {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'extLength' field " + _err.Error())
        }

        extLength = &_val
    }

    // Array field (propertyIdentifier)
    // Length array
    propertyIdentifier := make([]uint8, 0)
    _propertyIdentifierLength := utils.InlineIf(bool(bool((lengthValueType) == ((5)))), uint16((*extLength)), uint16(lengthValueType))
    _propertyIdentifierEndPos := io.GetPos() + uint16(_propertyIdentifierLength)
    for ;io.GetPos() < _propertyIdentifierEndPos; {
        _item, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'propertyIdentifier' field " + _err.Error())
        }
        propertyIdentifier = append(propertyIdentifier, _item)
    }

    // Const Field (openTag)
    openTag, _openTagErr := io.ReadUint8(8)
    if _openTagErr != nil {
        return nil, errors.New("Error parsing 'openTag' field " + _openTagErr.Error())
    }
    if openTag != BACnetTagWithContent_OPENTAG {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetTagWithContent_OPENTAG)) + " but got " + strconv.Itoa(int(openTag)))
    }

    // Simple Field (value)
    value, _valueErr := BACnetTagParse(io)
    if _valueErr != nil {
        return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
    }

    // Const Field (closingTag)
    closingTag, _closingTagErr := io.ReadUint8(8)
    if _closingTagErr != nil {
        return nil, errors.New("Error parsing 'closingTag' field " + _closingTagErr.Error())
    }
    if closingTag != BACnetTagWithContent_CLOSINGTAG {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetTagWithContent_CLOSINGTAG)) + " but got " + strconv.Itoa(int(closingTag)))
    }

    // Create the instance
    return NewBACnetTagWithContent(typeOrTagNumber, contextSpecificTag, lengthValueType, extTagNumber, extLength, propertyIdentifier, value), nil
}

func (m *BACnetTagWithContent) Serialize(io utils.WriteBuffer) error {

    // Simple Field (typeOrTagNumber)
    typeOrTagNumber := uint8(m.TypeOrTagNumber)
    _typeOrTagNumberErr := io.WriteUint8(4, (typeOrTagNumber))
    if _typeOrTagNumberErr != nil {
        return errors.New("Error serializing 'typeOrTagNumber' field " + _typeOrTagNumberErr.Error())
    }

    // Simple Field (contextSpecificTag)
    contextSpecificTag := uint8(m.ContextSpecificTag)
    _contextSpecificTagErr := io.WriteUint8(1, (contextSpecificTag))
    if _contextSpecificTagErr != nil {
        return errors.New("Error serializing 'contextSpecificTag' field " + _contextSpecificTagErr.Error())
    }

    // Simple Field (lengthValueType)
    lengthValueType := uint8(m.LengthValueType)
    _lengthValueTypeErr := io.WriteUint8(3, (lengthValueType))
    if _lengthValueTypeErr != nil {
        return errors.New("Error serializing 'lengthValueType' field " + _lengthValueTypeErr.Error())
    }

    // Optional Field (extTagNumber) (Can be skipped, if the value is null)
    var extTagNumber *uint8 = nil
    if m.ExtTagNumber != nil {
        extTagNumber = m.ExtTagNumber
        _extTagNumberErr := io.WriteUint8(8, *(extTagNumber))
        if _extTagNumberErr != nil {
            return errors.New("Error serializing 'extTagNumber' field " + _extTagNumberErr.Error())
        }
    }

    // Optional Field (extLength) (Can be skipped, if the value is null)
    var extLength *uint8 = nil
    if m.ExtLength != nil {
        extLength = m.ExtLength
        _extLengthErr := io.WriteUint8(8, *(extLength))
        if _extLengthErr != nil {
            return errors.New("Error serializing 'extLength' field " + _extLengthErr.Error())
        }
    }

    // Array Field (propertyIdentifier)
    if m.PropertyIdentifier != nil {
        for _, _element := range m.PropertyIdentifier {
            _elementErr := io.WriteUint8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'propertyIdentifier' field " + _elementErr.Error())
            }
        }
    }

    // Const Field (openTag)
    _openTagErr := io.WriteUint8(8, 0x2e)
    if _openTagErr != nil {
        return errors.New("Error serializing 'openTag' field " + _openTagErr.Error())
    }

    // Simple Field (value)
    _valueErr := m.Value.Serialize(io)
    if _valueErr != nil {
        return errors.New("Error serializing 'value' field " + _valueErr.Error())
    }

    // Const Field (closingTag)
    _closingTagErr := io.WriteUint8(8, 0x2f)
    if _closingTagErr != nil {
        return errors.New("Error serializing 'closingTag' field " + _closingTagErr.Error())
    }

    return nil
}

func (m *BACnetTagWithContent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "typeOrTagNumber":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TypeOrTagNumber = data
            case "contextSpecificTag":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ContextSpecificTag = data
            case "lengthValueType":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.LengthValueType = data
            case "extTagNumber":
                var data *uint8
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ExtTagNumber = data
            case "extLength":
                var data *uint8
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ExtLength = data
            case "propertyIdentifier":
                var data []uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.PropertyIdentifier = data
            case "value":
                var dt *BACnetTag
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.Value = dt
            }
        }
    }
}

func (m *BACnetTagWithContent) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.bacnetip.readwrite.BACnetTagWithContent"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.TypeOrTagNumber, xml.StartElement{Name: xml.Name{Local: "typeOrTagNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ContextSpecificTag, xml.StartElement{Name: xml.Name{Local: "contextSpecificTag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.LengthValueType, xml.StartElement{Name: xml.Name{Local: "lengthValueType"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ExtTagNumber, xml.StartElement{Name: xml.Name{Local: "extTagNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ExtLength, xml.StartElement{Name: xml.Name{Local: "extLength"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "propertyIdentifier"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.PropertyIdentifier, xml.StartElement{Name: xml.Name{Local: "propertyIdentifier"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "propertyIdentifier"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Value, xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

