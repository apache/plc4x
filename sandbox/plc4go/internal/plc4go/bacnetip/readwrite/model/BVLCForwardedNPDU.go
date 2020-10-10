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
	npdu NPDU
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

func NewBVLCForwardedNPDU(ip []uint8, port uint16, npdu NPDU) BVLCInitializer {
	return &BVLCForwardedNPDU{ip: ip, port: port, npdu: npdu}
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

func BVLCForwardedNPDUParse(io spi.ReadBuffer, bvlcLength uint16) (BVLCInitializer, error) {

	// Array field (ip)
	var ip []uint8
	// Count array
	{
		ip := make([]uint8, 4)
		for curItem := uint16(0); curItem < uint16(4); curItem++ {

			ip = append(ip, io.ReadUint8(8))
		}
	}

	// Simple Field (port)
	var port uint16 = io.ReadUint16(16)

	// Simple Field (npdu)
	_npduMessage, _err := NPDUParse(io, uint16((bvlcLength)-(10)))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'npdu'. " + _err.Error())
	}
	var npdu NPDU
	npdu, _npduOk := _npduMessage.(NPDU)
	if !_npduOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_npduMessage).Name() + " to NPDU")
	}

	// Create the instance
	return NewBVLCForwardedNPDU(ip, port, npdu), nil
}

func (m BVLCForwardedNPDU) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IBVLCForwardedNPDU); ok {

			// Array Field (ip)
			if m.ip != nil {
				for _, _element := range m.ip {
					io.WriteUint8(8, _element)
				}
			}

			// Simple Field (port)
			var port uint16 = m.port
			io.WriteUint16(16, (port))

			// Simple Field (npdu)
			var npdu NPDU = m.npdu
			npdu.Serialize(io)
		}
	}
	serializeFunc(m)
}
