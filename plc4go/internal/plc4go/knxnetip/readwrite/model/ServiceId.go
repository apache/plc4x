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
    "reflect"
    "strings"
)

// The data-structure of this message
type ServiceId struct {
    Child IServiceIdChild
    IServiceId
    IServiceIdParent
}

// The corresponding interface
type IServiceId interface {
    ServiceType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type IServiceIdParent interface {
    SerializeParent(io utils.WriteBuffer, child IServiceId, serializeChildFunction func() error) error
    GetTypeName() string
}

type IServiceIdChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *ServiceId)
    GetTypeName() string
    IServiceId
}

func NewServiceId() *ServiceId {
    return &ServiceId{}
}

func CastServiceId(structType interface{}) *ServiceId {
    castFunc := func(typ interface{}) *ServiceId {
        if casted, ok := typ.(ServiceId); ok {
            return &casted
        }
        if casted, ok := typ.(*ServiceId); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ServiceId) GetTypeName() string {
    return "ServiceId"
}

func (m *ServiceId) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (serviceType)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *ServiceId) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ServiceIdParse(io *utils.ReadBuffer) (*ServiceId, error) {

    // Discriminator Field (serviceType) (Used as input to a switch field)
    serviceType, _serviceTypeErr := io.ReadUint8(8)
    if _serviceTypeErr != nil {
        return nil, errors.New("Error parsing 'serviceType' field " + _serviceTypeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *ServiceId
    var typeSwitchError error
    switch {
    case serviceType == 0x02:
        _parent, typeSwitchError = KnxNetIpCoreParse(io)
    case serviceType == 0x03:
        _parent, typeSwitchError = KnxNetIpDeviceManagementParse(io)
    case serviceType == 0x04:
        _parent, typeSwitchError = KnxNetIpTunnelingParse(io)
    case serviceType == 0x06:
        _parent, typeSwitchError = KnxNetRemoteLoggingParse(io)
    case serviceType == 0x07:
        _parent, typeSwitchError = KnxNetRemoteConfigurationAndDiagnosisParse(io)
    case serviceType == 0x08:
        _parent, typeSwitchError = KnxNetObjectServerParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *ServiceId) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *ServiceId) SerializeParent(io utils.WriteBuffer, child IServiceId, serializeChildFunction func() error) error {

    // Discriminator Field (serviceType) (Used as input to a switch field)
    serviceType := uint8(child.ServiceType())
    _serviceTypeErr := io.WriteUint8(8, (serviceType))
    if _serviceTypeErr != nil {
        return errors.New("Error serializing 'serviceType' field " + _serviceTypeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *ServiceId) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpCore":
                        var dt *KnxNetIpCore
                        if m.Child != nil {
                            dt = m.Child.(*KnxNetIpCore)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpDeviceManagement":
                        var dt *KnxNetIpDeviceManagement
                        if m.Child != nil {
                            dt = m.Child.(*KnxNetIpDeviceManagement)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpTunneling":
                        var dt *KnxNetIpTunneling
                        if m.Child != nil {
                            dt = m.Child.(*KnxNetIpTunneling)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetRemoteLogging":
                        var dt *KnxNetRemoteLogging
                        if m.Child != nil {
                            dt = m.Child.(*KnxNetRemoteLogging)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetRemoteConfigurationAndDiagnosis":
                        var dt *KnxNetRemoteConfigurationAndDiagnosis
                        if m.Child != nil {
                            dt = m.Child.(*KnxNetRemoteConfigurationAndDiagnosis)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.KnxNetObjectServer":
                        var dt *KnxNetObjectServer
                        if m.Child != nil {
                            dt = m.Child.(*KnxNetObjectServer)
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

func (m *ServiceId) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

