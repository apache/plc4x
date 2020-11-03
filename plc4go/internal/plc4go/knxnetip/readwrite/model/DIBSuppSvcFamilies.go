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
type DIBSuppSvcFamilies struct {
    DescriptionType uint8
    ServiceIds []IServiceId

}

// The corresponding interface
type IDIBSuppSvcFamilies interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewDIBSuppSvcFamilies(descriptionType uint8, serviceIds []IServiceId) spi.Message {
    return &DIBSuppSvcFamilies{DescriptionType: descriptionType, ServiceIds: serviceIds}
}

func CastIDIBSuppSvcFamilies(structType interface{}) IDIBSuppSvcFamilies {
    castFunc := func(typ interface{}) IDIBSuppSvcFamilies {
        if iDIBSuppSvcFamilies, ok := typ.(IDIBSuppSvcFamilies); ok {
            return iDIBSuppSvcFamilies
        }
        return nil
    }
    return castFunc(structType)
}

func CastDIBSuppSvcFamilies(structType interface{}) DIBSuppSvcFamilies {
    castFunc := func(typ interface{}) DIBSuppSvcFamilies {
        if sDIBSuppSvcFamilies, ok := typ.(DIBSuppSvcFamilies); ok {
            return sDIBSuppSvcFamilies
        }
        if sDIBSuppSvcFamilies, ok := typ.(*DIBSuppSvcFamilies); ok {
            return *sDIBSuppSvcFamilies
        }
        return DIBSuppSvcFamilies{}
    }
    return castFunc(structType)
}

func (m DIBSuppSvcFamilies) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Implicit Field (structureLength)
    lengthInBits += 8

    // Simple field (descriptionType)
    lengthInBits += 8

    // Array field
    if len(m.ServiceIds) > 0 {
        for _, element := range m.ServiceIds {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m DIBSuppSvcFamilies) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DIBSuppSvcFamiliesParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _structureLengthErr := io.ReadUint8(8)
    if _structureLengthErr != nil {
        return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Simple Field (descriptionType)
    descriptionType, _descriptionTypeErr := io.ReadUint8(8)
    if _descriptionTypeErr != nil {
        return nil, errors.New("Error parsing 'descriptionType' field " + _descriptionTypeErr.Error())
    }

    // Array field (serviceIds)
    // Count array
    serviceIds := make([]IServiceId, uint16(3))
    for curItem := uint16(0); curItem < uint16(uint16(3)); curItem++ {

        _message, _err := ServiceIdParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'serviceIds' field " + _err.Error())
        }
        var _item IServiceId
        _item, _ok := _message.(IServiceId)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ServiceId")
        }
        serviceIds[curItem] = _item
    }

    // Create the instance
    return NewDIBSuppSvcFamilies(descriptionType, serviceIds), nil
}

func (m DIBSuppSvcFamilies) Serialize(io utils.WriteBuffer) error {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    structureLength := uint8(uint8(m.LengthInBytes()))
    _structureLengthErr := io.WriteUint8(8, (structureLength))
    if _structureLengthErr != nil {
        return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Simple Field (descriptionType)
    descriptionType := uint8(m.DescriptionType)
    _descriptionTypeErr := io.WriteUint8(8, (descriptionType))
    if _descriptionTypeErr != nil {
        return errors.New("Error serializing 'descriptionType' field " + _descriptionTypeErr.Error())
    }

    // Array Field (serviceIds)
    if m.ServiceIds != nil {
        for _, _element := range m.ServiceIds {
            _elementErr := _element.Serialize(io)
            if _elementErr != nil {
                return errors.New("Error serializing 'serviceIds' field " + _elementErr.Error())
            }
        }
    }

    return nil
}

func (m *DIBSuppSvcFamilies) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "descriptionType":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DescriptionType = data
            case "serviceIds":
                var _values []IServiceId
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpCore":
                        var dt *KnxNetIpCore
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        _values = append(_values, dt)
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpDeviceManagement":
                        var dt *KnxNetIpDeviceManagement
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        _values = append(_values, dt)
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpTunneling":
                        var dt *KnxNetIpTunneling
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        _values = append(_values, dt)
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetRemoteLogging":
                        var dt *KnxNetRemoteLogging
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        _values = append(_values, dt)
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetRemoteConfigurationAndDiagnosis":
                        var dt *KnxNetRemoteConfigurationAndDiagnosis
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        _values = append(_values, dt)
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetObjectServer":
                        var dt *KnxNetObjectServer
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        _values = append(_values, dt)
                    }
                    m.ServiceIds = _values
            }
        }
    }
}

func (m DIBSuppSvcFamilies) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.DIBSuppSvcFamilies"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DescriptionType, xml.StartElement{Name: xml.Name{Local: "descriptionType"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "serviceIds"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ServiceIds, xml.StartElement{Name: xml.Name{Local: "serviceIds"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "serviceIds"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

