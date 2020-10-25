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
type DeviceConfigurationRequest struct {
    DeviceConfigurationRequestDataBlock IDeviceConfigurationRequestDataBlock
    Cemi ICEMI
    KNXNetIPMessage
}

// The corresponding interface
type IDeviceConfigurationRequest interface {
    IKNXNetIPMessage
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m DeviceConfigurationRequest) MsgType() uint16 {
    return 0x0310
}

func (m DeviceConfigurationRequest) initialize() spi.Message {
    return m
}

func NewDeviceConfigurationRequest(deviceConfigurationRequestDataBlock IDeviceConfigurationRequestDataBlock, cemi ICEMI) KNXNetIPMessageInitializer {
    return &DeviceConfigurationRequest{DeviceConfigurationRequestDataBlock: deviceConfigurationRequestDataBlock, Cemi: cemi}
}

func CastIDeviceConfigurationRequest(structType interface{}) IDeviceConfigurationRequest {
    castFunc := func(typ interface{}) IDeviceConfigurationRequest {
        if iDeviceConfigurationRequest, ok := typ.(IDeviceConfigurationRequest); ok {
            return iDeviceConfigurationRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastDeviceConfigurationRequest(structType interface{}) DeviceConfigurationRequest {
    castFunc := func(typ interface{}) DeviceConfigurationRequest {
        if sDeviceConfigurationRequest, ok := typ.(DeviceConfigurationRequest); ok {
            return sDeviceConfigurationRequest
        }
        if sDeviceConfigurationRequest, ok := typ.(*DeviceConfigurationRequest); ok {
            return *sDeviceConfigurationRequest
        }
        return DeviceConfigurationRequest{}
    }
    return castFunc(structType)
}

func (m DeviceConfigurationRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

    // Simple field (deviceConfigurationRequestDataBlock)
    lengthInBits += m.DeviceConfigurationRequestDataBlock.LengthInBits()

    // Simple field (cemi)
    lengthInBits += m.Cemi.LengthInBits()

    return lengthInBits
}

func (m DeviceConfigurationRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DeviceConfigurationRequestParse(io *utils.ReadBuffer, totalLength uint16) (KNXNetIPMessageInitializer, error) {

    // Simple Field (deviceConfigurationRequestDataBlock)
    _deviceConfigurationRequestDataBlockMessage, _err := DeviceConfigurationRequestDataBlockParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'deviceConfigurationRequestDataBlock'. " + _err.Error())
    }
    var deviceConfigurationRequestDataBlock IDeviceConfigurationRequestDataBlock
    deviceConfigurationRequestDataBlock, _deviceConfigurationRequestDataBlockOk := _deviceConfigurationRequestDataBlockMessage.(IDeviceConfigurationRequestDataBlock)
    if !_deviceConfigurationRequestDataBlockOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_deviceConfigurationRequestDataBlockMessage).Name() + " to IDeviceConfigurationRequestDataBlock")
    }

    // Simple Field (cemi)
    _cemiMessage, _err := CEMIParse(io, uint8(totalLength) - uint8(uint8(uint8(uint8(6)) + uint8(deviceConfigurationRequestDataBlock.LengthInBytes()))))
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'cemi'. " + _err.Error())
    }
    var cemi ICEMI
    cemi, _cemiOk := _cemiMessage.(ICEMI)
    if !_cemiOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_cemiMessage).Name() + " to ICEMI")
    }

    // Create the instance
    return NewDeviceConfigurationRequest(deviceConfigurationRequestDataBlock, cemi), nil
}

func (m DeviceConfigurationRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (deviceConfigurationRequestDataBlock)
    deviceConfigurationRequestDataBlock := CastIDeviceConfigurationRequestDataBlock(m.DeviceConfigurationRequestDataBlock)
    _deviceConfigurationRequestDataBlockErr := deviceConfigurationRequestDataBlock.Serialize(io)
    if _deviceConfigurationRequestDataBlockErr != nil {
        return errors.New("Error serializing 'deviceConfigurationRequestDataBlock' field " + _deviceConfigurationRequestDataBlockErr.Error())
    }

    // Simple Field (cemi)
    cemi := CastICEMI(m.Cemi)
    _cemiErr := cemi.Serialize(io)
    if _cemiErr != nil {
        return errors.New("Error serializing 'cemi' field " + _cemiErr.Error())
    }

        return nil
    }
    return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}

func (m *DeviceConfigurationRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "deviceConfigurationRequestDataBlock":
                var data *DeviceConfigurationRequestDataBlock
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DeviceConfigurationRequestDataBlock = CastIDeviceConfigurationRequestDataBlock(data)
            case "cemi":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataReq":
                        var dt *CEMIDataReq
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataCon":
                        var dt *CEMIDataCon
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIDataInd":
                        var dt *CEMIDataInd
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIRawReq":
                        var dt *CEMIRawReq
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIRawCon":
                        var dt *CEMIRawCon
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIRawInd":
                        var dt *CEMIRawInd
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIPollDataReq":
                        var dt *CEMIPollDataReq
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIPollDataCon":
                        var dt *CEMIPollDataCon
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIBusmonInd":
                        var dt *CEMIBusmonInd
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIMPropReadReq":
                        var dt *CEMIMPropReadReq
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.CEMIMPropReadCon":
                        var dt *CEMIMPropReadCon
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.Cemi = dt
                    }
            }
        }
    }
}

func (m DeviceConfigurationRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.DeviceConfigurationRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DeviceConfigurationRequestDataBlock, xml.StartElement{Name: xml.Name{Local: "deviceConfigurationRequestDataBlock"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Cemi, xml.StartElement{Name: xml.Name{Local: "cemi"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

