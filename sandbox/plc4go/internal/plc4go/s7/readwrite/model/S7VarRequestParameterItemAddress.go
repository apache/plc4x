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
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
)

// The data-structure of this message
type S7VarRequestParameterItemAddress struct {
	address S7Address
	S7VarRequestParameterItem
}

// The corresponding interface
type IS7VarRequestParameterItemAddress interface {
	IS7VarRequestParameterItem
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7VarRequestParameterItemAddress) ItemType() uint8 {
	return 0x12
}

func (m S7VarRequestParameterItemAddress) initialize() spi.Message {
	return m
}

func NewS7VarRequestParameterItemAddress(address S7Address) S7VarRequestParameterItemInitializer {
	return &S7VarRequestParameterItemAddress{address: address}
}

func (m S7VarRequestParameterItemAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7VarRequestParameterItem.LengthInBits()

	// Implicit Field (itemLength)
	lengthInBits += 8

	// Simple field (address)
	lengthInBits += m.address.LengthInBits()

	return lengthInBits
}

func (m S7VarRequestParameterItemAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7VarRequestParameterItemAddressParse(io spi.ReadBuffer) (S7VarRequestParameterItemInitializer, error) {

	// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint8 = io.ReadUint8(8)

	// Simple Field (address)
	_addressMessage, _err := S7AddressParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'address'. " + _err.Error())
	}
	var address S7Address
	address, _addressOk := _addressMessage.(S7Address)
	if !_addressOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_addressMessage).Name() + " to S7Address")
	}

	// Create the instance
	return NewS7VarRequestParameterItemAddress(address), nil
}

func (m S7VarRequestParameterItemAddress) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IS7VarRequestParameterItemAddress); ok {

			// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			itemLength := uint8(m.address.LengthInBytes())
			io.WriteUint8(8, (itemLength))

			// Simple Field (address)
			var address S7Address = m.address
			address.Serialize(io)
		}
	}
	serializeFunc(m)
}
