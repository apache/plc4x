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
	"io"
	"reflect"
	"strings"
)

// The data-structure of this message
type S7Address struct {
	Child IS7AddressChild
	IS7Address
	IS7AddressParent
}

// The corresponding interface
type IS7Address interface {
	AddressType() uint8
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

type IS7AddressParent interface {
	SerializeParent(io utils.WriteBuffer, child IS7Address, serializeChildFunction func() error) error
	GetTypeName() string
}

type IS7AddressChild interface {
	Serialize(io utils.WriteBuffer) error
	InitializeParent(parent *S7Address)
	GetTypeName() string
	IS7Address
}

func NewS7Address() *S7Address {
	return &S7Address{}
}

func CastS7Address(structType interface{}) *S7Address {
	castFunc := func(typ interface{}) *S7Address {
		if casted, ok := typ.(S7Address); ok {
			return &casted
		}
		if casted, ok := typ.(*S7Address); ok {
			return casted
		}
		return nil
	}
	return castFunc(structType)
}

func (m *S7Address) GetTypeName() string {
	return "S7Address"
}

func (m *S7Address) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Discriminator Field (addressType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...
	lengthInBits += m.Child.LengthInBits()

	return lengthInBits
}

func (m *S7Address) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7AddressParse(io *utils.ReadBuffer) (*S7Address, error) {

	// Discriminator Field (addressType) (Used as input to a switch field)
	addressType, _addressTypeErr := io.ReadUint8(8)
	if _addressTypeErr != nil {
		return nil, errors.New("Error parsing 'addressType' field " + _addressTypeErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _parent *S7Address
	var typeSwitchError error
	switch {
	case addressType == 0x10:
		_parent, typeSwitchError = S7AddressAnyParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Finish initializing
	_parent.Child.InitializeParent(_parent)
	return _parent, nil
}

func (m *S7Address) Serialize(io utils.WriteBuffer) error {
	return m.Child.Serialize(io)
}

func (m *S7Address) SerializeParent(io utils.WriteBuffer, child IS7Address, serializeChildFunction func() error) error {

	// Discriminator Field (addressType) (Used as input to a switch field)
	addressType := uint8(child.AddressType())
	_addressTypeErr := io.WriteUint8(8, (addressType))
	if _addressTypeErr != nil {
		return errors.New("Error serializing 'addressType' field " + _addressTypeErr.Error())
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	_typeSwitchErr := serializeChildFunction()
	if _typeSwitchErr != nil {
		return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
	}

	return nil
}

func (m *S7Address) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				case "org.apache.plc4x.java.s7.readwrite.S7AddressAny":
					var dt *S7AddressAny
					if m.Child != nil {
						dt = m.Child.(*S7AddressAny)
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

func (m *S7Address) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	className := reflect.TypeOf(m.Child).String()
	className = "org.apache.plc4x.java.s7.readwrite." + className[strings.LastIndex(className, ".")+1:]
	if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
		{Name: xml.Name{Local: "className"}, Value: className},
	}}); err != nil {
		return err
	}
	marshaller, ok := m.Child.(xml.Marshaler)
	if !ok {
		return errors.New("child is not castable to Marshaler")
	}
	if err := marshaller.MarshalXML(e, start); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}
