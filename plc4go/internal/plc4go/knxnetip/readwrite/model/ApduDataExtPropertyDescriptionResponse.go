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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    log "github.com/sirupsen/logrus"
    "io"
)

// The data-structure of this message
type ApduDataExtPropertyDescriptionResponse struct {
    ObjectIndex uint8
    PropertyId uint8
    Index uint8
    WriteEnabled bool
    PropertyDataType KnxPropertyDataType
    MaxNrOfElements uint16
    ReadLevel AccessLevel
    WriteLevel AccessLevel
    Parent *ApduDataExt
    IApduDataExtPropertyDescriptionResponse
}

// The corresponding interface
type IApduDataExtPropertyDescriptionResponse interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ApduDataExtPropertyDescriptionResponse) ExtApciType() uint8 {
    return 0x19
}


func (m *ApduDataExtPropertyDescriptionResponse) InitializeParent(parent *ApduDataExt) {
}

func NewApduDataExtPropertyDescriptionResponse(objectIndex uint8, propertyId uint8, index uint8, writeEnabled bool, propertyDataType KnxPropertyDataType, maxNrOfElements uint16, readLevel AccessLevel, writeLevel AccessLevel) *ApduDataExt {
    child := &ApduDataExtPropertyDescriptionResponse{
        ObjectIndex: objectIndex,
        PropertyId: propertyId,
        Index: index,
        WriteEnabled: writeEnabled,
        PropertyDataType: propertyDataType,
        MaxNrOfElements: maxNrOfElements,
        ReadLevel: readLevel,
        WriteLevel: writeLevel,
        Parent: NewApduDataExt(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastApduDataExtPropertyDescriptionResponse(structType interface{}) *ApduDataExtPropertyDescriptionResponse {
    castFunc := func(typ interface{}) *ApduDataExtPropertyDescriptionResponse {
        if casted, ok := typ.(ApduDataExtPropertyDescriptionResponse); ok {
            return &casted
        }
        if casted, ok := typ.(*ApduDataExtPropertyDescriptionResponse); ok {
            return casted
        }
        if casted, ok := typ.(ApduDataExt); ok {
            return CastApduDataExtPropertyDescriptionResponse(casted.Child)
        }
        if casted, ok := typ.(*ApduDataExt); ok {
            return CastApduDataExtPropertyDescriptionResponse(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ApduDataExtPropertyDescriptionResponse) GetTypeName() string {
    return "ApduDataExtPropertyDescriptionResponse"
}

func (m *ApduDataExtPropertyDescriptionResponse) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (objectIndex)
    lengthInBits += 8

    // Simple field (propertyId)
    lengthInBits += 8

    // Simple field (index)
    lengthInBits += 8

    // Simple field (writeEnabled)
    lengthInBits += 1

    // Reserved Field (reserved)
    lengthInBits += 1

    // Simple field (propertyDataType)
    lengthInBits += 6

    // Reserved Field (reserved)
    lengthInBits += 4

    // Simple field (maxNrOfElements)
    lengthInBits += 12

    // Simple field (readLevel)
    lengthInBits += 4

    // Simple field (writeLevel)
    lengthInBits += 4

    return lengthInBits
}

func (m *ApduDataExtPropertyDescriptionResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ApduDataExtPropertyDescriptionResponseParse(io *utils.ReadBuffer) (*ApduDataExt, error) {

    // Simple Field (objectIndex)
    objectIndex, _objectIndexErr := io.ReadUint8(8)
    if _objectIndexErr != nil {
        return nil, errors.New("Error parsing 'objectIndex' field " + _objectIndexErr.Error())
    }

    // Simple Field (propertyId)
    propertyId, _propertyIdErr := io.ReadUint8(8)
    if _propertyIdErr != nil {
        return nil, errors.New("Error parsing 'propertyId' field " + _propertyIdErr.Error())
    }

    // Simple Field (index)
    index, _indexErr := io.ReadUint8(8)
    if _indexErr != nil {
        return nil, errors.New("Error parsing 'index' field " + _indexErr.Error())
    }

    // Simple Field (writeEnabled)
    writeEnabled, _writeEnabledErr := io.ReadBit()
    if _writeEnabledErr != nil {
        return nil, errors.New("Error parsing 'writeEnabled' field " + _writeEnabledErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(1)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0x0) {
            log.WithFields(log.Fields{
                "expected value": uint8(0x0),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (propertyDataType)
    propertyDataType, _propertyDataTypeErr := KnxPropertyDataTypeParse(io)
    if _propertyDataTypeErr != nil {
        return nil, errors.New("Error parsing 'propertyDataType' field " + _propertyDataTypeErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(4)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0x0) {
            log.WithFields(log.Fields{
                "expected value": uint8(0x0),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (maxNrOfElements)
    maxNrOfElements, _maxNrOfElementsErr := io.ReadUint16(12)
    if _maxNrOfElementsErr != nil {
        return nil, errors.New("Error parsing 'maxNrOfElements' field " + _maxNrOfElementsErr.Error())
    }

    // Simple Field (readLevel)
    readLevel, _readLevelErr := AccessLevelParse(io)
    if _readLevelErr != nil {
        return nil, errors.New("Error parsing 'readLevel' field " + _readLevelErr.Error())
    }

    // Simple Field (writeLevel)
    writeLevel, _writeLevelErr := AccessLevelParse(io)
    if _writeLevelErr != nil {
        return nil, errors.New("Error parsing 'writeLevel' field " + _writeLevelErr.Error())
    }

    // Create a partially initialized instance
    _child := &ApduDataExtPropertyDescriptionResponse{
        ObjectIndex: objectIndex,
        PropertyId: propertyId,
        Index: index,
        WriteEnabled: writeEnabled,
        PropertyDataType: propertyDataType,
        MaxNrOfElements: maxNrOfElements,
        ReadLevel: readLevel,
        WriteLevel: writeLevel,
        Parent: &ApduDataExt{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ApduDataExtPropertyDescriptionResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (objectIndex)
    objectIndex := uint8(m.ObjectIndex)
    _objectIndexErr := io.WriteUint8(8, (objectIndex))
    if _objectIndexErr != nil {
        return errors.New("Error serializing 'objectIndex' field " + _objectIndexErr.Error())
    }

    // Simple Field (propertyId)
    propertyId := uint8(m.PropertyId)
    _propertyIdErr := io.WriteUint8(8, (propertyId))
    if _propertyIdErr != nil {
        return errors.New("Error serializing 'propertyId' field " + _propertyIdErr.Error())
    }

    // Simple Field (index)
    index := uint8(m.Index)
    _indexErr := io.WriteUint8(8, (index))
    if _indexErr != nil {
        return errors.New("Error serializing 'index' field " + _indexErr.Error())
    }

    // Simple Field (writeEnabled)
    writeEnabled := bool(m.WriteEnabled)
    _writeEnabledErr := io.WriteBit((writeEnabled))
    if _writeEnabledErr != nil {
        return errors.New("Error serializing 'writeEnabled' field " + _writeEnabledErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(1, uint8(0x0))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (propertyDataType)
    _propertyDataTypeErr := m.PropertyDataType.Serialize(io)
    if _propertyDataTypeErr != nil {
        return errors.New("Error serializing 'propertyDataType' field " + _propertyDataTypeErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(4, uint8(0x0))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (maxNrOfElements)
    maxNrOfElements := uint16(m.MaxNrOfElements)
    _maxNrOfElementsErr := io.WriteUint16(12, (maxNrOfElements))
    if _maxNrOfElementsErr != nil {
        return errors.New("Error serializing 'maxNrOfElements' field " + _maxNrOfElementsErr.Error())
    }

    // Simple Field (readLevel)
    _readLevelErr := m.ReadLevel.Serialize(io)
    if _readLevelErr != nil {
        return errors.New("Error serializing 'readLevel' field " + _readLevelErr.Error())
    }

    // Simple Field (writeLevel)
    _writeLevelErr := m.WriteLevel.Serialize(io)
    if _writeLevelErr != nil {
        return errors.New("Error serializing 'writeLevel' field " + _writeLevelErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ApduDataExtPropertyDescriptionResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "objectIndex":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ObjectIndex = data
            case "propertyId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.PropertyId = data
            case "index":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Index = data
            case "writeEnabled":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.WriteEnabled = data
            case "propertyDataType":
                var data KnxPropertyDataType
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.PropertyDataType = data
            case "maxNrOfElements":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MaxNrOfElements = data
            case "readLevel":
                var data AccessLevel
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ReadLevel = data
            case "writeLevel":
                var data AccessLevel
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.WriteLevel = data
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

func (m *ApduDataExtPropertyDescriptionResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.ObjectIndex, xml.StartElement{Name: xml.Name{Local: "objectIndex"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.PropertyId, xml.StartElement{Name: xml.Name{Local: "propertyId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Index, xml.StartElement{Name: xml.Name{Local: "index"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.WriteEnabled, xml.StartElement{Name: xml.Name{Local: "writeEnabled"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.PropertyDataType, xml.StartElement{Name: xml.Name{Local: "propertyDataType"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MaxNrOfElements, xml.StartElement{Name: xml.Name{Local: "maxNrOfElements"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ReadLevel, xml.StartElement{Name: xml.Name{Local: "readLevel"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.WriteLevel, xml.StartElement{Name: xml.Name{Local: "writeLevel"}}); err != nil {
        return err
    }
    return nil
}

