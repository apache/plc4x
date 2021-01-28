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
type ApduDataExtPropertyValueWrite struct {
    Parent *ApduDataExt
    IApduDataExtPropertyValueWrite
}

// The corresponding interface
type IApduDataExtPropertyValueWrite interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ApduDataExtPropertyValueWrite) ExtApciType() uint8 {
    return 0x17
}


func (m *ApduDataExtPropertyValueWrite) InitializeParent(parent *ApduDataExt) {
}

func NewApduDataExtPropertyValueWrite() *ApduDataExt {
    child := &ApduDataExtPropertyValueWrite{
        Parent: NewApduDataExt(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastApduDataExtPropertyValueWrite(structType interface{}) *ApduDataExtPropertyValueWrite {
    castFunc := func(typ interface{}) *ApduDataExtPropertyValueWrite {
        if casted, ok := typ.(ApduDataExtPropertyValueWrite); ok {
            return &casted
        }
        if casted, ok := typ.(*ApduDataExtPropertyValueWrite); ok {
            return casted
        }
        if casted, ok := typ.(ApduDataExt); ok {
            return CastApduDataExtPropertyValueWrite(casted.Child)
        }
        if casted, ok := typ.(*ApduDataExt); ok {
            return CastApduDataExtPropertyValueWrite(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ApduDataExtPropertyValueWrite) GetTypeName() string {
    return "ApduDataExtPropertyValueWrite"
}

func (m *ApduDataExtPropertyValueWrite) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *ApduDataExtPropertyValueWrite) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ApduDataExtPropertyValueWriteParse(io *utils.ReadBuffer) (*ApduDataExt, error) {

    // Create a partially initialized instance
    _child := &ApduDataExtPropertyValueWrite{
        Parent: &ApduDataExt{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ApduDataExtPropertyValueWrite) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ApduDataExtPropertyValueWrite) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *ApduDataExtPropertyValueWrite) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

