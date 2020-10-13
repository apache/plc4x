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
	Serialize(io spi.WriteBuffer) error
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

func S7ParameterUserDataItemCPUFunctionsParse(io *spi.ReadBuffer) (S7ParameterUserDataItemInitializer, error) {

	// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _itemLengthErr := io.ReadUint8(8)
	if _itemLengthErr != nil {
		return nil, errors.New("Error parsing 'itemLength' field " + _itemLengthErr.Error())
	}

	// Simple Field (method)
	method, _methodErr := io.ReadUint8(8)
	if _methodErr != nil {
		return nil, errors.New("Error parsing 'method' field " + _methodErr.Error())
	}

	// Simple Field (cpuFunctionType)
	cpuFunctionType, _cpuFunctionTypeErr := io.ReadUint8(4)
	if _cpuFunctionTypeErr != nil {
		return nil, errors.New("Error parsing 'cpuFunctionType' field " + _cpuFunctionTypeErr.Error())
	}

	// Simple Field (cpuFunctionGroup)
	cpuFunctionGroup, _cpuFunctionGroupErr := io.ReadUint8(4)
	if _cpuFunctionGroupErr != nil {
		return nil, errors.New("Error parsing 'cpuFunctionGroup' field " + _cpuFunctionGroupErr.Error())
	}

	// Simple Field (cpuSubfunction)
	cpuSubfunction, _cpuSubfunctionErr := io.ReadUint8(8)
	if _cpuSubfunctionErr != nil {
		return nil, errors.New("Error parsing 'cpuSubfunction' field " + _cpuSubfunctionErr.Error())
	}

	// Simple Field (sequenceNumber)
	sequenceNumber, _sequenceNumberErr := io.ReadUint8(8)
	if _sequenceNumberErr != nil {
		return nil, errors.New("Error parsing 'sequenceNumber' field " + _sequenceNumberErr.Error())
	}

	// Optional Field (dataUnitReferenceNumber) (Can be skipped, if a given expression evaluates to false)
	var dataUnitReferenceNumber *uint8 = nil
	if bool((cpuFunctionType) == (8)) {
		_val, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'dataUnitReferenceNumber' field " + _err.Error())
		}

		dataUnitReferenceNumber = &_val
	}

	// Optional Field (lastDataUnit) (Can be skipped, if a given expression evaluates to false)
	var lastDataUnit *uint8 = nil
	if bool((cpuFunctionType) == (8)) {
		_val, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'lastDataUnit' field " + _err.Error())
		}

		lastDataUnit = &_val
	}

	// Optional Field (errorCode) (Can be skipped, if a given expression evaluates to false)
	var errorCode *uint16 = nil
	if bool((cpuFunctionType) == (8)) {
		_val, _err := io.ReadUint16(16)
		if _err != nil {
			return nil, errors.New("Error parsing 'errorCode' field " + _err.Error())
		}

		errorCode = &_val
	}

	// Create the instance
	return NewS7ParameterUserDataItemCPUFunctions(method, cpuFunctionType, cpuFunctionGroup, cpuSubfunction, sequenceNumber, dataUnitReferenceNumber, lastDataUnit, errorCode), nil
}

func (m S7ParameterUserDataItemCPUFunctions) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Implicit Field (itemLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		itemLength := uint8(uint8(uint8(m.LengthInBytes())) - uint8(uint8(2)))
		_itemLengthErr := io.WriteUint8(8, (itemLength))
		if _itemLengthErr != nil {
			return errors.New("Error serializing 'itemLength' field " + _itemLengthErr.Error())
		}

		// Simple Field (method)
		method := uint8(m.method)
		_methodErr := io.WriteUint8(8, (method))
		if _methodErr != nil {
			return errors.New("Error serializing 'method' field " + _methodErr.Error())
		}

		// Simple Field (cpuFunctionType)
		cpuFunctionType := uint8(m.cpuFunctionType)
		_cpuFunctionTypeErr := io.WriteUint8(4, (cpuFunctionType))
		if _cpuFunctionTypeErr != nil {
			return errors.New("Error serializing 'cpuFunctionType' field " + _cpuFunctionTypeErr.Error())
		}

		// Simple Field (cpuFunctionGroup)
		cpuFunctionGroup := uint8(m.cpuFunctionGroup)
		_cpuFunctionGroupErr := io.WriteUint8(4, (cpuFunctionGroup))
		if _cpuFunctionGroupErr != nil {
			return errors.New("Error serializing 'cpuFunctionGroup' field " + _cpuFunctionGroupErr.Error())
		}

		// Simple Field (cpuSubfunction)
		cpuSubfunction := uint8(m.cpuSubfunction)
		_cpuSubfunctionErr := io.WriteUint8(8, (cpuSubfunction))
		if _cpuSubfunctionErr != nil {
			return errors.New("Error serializing 'cpuSubfunction' field " + _cpuSubfunctionErr.Error())
		}

		// Simple Field (sequenceNumber)
		sequenceNumber := uint8(m.sequenceNumber)
		_sequenceNumberErr := io.WriteUint8(8, (sequenceNumber))
		if _sequenceNumberErr != nil {
			return errors.New("Error serializing 'sequenceNumber' field " + _sequenceNumberErr.Error())
		}

		// Optional Field (dataUnitReferenceNumber) (Can be skipped, if the value is null)
		var dataUnitReferenceNumber *uint8 = nil
		if m.dataUnitReferenceNumber != nil {
			dataUnitReferenceNumber = m.dataUnitReferenceNumber
			_dataUnitReferenceNumberErr := io.WriteUint8(8, *(dataUnitReferenceNumber))
			if _dataUnitReferenceNumberErr != nil {
				return errors.New("Error serializing 'dataUnitReferenceNumber' field " + _dataUnitReferenceNumberErr.Error())
			}
		}

		// Optional Field (lastDataUnit) (Can be skipped, if the value is null)
		var lastDataUnit *uint8 = nil
		if m.lastDataUnit != nil {
			lastDataUnit = m.lastDataUnit
			_lastDataUnitErr := io.WriteUint8(8, *(lastDataUnit))
			if _lastDataUnitErr != nil {
				return errors.New("Error serializing 'lastDataUnit' field " + _lastDataUnitErr.Error())
			}
		}

		// Optional Field (errorCode) (Can be skipped, if the value is null)
		var errorCode *uint16 = nil
		if m.errorCode != nil {
			errorCode = m.errorCode
			_errorCodeErr := io.WriteUint16(16, *(errorCode))
			if _errorCodeErr != nil {
				return errors.New("Error serializing 'errorCode' field " + _errorCodeErr.Error())
			}
		}

		return nil
	}
	return S7ParameterUserDataItemSerialize(io, m.S7ParameterUserDataItem, CastIS7ParameterUserDataItem(m), ser)
}
