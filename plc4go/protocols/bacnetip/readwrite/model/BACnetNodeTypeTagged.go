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

// BACnetNodeTypeTagged is the corresponding interface of BACnetNodeTypeTagged
type BACnetNodeTypeTagged interface {
	utils.LengthAware
	utils.Serializable
	// GetHeader returns Header (property field)
	GetHeader() BACnetTagHeader
	// GetValue returns Value (property field)
	GetValue() BACnetNodeType
}

// BACnetNodeTypeTaggedExactly can be used when we want exactly this type and not a type which fulfills BACnetNodeTypeTagged.
// This is useful for switch cases.
type BACnetNodeTypeTaggedExactly interface {
	BACnetNodeTypeTagged
	isBACnetNodeTypeTagged() bool
}

// _BACnetNodeTypeTagged is the data-structure of this message
type _BACnetNodeTypeTagged struct {
	Header BACnetTagHeader
	Value  BACnetNodeType

	// Arguments.
	TagNumber uint8
	TagClass  TagClass
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNodeTypeTagged) GetHeader() BACnetTagHeader {
	return m.Header
}

func (m *_BACnetNodeTypeTagged) GetValue() BACnetNodeType {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetNodeTypeTagged factory function for _BACnetNodeTypeTagged
func NewBACnetNodeTypeTagged(header BACnetTagHeader, value BACnetNodeType, tagNumber uint8, tagClass TagClass) *_BACnetNodeTypeTagged {
	return &_BACnetNodeTypeTagged{Header: header, Value: value, TagNumber: tagNumber, TagClass: tagClass}
}

// Deprecated: use the interface for direct cast
func CastBACnetNodeTypeTagged(structType interface{}) BACnetNodeTypeTagged {
	if casted, ok := structType.(BACnetNodeTypeTagged); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNodeTypeTagged); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNodeTypeTagged) GetTypeName() string {
	return "BACnetNodeTypeTagged"
}

func (m *_BACnetNodeTypeTagged) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetNodeTypeTagged) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits()

	// Manual Field (value)
	lengthInBits += uint16(int32(m.GetHeader().GetActualLength()) * int32(int32(8)))

	return lengthInBits
}

func (m *_BACnetNodeTypeTagged) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetNodeTypeTaggedParse(theBytes []byte, tagNumber uint8, tagClass TagClass) (BACnetNodeTypeTagged, error) {
	return BACnetNodeTypeTaggedParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, tagClass)
}

func BACnetNodeTypeTaggedParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, tagClass TagClass) (BACnetNodeTypeTagged, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNodeTypeTagged"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNodeTypeTagged")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (header)
	if pullErr := readBuffer.PullContext("header"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for header")
	}
	_header, _headerErr := BACnetTagHeaderParseWithBuffer(readBuffer)
	if _headerErr != nil {
		return nil, errors.Wrap(_headerErr, "Error parsing 'header' field of BACnetNodeTypeTagged")
	}
	header := _header.(BACnetTagHeader)
	if closeErr := readBuffer.CloseContext("header"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for header")
	}

	// Validation
	if !(bool((header.GetTagClass()) == (tagClass))) {
		return nil, errors.WithStack(utils.ParseValidationError{"tag class doesn't match"})
	}

	// Validation
	if !(bool((bool((header.GetTagClass()) == (TagClass_APPLICATION_TAGS)))) || bool((bool((header.GetActualTagNumber()) == (tagNumber))))) {
		return nil, errors.WithStack(utils.ParseAssertError{"tagnumber doesn't match"})
	}

	// Manual Field (value)
	_value, _valueErr := ReadEnumGenericFailing(readBuffer, header.GetActualLength(), BACnetNodeType_UNKNOWN)
	if _valueErr != nil {
		return nil, errors.Wrap(_valueErr, "Error parsing 'value' field of BACnetNodeTypeTagged")
	}
	var value BACnetNodeType
	if _value != nil {
		value = _value.(BACnetNodeType)
	}

	if closeErr := readBuffer.CloseContext("BACnetNodeTypeTagged"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNodeTypeTagged")
	}

	// Create the instance
	return &_BACnetNodeTypeTagged{
		TagNumber: tagNumber,
		TagClass:  tagClass,
		Header:    header,
		Value:     value,
	}, nil
}

func (m *_BACnetNodeTypeTagged) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNodeTypeTagged) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetNodeTypeTagged"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetNodeTypeTagged")
	}

	// Simple Field (header)
	if pushErr := writeBuffer.PushContext("header"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for header")
	}
	_headerErr := writeBuffer.WriteSerializable(m.GetHeader())
	if popErr := writeBuffer.PopContext("header"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for header")
	}
	if _headerErr != nil {
		return errors.Wrap(_headerErr, "Error serializing 'header' field")
	}

	// Manual Field (value)
	_valueErr := WriteEnumGeneric(writeBuffer, m.GetValue())
	if _valueErr != nil {
		return errors.Wrap(_valueErr, "Error serializing 'value' field")
	}

	if popErr := writeBuffer.PopContext("BACnetNodeTypeTagged"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetNodeTypeTagged")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetNodeTypeTagged) GetTagNumber() uint8 {
	return m.TagNumber
}
func (m *_BACnetNodeTypeTagged) GetTagClass() TagClass {
	return m.TagClass
}

//
////

func (m *_BACnetNodeTypeTagged) isBACnetNodeTypeTagged() bool {
	return true
}

func (m *_BACnetNodeTypeTagged) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
