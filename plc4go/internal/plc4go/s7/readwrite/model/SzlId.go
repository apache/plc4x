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
type SzlId struct {
    TypeClass SzlModuleTypeClass
    SublistExtract uint8
    SublistList SzlSublist
    ISzlId
}

// The corresponding interface
type ISzlId interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

func NewSzlId(typeClass SzlModuleTypeClass, sublistExtract uint8, sublistList SzlSublist) *SzlId {
    return &SzlId{TypeClass: typeClass, SublistExtract: sublistExtract, SublistList: sublistList}
}

func CastSzlId(structType interface{}) SzlId {
    castFunc := func(typ interface{}) SzlId {
        if casted, ok := typ.(SzlId); ok {
            return casted
        }
        if casted, ok := typ.(*SzlId); ok {
            return *casted
        }
        return SzlId{}
    }
    return castFunc(structType)
}

func (m *SzlId) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Enum Field (typeClass)
    lengthInBits += 4

    // Simple field (sublistExtract)
    lengthInBits += 4

    // Enum Field (sublistList)
    lengthInBits += 8

    return lengthInBits
}

func (m *SzlId) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func SzlIdParse(io *utils.ReadBuffer) (*SzlId, error) {

    // Enum field (typeClass)
    typeClass, _typeClassErr := SzlModuleTypeClassParse(io)
    if _typeClassErr != nil {
        return nil, errors.New("Error parsing 'typeClass' field " + _typeClassErr.Error())
    }

    // Simple Field (sublistExtract)
    sublistExtract, _sublistExtractErr := io.ReadUint8(4)
    if _sublistExtractErr != nil {
        return nil, errors.New("Error parsing 'sublistExtract' field " + _sublistExtractErr.Error())
    }

    // Enum field (sublistList)
    sublistList, _sublistListErr := SzlSublistParse(io)
    if _sublistListErr != nil {
        return nil, errors.New("Error parsing 'sublistList' field " + _sublistListErr.Error())
    }

    // Create the instance
    return NewSzlId(typeClass, sublistExtract, sublistList), nil
}

func (m *SzlId) Serialize(io utils.WriteBuffer) error {

    // Enum field (typeClass)
    typeClass := CastSzlModuleTypeClass(m.TypeClass)
    _typeClassErr := typeClass.Serialize(io)
    if _typeClassErr != nil {
        return errors.New("Error serializing 'typeClass' field " + _typeClassErr.Error())
    }

    // Simple Field (sublistExtract)
    sublistExtract := uint8(m.SublistExtract)
    _sublistExtractErr := io.WriteUint8(4, (sublistExtract))
    if _sublistExtractErr != nil {
        return errors.New("Error serializing 'sublistExtract' field " + _sublistExtractErr.Error())
    }

    // Enum field (sublistList)
    sublistList := CastSzlSublist(m.SublistList)
    _sublistListErr := sublistList.Serialize(io)
    if _sublistListErr != nil {
        return errors.New("Error serializing 'sublistList' field " + _sublistListErr.Error())
    }

    return nil
}

func (m *SzlId) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "typeClass":
                var data SzlModuleTypeClass
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TypeClass = data
            case "sublistExtract":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SublistExtract = data
            case "sublistList":
                var data SzlSublist
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SublistList = data
            }
        }
    }
}

func (m *SzlId) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.SzlId"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.TypeClass, xml.StartElement{Name: xml.Name{Local: "typeClass"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SublistExtract, xml.StartElement{Name: xml.Name{Local: "sublistExtract"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SublistList, xml.StartElement{Name: xml.Name{Local: "sublistList"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

