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
    "reflect"
    "strings"
)

// The data-structure of this message
type S7ParameterUserDataItem struct {
    Child IS7ParameterUserDataItemChild
    IS7ParameterUserDataItem
    IS7ParameterUserDataItemParent
}

// The corresponding interface
type IS7ParameterUserDataItem interface {
    ItemType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type IS7ParameterUserDataItemParent interface {
    SerializeParent(io utils.WriteBuffer, child IS7ParameterUserDataItem, serializeChildFunction func() error) error
}

type IS7ParameterUserDataItemChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *S7ParameterUserDataItem)
    IS7ParameterUserDataItem
}

func NewS7ParameterUserDataItem() *S7ParameterUserDataItem {
    return &S7ParameterUserDataItem{}
}

func CastS7ParameterUserDataItem(structType interface{}) S7ParameterUserDataItem {
    castFunc := func(typ interface{}) S7ParameterUserDataItem {
        if casted, ok := typ.(S7ParameterUserDataItem); ok {
            return casted
        }
        if casted, ok := typ.(*S7ParameterUserDataItem); ok {
            return *casted
        }
        return S7ParameterUserDataItem{}
    }
    return castFunc(structType)
}

func (m *S7ParameterUserDataItem) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (itemType)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *S7ParameterUserDataItem) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7ParameterUserDataItemParse(io *utils.ReadBuffer) (*S7ParameterUserDataItem, error) {

    // Discriminator Field (itemType) (Used as input to a switch field)
    itemType, _itemTypeErr := io.ReadUint8(8)
    if _itemTypeErr != nil {
        return nil, errors.New("Error parsing 'itemType' field " + _itemTypeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *S7ParameterUserDataItem
    var typeSwitchError error
    switch {
    case itemType == 0x12:
        _parent, typeSwitchError = S7ParameterUserDataItemCPUFunctionsParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *S7ParameterUserDataItem) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *S7ParameterUserDataItem) SerializeParent(io utils.WriteBuffer, child IS7ParameterUserDataItem, serializeChildFunction func() error) error {

    // Discriminator Field (itemType) (Used as input to a switch field)
    itemType := uint8(child.ItemType())
    _itemTypeErr := io.WriteUint8(8, (itemType))
    if _itemTypeErr != nil {
        return errors.New("Error serializing 'itemType' field " + _itemTypeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *S7ParameterUserDataItem) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                default:
                    switch start.Attr[0].Value {
                        case "org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions":
                            var dt *S7ParameterUserDataItemCPUFunctions
                            if m.Child != nil {
                                dt = m.Child.(*S7ParameterUserDataItemCPUFunctions)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                    }
            }
        }
    }
}

func (m *S7ParameterUserDataItem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.s7.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

