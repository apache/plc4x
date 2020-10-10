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
type CEMIMPropReadReq struct {
	interfaceObjectType uint16
	objectInstance      uint8
	propertyId          uint8
	numberOfElements    uint8
	startIndex          uint16
	CEMI
}

// The corresponding interface
type ICEMIMPropReadReq interface {
	ICEMI
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m CEMIMPropReadReq) MessageCode() uint8 {
	return 0xFC
}

func (m CEMIMPropReadReq) initialize() spi.Message {
	return m
}

func NewCEMIMPropReadReq(interfaceObjectType uint16, objectInstance uint8, propertyId uint8, numberOfElements uint8, startIndex uint16) CEMIInitializer {
	return &CEMIMPropReadReq{interfaceObjectType: interfaceObjectType, objectInstance: objectInstance, propertyId: propertyId, numberOfElements: numberOfElements, startIndex: startIndex}
}

func (m CEMIMPropReadReq) LengthInBits() uint16 {
	var lengthInBits uint16 = m.CEMI.LengthInBits()

	// Simple field (interfaceObjectType)
	lengthInBits += 16

	// Simple field (objectInstance)
	lengthInBits += 8

	// Simple field (propertyId)
	lengthInBits += 8

	// Simple field (numberOfElements)
	lengthInBits += 4

	// Simple field (startIndex)
	lengthInBits += 12

	return lengthInBits
}

func (m CEMIMPropReadReq) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIMPropReadReqParse(io spi.ReadBuffer) (CEMIInitializer, error) {

	// Simple Field (interfaceObjectType)
	var interfaceObjectType uint16 = io.ReadUint16(16)

	// Simple Field (objectInstance)
	var objectInstance uint8 = io.ReadUint8(8)

	// Simple Field (propertyId)
	var propertyId uint8 = io.ReadUint8(8)

	// Simple Field (numberOfElements)
	var numberOfElements uint8 = io.ReadUint8(4)

	// Simple Field (startIndex)
	var startIndex uint16 = io.ReadUint16(12)

	// Create the instance
	return NewCEMIMPropReadReq(interfaceObjectType, objectInstance, propertyId, numberOfElements, startIndex), nil
}

func (m CEMIMPropReadReq) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ICEMIMPropReadReq); ok {

			// Simple Field (interfaceObjectType)
			var interfaceObjectType uint16 = m.interfaceObjectType
			io.WriteUint16(16, (interfaceObjectType))

			// Simple Field (objectInstance)
			var objectInstance uint8 = m.objectInstance
			io.WriteUint8(8, (objectInstance))

			// Simple Field (propertyId)
			var propertyId uint8 = m.propertyId
			io.WriteUint8(8, (propertyId))

			// Simple Field (numberOfElements)
			var numberOfElements uint8 = m.numberOfElements
			io.WriteUint8(4, (numberOfElements))

			// Simple Field (startIndex)
			var startIndex uint16 = m.startIndex
			io.WriteUint16(12, (startIndex))
		}
	}
	serializeFunc(m)
}
