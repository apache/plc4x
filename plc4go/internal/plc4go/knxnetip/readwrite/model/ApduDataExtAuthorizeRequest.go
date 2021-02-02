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
type ApduDataExtAuthorizeRequest struct {
    Parent *ApduDataExt
    IApduDataExtAuthorizeRequest
}

// The corresponding interface
type IApduDataExtAuthorizeRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ApduDataExtAuthorizeRequest) ExtApciType() uint8 {
    return 0x11
}


func (m *ApduDataExtAuthorizeRequest) InitializeParent(parent *ApduDataExt) {
}

func NewApduDataExtAuthorizeRequest() *ApduDataExt {
    child := &ApduDataExtAuthorizeRequest{
        Parent: NewApduDataExt(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastApduDataExtAuthorizeRequest(structType interface{}) *ApduDataExtAuthorizeRequest {
    castFunc := func(typ interface{}) *ApduDataExtAuthorizeRequest {
        if casted, ok := typ.(ApduDataExtAuthorizeRequest); ok {
            return &casted
        }
        if casted, ok := typ.(*ApduDataExtAuthorizeRequest); ok {
            return casted
        }
        if casted, ok := typ.(ApduDataExt); ok {
            return CastApduDataExtAuthorizeRequest(casted.Child)
        }
        if casted, ok := typ.(*ApduDataExt); ok {
            return CastApduDataExtAuthorizeRequest(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ApduDataExtAuthorizeRequest) GetTypeName() string {
    return "ApduDataExtAuthorizeRequest"
}

func (m *ApduDataExtAuthorizeRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *ApduDataExtAuthorizeRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ApduDataExtAuthorizeRequestParse(io *utils.ReadBuffer) (*ApduDataExt, error) {

    // Create a partially initialized instance
    _child := &ApduDataExtAuthorizeRequest{
        Parent: &ApduDataExt{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ApduDataExtAuthorizeRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ApduDataExtAuthorizeRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *ApduDataExtAuthorizeRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

