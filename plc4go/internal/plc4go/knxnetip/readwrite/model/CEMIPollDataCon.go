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
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMIPollDataCon struct {
    CEMI
}

// The corresponding interface
type ICEMIPollDataCon interface {
    ICEMI
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIPollDataCon) MessageCode() uint8 {
    return 0x25
}

func (m CEMIPollDataCon) initialize() spi.Message {
    return m
}

func NewCEMIPollDataCon() CEMIInitializer {
    return &CEMIPollDataCon{}
}

func CastICEMIPollDataCon(structType interface{}) ICEMIPollDataCon {
    castFunc := func(typ interface{}) ICEMIPollDataCon {
        if iCEMIPollDataCon, ok := typ.(ICEMIPollDataCon); ok {
            return iCEMIPollDataCon
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIPollDataCon(structType interface{}) CEMIPollDataCon {
    castFunc := func(typ interface{}) CEMIPollDataCon {
        if sCEMIPollDataCon, ok := typ.(CEMIPollDataCon); ok {
            return sCEMIPollDataCon
        }
        if sCEMIPollDataCon, ok := typ.(*CEMIPollDataCon); ok {
            return *sCEMIPollDataCon
        }
        return CEMIPollDataCon{}
    }
    return castFunc(structType)
}

func (m CEMIPollDataCon) LengthInBits() uint16 {
    var lengthInBits uint16 = m.CEMI.LengthInBits()

    return lengthInBits
}

func (m CEMIPollDataCon) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIPollDataConParse(io *utils.ReadBuffer) (CEMIInitializer, error) {

    // Create the instance
    return NewCEMIPollDataCon(), nil
}

func (m CEMIPollDataCon) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return CEMISerialize(io, m.CEMI, CastICEMI(m), ser)
}

func (m *CEMIPollDataCon) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            }
        }
    }
}

func (m CEMIPollDataCon) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.CEMIPollDataCon"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

