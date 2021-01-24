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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "io"
)

// The data-structure of this message
type ApduDataMemoryRead struct {
    Parent *ApduData
    IApduDataMemoryRead
}

// The corresponding interface
type IApduDataMemoryRead interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ApduDataMemoryRead) ApciType() uint8 {
    return 0x8
}


func (m *ApduDataMemoryRead) InitializeParent(parent *ApduData) {
}

func NewApduDataMemoryRead() *ApduData {
    child := &ApduDataMemoryRead{
        Parent: NewApduData(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastApduDataMemoryRead(structType interface{}) *ApduDataMemoryRead {
    castFunc := func(typ interface{}) *ApduDataMemoryRead {
        if casted, ok := typ.(ApduDataMemoryRead); ok {
            return &casted
        }
        if casted, ok := typ.(*ApduDataMemoryRead); ok {
            return casted
        }
        if casted, ok := typ.(ApduData); ok {
            return CastApduDataMemoryRead(casted.Child)
        }
        if casted, ok := typ.(*ApduData); ok {
            return CastApduDataMemoryRead(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ApduDataMemoryRead) GetTypeName() string {
    return "ApduDataMemoryRead"
}

func (m *ApduDataMemoryRead) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *ApduDataMemoryRead) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ApduDataMemoryReadParse(io *utils.ReadBuffer) (*ApduData, error) {

    // Create a partially initialized instance
    _child := &ApduDataMemoryRead{
        Parent: &ApduData{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ApduDataMemoryRead) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ApduDataMemoryRead) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
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

func (m *ApduDataMemoryRead) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

