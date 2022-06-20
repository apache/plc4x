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

// BACnetConstructedDataReferencePort is the corresponding interface of BACnetConstructedDataReferencePort
type BACnetConstructedDataReferencePort interface {
	BACnetConstructedData
	// GetReferencePort returns ReferencePort (property field)
	GetReferencePort() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataReferencePort is the data-structure of this message
type _BACnetConstructedDataReferencePort struct {
	*_BACnetConstructedData
	ReferencePort BACnetApplicationTagUnsignedInteger

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataReferencePort) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataReferencePort) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_REFERENCE_PORT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataReferencePort) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataReferencePort) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataReferencePort) GetReferencePort() BACnetApplicationTagUnsignedInteger {
	return m.ReferencePort
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataReferencePort) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetReferencePort())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataReferencePort factory function for _BACnetConstructedDataReferencePort
func NewBACnetConstructedDataReferencePort(referencePort BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataReferencePort {
	_result := &_BACnetConstructedDataReferencePort{
		ReferencePort:          referencePort,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataReferencePort(structType interface{}) BACnetConstructedDataReferencePort {
	if casted, ok := structType.(BACnetConstructedDataReferencePort); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataReferencePort); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataReferencePort) GetTypeName() string {
	return "BACnetConstructedDataReferencePort"
}

func (m *_BACnetConstructedDataReferencePort) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataReferencePort) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (referencePort)
	lengthInBits += m.ReferencePort.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataReferencePort) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataReferencePortParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataReferencePort, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataReferencePort"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataReferencePort")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (referencePort)
	if pullErr := readBuffer.PullContext("referencePort"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for referencePort")
	}
	_referencePort, _referencePortErr := BACnetApplicationTagParse(readBuffer)
	if _referencePortErr != nil {
		return nil, errors.Wrap(_referencePortErr, "Error parsing 'referencePort' field")
	}
	referencePort := _referencePort.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("referencePort"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for referencePort")
	}

	// Virtual field
	_actualValue := referencePort
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataReferencePort"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataReferencePort")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataReferencePort{
		ReferencePort:          referencePort,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataReferencePort) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataReferencePort"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataReferencePort")
		}

		// Simple Field (referencePort)
		if pushErr := writeBuffer.PushContext("referencePort"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for referencePort")
		}
		_referencePortErr := writeBuffer.WriteSerializable(m.GetReferencePort())
		if popErr := writeBuffer.PopContext("referencePort"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for referencePort")
		}
		if _referencePortErr != nil {
			return errors.Wrap(_referencePortErr, "Error serializing 'referencePort' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataReferencePort"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataReferencePort")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataReferencePort) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
