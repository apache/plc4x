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

// BACnetActionTagged is the corresponding interface of BACnetActionTagged
type BACnetActionTagged interface {
	utils.LengthAware
	utils.Serializable
	// GetHeader returns Header (property field)
	GetHeader() BACnetTagHeader
	// GetValue returns Value (property field)
	GetValue() BACnetAction
}

// BACnetActionTaggedExactly can be used when we want exactly this type and not a type which fulfills BACnetActionTagged.
// This is useful for switch cases.
type BACnetActionTaggedExactly interface {
	BACnetActionTagged
	isBACnetActionTagged() bool
}

// _BACnetActionTagged is the data-structure of this message
type _BACnetActionTagged struct {
	Header BACnetTagHeader
	Value  BACnetAction

	// Arguments.
	TagNumber uint8
	TagClass  TagClass
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetActionTagged) GetHeader() BACnetTagHeader {
	return m.Header
}

func (m *_BACnetActionTagged) GetValue() BACnetAction {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetActionTagged factory function for _BACnetActionTagged
func NewBACnetActionTagged(header BACnetTagHeader, value BACnetAction, tagNumber uint8, tagClass TagClass) *_BACnetActionTagged {
	return &_BACnetActionTagged{Header: header, Value: value, TagNumber: tagNumber, TagClass: tagClass}
}

// Deprecated: use the interface for direct cast
func CastBACnetActionTagged(structType interface{}) BACnetActionTagged {
	if casted, ok := structType.(BACnetActionTagged); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetActionTagged); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetActionTagged) GetTypeName() string {
	return "BACnetActionTagged"
}

func (m *_BACnetActionTagged) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetActionTagged) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits()

	// Manual Field (value)
	lengthInBits += uint16(int32(m.GetHeader().GetActualLength()) * int32(int32(8)))

	return lengthInBits
}

func (m *_BACnetActionTagged) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetActionTaggedParse(readBuffer utils.ReadBuffer, tagNumber uint8, tagClass TagClass) (BACnetActionTagged, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetActionTagged"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetActionTagged")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (header)
	if pullErr := readBuffer.PullContext("header"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for header")
	}
	_header, _headerErr := BACnetTagHeaderParse(readBuffer)
	if _headerErr != nil {
		return nil, errors.Wrap(_headerErr, "Error parsing 'header' field of BACnetActionTagged")
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
	_value, _valueErr := ReadEnumGenericFailing(readBuffer, header.GetActualLength(), BACnetAction_DIRECT)
	if _valueErr != nil {
		return nil, errors.Wrap(_valueErr, "Error parsing 'value' field of BACnetActionTagged")
	}
	var value BACnetAction
	if _value != nil {
		value = _value.(BACnetAction)
	}

	if closeErr := readBuffer.CloseContext("BACnetActionTagged"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetActionTagged")
	}

	// Create the instance
	return &_BACnetActionTagged{
		TagNumber: tagNumber,
		TagClass:  tagClass,
		Header:    header,
		Value:     value,
	}, nil
}

func (m *_BACnetActionTagged) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetActionTagged"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetActionTagged")
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

	if popErr := writeBuffer.PopContext("BACnetActionTagged"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetActionTagged")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetActionTagged) GetTagNumber() uint8 {
	return m.TagNumber
}
func (m *_BACnetActionTagged) GetTagClass() TagClass {
	return m.TagClass
}

//
////

func (m *_BACnetActionTagged) isBACnetActionTagged() bool {
	return true
}

func (m *_BACnetActionTagged) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
