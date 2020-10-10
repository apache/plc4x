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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type MACAddress struct {
	addr []int8
}

// The corresponding interface
type IMACAddress interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewMACAddress(addr []int8) spi.Message {
	return &MACAddress{addr: addr}
}

func (m MACAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Array field
	if len(m.addr) > 0 {
		lengthInBits += 8 * uint16(len(m.addr))
	}

	return lengthInBits
}

func (m MACAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func MACAddressParse(io spi.ReadBuffer) (spi.Message, error) {

	// Array field (addr)
	var addr []int8
	// Count array
	{
		addr := make([]int8, 6)
		for curItem := uint16(0); curItem < uint16(6); curItem++ {

			addr = append(addr, io.ReadInt8(8))
		}
	}

	// Create the instance
	return NewMACAddress(addr), nil
}

func (m MACAddress) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IMACAddress); ok {

			// Array Field (addr)
			if m.addr != nil {
				for _, _element := range m.addr {
					io.WriteInt8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
