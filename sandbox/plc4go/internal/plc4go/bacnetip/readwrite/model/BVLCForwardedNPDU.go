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
type BVLCForwardedNPDU struct {
	ip   []uint8
	port uint16
	npdu INPDU
	BVLC
}

// The corresponding interface
type IBVLCForwardedNPDU interface {
	IBVLC
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BVLCForwardedNPDU) BvlcFunction() uint8 {
	return 0x04
}

func (m BVLCForwardedNPDU) initialize() spi.Message {
	return m
}

func NewBVLCForwardedNPDU(ip []uint8, port uint16, npdu INPDU) BVLCInitializer {
	return &BVLCForwardedNPDU{ip: ip, port: port, npdu: npdu}
}

func CastIBVLCForwardedNPDU(structType interface{}) IBVLCForwardedNPDU {
	castFunc := func(typ interface{}) IBVLCForwardedNPDU {
		if iBVLCForwardedNPDU, ok := typ.(IBVLCForwardedNPDU); ok {
			return iBVLCForwardedNPDU
		}
		return nil
	}
	return castFunc(structType)
}

func CastBVLCForwardedNPDU(structType interface{}) BVLCForwardedNPDU {
	castFunc := func(typ interface{}) BVLCForwardedNPDU {
		if sBVLCForwardedNPDU, ok := typ.(BVLCForwardedNPDU); ok {
			return sBVLCForwardedNPDU
		}
		return BVLCForwardedNPDU{}
	}
	return castFunc(structType)
}

func (m BVLCForwardedNPDU) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BVLC.LengthInBits()

	// Array field
	if len(m.ip) > 0 {
		lengthInBits += 8 * uint16(len(m.ip))
	}

	// Simple field (port)
	lengthInBits += 16

	// Simple field (npdu)
	lengthInBits += m.npdu.LengthInBits()

	return lengthInBits
}

func (m BVLCForwardedNPDU) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BVLCForwardedNPDUParse(io *spi.ReadBuffer, bvlcLength uint16) (BVLCInitializer, error) {

	// Array field (ip)
	// Count array
	ip := make([]uint8, uint16(4))
	for curItem := uint16(0); curItem < uint16(uint16(4)); curItem++ {

		_item, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'ip' field " + _err.Error())
		}
		ip[curItem] = _item
	}

	// Simple Field (port)
	port, _portErr := io.ReadUint16(16)
	if _portErr != nil {
		return nil, errors.New("Error parsing 'port' field " + _portErr.Error())
	}

	// Simple Field (npdu)
	_npduMessage, _err := NPDUParse(io, uint16(bvlcLength)-uint16(uint16(10)))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'npdu'. " + _err.Error())
	}
	var npdu INPDU
	npdu, _npduOk := _npduMessage.(INPDU)
	if !_npduOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_npduMessage).Name() + " to INPDU")
	}

	// Create the instance
	return NewBVLCForwardedNPDU(ip, port, npdu), nil
}

func (m BVLCForwardedNPDU) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Array Field (ip)
		if m.ip != nil {
			for _, _element := range m.ip {
				io.WriteUint8(8, _element)
			}
		}

		// Simple Field (port)
		port := uint16(m.port)
		io.WriteUint16(16, (port))

		// Simple Field (npdu)
		npdu := CastINPDU(m.npdu)
		npdu.Serialize(io)

	}
	BVLCSerialize(io, m.BVLC, CastIBVLC(m), ser)
}
