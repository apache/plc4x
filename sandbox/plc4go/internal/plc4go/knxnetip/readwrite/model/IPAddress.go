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
type IPAddress struct {
	addr []int8
}

// The corresponding interface
type IIPAddress interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewIPAddress(addr []int8) spi.Message {
	return &IPAddress{addr: addr}
}

func CastIIPAddress(structType interface{}) IIPAddress {
	castFunc := func(typ interface{}) IIPAddress {
		if iIPAddress, ok := typ.(IIPAddress); ok {
			return iIPAddress
		}
		return nil
	}
	return castFunc(structType)
}

func CastIPAddress(structType interface{}) IPAddress {
	castFunc := func(typ interface{}) IPAddress {
		if sIPAddress, ok := typ.(IPAddress); ok {
			return sIPAddress
		}
		return IPAddress{}
	}
	return castFunc(structType)
}

func (m IPAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Array field
	if len(m.addr) > 0 {
		lengthInBits += 8 * uint16(len(m.addr))
	}

	return lengthInBits
}

func (m IPAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func IPAddressParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Array field (addr)
	// Count array
	addr := make([]int8, uint16(4))
	for curItem := uint16(0); curItem < uint16(uint16(4)); curItem++ {

		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'addr' field " + _err.Error())
		}
		addr[curItem] = _item
	}

	// Create the instance
	return NewIPAddress(addr), nil
}

func (m IPAddress) Serialize(io spi.WriteBuffer) {

	// Array Field (addr)
	if m.addr != nil {
		for _, _element := range m.addr {
			io.WriteInt8(8, _element)
		}
	}

}
