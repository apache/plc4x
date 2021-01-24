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
    "io"
)

// The data-structure of this message
type DeviceDescriptorType2 struct {
    ManufacturerId uint16
    DeviceType uint16
    Version uint8
    ReadSupported bool
    WriteSupported bool
    LogicalTagBase uint8
    ChannelInfo1 *ChannelInformation
    ChannelInfo2 *ChannelInformation
    ChannelInfo3 *ChannelInformation
    ChannelInfo4 *ChannelInformation
    IDeviceDescriptorType2
}

// The corresponding interface
type IDeviceDescriptorType2 interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewDeviceDescriptorType2(manufacturerId uint16, deviceType uint16, version uint8, readSupported bool, writeSupported bool, logicalTagBase uint8, channelInfo1 *ChannelInformation, channelInfo2 *ChannelInformation, channelInfo3 *ChannelInformation, channelInfo4 *ChannelInformation) *DeviceDescriptorType2 {
    return &DeviceDescriptorType2{ManufacturerId: manufacturerId, DeviceType: deviceType, Version: version, ReadSupported: readSupported, WriteSupported: writeSupported, LogicalTagBase: logicalTagBase, ChannelInfo1: channelInfo1, ChannelInfo2: channelInfo2, ChannelInfo3: channelInfo3, ChannelInfo4: channelInfo4}
}

func CastDeviceDescriptorType2(structType interface{}) *DeviceDescriptorType2 {
    castFunc := func(typ interface{}) *DeviceDescriptorType2 {
        if casted, ok := typ.(DeviceDescriptorType2); ok {
            return &casted
        }
        if casted, ok := typ.(*DeviceDescriptorType2); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *DeviceDescriptorType2) GetTypeName() string {
    return "DeviceDescriptorType2"
}

func (m *DeviceDescriptorType2) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (manufacturerId)
    lengthInBits += 16

    // Simple field (deviceType)
    lengthInBits += 16

    // Simple field (version)
    lengthInBits += 8

    // Simple field (readSupported)
    lengthInBits += 1

    // Simple field (writeSupported)
    lengthInBits += 1

    // Simple field (logicalTagBase)
    lengthInBits += 6

    // Simple field (channelInfo1)
    lengthInBits += m.ChannelInfo1.LengthInBits()

    // Simple field (channelInfo2)
    lengthInBits += m.ChannelInfo2.LengthInBits()

    // Simple field (channelInfo3)
    lengthInBits += m.ChannelInfo3.LengthInBits()

    // Simple field (channelInfo4)
    lengthInBits += m.ChannelInfo4.LengthInBits()

    return lengthInBits
}

func (m *DeviceDescriptorType2) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DeviceDescriptorType2Parse(io *utils.ReadBuffer) (*DeviceDescriptorType2, error) {

    // Simple Field (manufacturerId)
    manufacturerId, _manufacturerIdErr := io.ReadUint16(16)
    if _manufacturerIdErr != nil {
        return nil, errors.New("Error parsing 'manufacturerId' field " + _manufacturerIdErr.Error())
    }

    // Simple Field (deviceType)
    deviceType, _deviceTypeErr := io.ReadUint16(16)
    if _deviceTypeErr != nil {
        return nil, errors.New("Error parsing 'deviceType' field " + _deviceTypeErr.Error())
    }

    // Simple Field (version)
    version, _versionErr := io.ReadUint8(8)
    if _versionErr != nil {
        return nil, errors.New("Error parsing 'version' field " + _versionErr.Error())
    }

    // Simple Field (readSupported)
    readSupported, _readSupportedErr := io.ReadBit()
    if _readSupportedErr != nil {
        return nil, errors.New("Error parsing 'readSupported' field " + _readSupportedErr.Error())
    }

    // Simple Field (writeSupported)
    writeSupported, _writeSupportedErr := io.ReadBit()
    if _writeSupportedErr != nil {
        return nil, errors.New("Error parsing 'writeSupported' field " + _writeSupportedErr.Error())
    }

    // Simple Field (logicalTagBase)
    logicalTagBase, _logicalTagBaseErr := io.ReadUint8(6)
    if _logicalTagBaseErr != nil {
        return nil, errors.New("Error parsing 'logicalTagBase' field " + _logicalTagBaseErr.Error())
    }

    // Simple Field (channelInfo1)
    channelInfo1, _channelInfo1Err := ChannelInformationParse(io)
    if _channelInfo1Err != nil {
        return nil, errors.New("Error parsing 'channelInfo1' field " + _channelInfo1Err.Error())
    }

    // Simple Field (channelInfo2)
    channelInfo2, _channelInfo2Err := ChannelInformationParse(io)
    if _channelInfo2Err != nil {
        return nil, errors.New("Error parsing 'channelInfo2' field " + _channelInfo2Err.Error())
    }

    // Simple Field (channelInfo3)
    channelInfo3, _channelInfo3Err := ChannelInformationParse(io)
    if _channelInfo3Err != nil {
        return nil, errors.New("Error parsing 'channelInfo3' field " + _channelInfo3Err.Error())
    }

    // Simple Field (channelInfo4)
    channelInfo4, _channelInfo4Err := ChannelInformationParse(io)
    if _channelInfo4Err != nil {
        return nil, errors.New("Error parsing 'channelInfo4' field " + _channelInfo4Err.Error())
    }

    // Create the instance
    return NewDeviceDescriptorType2(manufacturerId, deviceType, version, readSupported, writeSupported, logicalTagBase, channelInfo1, channelInfo2, channelInfo3, channelInfo4), nil
}

func (m *DeviceDescriptorType2) Serialize(io utils.WriteBuffer) error {

    // Simple Field (manufacturerId)
    manufacturerId := uint16(m.ManufacturerId)
    _manufacturerIdErr := io.WriteUint16(16, (manufacturerId))
    if _manufacturerIdErr != nil {
        return errors.New("Error serializing 'manufacturerId' field " + _manufacturerIdErr.Error())
    }

    // Simple Field (deviceType)
    deviceType := uint16(m.DeviceType)
    _deviceTypeErr := io.WriteUint16(16, (deviceType))
    if _deviceTypeErr != nil {
        return errors.New("Error serializing 'deviceType' field " + _deviceTypeErr.Error())
    }

    // Simple Field (version)
    version := uint8(m.Version)
    _versionErr := io.WriteUint8(8, (version))
    if _versionErr != nil {
        return errors.New("Error serializing 'version' field " + _versionErr.Error())
    }

    // Simple Field (readSupported)
    readSupported := bool(m.ReadSupported)
    _readSupportedErr := io.WriteBit((readSupported))
    if _readSupportedErr != nil {
        return errors.New("Error serializing 'readSupported' field " + _readSupportedErr.Error())
    }

    // Simple Field (writeSupported)
    writeSupported := bool(m.WriteSupported)
    _writeSupportedErr := io.WriteBit((writeSupported))
    if _writeSupportedErr != nil {
        return errors.New("Error serializing 'writeSupported' field " + _writeSupportedErr.Error())
    }

    // Simple Field (logicalTagBase)
    logicalTagBase := uint8(m.LogicalTagBase)
    _logicalTagBaseErr := io.WriteUint8(6, (logicalTagBase))
    if _logicalTagBaseErr != nil {
        return errors.New("Error serializing 'logicalTagBase' field " + _logicalTagBaseErr.Error())
    }

    // Simple Field (channelInfo1)
    _channelInfo1Err := m.ChannelInfo1.Serialize(io)
    if _channelInfo1Err != nil {
        return errors.New("Error serializing 'channelInfo1' field " + _channelInfo1Err.Error())
    }

    // Simple Field (channelInfo2)
    _channelInfo2Err := m.ChannelInfo2.Serialize(io)
    if _channelInfo2Err != nil {
        return errors.New("Error serializing 'channelInfo2' field " + _channelInfo2Err.Error())
    }

    // Simple Field (channelInfo3)
    _channelInfo3Err := m.ChannelInfo3.Serialize(io)
    if _channelInfo3Err != nil {
        return errors.New("Error serializing 'channelInfo3' field " + _channelInfo3Err.Error())
    }

    // Simple Field (channelInfo4)
    _channelInfo4Err := m.ChannelInfo4.Serialize(io)
    if _channelInfo4Err != nil {
        return errors.New("Error serializing 'channelInfo4' field " + _channelInfo4Err.Error())
    }

    return nil
}

func (m *DeviceDescriptorType2) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "manufacturerId":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ManufacturerId = data
            case "deviceType":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DeviceType = data
            case "version":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Version = data
            case "readSupported":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ReadSupported = data
            case "writeSupported":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.WriteSupported = data
            case "logicalTagBase":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.LogicalTagBase = data
            case "channelInfo1":
                var data *ChannelInformation
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ChannelInfo1 = data
            case "channelInfo2":
                var data *ChannelInformation
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ChannelInfo2 = data
            case "channelInfo3":
                var data *ChannelInformation
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ChannelInfo3 = data
            case "channelInfo4":
                var data *ChannelInformation
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ChannelInfo4 = data
            }
        }
    }
}

func (m *DeviceDescriptorType2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.knxnetip.readwrite.DeviceDescriptorType2"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ManufacturerId, xml.StartElement{Name: xml.Name{Local: "manufacturerId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DeviceType, xml.StartElement{Name: xml.Name{Local: "deviceType"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Version, xml.StartElement{Name: xml.Name{Local: "version"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ReadSupported, xml.StartElement{Name: xml.Name{Local: "readSupported"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.WriteSupported, xml.StartElement{Name: xml.Name{Local: "writeSupported"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.LogicalTagBase, xml.StartElement{Name: xml.Name{Local: "logicalTagBase"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ChannelInfo1, xml.StartElement{Name: xml.Name{Local: "channelInfo1"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ChannelInfo2, xml.StartElement{Name: xml.Name{Local: "channelInfo2"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ChannelInfo3, xml.StartElement{Name: xml.Name{Local: "channelInfo3"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ChannelInfo4, xml.StartElement{Name: xml.Name{Local: "channelInfo4"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

