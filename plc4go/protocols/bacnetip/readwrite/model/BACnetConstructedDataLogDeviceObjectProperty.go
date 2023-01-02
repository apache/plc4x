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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataLogDeviceObjectProperty is the corresponding interface of BACnetConstructedDataLogDeviceObjectProperty
type BACnetConstructedDataLogDeviceObjectProperty interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetLogDeviceObjectProperty returns LogDeviceObjectProperty (property field)
	GetLogDeviceObjectProperty() BACnetDeviceObjectPropertyReference
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetDeviceObjectPropertyReference
}

// BACnetConstructedDataLogDeviceObjectPropertyExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataLogDeviceObjectProperty.
// This is useful for switch cases.
type BACnetConstructedDataLogDeviceObjectPropertyExactly interface {
	BACnetConstructedDataLogDeviceObjectProperty
	isBACnetConstructedDataLogDeviceObjectProperty() bool
}

// _BACnetConstructedDataLogDeviceObjectProperty is the data-structure of this message
type _BACnetConstructedDataLogDeviceObjectProperty struct {
	*_BACnetConstructedData
	LogDeviceObjectProperty BACnetDeviceObjectPropertyReference
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_LOG_DEVICE_OBJECT_PROPERTY
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLogDeviceObjectProperty) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetLogDeviceObjectProperty() BACnetDeviceObjectPropertyReference {
	return m.LogDeviceObjectProperty
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetActualValue() BACnetDeviceObjectPropertyReference {
	return CastBACnetDeviceObjectPropertyReference(m.GetLogDeviceObjectProperty())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataLogDeviceObjectProperty factory function for _BACnetConstructedDataLogDeviceObjectProperty
func NewBACnetConstructedDataLogDeviceObjectProperty(logDeviceObjectProperty BACnetDeviceObjectPropertyReference, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLogDeviceObjectProperty {
	_result := &_BACnetConstructedDataLogDeviceObjectProperty{
		LogDeviceObjectProperty: logDeviceObjectProperty,
		_BACnetConstructedData:  NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLogDeviceObjectProperty(structType interface{}) BACnetConstructedDataLogDeviceObjectProperty {
	if casted, ok := structType.(BACnetConstructedDataLogDeviceObjectProperty); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLogDeviceObjectProperty); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetTypeName() string {
	return "BACnetConstructedDataLogDeviceObjectProperty"
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (logDeviceObjectProperty)
	lengthInBits += m.LogDeviceObjectProperty.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataLogDeviceObjectPropertyParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataLogDeviceObjectProperty, error) {
	return BACnetConstructedDataLogDeviceObjectPropertyParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataLogDeviceObjectPropertyParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataLogDeviceObjectProperty, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLogDeviceObjectProperty"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLogDeviceObjectProperty")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (logDeviceObjectProperty)
	if pullErr := readBuffer.PullContext("logDeviceObjectProperty"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for logDeviceObjectProperty")
	}
	_logDeviceObjectProperty, _logDeviceObjectPropertyErr := BACnetDeviceObjectPropertyReferenceParseWithBuffer(readBuffer)
	if _logDeviceObjectPropertyErr != nil {
		return nil, errors.Wrap(_logDeviceObjectPropertyErr, "Error parsing 'logDeviceObjectProperty' field of BACnetConstructedDataLogDeviceObjectProperty")
	}
	logDeviceObjectProperty := _logDeviceObjectProperty.(BACnetDeviceObjectPropertyReference)
	if closeErr := readBuffer.CloseContext("logDeviceObjectProperty"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for logDeviceObjectProperty")
	}

	// Virtual field
	_actualValue := logDeviceObjectProperty
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLogDeviceObjectProperty"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLogDeviceObjectProperty")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataLogDeviceObjectProperty{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		LogDeviceObjectProperty: logDeviceObjectProperty,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLogDeviceObjectProperty"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLogDeviceObjectProperty")
		}

		// Simple Field (logDeviceObjectProperty)
		if pushErr := writeBuffer.PushContext("logDeviceObjectProperty"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for logDeviceObjectProperty")
		}
		_logDeviceObjectPropertyErr := writeBuffer.WriteSerializable(m.GetLogDeviceObjectProperty())
		if popErr := writeBuffer.PopContext("logDeviceObjectProperty"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for logDeviceObjectProperty")
		}
		if _logDeviceObjectPropertyErr != nil {
			return errors.Wrap(_logDeviceObjectPropertyErr, "Error serializing 'logDeviceObjectProperty' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLogDeviceObjectProperty"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLogDeviceObjectProperty")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) isBACnetConstructedDataLogDeviceObjectProperty() bool {
	return true
}

func (m *_BACnetConstructedDataLogDeviceObjectProperty) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
