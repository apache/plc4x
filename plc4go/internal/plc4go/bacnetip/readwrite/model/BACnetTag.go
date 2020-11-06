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
type BACnetTag struct {
    TypeOrTagNumber uint8
    LengthValueType uint8
    ExtTagNumber *uint8
    ExtLength *uint8
    Child IBACnetTagChild
    IBACnetTag
    IBACnetTagParent
}

// The corresponding interface
type IBACnetTag interface {
    ContextSpecificTag() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IBACnetTagParent interface {
    SerializeParent(io utils.WriteBuffer, child IBACnetTag, serializeChildFunction func() error) error
}

type IBACnetTagChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *BACnetTag, typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8)
    IBACnetTag
}

func NewBACnetTag(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) *BACnetTag {
    return &BACnetTag{TypeOrTagNumber: typeOrTagNumber, LengthValueType: lengthValueType, ExtTagNumber: extTagNumber, ExtLength: extLength}
}

func CastBACnetTag(structType interface{}) BACnetTag {
    castFunc := func(typ interface{}) BACnetTag {
        if casted, ok := typ.(BACnetTag); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetTag); ok {
            return *casted
        }
        return BACnetTag{}
    }
    return castFunc(structType)
}

func (m *BACnetTag) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (typeOrTagNumber)
    lengthInBits += 4

    // Discriminator Field (contextSpecificTag)
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

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *BACnetTag) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagParse(io *utils.ReadBuffer) (*BACnetTag, error) {

    // Simple Field (typeOrTagNumber)
    typeOrTagNumber, _typeOrTagNumberErr := io.ReadUint8(4)
    if _typeOrTagNumberErr != nil {
        return nil, errors.New("Error parsing 'typeOrTagNumber' field " + _typeOrTagNumberErr.Error())
    }

    // Discriminator Field (contextSpecificTag) (Used as input to a switch field)
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

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *BACnetTag
    var typeSwitchError error
    switch {
    case contextSpecificTag == 0 && typeOrTagNumber == 0x0:
        _parent, typeSwitchError = BACnetTagApplicationNullParse(io)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x1:
        _parent, typeSwitchError = BACnetTagApplicationBooleanParse(io)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x2:
        _parent, typeSwitchError = BACnetTagApplicationUnsignedIntegerParse(io, lengthValueType, *extLength)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x3:
        _parent, typeSwitchError = BACnetTagApplicationSignedIntegerParse(io, lengthValueType, *extLength)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x4:
        _parent, typeSwitchError = BACnetTagApplicationRealParse(io, lengthValueType, *extLength)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x5:
        _parent, typeSwitchError = BACnetTagApplicationDoubleParse(io, lengthValueType, *extLength)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x6:
        _parent, typeSwitchError = BACnetTagApplicationOctetStringParse(io)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x7:
        _parent, typeSwitchError = BACnetTagApplicationCharacterStringParse(io)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x8:
        _parent, typeSwitchError = BACnetTagApplicationBitStringParse(io, lengthValueType, *extLength)
    case contextSpecificTag == 0 && typeOrTagNumber == 0x9:
        _parent, typeSwitchError = BACnetTagApplicationEnumeratedParse(io, lengthValueType, *extLength)
    case contextSpecificTag == 0 && typeOrTagNumber == 0xA:
        _parent, typeSwitchError = BACnetTagApplicationDateParse(io)
    case contextSpecificTag == 0 && typeOrTagNumber == 0xB:
        _parent, typeSwitchError = BACnetTagApplicationTimeParse(io)
    case contextSpecificTag == 0 && typeOrTagNumber == 0xC:
        _parent, typeSwitchError = BACnetTagApplicationObjectIdentifierParse(io)
    case contextSpecificTag == 1:
        _parent, typeSwitchError = BACnetTagContextParse(io, typeOrTagNumber, *extTagNumber, lengthValueType, *extLength)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent, typeOrTagNumber, lengthValueType, extTagNumber, extLength)
    return _parent, nil
}

func (m *BACnetTag) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *BACnetTag) SerializeParent(io utils.WriteBuffer, child IBACnetTag, serializeChildFunction func() error) error {

    // Simple Field (typeOrTagNumber)
    typeOrTagNumber := uint8(m.TypeOrTagNumber)
    _typeOrTagNumberErr := io.WriteUint8(4, (typeOrTagNumber))
    if _typeOrTagNumberErr != nil {
        return errors.New("Error serializing 'typeOrTagNumber' field " + _typeOrTagNumberErr.Error())
    }

    // Discriminator Field (contextSpecificTag) (Used as input to a switch field)
    contextSpecificTag := uint8(child.ContextSpecificTag())
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

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *BACnetTag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "typeOrTagNumber":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TypeOrTagNumber = data
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
            }
        }
    }
}

func (m *BACnetTag) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetTag"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.TypeOrTagNumber, xml.StartElement{Name: xml.Name{Local: "typeOrTagNumber"}}); err != nil {
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
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

