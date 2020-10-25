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
    "encoding/base64"
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type MACAddress struct {
    Addr []int8

}

// The corresponding interface
type IMACAddress interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewMACAddress(addr []int8) spi.Message {
    return &MACAddress{Addr: addr}
}

func CastIMACAddress(structType interface{}) IMACAddress {
    castFunc := func(typ interface{}) IMACAddress {
        if iMACAddress, ok := typ.(IMACAddress); ok {
            return iMACAddress
        }
        return nil
    }
    return castFunc(structType)
}

func CastMACAddress(structType interface{}) MACAddress {
    castFunc := func(typ interface{}) MACAddress {
        if sMACAddress, ok := typ.(MACAddress); ok {
            return sMACAddress
        }
        if sMACAddress, ok := typ.(*MACAddress); ok {
            return *sMACAddress
        }
        return MACAddress{}
    }
    return castFunc(structType)
}

func (m MACAddress) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Array field
    if len(m.Addr) > 0 {
        lengthInBits += 8 * uint16(len(m.Addr))
    }

    return lengthInBits
}

func (m MACAddress) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func MACAddressParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Array field (addr)
    // Count array
    addr := make([]int8, uint16(6))
    for curItem := uint16(0); curItem < uint16(uint16(6)); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'addr' field " + _err.Error())
        }
        addr[curItem] = _item
    }

    // Create the instance
    return NewMACAddress(addr), nil
}

func (m MACAddress) Serialize(io utils.WriteBuffer) error {

    // Array Field (addr)
    if m.Addr != nil {
        for _, _element := range m.Addr {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'addr' field " + _elementErr.Error())
            }
        }
    }

    return nil
}

func (m *MACAddress) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "addr":
                var data []int8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Addr = data
            }
        }
    }
}

func (m MACAddress) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.MACAddress"},
        }}); err != nil {
        return err
    }
    _encodedAddr := make([]byte, base64.StdEncoding.EncodedLen(len(m.Addr)))
    base64.StdEncoding.Encode(_encodedAddr, utils.Int8ToByte(m.Addr))
    if err := e.EncodeElement(_encodedAddr, xml.StartElement{Name: xml.Name{Local: "addr"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

