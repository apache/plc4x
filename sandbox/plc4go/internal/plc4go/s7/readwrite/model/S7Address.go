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
)

// The data-structure of this message
type S7Address struct {
}

// The corresponding interface
type IS7Address interface {
	spi.Message
	AddressType() uint8
	Serialize(io spi.WriteBuffer)
}

type S7AddressInitializer interface {
	initialize() spi.Message
}

func S7AddressAddressType(m IS7Address) uint8 {
	return m.AddressType()
}

func CastIS7Address(structType interface{}) IS7Address {
	castFunc := func(typ interface{}) IS7Address {
		if iS7Address, ok := typ.(IS7Address); ok {
			return iS7Address
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7Address(structType interface{}) S7Address {
	castFunc := func(typ interface{}) S7Address {
		if sS7Address, ok := typ.(S7Address); ok {
			return sS7Address
		}
		return S7Address{}
	}
	return castFunc(structType)
}

func (m S7Address) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (addressType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m S7Address) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7AddressParse(io spi.ReadBuffer) (spi.Message, error) {

	// Discriminator Field (addressType) (Used as input to a switch field)
	addressType, _addressTypeErr := io.ReadUint8(8)
	if _addressTypeErr != nil {
		return nil, errors.New("Error parsing 'addressType' field " + _addressTypeErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer S7AddressInitializer
	var typeSwitchError error
	switch {
	case addressType == 0x10:
		initializer, typeSwitchError = S7AddressAnyParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func S7AddressSerialize(io spi.WriteBuffer, m S7Address, i IS7Address, childSerialize func()) {

	// Discriminator Field (addressType) (Used as input to a switch field)
	addressType := uint8(i.AddressType())
	io.WriteUint8(8, (addressType))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	childSerialize()

}
