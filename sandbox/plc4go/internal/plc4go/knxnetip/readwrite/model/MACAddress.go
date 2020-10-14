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
type MACAddress struct {
	Addr []int8
}

// The corresponding interface
type IMACAddress interface {
	spi.Message
	Serialize(io spi.WriteBuffer) error
}

func NewMACAddress(addr []int8) spi.Message {
	return &MACAddress{Addr: addr}
}

func CastIMACAddress(structType interface{}) IMACAddress {
	castFunc := func(typ interface{}) IMACAddress {
		if iMACAddress, ok := typ.(IMACAddress); ok {
			return iMACAddress
		}
		return nil
	}
	return castFunc(structType)
}

func CastMACAddress(structType interface{}) MACAddress {
	castFunc := func(typ interface{}) MACAddress {
		if sMACAddress, ok := typ.(MACAddress); ok {
			return sMACAddress
		}
		return MACAddress{}
	}
	return castFunc(structType)
}

func (m MACAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Array field
	if len(m.Addr) > 0 {
		lengthInBits += 8 * uint16(len(m.Addr))
	}

	return lengthInBits
}

func (m MACAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func MACAddressParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Array field (addr)
	// Count array
	addr := make([]int8, uint16(6))
	for curItem := uint16(0); curItem < uint16(uint16(6)); curItem++ {

		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'addr' field " + _err.Error())
		}
		addr[curItem] = _item
	}

	// Create the instance
	return NewMACAddress(addr), nil
}

func (m MACAddress) Serialize(io spi.WriteBuffer) error {

	// Array Field (addr)
	if m.Addr != nil {
		for _, _element := range m.Addr {
			_elementErr := io.WriteInt8(8, _element)
			if _elementErr != nil {
				return errors.New("Error serializing 'addr' field " + _elementErr.Error())
			}
		}
	}

	return nil
}
