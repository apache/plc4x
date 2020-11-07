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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "strconv"
)

// Constant values.
const ModbusConstants_MODBUSTCPDEFAULTPORT uint16 = 502

// The data-structure of this message
type ModbusConstants struct {
    IModbusConstants
}

// The corresponding interface
type IModbusConstants interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewModbusConstants() *ModbusConstants {
    return &ModbusConstants{}
}

func CastModbusConstants(structType interface{}) ModbusConstants {
    castFunc := func(typ interface{}) ModbusConstants {
        if casted, ok := typ.(ModbusConstants); ok {
            return casted
        }
        if casted, ok := typ.(*ModbusConstants); ok {
            return *casted
        }
        return ModbusConstants{}
    }
    return castFunc(structType)
}

func (m *ModbusConstants) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Const Field (modbusTcpDefaultPort)
    lengthInBits += 16

    return lengthInBits
}

func (m *ModbusConstants) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusConstantsParse(io *utils.ReadBuffer) (*ModbusConstants, error) {

    // Const Field (modbusTcpDefaultPort)
    modbusTcpDefaultPort, _modbusTcpDefaultPortErr := io.ReadUint16(16)
    if _modbusTcpDefaultPortErr != nil {
        return nil, errors.New("Error parsing 'modbusTcpDefaultPort' field " + _modbusTcpDefaultPortErr.Error())
    }
    if modbusTcpDefaultPort != ModbusConstants_MODBUSTCPDEFAULTPORT {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(ModbusConstants_MODBUSTCPDEFAULTPORT)) + " but got " + strconv.Itoa(int(modbusTcpDefaultPort)))
    }

    // Create the instance
    return NewModbusConstants(), nil
}

func (m *ModbusConstants) Serialize(io utils.WriteBuffer) error {

    // Const Field (modbusTcpDefaultPort)
    _modbusTcpDefaultPortErr := io.WriteUint16(16, 502)
    if _modbusTcpDefaultPortErr != nil {
        return errors.New("Error serializing 'modbusTcpDefaultPort' field " + _modbusTcpDefaultPortErr.Error())
    }

    return nil
}

func (m *ModbusConstants) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            }
        }
    }
}

func (m *ModbusConstants) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.modbus.readwrite.ModbusConstants"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

