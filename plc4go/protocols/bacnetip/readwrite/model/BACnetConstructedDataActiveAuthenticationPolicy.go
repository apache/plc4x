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

// BACnetConstructedDataActiveAuthenticationPolicy is the corresponding interface of BACnetConstructedDataActiveAuthenticationPolicy
type BACnetConstructedDataActiveAuthenticationPolicy interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetActiveAuthenticationPolicy returns ActiveAuthenticationPolicy (property field)
	GetActiveAuthenticationPolicy() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
}

// BACnetConstructedDataActiveAuthenticationPolicyExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataActiveAuthenticationPolicy.
// This is useful for switch cases.
type BACnetConstructedDataActiveAuthenticationPolicyExactly interface {
	BACnetConstructedDataActiveAuthenticationPolicy
	isBACnetConstructedDataActiveAuthenticationPolicy() bool
}

// _BACnetConstructedDataActiveAuthenticationPolicy is the data-structure of this message
type _BACnetConstructedDataActiveAuthenticationPolicy struct {
	*_BACnetConstructedData
	ActiveAuthenticationPolicy BACnetApplicationTagUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_ACTIVE_AUTHENTICATION_POLICY
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetActiveAuthenticationPolicy() BACnetApplicationTagUnsignedInteger {
	return m.ActiveAuthenticationPolicy
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetActiveAuthenticationPolicy())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataActiveAuthenticationPolicy factory function for _BACnetConstructedDataActiveAuthenticationPolicy
func NewBACnetConstructedDataActiveAuthenticationPolicy(activeAuthenticationPolicy BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataActiveAuthenticationPolicy {
	_result := &_BACnetConstructedDataActiveAuthenticationPolicy{
		ActiveAuthenticationPolicy: activeAuthenticationPolicy,
		_BACnetConstructedData:     NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataActiveAuthenticationPolicy(structType interface{}) BACnetConstructedDataActiveAuthenticationPolicy {
	if casted, ok := structType.(BACnetConstructedDataActiveAuthenticationPolicy); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataActiveAuthenticationPolicy); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetTypeName() string {
	return "BACnetConstructedDataActiveAuthenticationPolicy"
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (activeAuthenticationPolicy)
	lengthInBits += m.ActiveAuthenticationPolicy.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataActiveAuthenticationPolicyParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataActiveAuthenticationPolicy, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataActiveAuthenticationPolicy"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataActiveAuthenticationPolicy")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (activeAuthenticationPolicy)
	if pullErr := readBuffer.PullContext("activeAuthenticationPolicy"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for activeAuthenticationPolicy")
	}
	_activeAuthenticationPolicy, _activeAuthenticationPolicyErr := BACnetApplicationTagParse(readBuffer)
	if _activeAuthenticationPolicyErr != nil {
		return nil, errors.Wrap(_activeAuthenticationPolicyErr, "Error parsing 'activeAuthenticationPolicy' field of BACnetConstructedDataActiveAuthenticationPolicy")
	}
	activeAuthenticationPolicy := _activeAuthenticationPolicy.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("activeAuthenticationPolicy"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for activeAuthenticationPolicy")
	}

	// Virtual field
	_actualValue := activeAuthenticationPolicy
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataActiveAuthenticationPolicy"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataActiveAuthenticationPolicy")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataActiveAuthenticationPolicy{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		ActiveAuthenticationPolicy: activeAuthenticationPolicy,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataActiveAuthenticationPolicy"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataActiveAuthenticationPolicy")
		}

		// Simple Field (activeAuthenticationPolicy)
		if pushErr := writeBuffer.PushContext("activeAuthenticationPolicy"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for activeAuthenticationPolicy")
		}
		_activeAuthenticationPolicyErr := writeBuffer.WriteSerializable(m.GetActiveAuthenticationPolicy())
		if popErr := writeBuffer.PopContext("activeAuthenticationPolicy"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for activeAuthenticationPolicy")
		}
		if _activeAuthenticationPolicyErr != nil {
			return errors.Wrap(_activeAuthenticationPolicyErr, "Error serializing 'activeAuthenticationPolicy' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataActiveAuthenticationPolicy"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataActiveAuthenticationPolicy")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) isBACnetConstructedDataActiveAuthenticationPolicy() bool {
	return true
}

func (m *_BACnetConstructedDataActiveAuthenticationPolicy) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
