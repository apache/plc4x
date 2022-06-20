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

// BACnetApplicationTag is the corresponding interface of BACnetApplicationTag
type BACnetApplicationTag interface {
	// GetHeader returns Header (property field)
	GetHeader() BACnetTagHeader
	// GetActualTagNumber returns ActualTagNumber (virtual field)
	GetActualTagNumber() uint8
	// GetActualLength returns ActualLength (virtual field)
	GetActualLength() uint32
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetApplicationTag is the data-structure of this message
type _BACnetApplicationTag struct {
	_BACnetApplicationTagChildRequirements
	Header BACnetTagHeader
}

type _BACnetApplicationTagChildRequirements interface {
	GetLengthInBits() uint16
	GetLengthInBitsConditional(lastItem bool) uint16
	Serialize(writeBuffer utils.WriteBuffer) error
}

type BACnetApplicationTagParent interface {
	SerializeParent(writeBuffer utils.WriteBuffer, child BACnetApplicationTag, serializeChildFunction func() error) error
	GetTypeName() string
}

type BACnetApplicationTagChild interface {
	Serialize(writeBuffer utils.WriteBuffer) error
	InitializeParent(parent BACnetApplicationTag, header BACnetTagHeader)
	GetParent() *BACnetApplicationTag

	GetTypeName() string
	BACnetApplicationTag
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetApplicationTag) GetHeader() BACnetTagHeader {
	return m.Header
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetApplicationTag) GetActualTagNumber() uint8 {
	return uint8(m.GetHeader().GetActualTagNumber())
}

func (m *_BACnetApplicationTag) GetActualLength() uint32 {
	return uint32(m.GetHeader().GetActualLength())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetApplicationTag factory function for _BACnetApplicationTag
func NewBACnetApplicationTag(header BACnetTagHeader) *_BACnetApplicationTag {
	return &_BACnetApplicationTag{Header: header}
}

// Deprecated: use the interface for direct cast
func CastBACnetApplicationTag(structType interface{}) BACnetApplicationTag {
	if casted, ok := structType.(BACnetApplicationTag); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetApplicationTag); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetApplicationTag) GetTypeName() string {
	return "BACnetApplicationTag"
}

func (m *_BACnetApplicationTag) GetParentLengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetApplicationTag) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetApplicationTagParse(readBuffer utils.ReadBuffer) (BACnetApplicationTag, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetApplicationTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetApplicationTag")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (header)
	if pullErr := readBuffer.PullContext("header"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for header")
	}
	_header, _headerErr := BACnetTagHeaderParse(readBuffer)
	if _headerErr != nil {
		return nil, errors.Wrap(_headerErr, "Error parsing 'header' field")
	}
	header := _header.(BACnetTagHeader)
	if closeErr := readBuffer.CloseContext("header"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for header")
	}

	// Validation
	if !(bool((header.GetTagClass()) == (TagClass_APPLICATION_TAGS))) {
		return nil, errors.WithStack(utils.ParseValidationError{"should be a application tag"})
	}

	// Virtual field
	_actualTagNumber := header.GetActualTagNumber()
	actualTagNumber := uint8(_actualTagNumber)
	_ = actualTagNumber

	// Virtual field
	_actualLength := header.GetActualLength()
	actualLength := uint32(_actualLength)
	_ = actualLength

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type BACnetApplicationTagChildSerializeRequirement interface {
		BACnetApplicationTag
		InitializeParent(BACnetApplicationTag, BACnetTagHeader)
		GetParent() BACnetApplicationTag
	}
	var _childTemp interface{}
	var _child BACnetApplicationTagChildSerializeRequirement
	var typeSwitchError error
	switch {
	case actualTagNumber == 0x0: // BACnetApplicationTagNull
		_childTemp, typeSwitchError = BACnetApplicationTagNullParse(readBuffer)
	case actualTagNumber == 0x1: // BACnetApplicationTagBoolean
		_childTemp, typeSwitchError = BACnetApplicationTagBooleanParse(readBuffer, header)
	case actualTagNumber == 0x2: // BACnetApplicationTagUnsignedInteger
		_childTemp, typeSwitchError = BACnetApplicationTagUnsignedIntegerParse(readBuffer, header)
	case actualTagNumber == 0x3: // BACnetApplicationTagSignedInteger
		_childTemp, typeSwitchError = BACnetApplicationTagSignedIntegerParse(readBuffer, header)
	case actualTagNumber == 0x4: // BACnetApplicationTagReal
		_childTemp, typeSwitchError = BACnetApplicationTagRealParse(readBuffer)
	case actualTagNumber == 0x5: // BACnetApplicationTagDouble
		_childTemp, typeSwitchError = BACnetApplicationTagDoubleParse(readBuffer)
	case actualTagNumber == 0x6: // BACnetApplicationTagOctetString
		_childTemp, typeSwitchError = BACnetApplicationTagOctetStringParse(readBuffer, header)
	case actualTagNumber == 0x7: // BACnetApplicationTagCharacterString
		_childTemp, typeSwitchError = BACnetApplicationTagCharacterStringParse(readBuffer, header)
	case actualTagNumber == 0x8: // BACnetApplicationTagBitString
		_childTemp, typeSwitchError = BACnetApplicationTagBitStringParse(readBuffer, header)
	case actualTagNumber == 0x9: // BACnetApplicationTagEnumerated
		_childTemp, typeSwitchError = BACnetApplicationTagEnumeratedParse(readBuffer, header)
	case actualTagNumber == 0xA: // BACnetApplicationTagDate
		_childTemp, typeSwitchError = BACnetApplicationTagDateParse(readBuffer)
	case actualTagNumber == 0xB: // BACnetApplicationTagTime
		_childTemp, typeSwitchError = BACnetApplicationTagTimeParse(readBuffer)
	case actualTagNumber == 0xC: // BACnetApplicationTagObjectIdentifier
		_childTemp, typeSwitchError = BACnetApplicationTagObjectIdentifierParse(readBuffer)
	default:
		// TODO: return actual type
		typeSwitchError = errors.New("Unmapped type")
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch.")
	}
	_child = _childTemp.(BACnetApplicationTagChildSerializeRequirement)

	if closeErr := readBuffer.CloseContext("BACnetApplicationTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetApplicationTag")
	}

	// Finish initializing
	_child.InitializeParent(_child, header)
	return _child, nil
}

func (pm *_BACnetApplicationTag) SerializeParent(writeBuffer utils.WriteBuffer, child BACnetApplicationTag, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetApplicationTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetApplicationTag")
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
	// Virtual field
	if _actualTagNumberErr := writeBuffer.WriteVirtual("actualTagNumber", m.GetActualTagNumber()); _actualTagNumberErr != nil {
		return errors.Wrap(_actualTagNumberErr, "Error serializing 'actualTagNumber' field")
	}
	// Virtual field
	if _actualLengthErr := writeBuffer.WriteVirtual("actualLength", m.GetActualLength()); _actualLengthErr != nil {
		return errors.Wrap(_actualLengthErr, "Error serializing 'actualLength' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("BACnetApplicationTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetApplicationTag")
	}
	return nil
}

func (m *_BACnetApplicationTag) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
