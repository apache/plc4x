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
type AmsTCPPacket struct {
    Userdata *AmsPacket
    IAmsTCPPacket
}

// The corresponding interface
type IAmsTCPPacket interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewAmsTCPPacket(userdata *AmsPacket) *AmsTCPPacket {
    return &AmsTCPPacket{Userdata: userdata}
}

func CastAmsTCPPacket(structType interface{}) *AmsTCPPacket {
    castFunc := func(typ interface{}) *AmsTCPPacket {
        if casted, ok := typ.(AmsTCPPacket); ok {
            return &casted
        }
        if casted, ok := typ.(*AmsTCPPacket); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *AmsTCPPacket) GetTypeName() string {
    return "AmsTCPPacket"
}

func (m *AmsTCPPacket) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Reserved Field (reserved)
    lengthInBits += 16

    // Implicit Field (length)
    lengthInBits += 32

    // Simple field (userdata)
    lengthInBits += m.Userdata.LengthInBits()

    return lengthInBits
}

func (m *AmsTCPPacket) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func AmsTCPPacketParse(io *utils.ReadBuffer) (*AmsTCPPacket, error) {

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint16(16)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint16(0x0000) {
            log.WithFields(log.Fields{
                "expected value": uint16(0x0000),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _lengthErr := io.ReadUint32(32)
    if _lengthErr != nil {
        return nil, errors.New("Error parsing 'length' field " + _lengthErr.Error())
    }

    // Simple Field (userdata)
    userdata, _userdataErr := AmsPacketParse(io)
    if _userdataErr != nil {
        return nil, errors.New("Error parsing 'userdata' field " + _userdataErr.Error())
    }

    // Create the instance
    return NewAmsTCPPacket(userdata), nil
}

func (m *AmsTCPPacket) Serialize(io utils.WriteBuffer) error {

    // Reserved Field (reserved)
    {
        _err := io.WriteUint16(16, uint16(0x0000))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    length := uint32(m.Userdata.LengthInBytes())
    _lengthErr := io.WriteUint32(32, (length))
    if _lengthErr != nil {
        return errors.New("Error serializing 'length' field " + _lengthErr.Error())
    }

    // Simple Field (userdata)
    _userdataErr := m.Userdata.Serialize(io)
    if _userdataErr != nil {
        return errors.New("Error serializing 'userdata' field " + _userdataErr.Error())
    }

    return nil
}

func (m *AmsTCPPacket) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "userdata":
                var data *AmsPacket
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.Userdata = data
            }
        }
    }
}

func (m *AmsTCPPacket) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.ads.readwrite.AmsTCPPacket"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Userdata, xml.StartElement{Name: xml.Name{Local: "userdata"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

