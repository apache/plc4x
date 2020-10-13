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
type BACnetAddress struct {
	address []uint8
	port    uint16
}

// The corresponding interface
type IBACnetAddress interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewBACnetAddress(address []uint8, port uint16) spi.Message {
	return &BACnetAddress{address: address, port: port}
}

func CastIBACnetAddress(structType interface{}) IBACnetAddress {
	castFunc := func(typ interface{}) IBACnetAddress {
		if iBACnetAddress, ok := typ.(IBACnetAddress); ok {
			return iBACnetAddress
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetAddress(structType interface{}) BACnetAddress {
	castFunc := func(typ interface{}) BACnetAddress {
		if sBACnetAddress, ok := typ.(BACnetAddress); ok {
			return sBACnetAddress
		}
		return BACnetAddress{}
	}
	return castFunc(structType)
}

func (m BACnetAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Array field
	if len(m.address) > 0 {
		lengthInBits += 8 * uint16(len(m.address))
	}

	// Simple field (port)
	lengthInBits += 16

	return lengthInBits
}

func (m BACnetAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetAddressParse(io spi.ReadBuffer) (spi.Message, error) {

	// Array field (address)
	var address []uint8
	// Count array
	{
		address := make([]uint8, uint16(4))
		for curItem := uint16(0); curItem < uint16(uint16(4)); curItem++ {

			_addressVal, _err := io.ReadUint8(8)
			if _err != nil {
				return nil, errors.New("Error parsing 'address' field " + _err.Error())
			}
			address = append(address, _addressVal)
		}
	}

	// Simple Field (port)
	port, _portErr := io.ReadUint16(16)
	if _portErr != nil {
		return nil, errors.New("Error parsing 'port' field " + _portErr.Error())
	}

	// Create the instance
	return NewBACnetAddress(address, port), nil
}

func (m BACnetAddress) Serialize(io spi.WriteBuffer) {

	// Array Field (address)
	if m.address != nil {
		for _, _element := range m.address {
			io.WriteUint8(8, _element)
		}
	}

	// Simple Field (port)
	port := uint16(m.port)
	io.WriteUint16(16, (port))
}
