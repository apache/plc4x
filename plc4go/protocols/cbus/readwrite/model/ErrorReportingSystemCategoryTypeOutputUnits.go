/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// ErrorReportingSystemCategoryTypeOutputUnits is the corresponding interface of ErrorReportingSystemCategoryTypeOutputUnits
type ErrorReportingSystemCategoryTypeOutputUnits interface {
	utils.LengthAware
	utils.Serializable
	ErrorReportingSystemCategoryType
	// GetCategoryForType returns CategoryForType (property field)
	GetCategoryForType() ErrorReportingSystemCategoryTypeForOutputUnits
}

// ErrorReportingSystemCategoryTypeOutputUnitsExactly can be used when we want exactly this type and not a type which fulfills ErrorReportingSystemCategoryTypeOutputUnits.
// This is useful for switch cases.
type ErrorReportingSystemCategoryTypeOutputUnitsExactly interface {
	ErrorReportingSystemCategoryTypeOutputUnits
	isErrorReportingSystemCategoryTypeOutputUnits() bool
}

// _ErrorReportingSystemCategoryTypeOutputUnits is the data-structure of this message
type _ErrorReportingSystemCategoryTypeOutputUnits struct {
	*_ErrorReportingSystemCategoryType
	CategoryForType ErrorReportingSystemCategoryTypeForOutputUnits
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetErrorReportingSystemCategoryClass() ErrorReportingSystemCategoryClass {
	return ErrorReportingSystemCategoryClass_OUTPUT_UNITS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) InitializeParent(parent ErrorReportingSystemCategoryType) {
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetParent() ErrorReportingSystemCategoryType {
	return m._ErrorReportingSystemCategoryType
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetCategoryForType() ErrorReportingSystemCategoryTypeForOutputUnits {
	return m.CategoryForType
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewErrorReportingSystemCategoryTypeOutputUnits factory function for _ErrorReportingSystemCategoryTypeOutputUnits
func NewErrorReportingSystemCategoryTypeOutputUnits(categoryForType ErrorReportingSystemCategoryTypeForOutputUnits) *_ErrorReportingSystemCategoryTypeOutputUnits {
	_result := &_ErrorReportingSystemCategoryTypeOutputUnits{
		CategoryForType:                   categoryForType,
		_ErrorReportingSystemCategoryType: NewErrorReportingSystemCategoryType(),
	}
	_result._ErrorReportingSystemCategoryType._ErrorReportingSystemCategoryTypeChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastErrorReportingSystemCategoryTypeOutputUnits(structType interface{}) ErrorReportingSystemCategoryTypeOutputUnits {
	if casted, ok := structType.(ErrorReportingSystemCategoryTypeOutputUnits); ok {
		return casted
	}
	if casted, ok := structType.(*ErrorReportingSystemCategoryTypeOutputUnits); ok {
		return *casted
	}
	return nil
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetTypeName() string {
	return "ErrorReportingSystemCategoryTypeOutputUnits"
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (categoryForType)
	lengthInBits += 4

	return lengthInBits
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ErrorReportingSystemCategoryTypeOutputUnitsParse(readBuffer utils.ReadBuffer, errorReportingSystemCategoryClass ErrorReportingSystemCategoryClass) (ErrorReportingSystemCategoryTypeOutputUnits, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ErrorReportingSystemCategoryTypeOutputUnits"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ErrorReportingSystemCategoryTypeOutputUnits")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (categoryForType)
	if pullErr := readBuffer.PullContext("categoryForType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for categoryForType")
	}
	_categoryForType, _categoryForTypeErr := ErrorReportingSystemCategoryTypeForOutputUnitsParse(readBuffer)
	if _categoryForTypeErr != nil {
		return nil, errors.Wrap(_categoryForTypeErr, "Error parsing 'categoryForType' field of ErrorReportingSystemCategoryTypeOutputUnits")
	}
	categoryForType := _categoryForType
	if closeErr := readBuffer.CloseContext("categoryForType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for categoryForType")
	}

	if closeErr := readBuffer.CloseContext("ErrorReportingSystemCategoryTypeOutputUnits"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ErrorReportingSystemCategoryTypeOutputUnits")
	}

	// Create a partially initialized instance
	_child := &_ErrorReportingSystemCategoryTypeOutputUnits{
		_ErrorReportingSystemCategoryType: &_ErrorReportingSystemCategoryType{},
		CategoryForType:                   categoryForType,
	}
	_child._ErrorReportingSystemCategoryType._ErrorReportingSystemCategoryTypeChildRequirements = _child
	return _child, nil
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ErrorReportingSystemCategoryTypeOutputUnits"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ErrorReportingSystemCategoryTypeOutputUnits")
		}

		// Simple Field (categoryForType)
		if pushErr := writeBuffer.PushContext("categoryForType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for categoryForType")
		}
		_categoryForTypeErr := writeBuffer.WriteSerializable(m.GetCategoryForType())
		if popErr := writeBuffer.PopContext("categoryForType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for categoryForType")
		}
		if _categoryForTypeErr != nil {
			return errors.Wrap(_categoryForTypeErr, "Error serializing 'categoryForType' field")
		}

		if popErr := writeBuffer.PopContext("ErrorReportingSystemCategoryTypeOutputUnits"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ErrorReportingSystemCategoryTypeOutputUnits")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) isErrorReportingSystemCategoryTypeOutputUnits() bool {
	return true
}

func (m *_ErrorReportingSystemCategoryTypeOutputUnits) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
