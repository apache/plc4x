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
	Address IS7Address
	S7VarRequestParameterItem
}

// The corresponding interface
type IS7VarRequestParameterItemAddress interface {
	IS7VarRequestParameterItem
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7VarRequestParameterItemAddress) ItemType() uint8 {
	return 0x12
}

func (m S7VarRequestParameterItemAddress) initialize() spi.Message {
	return m
}

func NewS7VarRequestParameterItemAddress(address IS7Address) S7VarRequestParameterItemInitializer {
	return &S7VarRequestParameterItemAddress{Address: address}
}

func CastIS7VarRequestParameterItemAddress(structType interface{}) IS7VarRequestParameterItemAddress {
	castFunc := func(typ interface{}) IS7VarRequestParameterItemAddress {
		if iS7VarRequestParameterItemAddress, ok := typ.(IS7VarRequestParameterItemAddress); ok {
			return iS7VarRequestParameterItemAddress
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7VarRequestParameterItemAddress(structType interface{}) S7VarRequestParameterItemAddress {
	castFunc := func(typ interface{}) S7VarRequestParameterItemAddress {
		if sS7VarRequestParameterItemAddress, ok := typ.(S7VarRequestParameterItemAddress); ok {
			return sS7VarRequestParameterItemAddress
		}
		return S7VarRequestParameterItemAddress{}
	}
	return castFunc(structType)
}

func (m S7VarRequestParameterItemAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7VarRequestParameterItem.LengthInBits()

	// Implicit Field (itemLength)
	lengthInBits += 8

	// Simple field (address)
	lengthInBits += m.Address.LengthInBits()

	return lengthInBits
}

func (m S7VarRequestParameterItemAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7VarRequestParameterItemAddressParse(io *spi.ReadBuffer) (S7VarRequestParameterItemInitializer, error) {

	// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _itemLengthErr := io.ReadUint8(8)
	if _itemLengthErr != nil {
		return nil, errors.New("Error parsing 'itemLength' field " + _itemLengthErr.Error())
	}

	// Simple Field (address)
	_addressMessage, _err := S7AddressParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'address'. " + _err.Error())
	}
	var address IS7Address
	address, _addressOk := _addressMessage.(IS7Address)
	if !_addressOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_addressMessage).Name() + " to IS7Address")
	}

	// Create the instance
	return NewS7VarRequestParameterItemAddress(address), nil
}

func (m S7VarRequestParameterItemAddress) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		itemLength := uint8(m.Address.LengthInBytes())
		_itemLengthErr := io.WriteUint8(8, (itemLength))
		if _itemLengthErr != nil {
			return errors.New("Error serializing 'itemLength' field " + _itemLengthErr.Error())
		}

		// Simple Field (address)
		address := CastIS7Address(m.Address)
		_addressErr := address.Serialize(io)
		if _addressErr != nil {
			return errors.New("Error serializing 'address' field " + _addressErr.Error())
		}

		return nil
	}
	return S7VarRequestParameterItemSerialize(io, m.S7VarRequestParameterItem, CastIS7VarRequestParameterItem(m), ser)
}
