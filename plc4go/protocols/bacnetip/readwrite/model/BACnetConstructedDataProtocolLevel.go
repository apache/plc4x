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

// BACnetConstructedDataProtocolLevel is the corresponding interface of BACnetConstructedDataProtocolLevel
type BACnetConstructedDataProtocolLevel interface {
	BACnetConstructedData
	// GetProtocolLevel returns ProtocolLevel (property field)
	GetProtocolLevel() BACnetProtocolLevelTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetProtocolLevelTagged
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataProtocolLevel is the data-structure of this message
type _BACnetConstructedDataProtocolLevel struct {
	*_BACnetConstructedData
	ProtocolLevel BACnetProtocolLevelTagged

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataProtocolLevel) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataProtocolLevel) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_PROTOCOL_LEVEL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataProtocolLevel) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataProtocolLevel) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataProtocolLevel) GetProtocolLevel() BACnetProtocolLevelTagged {
	return m.ProtocolLevel
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataProtocolLevel) GetActualValue() BACnetProtocolLevelTagged {
	return CastBACnetProtocolLevelTagged(m.GetProtocolLevel())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataProtocolLevel factory function for _BACnetConstructedDataProtocolLevel
func NewBACnetConstructedDataProtocolLevel(protocolLevel BACnetProtocolLevelTagged, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataProtocolLevel {
	_result := &_BACnetConstructedDataProtocolLevel{
		ProtocolLevel:          protocolLevel,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataProtocolLevel(structType interface{}) BACnetConstructedDataProtocolLevel {
	if casted, ok := structType.(BACnetConstructedDataProtocolLevel); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataProtocolLevel); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataProtocolLevel) GetTypeName() string {
	return "BACnetConstructedDataProtocolLevel"
}

func (m *_BACnetConstructedDataProtocolLevel) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataProtocolLevel) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (protocolLevel)
	lengthInBits += m.ProtocolLevel.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataProtocolLevel) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataProtocolLevelParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataProtocolLevel, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataProtocolLevel"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataProtocolLevel")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (protocolLevel)
	if pullErr := readBuffer.PullContext("protocolLevel"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for protocolLevel")
	}
	_protocolLevel, _protocolLevelErr := BACnetProtocolLevelTaggedParse(readBuffer, uint8(uint8(0)), TagClass(TagClass_APPLICATION_TAGS))
	if _protocolLevelErr != nil {
		return nil, errors.Wrap(_protocolLevelErr, "Error parsing 'protocolLevel' field")
	}
	protocolLevel := _protocolLevel.(BACnetProtocolLevelTagged)
	if closeErr := readBuffer.CloseContext("protocolLevel"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for protocolLevel")
	}

	// Virtual field
	_actualValue := protocolLevel
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataProtocolLevel"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataProtocolLevel")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataProtocolLevel{
		ProtocolLevel:          protocolLevel,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataProtocolLevel) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataProtocolLevel"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataProtocolLevel")
		}

		// Simple Field (protocolLevel)
		if pushErr := writeBuffer.PushContext("protocolLevel"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for protocolLevel")
		}
		_protocolLevelErr := writeBuffer.WriteSerializable(m.GetProtocolLevel())
		if popErr := writeBuffer.PopContext("protocolLevel"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for protocolLevel")
		}
		if _protocolLevelErr != nil {
			return errors.Wrap(_protocolLevelErr, "Error serializing 'protocolLevel' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataProtocolLevel"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataProtocolLevel")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataProtocolLevel) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
