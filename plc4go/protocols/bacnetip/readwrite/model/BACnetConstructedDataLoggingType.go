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

// BACnetConstructedDataLoggingType is the corresponding interface of BACnetConstructedDataLoggingType
type BACnetConstructedDataLoggingType interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetLoggingType returns LoggingType (property field)
	GetLoggingType() BACnetLoggingTypeTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetLoggingTypeTagged
}

// BACnetConstructedDataLoggingTypeExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataLoggingType.
// This is useful for switch cases.
type BACnetConstructedDataLoggingTypeExactly interface {
	BACnetConstructedDataLoggingType
	isBACnetConstructedDataLoggingType() bool
}

// _BACnetConstructedDataLoggingType is the data-structure of this message
type _BACnetConstructedDataLoggingType struct {
	*_BACnetConstructedData
	LoggingType BACnetLoggingTypeTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLoggingType) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataLoggingType) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_LOGGING_TYPE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLoggingType) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataLoggingType) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLoggingType) GetLoggingType() BACnetLoggingTypeTagged {
	return m.LoggingType
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataLoggingType) GetActualValue() BACnetLoggingTypeTagged {
	return CastBACnetLoggingTypeTagged(m.GetLoggingType())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataLoggingType factory function for _BACnetConstructedDataLoggingType
func NewBACnetConstructedDataLoggingType(loggingType BACnetLoggingTypeTagged, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLoggingType {
	_result := &_BACnetConstructedDataLoggingType{
		LoggingType:            loggingType,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLoggingType(structType interface{}) BACnetConstructedDataLoggingType {
	if casted, ok := structType.(BACnetConstructedDataLoggingType); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLoggingType); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLoggingType) GetTypeName() string {
	return "BACnetConstructedDataLoggingType"
}

func (m *_BACnetConstructedDataLoggingType) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataLoggingType) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (loggingType)
	lengthInBits += m.LoggingType.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataLoggingType) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataLoggingTypeParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataLoggingType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLoggingType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLoggingType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (loggingType)
	if pullErr := readBuffer.PullContext("loggingType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for loggingType")
	}
	_loggingType, _loggingTypeErr := BACnetLoggingTypeTaggedParse(readBuffer, uint8(uint8(0)), TagClass(TagClass_APPLICATION_TAGS))
	if _loggingTypeErr != nil {
		return nil, errors.Wrap(_loggingTypeErr, "Error parsing 'loggingType' field of BACnetConstructedDataLoggingType")
	}
	loggingType := _loggingType.(BACnetLoggingTypeTagged)
	if closeErr := readBuffer.CloseContext("loggingType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for loggingType")
	}

	// Virtual field
	_actualValue := loggingType
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLoggingType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLoggingType")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataLoggingType{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		LoggingType: loggingType,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataLoggingType) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLoggingType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLoggingType")
		}

		// Simple Field (loggingType)
		if pushErr := writeBuffer.PushContext("loggingType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for loggingType")
		}
		_loggingTypeErr := writeBuffer.WriteSerializable(m.GetLoggingType())
		if popErr := writeBuffer.PopContext("loggingType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for loggingType")
		}
		if _loggingTypeErr != nil {
			return errors.Wrap(_loggingTypeErr, "Error serializing 'loggingType' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLoggingType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLoggingType")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLoggingType) isBACnetConstructedDataLoggingType() bool {
	return true
}

func (m *_BACnetConstructedDataLoggingType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
