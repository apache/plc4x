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
type BACnetTagApplicationEnumerated struct {
	data []int8
	BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationEnumerated interface {
	IBACnetTag
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetTagApplicationEnumerated) ContextSpecificTag() uint8 {
	return 0
}

func (m BACnetTagApplicationEnumerated) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
	m.typeOrTagNumber = typeOrTagNumber
	m.lengthValueType = lengthValueType
	m.extTagNumber = extTagNumber
	m.extLength = extLength
	return m
}

func NewBACnetTagApplicationEnumerated(data []int8) BACnetTagInitializer {
	return &BACnetTagApplicationEnumerated{data: data}
}

func CastIBACnetTagApplicationEnumerated(structType interface{}) IBACnetTagApplicationEnumerated {
	castFunc := func(typ interface{}) IBACnetTagApplicationEnumerated {
		if iBACnetTagApplicationEnumerated, ok := typ.(IBACnetTagApplicationEnumerated); ok {
			return iBACnetTagApplicationEnumerated
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetTagApplicationEnumerated(structType interface{}) BACnetTagApplicationEnumerated {
	castFunc := func(typ interface{}) BACnetTagApplicationEnumerated {
		if sBACnetTagApplicationEnumerated, ok := typ.(BACnetTagApplicationEnumerated); ok {
			return sBACnetTagApplicationEnumerated
		}
		return BACnetTagApplicationEnumerated{}
	}
	return castFunc(structType)
}

func (m BACnetTagApplicationEnumerated) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetTag.LengthInBits()

	// Array field
	if len(m.data) > 0 {
		lengthInBits += 8 * uint16(len(m.data))
	}

	return lengthInBits
}

func (m BACnetTagApplicationEnumerated) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetTagApplicationEnumeratedParse(io *spi.ReadBuffer, lengthValueType uint8, extLength uint8) (BACnetTagInitializer, error) {

	// Array field (data)
	// Length array
	data := make([]int8, 0)
	_dataLength := spi.InlineIf(bool(bool((lengthValueType) == (5))), uint16(extLength), uint16(lengthValueType))
	_dataEndPos := io.GetPos() + uint16(_dataLength)
	for io.GetPos() < _dataEndPos {
		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'data' field " + _err.Error())
		}
		data = append(data, _item)
	}

	// Create the instance
	return NewBACnetTagApplicationEnumerated(data), nil
}

func (m BACnetTagApplicationEnumerated) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Array Field (data)
		if m.data != nil {
			for _, _element := range m.data {
				io.WriteInt8(8, _element)
			}
		}

	}
	BACnetTagSerialize(io, m.BACnetTag, CastIBACnetTag(m), ser)
}
