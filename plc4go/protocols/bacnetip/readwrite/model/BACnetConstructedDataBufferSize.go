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

// BACnetConstructedDataBufferSize is the corresponding interface of BACnetConstructedDataBufferSize
type BACnetConstructedDataBufferSize interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetBufferSize returns BufferSize (property field)
	GetBufferSize() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
}

// BACnetConstructedDataBufferSizeExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataBufferSize.
// This is useful for switch cases.
type BACnetConstructedDataBufferSizeExactly interface {
	BACnetConstructedDataBufferSize
	isBACnetConstructedDataBufferSize() bool
}

// _BACnetConstructedDataBufferSize is the data-structure of this message
type _BACnetConstructedDataBufferSize struct {
	*_BACnetConstructedData
	BufferSize BACnetApplicationTagUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataBufferSize) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataBufferSize) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_BUFFER_SIZE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataBufferSize) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataBufferSize) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataBufferSize) GetBufferSize() BACnetApplicationTagUnsignedInteger {
	return m.BufferSize
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataBufferSize) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetBufferSize())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataBufferSize factory function for _BACnetConstructedDataBufferSize
func NewBACnetConstructedDataBufferSize(bufferSize BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataBufferSize {
	_result := &_BACnetConstructedDataBufferSize{
		BufferSize:             bufferSize,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataBufferSize(structType interface{}) BACnetConstructedDataBufferSize {
	if casted, ok := structType.(BACnetConstructedDataBufferSize); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataBufferSize); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataBufferSize) GetTypeName() string {
	return "BACnetConstructedDataBufferSize"
}

func (m *_BACnetConstructedDataBufferSize) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataBufferSize) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (bufferSize)
	lengthInBits += m.BufferSize.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataBufferSize) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataBufferSizeParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataBufferSize, error) {
	return BACnetConstructedDataBufferSizeParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataBufferSizeParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataBufferSize, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataBufferSize"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataBufferSize")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (bufferSize)
	if pullErr := readBuffer.PullContext("bufferSize"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for bufferSize")
	}
	_bufferSize, _bufferSizeErr := BACnetApplicationTagParseWithBuffer(readBuffer)
	if _bufferSizeErr != nil {
		return nil, errors.Wrap(_bufferSizeErr, "Error parsing 'bufferSize' field of BACnetConstructedDataBufferSize")
	}
	bufferSize := _bufferSize.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("bufferSize"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for bufferSize")
	}

	// Virtual field
	_actualValue := bufferSize
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataBufferSize"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataBufferSize")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataBufferSize{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		BufferSize: bufferSize,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataBufferSize) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataBufferSize) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataBufferSize"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataBufferSize")
		}

		// Simple Field (bufferSize)
		if pushErr := writeBuffer.PushContext("bufferSize"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for bufferSize")
		}
		_bufferSizeErr := writeBuffer.WriteSerializable(m.GetBufferSize())
		if popErr := writeBuffer.PopContext("bufferSize"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for bufferSize")
		}
		if _bufferSizeErr != nil {
			return errors.Wrap(_bufferSizeErr, "Error serializing 'bufferSize' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataBufferSize"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataBufferSize")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataBufferSize) isBACnetConstructedDataBufferSize() bool {
	return true
}

func (m *_BACnetConstructedDataBufferSize) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
