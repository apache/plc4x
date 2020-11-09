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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

// The data-structure of this message
type DeviceConfigurationRequest struct {
    DeviceConfigurationRequestDataBlock *DeviceConfigurationRequestDataBlock
    Cemi *CEMI
    Parent *KnxNetIpMessage
    IDeviceConfigurationRequest
}

// The corresponding interface
type IDeviceConfigurationRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *DeviceConfigurationRequest) MsgType() uint16 {
    return 0x0310
}


func (m *DeviceConfigurationRequest) InitializeParent(parent *KnxNetIpMessage) {
}

func NewDeviceConfigurationRequest(deviceConfigurationRequestDataBlock *DeviceConfigurationRequestDataBlock, cemi *CEMI, ) *KnxNetIpMessage {
    child := &DeviceConfigurationRequest{
        DeviceConfigurationRequestDataBlock: deviceConfigurationRequestDataBlock,
        Cemi: cemi,
        Parent: NewKnxNetIpMessage(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastDeviceConfigurationRequest(structType interface{}) *DeviceConfigurationRequest {
    castFunc := func(typ interface{}) *DeviceConfigurationRequest {
        if casted, ok := typ.(DeviceConfigurationRequest); ok {
            return &casted
        }
        if casted, ok := typ.(*DeviceConfigurationRequest); ok {
            return casted
        }
        if casted, ok := typ.(KnxNetIpMessage); ok {
            return CastDeviceConfigurationRequest(casted.Child)
        }
        if casted, ok := typ.(*KnxNetIpMessage); ok {
            return CastDeviceConfigurationRequest(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *DeviceConfigurationRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (deviceConfigurationRequestDataBlock)
    lengthInBits += m.DeviceConfigurationRequestDataBlock.LengthInBits()

    // Simple field (cemi)
    lengthInBits += m.Cemi.LengthInBits()

    return lengthInBits
}

func (m *DeviceConfigurationRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DeviceConfigurationRequestParse(io *utils.ReadBuffer, totalLength uint16) (*KnxNetIpMessage, error) {

    // Simple Field (deviceConfigurationRequestDataBlock)
    deviceConfigurationRequestDataBlock, _deviceConfigurationRequestDataBlockErr := DeviceConfigurationRequestDataBlockParse(io)
    if _deviceConfigurationRequestDataBlockErr != nil {
        return nil, errors.New("Error parsing 'deviceConfigurationRequestDataBlock' field " + _deviceConfigurationRequestDataBlockErr.Error())
    }

    // Simple Field (cemi)
    cemi, _cemiErr := CEMIParse(io, uint8(totalLength) - uint8(uint8(uint8(uint8(6)) + uint8(deviceConfigurationRequestDataBlock.LengthInBytes()))))
    if _cemiErr != nil {
        return nil, errors.New("Error parsing 'cemi' field " + _cemiErr.Error())
    }

    // Create a partially initialized instance
    _child := &DeviceConfigurationRequest{
        DeviceConfigurationRequestDataBlock: deviceConfigurationRequestDataBlock,
        Cemi: cemi,
        Parent: &KnxNetIpMessage{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *DeviceConfigurationRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (deviceConfigurationRequestDataBlock)
    _deviceConfigurationRequestDataBlockErr := m.DeviceConfigurationRequestDataBlock.Serialize(io)
    if _deviceConfigurationRequestDataBlockErr != nil {
        return errors.New("Error serializing 'deviceConfigurationRequestDataBlock' field " + _deviceConfigurationRequestDataBlockErr.Error())
    }

    // Simple Field (cemi)
    _cemiErr := m.Cemi.Serialize(io)
    if _cemiErr != nil {
        return errors.New("Error serializing 'cemi' field " + _cemiErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *DeviceConfigurationRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "deviceConfigurationRequestDataBlock":
                var data *DeviceConfigurationRequestDataBlock
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.DeviceConfigurationRequestDataBlock = data
            case "cemi":
                var dt *CEMI
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.Cemi = dt
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

func (m *DeviceConfigurationRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.DeviceConfigurationRequestDataBlock, xml.StartElement{Name: xml.Name{Local: "deviceConfigurationRequestDataBlock"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Cemi, xml.StartElement{Name: xml.Name{Local: "cemi"}}); err != nil {
        return err
    }
    return nil
}

