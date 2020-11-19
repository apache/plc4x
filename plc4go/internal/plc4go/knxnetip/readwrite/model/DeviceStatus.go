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
    log "github.com/sirupsen/logrus"
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

// The data-structure of this message
type DeviceStatus struct {
    ProgramMode bool
    IDeviceStatus
}

// The corresponding interface
type IDeviceStatus interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewDeviceStatus(programMode bool) *DeviceStatus {
    return &DeviceStatus{ProgramMode: programMode}
}

func CastDeviceStatus(structType interface{}) *DeviceStatus {
    castFunc := func(typ interface{}) *DeviceStatus {
        if casted, ok := typ.(DeviceStatus); ok {
            return &casted
        }
        if casted, ok := typ.(*DeviceStatus); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *DeviceStatus) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Reserved Field (reserved)
    lengthInBits += 7

    // Simple field (programMode)
    lengthInBits += 1

    return lengthInBits
}

func (m *DeviceStatus) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DeviceStatusParse(io *utils.ReadBuffer) (*DeviceStatus, error) {

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(7)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0x00) {
            log.WithFields(log.Fields{
                "expected value": uint8(0x00),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (programMode)
    programMode, _programModeErr := io.ReadBit()
    if _programModeErr != nil {
        return nil, errors.New("Error parsing 'programMode' field " + _programModeErr.Error())
    }

    // Create the instance
    return NewDeviceStatus(programMode), nil
}

func (m *DeviceStatus) Serialize(io utils.WriteBuffer) error {

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(7, uint8(0x00))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (programMode)
    programMode := bool(m.ProgramMode)
    _programModeErr := io.WriteBit((programMode))
    if _programModeErr != nil {
        return errors.New("Error serializing 'programMode' field " + _programModeErr.Error())
    }

    return nil
}

func (m *DeviceStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "programMode":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ProgramMode = data
            }
        }
    }
}

func (m *DeviceStatus) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.knxnetip.readwrite.DeviceStatus"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ProgramMode, xml.StartElement{Name: xml.Name{Local: "programMode"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

