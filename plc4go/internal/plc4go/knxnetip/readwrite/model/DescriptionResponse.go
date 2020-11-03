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
    "reflect"
)

// The data-structure of this message
type DescriptionResponse struct {
    DibDeviceInfo IDIBDeviceInfo
    DibSuppSvcFamilies IDIBSuppSvcFamilies
    KNXNetIPMessage
}

// The corresponding interface
type IDescriptionResponse interface {
    IKNXNetIPMessage
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m DescriptionResponse) MsgType() uint16 {
    return 0x0204
}

func (m DescriptionResponse) initialize() spi.Message {
    return m
}

func NewDescriptionResponse(dibDeviceInfo IDIBDeviceInfo, dibSuppSvcFamilies IDIBSuppSvcFamilies) KNXNetIPMessageInitializer {
    return &DescriptionResponse{DibDeviceInfo: dibDeviceInfo, DibSuppSvcFamilies: dibSuppSvcFamilies}
}

func CastIDescriptionResponse(structType interface{}) IDescriptionResponse {
    castFunc := func(typ interface{}) IDescriptionResponse {
        if iDescriptionResponse, ok := typ.(IDescriptionResponse); ok {
            return iDescriptionResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastDescriptionResponse(structType interface{}) DescriptionResponse {
    castFunc := func(typ interface{}) DescriptionResponse {
        if sDescriptionResponse, ok := typ.(DescriptionResponse); ok {
            return sDescriptionResponse
        }
        if sDescriptionResponse, ok := typ.(*DescriptionResponse); ok {
            return *sDescriptionResponse
        }
        return DescriptionResponse{}
    }
    return castFunc(structType)
}

func (m DescriptionResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

    // Simple field (dibDeviceInfo)
    lengthInBits += m.DibDeviceInfo.LengthInBits()

    // Simple field (dibSuppSvcFamilies)
    lengthInBits += m.DibSuppSvcFamilies.LengthInBits()

    return lengthInBits
}

func (m DescriptionResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DescriptionResponseParse(io *utils.ReadBuffer) (KNXNetIPMessageInitializer, error) {

    // Simple Field (dibDeviceInfo)
    _dibDeviceInfoMessage, _err := DIBDeviceInfoParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'dibDeviceInfo'. " + _err.Error())
    }
    var dibDeviceInfo IDIBDeviceInfo
    dibDeviceInfo, _dibDeviceInfoOk := _dibDeviceInfoMessage.(IDIBDeviceInfo)
    if !_dibDeviceInfoOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_dibDeviceInfoMessage).Name() + " to IDIBDeviceInfo")
    }

    // Simple Field (dibSuppSvcFamilies)
    _dibSuppSvcFamiliesMessage, _err := DIBSuppSvcFamiliesParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'dibSuppSvcFamilies'. " + _err.Error())
    }
    var dibSuppSvcFamilies IDIBSuppSvcFamilies
    dibSuppSvcFamilies, _dibSuppSvcFamiliesOk := _dibSuppSvcFamiliesMessage.(IDIBSuppSvcFamilies)
    if !_dibSuppSvcFamiliesOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_dibSuppSvcFamiliesMessage).Name() + " to IDIBSuppSvcFamilies")
    }

    // Create the instance
    return NewDescriptionResponse(dibDeviceInfo, dibSuppSvcFamilies), nil
}

func (m DescriptionResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (dibDeviceInfo)
    dibDeviceInfo := CastIDIBDeviceInfo(m.DibDeviceInfo)
    _dibDeviceInfoErr := dibDeviceInfo.Serialize(io)
    if _dibDeviceInfoErr != nil {
        return errors.New("Error serializing 'dibDeviceInfo' field " + _dibDeviceInfoErr.Error())
    }

    // Simple Field (dibSuppSvcFamilies)
    dibSuppSvcFamilies := CastIDIBSuppSvcFamilies(m.DibSuppSvcFamilies)
    _dibSuppSvcFamiliesErr := dibSuppSvcFamilies.Serialize(io)
    if _dibSuppSvcFamiliesErr != nil {
        return errors.New("Error serializing 'dibSuppSvcFamilies' field " + _dibSuppSvcFamiliesErr.Error())
    }

        return nil
    }
    return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}

func (m *DescriptionResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "dibDeviceInfo":
                var data *DIBDeviceInfo
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DibDeviceInfo = CastIDIBDeviceInfo(data)
            case "dibSuppSvcFamilies":
                var data *DIBSuppSvcFamilies
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DibSuppSvcFamilies = CastIDIBSuppSvcFamilies(data)
            }
        }
    }
}

func (m DescriptionResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.DescriptionResponse"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DibDeviceInfo, xml.StartElement{Name: xml.Name{Local: "dibDeviceInfo"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DibSuppSvcFamilies, xml.StartElement{Name: xml.Name{Local: "dibSuppSvcFamilies"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

