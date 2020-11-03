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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMIMPropReadCon struct {
    InterfaceObjectType uint16
    ObjectInstance uint8
    PropertyId uint8
    NumberOfElements uint8
    StartIndex uint16
    Unknown uint16
    CEMI
}

// The corresponding interface
type ICEMIMPropReadCon interface {
    ICEMI
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIMPropReadCon) MessageCode() uint8 {
    return 0xFB
}

func (m CEMIMPropReadCon) initialize() spi.Message {
    return m
}

func NewCEMIMPropReadCon(interfaceObjectType uint16, objectInstance uint8, propertyId uint8, numberOfElements uint8, startIndex uint16, unknown uint16) CEMIInitializer {
    return &CEMIMPropReadCon{InterfaceObjectType: interfaceObjectType, ObjectInstance: objectInstance, PropertyId: propertyId, NumberOfElements: numberOfElements, StartIndex: startIndex, Unknown: unknown}
}

func CastICEMIMPropReadCon(structType interface{}) ICEMIMPropReadCon {
    castFunc := func(typ interface{}) ICEMIMPropReadCon {
        if iCEMIMPropReadCon, ok := typ.(ICEMIMPropReadCon); ok {
            return iCEMIMPropReadCon
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIMPropReadCon(structType interface{}) CEMIMPropReadCon {
    castFunc := func(typ interface{}) CEMIMPropReadCon {
        if sCEMIMPropReadCon, ok := typ.(CEMIMPropReadCon); ok {
            return sCEMIMPropReadCon
        }
        if sCEMIMPropReadCon, ok := typ.(*CEMIMPropReadCon); ok {
            return *sCEMIMPropReadCon
        }
        return CEMIMPropReadCon{}
    }
    return castFunc(structType)
}

func (m CEMIMPropReadCon) LengthInBits() uint16 {
    var lengthInBits uint16 = m.CEMI.LengthInBits()

    // Simple field (interfaceObjectType)
    lengthInBits += 16

    // Simple field (objectInstance)
    lengthInBits += 8

    // Simple field (propertyId)
    lengthInBits += 8

    // Simple field (numberOfElements)
    lengthInBits += 4

    // Simple field (startIndex)
    lengthInBits += 12

    // Simple field (unknown)
    lengthInBits += 16

    return lengthInBits
}

func (m CEMIMPropReadCon) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIMPropReadConParse(io *utils.ReadBuffer) (CEMIInitializer, error) {

    // Simple Field (interfaceObjectType)
    interfaceObjectType, _interfaceObjectTypeErr := io.ReadUint16(16)
    if _interfaceObjectTypeErr != nil {
        return nil, errors.New("Error parsing 'interfaceObjectType' field " + _interfaceObjectTypeErr.Error())
    }

    // Simple Field (objectInstance)
    objectInstance, _objectInstanceErr := io.ReadUint8(8)
    if _objectInstanceErr != nil {
        return nil, errors.New("Error parsing 'objectInstance' field " + _objectInstanceErr.Error())
    }

    // Simple Field (propertyId)
    propertyId, _propertyIdErr := io.ReadUint8(8)
    if _propertyIdErr != nil {
        return nil, errors.New("Error parsing 'propertyId' field " + _propertyIdErr.Error())
    }

    // Simple Field (numberOfElements)
    numberOfElements, _numberOfElementsErr := io.ReadUint8(4)
    if _numberOfElementsErr != nil {
        return nil, errors.New("Error parsing 'numberOfElements' field " + _numberOfElementsErr.Error())
    }

    // Simple Field (startIndex)
    startIndex, _startIndexErr := io.ReadUint16(12)
    if _startIndexErr != nil {
        return nil, errors.New("Error parsing 'startIndex' field " + _startIndexErr.Error())
    }

    // Simple Field (unknown)
    unknown, _unknownErr := io.ReadUint16(16)
    if _unknownErr != nil {
        return nil, errors.New("Error parsing 'unknown' field " + _unknownErr.Error())
    }

    // Create the instance
    return NewCEMIMPropReadCon(interfaceObjectType, objectInstance, propertyId, numberOfElements, startIndex, unknown), nil
}

func (m CEMIMPropReadCon) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (interfaceObjectType)
    interfaceObjectType := uint16(m.InterfaceObjectType)
    _interfaceObjectTypeErr := io.WriteUint16(16, (interfaceObjectType))
    if _interfaceObjectTypeErr != nil {
        return errors.New("Error serializing 'interfaceObjectType' field " + _interfaceObjectTypeErr.Error())
    }

    // Simple Field (objectInstance)
    objectInstance := uint8(m.ObjectInstance)
    _objectInstanceErr := io.WriteUint8(8, (objectInstance))
    if _objectInstanceErr != nil {
        return errors.New("Error serializing 'objectInstance' field " + _objectInstanceErr.Error())
    }

    // Simple Field (propertyId)
    propertyId := uint8(m.PropertyId)
    _propertyIdErr := io.WriteUint8(8, (propertyId))
    if _propertyIdErr != nil {
        return errors.New("Error serializing 'propertyId' field " + _propertyIdErr.Error())
    }

    // Simple Field (numberOfElements)
    numberOfElements := uint8(m.NumberOfElements)
    _numberOfElementsErr := io.WriteUint8(4, (numberOfElements))
    if _numberOfElementsErr != nil {
        return errors.New("Error serializing 'numberOfElements' field " + _numberOfElementsErr.Error())
    }

    // Simple Field (startIndex)
    startIndex := uint16(m.StartIndex)
    _startIndexErr := io.WriteUint16(12, (startIndex))
    if _startIndexErr != nil {
        return errors.New("Error serializing 'startIndex' field " + _startIndexErr.Error())
    }

    // Simple Field (unknown)
    unknown := uint16(m.Unknown)
    _unknownErr := io.WriteUint16(16, (unknown))
    if _unknownErr != nil {
        return errors.New("Error serializing 'unknown' field " + _unknownErr.Error())
    }

        return nil
    }
    return CEMISerialize(io, m.CEMI, CastICEMI(m), ser)
}

func (m *CEMIMPropReadCon) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "interfaceObjectType":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.InterfaceObjectType = data
            case "objectInstance":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ObjectInstance = data
            case "propertyId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.PropertyId = data
            case "numberOfElements":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.NumberOfElements = data
            case "startIndex":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.StartIndex = data
            case "unknown":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Unknown = data
            }
        }
    }
}

func (m CEMIMPropReadCon) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.CEMIMPropReadCon"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.InterfaceObjectType, xml.StartElement{Name: xml.Name{Local: "interfaceObjectType"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ObjectInstance, xml.StartElement{Name: xml.Name{Local: "objectInstance"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.PropertyId, xml.StartElement{Name: xml.Name{Local: "propertyId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.NumberOfElements, xml.StartElement{Name: xml.Name{Local: "numberOfElements"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.StartIndex, xml.StartElement{Name: xml.Name{Local: "startIndex"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Unknown, xml.StartElement{Name: xml.Name{Local: "unknown"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

