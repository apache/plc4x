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

// BACnetConstructedDataFileAccessMethod is the corresponding interface of BACnetConstructedDataFileAccessMethod
type BACnetConstructedDataFileAccessMethod interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetFileAccessMethod returns FileAccessMethod (property field)
	GetFileAccessMethod() BACnetFileAccessMethodTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetFileAccessMethodTagged
}

// BACnetConstructedDataFileAccessMethodExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataFileAccessMethod.
// This is useful for switch cases.
type BACnetConstructedDataFileAccessMethodExactly interface {
	BACnetConstructedDataFileAccessMethod
	isBACnetConstructedDataFileAccessMethod() bool
}

// _BACnetConstructedDataFileAccessMethod is the data-structure of this message
type _BACnetConstructedDataFileAccessMethod struct {
	*_BACnetConstructedData
	FileAccessMethod BACnetFileAccessMethodTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataFileAccessMethod) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataFileAccessMethod) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_FILE_ACCESS_METHOD
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataFileAccessMethod) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataFileAccessMethod) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataFileAccessMethod) GetFileAccessMethod() BACnetFileAccessMethodTagged {
	return m.FileAccessMethod
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataFileAccessMethod) GetActualValue() BACnetFileAccessMethodTagged {
	return CastBACnetFileAccessMethodTagged(m.GetFileAccessMethod())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataFileAccessMethod factory function for _BACnetConstructedDataFileAccessMethod
func NewBACnetConstructedDataFileAccessMethod(fileAccessMethod BACnetFileAccessMethodTagged, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataFileAccessMethod {
	_result := &_BACnetConstructedDataFileAccessMethod{
		FileAccessMethod:       fileAccessMethod,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataFileAccessMethod(structType interface{}) BACnetConstructedDataFileAccessMethod {
	if casted, ok := structType.(BACnetConstructedDataFileAccessMethod); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataFileAccessMethod); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataFileAccessMethod) GetTypeName() string {
	return "BACnetConstructedDataFileAccessMethod"
}

func (m *_BACnetConstructedDataFileAccessMethod) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataFileAccessMethod) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (fileAccessMethod)
	lengthInBits += m.FileAccessMethod.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataFileAccessMethod) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataFileAccessMethodParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataFileAccessMethod, error) {
	return BACnetConstructedDataFileAccessMethodParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataFileAccessMethodParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataFileAccessMethod, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataFileAccessMethod"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataFileAccessMethod")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (fileAccessMethod)
	if pullErr := readBuffer.PullContext("fileAccessMethod"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for fileAccessMethod")
	}
	_fileAccessMethod, _fileAccessMethodErr := BACnetFileAccessMethodTaggedParseWithBuffer(readBuffer, uint8(uint8(0)), TagClass(TagClass_APPLICATION_TAGS))
	if _fileAccessMethodErr != nil {
		return nil, errors.Wrap(_fileAccessMethodErr, "Error parsing 'fileAccessMethod' field of BACnetConstructedDataFileAccessMethod")
	}
	fileAccessMethod := _fileAccessMethod.(BACnetFileAccessMethodTagged)
	if closeErr := readBuffer.CloseContext("fileAccessMethod"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for fileAccessMethod")
	}

	// Virtual field
	_actualValue := fileAccessMethod
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataFileAccessMethod"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataFileAccessMethod")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataFileAccessMethod{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		FileAccessMethod: fileAccessMethod,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataFileAccessMethod) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataFileAccessMethod) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataFileAccessMethod"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataFileAccessMethod")
		}

		// Simple Field (fileAccessMethod)
		if pushErr := writeBuffer.PushContext("fileAccessMethod"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for fileAccessMethod")
		}
		_fileAccessMethodErr := writeBuffer.WriteSerializable(m.GetFileAccessMethod())
		if popErr := writeBuffer.PopContext("fileAccessMethod"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for fileAccessMethod")
		}
		if _fileAccessMethodErr != nil {
			return errors.Wrap(_fileAccessMethodErr, "Error serializing 'fileAccessMethod' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataFileAccessMethod"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataFileAccessMethod")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataFileAccessMethod) isBACnetConstructedDataFileAccessMethod() bool {
	return true
}

func (m *_BACnetConstructedDataFileAccessMethod) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
