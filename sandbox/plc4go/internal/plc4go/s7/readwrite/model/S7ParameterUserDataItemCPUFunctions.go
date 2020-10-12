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
type S7ParameterUserDataItemCPUFunctions struct {
	method                  uint8
	cpuFunctionType         uint8
	cpuFunctionGroup        uint8
	cpuSubfunction          uint8
	sequenceNumber          uint8
	dataUnitReferenceNumber *uint8
	lastDataUnit            *uint8
	errorCode               *uint16
	S7ParameterUserDataItem
}

// The corresponding interface
type IS7ParameterUserDataItemCPUFunctions interface {
	IS7ParameterUserDataItem
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7ParameterUserDataItemCPUFunctions) ItemType() uint8 {
	return 0x12
}

func (m S7ParameterUserDataItemCPUFunctions) initialize() spi.Message {
	return m
}

func NewS7ParameterUserDataItemCPUFunctions(method uint8, cpuFunctionType uint8, cpuFunctionGroup uint8, cpuSubfunction uint8, sequenceNumber uint8, dataUnitReferenceNumber *uint8, lastDataUnit *uint8, errorCode *uint16) S7ParameterUserDataItemInitializer {
	return &S7ParameterUserDataItemCPUFunctions{method: method, cpuFunctionType: cpuFunctionType, cpuFunctionGroup: cpuFunctionGroup, cpuSubfunction: cpuSubfunction, sequenceNumber: sequenceNumber, dataUnitReferenceNumber: dataUnitReferenceNumber, lastDataUnit: lastDataUnit, errorCode: errorCode}
}

func CastIS7ParameterUserDataItemCPUFunctions(structType interface{}) IS7ParameterUserDataItemCPUFunctions {
	castFunc := func(typ interface{}) IS7ParameterUserDataItemCPUFunctions {
		if iS7ParameterUserDataItemCPUFunctions, ok := typ.(IS7ParameterUserDataItemCPUFunctions); ok {
			return iS7ParameterUserDataItemCPUFunctions
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7ParameterUserDataItemCPUFunctions(structType interface{}) S7ParameterUserDataItemCPUFunctions {
	castFunc := func(typ interface{}) S7ParameterUserDataItemCPUFunctions {
		if sS7ParameterUserDataItemCPUFunctions, ok := typ.(S7ParameterUserDataItemCPUFunctions); ok {
			return sS7ParameterUserDataItemCPUFunctions
		}
		return S7ParameterUserDataItemCPUFunctions{}
	}
	return castFunc(structType)
}

func (m S7ParameterUserDataItemCPUFunctions) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7ParameterUserDataItem.LengthInBits()

	// Implicit Field (itemLength)
	lengthInBits += 8

	// Simple field (method)
	lengthInBits += 8

	// Simple field (cpuFunctionType)
	lengthInBits += 4

	// Simple field (cpuFunctionGroup)
	lengthInBits += 4

	// Simple field (cpuSubfunction)
	lengthInBits += 8

	// Simple field (sequenceNumber)
	lengthInBits += 8

	// Optional Field (dataUnitReferenceNumber)
	if m.dataUnitReferenceNumber != nil {
		lengthInBits += 8
	}

	// Optional Field (lastDataUnit)
	if m.lastDataUnit != nil {
		lengthInBits += 8
	}

	// Optional Field (errorCode)
	if m.errorCode != nil {
		lengthInBits += 16
	}

	return lengthInBits
}

func (m S7ParameterUserDataItemCPUFunctions) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterUserDataItemCPUFunctionsParse(io spi.ReadBuffer) (S7ParameterUserDataItemInitializer, error) {

	// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint8 = io.ReadUint8(8)

	// Simple Field (method)
	var method uint8 = io.ReadUint8(8)

	// Simple Field (cpuFunctionType)
	var cpuFunctionType uint8 = io.ReadUint8(4)

	// Simple Field (cpuFunctionGroup)
	var cpuFunctionGroup uint8 = io.ReadUint8(4)

	// Simple Field (cpuSubfunction)
	var cpuSubfunction uint8 = io.ReadUint8(8)

	// Simple Field (sequenceNumber)
	var sequenceNumber uint8 = io.ReadUint8(8)

	// Optional Field (dataUnitReferenceNumber) (Can be skipped, if a given expression evaluates to false)
	var dataUnitReferenceNumber *uint8 = nil
	if bool((cpuFunctionType) == (8)) {
		_val := io.ReadUint8(8)
		dataUnitReferenceNumber = &_val
	}

	// Optional Field (lastDataUnit) (Can be skipped, if a given expression evaluates to false)
	var lastDataUnit *uint8 = nil
	if bool((cpuFunctionType) == (8)) {
		_val := io.ReadUint8(8)
		lastDataUnit = &_val
	}

	// Optional Field (errorCode) (Can be skipped, if a given expression evaluates to false)
	var errorCode *uint16 = nil
	if bool((cpuFunctionType) == (8)) {
		_val := io.ReadUint16(16)
		errorCode = &_val
	}

	// Create the instance
	return NewS7ParameterUserDataItemCPUFunctions(method, cpuFunctionType, cpuFunctionGroup, cpuSubfunction, sequenceNumber, dataUnitReferenceNumber, lastDataUnit, errorCode), nil
}

func (m S7ParameterUserDataItemCPUFunctions) Serialize(io spi.WriteBuffer) {

	// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	itemLength := uint8(uint8(uint8(m.LengthInBytes())) - uint8(uint8(2)))
	io.WriteUint8(8, (itemLength))

	// Simple Field (method)
	method := uint8(m.method)
	io.WriteUint8(8, (method))

	// Simple Field (cpuFunctionType)
	cpuFunctionType := uint8(m.cpuFunctionType)
	io.WriteUint8(4, (cpuFunctionType))

	// Simple Field (cpuFunctionGroup)
	cpuFunctionGroup := uint8(m.cpuFunctionGroup)
	io.WriteUint8(4, (cpuFunctionGroup))

	// Simple Field (cpuSubfunction)
	cpuSubfunction := uint8(m.cpuSubfunction)
	io.WriteUint8(8, (cpuSubfunction))

	// Simple Field (sequenceNumber)
	sequenceNumber := uint8(m.sequenceNumber)
	io.WriteUint8(8, (sequenceNumber))

	// Optional Field (dataUnitReferenceNumber) (Can be skipped, if the value is null)
	var dataUnitReferenceNumber *uint8 = nil
	if m.dataUnitReferenceNumber != nil {
		dataUnitReferenceNumber = m.dataUnitReferenceNumber
		io.WriteUint8(8, *(dataUnitReferenceNumber))
	}

	// Optional Field (lastDataUnit) (Can be skipped, if the value is null)
	var lastDataUnit *uint8 = nil
	if m.lastDataUnit != nil {
		lastDataUnit = m.lastDataUnit
		io.WriteUint8(8, *(lastDataUnit))
	}

	// Optional Field (errorCode) (Can be skipped, if the value is null)
	var errorCode *uint16 = nil
	if m.errorCode != nil {
		errorCode = m.errorCode
		io.WriteUint16(16, *(errorCode))
	}
}
