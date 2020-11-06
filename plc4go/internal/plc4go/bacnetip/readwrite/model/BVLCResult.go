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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BVLCResult struct {
    Parent *BVLC
    IBVLCResult
}

// The corresponding interface
type IBVLCResult interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *BVLCResult) BvlcFunction() uint8 {
    return 0x00
}


func (m *BVLCResult) InitializeParent(parent *BVLC) {
}

func NewBVLCResult() *BVLC {
    child := &BVLCResult{
        Parent: NewBVLC(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastBVLCResult(structType interface{}) BVLCResult {
    castFunc := func(typ interface{}) BVLCResult {
        if casted, ok := typ.(BVLCResult); ok {
            return casted
        }
        if casted, ok := typ.(*BVLCResult); ok {
            return *casted
        }
        if casted, ok := typ.(BVLC); ok {
            return CastBVLCResult(casted.Child)
        }
        if casted, ok := typ.(*BVLC); ok {
            return CastBVLCResult(casted.Child)
        }
        return BVLCResult{}
    }
    return castFunc(structType)
}

func (m *BVLCResult) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *BVLCResult) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BVLCResultParse(io *utils.ReadBuffer) (*BVLC, error) {

    // Create a partially initialized instance
    _child := &BVLCResult{
        Parent: &BVLC{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *BVLCResult) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *BVLCResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *BVLCResult) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BVLCResult"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

