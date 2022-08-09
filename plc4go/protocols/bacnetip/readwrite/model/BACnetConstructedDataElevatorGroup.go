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

// BACnetConstructedDataElevatorGroup is the corresponding interface of BACnetConstructedDataElevatorGroup
type BACnetConstructedDataElevatorGroup interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetElevatorGroup returns ElevatorGroup (property field)
	GetElevatorGroup() BACnetApplicationTagObjectIdentifier
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagObjectIdentifier
}

// BACnetConstructedDataElevatorGroupExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataElevatorGroup.
// This is useful for switch cases.
type BACnetConstructedDataElevatorGroupExactly interface {
	BACnetConstructedDataElevatorGroup
	isBACnetConstructedDataElevatorGroup() bool
}

// _BACnetConstructedDataElevatorGroup is the data-structure of this message
type _BACnetConstructedDataElevatorGroup struct {
	*_BACnetConstructedData
	ElevatorGroup BACnetApplicationTagObjectIdentifier
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataElevatorGroup) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataElevatorGroup) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_ELEVATOR_GROUP
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataElevatorGroup) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataElevatorGroup) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataElevatorGroup) GetElevatorGroup() BACnetApplicationTagObjectIdentifier {
	return m.ElevatorGroup
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataElevatorGroup) GetActualValue() BACnetApplicationTagObjectIdentifier {
	return CastBACnetApplicationTagObjectIdentifier(m.GetElevatorGroup())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataElevatorGroup factory function for _BACnetConstructedDataElevatorGroup
func NewBACnetConstructedDataElevatorGroup(elevatorGroup BACnetApplicationTagObjectIdentifier, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataElevatorGroup {
	_result := &_BACnetConstructedDataElevatorGroup{
		ElevatorGroup:          elevatorGroup,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataElevatorGroup(structType interface{}) BACnetConstructedDataElevatorGroup {
	if casted, ok := structType.(BACnetConstructedDataElevatorGroup); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataElevatorGroup); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataElevatorGroup) GetTypeName() string {
	return "BACnetConstructedDataElevatorGroup"
}

func (m *_BACnetConstructedDataElevatorGroup) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataElevatorGroup) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (elevatorGroup)
	lengthInBits += m.ElevatorGroup.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataElevatorGroup) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataElevatorGroupParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataElevatorGroup, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataElevatorGroup"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataElevatorGroup")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (elevatorGroup)
	if pullErr := readBuffer.PullContext("elevatorGroup"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for elevatorGroup")
	}
	_elevatorGroup, _elevatorGroupErr := BACnetApplicationTagParse(readBuffer)
	if _elevatorGroupErr != nil {
		return nil, errors.Wrap(_elevatorGroupErr, "Error parsing 'elevatorGroup' field of BACnetConstructedDataElevatorGroup")
	}
	elevatorGroup := _elevatorGroup.(BACnetApplicationTagObjectIdentifier)
	if closeErr := readBuffer.CloseContext("elevatorGroup"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for elevatorGroup")
	}

	// Virtual field
	_actualValue := elevatorGroup
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataElevatorGroup"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataElevatorGroup")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataElevatorGroup{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		ElevatorGroup: elevatorGroup,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataElevatorGroup) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataElevatorGroup"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataElevatorGroup")
		}

		// Simple Field (elevatorGroup)
		if pushErr := writeBuffer.PushContext("elevatorGroup"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for elevatorGroup")
		}
		_elevatorGroupErr := writeBuffer.WriteSerializable(m.GetElevatorGroup())
		if popErr := writeBuffer.PopContext("elevatorGroup"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for elevatorGroup")
		}
		if _elevatorGroupErr != nil {
			return errors.Wrap(_elevatorGroupErr, "Error serializing 'elevatorGroup' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataElevatorGroup"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataElevatorGroup")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataElevatorGroup) isBACnetConstructedDataElevatorGroup() bool {
	return true
}

func (m *_BACnetConstructedDataElevatorGroup) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
