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
)

// The data-structure of this message
type KNXGroupAddress struct {
    Child IKNXGroupAddressChild
    IKNXGroupAddress
    IKNXGroupAddressParent
}

// The corresponding interface
type IKNXGroupAddress interface {
    NumLevels() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IKNXGroupAddressParent interface {
    SerializeParent(io utils.WriteBuffer, child IKNXGroupAddress, serializeChildFunction func() error) error
}

type IKNXGroupAddressChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *KNXGroupAddress)
    IKNXGroupAddress
}

func NewKNXGroupAddress() *KNXGroupAddress {
    return &KNXGroupAddress{}
}

func CastKNXGroupAddress(structType interface{}) KNXGroupAddress {
    castFunc := func(typ interface{}) KNXGroupAddress {
        if casted, ok := typ.(KNXGroupAddress); ok {
            return casted
        }
        if casted, ok := typ.(*KNXGroupAddress); ok {
            return *casted
        }
        return KNXGroupAddress{}
    }
    return castFunc(structType)
}

func (m *KNXGroupAddress) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *KNXGroupAddress) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KNXGroupAddressParse(io *utils.ReadBuffer, numLevels uint8) (*KNXGroupAddress, error) {

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *KNXGroupAddress
    var typeSwitchError error
    switch {
    case numLevels == 1:
        _parent, typeSwitchError = KNXGroupAddressFreeLevelParse(io)
    case numLevels == 2:
        _parent, typeSwitchError = KNXGroupAddress2LevelParse(io)
    case numLevels == 3:
        _parent, typeSwitchError = KNXGroupAddress3LevelParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *KNXGroupAddress) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *KNXGroupAddress) SerializeParent(io utils.WriteBuffer, child IKNXGroupAddress, serializeChildFunction func() error) error {

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *KNXGroupAddress) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *KNXGroupAddress) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddress"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

