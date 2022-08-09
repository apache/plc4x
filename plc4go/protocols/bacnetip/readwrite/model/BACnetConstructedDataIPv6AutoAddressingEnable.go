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

// BACnetConstructedDataIPv6AutoAddressingEnable is the corresponding interface of BACnetConstructedDataIPv6AutoAddressingEnable
type BACnetConstructedDataIPv6AutoAddressingEnable interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetAutoAddressingEnable returns AutoAddressingEnable (property field)
	GetAutoAddressingEnable() BACnetApplicationTagBoolean
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagBoolean
}

// BACnetConstructedDataIPv6AutoAddressingEnableExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataIPv6AutoAddressingEnable.
// This is useful for switch cases.
type BACnetConstructedDataIPv6AutoAddressingEnableExactly interface {
	BACnetConstructedDataIPv6AutoAddressingEnable
	isBACnetConstructedDataIPv6AutoAddressingEnable() bool
}

// _BACnetConstructedDataIPv6AutoAddressingEnable is the data-structure of this message
type _BACnetConstructedDataIPv6AutoAddressingEnable struct {
	*_BACnetConstructedData
	AutoAddressingEnable BACnetApplicationTagBoolean
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_IPV6_AUTO_ADDRESSING_ENABLE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetAutoAddressingEnable() BACnetApplicationTagBoolean {
	return m.AutoAddressingEnable
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetActualValue() BACnetApplicationTagBoolean {
	return CastBACnetApplicationTagBoolean(m.GetAutoAddressingEnable())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataIPv6AutoAddressingEnable factory function for _BACnetConstructedDataIPv6AutoAddressingEnable
func NewBACnetConstructedDataIPv6AutoAddressingEnable(autoAddressingEnable BACnetApplicationTagBoolean, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataIPv6AutoAddressingEnable {
	_result := &_BACnetConstructedDataIPv6AutoAddressingEnable{
		AutoAddressingEnable:   autoAddressingEnable,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataIPv6AutoAddressingEnable(structType interface{}) BACnetConstructedDataIPv6AutoAddressingEnable {
	if casted, ok := structType.(BACnetConstructedDataIPv6AutoAddressingEnable); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataIPv6AutoAddressingEnable); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetTypeName() string {
	return "BACnetConstructedDataIPv6AutoAddressingEnable"
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (autoAddressingEnable)
	lengthInBits += m.AutoAddressingEnable.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataIPv6AutoAddressingEnableParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataIPv6AutoAddressingEnable, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataIPv6AutoAddressingEnable"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataIPv6AutoAddressingEnable")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (autoAddressingEnable)
	if pullErr := readBuffer.PullContext("autoAddressingEnable"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for autoAddressingEnable")
	}
	_autoAddressingEnable, _autoAddressingEnableErr := BACnetApplicationTagParse(readBuffer)
	if _autoAddressingEnableErr != nil {
		return nil, errors.Wrap(_autoAddressingEnableErr, "Error parsing 'autoAddressingEnable' field of BACnetConstructedDataIPv6AutoAddressingEnable")
	}
	autoAddressingEnable := _autoAddressingEnable.(BACnetApplicationTagBoolean)
	if closeErr := readBuffer.CloseContext("autoAddressingEnable"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for autoAddressingEnable")
	}

	// Virtual field
	_actualValue := autoAddressingEnable
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataIPv6AutoAddressingEnable"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataIPv6AutoAddressingEnable")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataIPv6AutoAddressingEnable{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		AutoAddressingEnable: autoAddressingEnable,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataIPv6AutoAddressingEnable"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataIPv6AutoAddressingEnable")
		}

		// Simple Field (autoAddressingEnable)
		if pushErr := writeBuffer.PushContext("autoAddressingEnable"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for autoAddressingEnable")
		}
		_autoAddressingEnableErr := writeBuffer.WriteSerializable(m.GetAutoAddressingEnable())
		if popErr := writeBuffer.PopContext("autoAddressingEnable"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for autoAddressingEnable")
		}
		if _autoAddressingEnableErr != nil {
			return errors.Wrap(_autoAddressingEnableErr, "Error serializing 'autoAddressingEnable' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataIPv6AutoAddressingEnable"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataIPv6AutoAddressingEnable")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) isBACnetConstructedDataIPv6AutoAddressingEnable() bool {
	return true
}

func (m *_BACnetConstructedDataIPv6AutoAddressingEnable) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
