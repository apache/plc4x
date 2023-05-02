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
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetTimerStateChangeValue is the corresponding interface of BACnetTimerStateChangeValue
type BACnetTimerStateChangeValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetPeekedTagHeader returns PeekedTagHeader (property field)
	GetPeekedTagHeader() BACnetTagHeader
	// GetPeekedTagNumber returns PeekedTagNumber (virtual field)
	GetPeekedTagNumber() uint8
	// GetPeekedIsContextTag returns PeekedIsContextTag (virtual field)
	GetPeekedIsContextTag() bool
}

// BACnetTimerStateChangeValueExactly can be used when we want exactly this type and not a type which fulfills BACnetTimerStateChangeValue.
// This is useful for switch cases.
type BACnetTimerStateChangeValueExactly interface {
	BACnetTimerStateChangeValue
	isBACnetTimerStateChangeValue() bool
}

// _BACnetTimerStateChangeValue is the data-structure of this message
type _BACnetTimerStateChangeValue struct {
	_BACnetTimerStateChangeValueChildRequirements
	PeekedTagHeader BACnetTagHeader

	// Arguments.
	ObjectTypeArgument BACnetObjectType
}

type _BACnetTimerStateChangeValueChildRequirements interface {
	utils.Serializable
	GetLengthInBits(ctx context.Context) uint16
	GetPeekedTagNumber() uint8
	GetPeekedIsContextTag() bool
}

type BACnetTimerStateChangeValueParent interface {
	SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetTimerStateChangeValue, serializeChildFunction func() error) error
	GetTypeName() string
}

type BACnetTimerStateChangeValueChild interface {
	utils.Serializable
	InitializeParent(parent BACnetTimerStateChangeValue, peekedTagHeader BACnetTagHeader)
	GetParent() *BACnetTimerStateChangeValue

	GetTypeName() string
	BACnetTimerStateChangeValue
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetTimerStateChangeValue) GetPeekedTagHeader() BACnetTagHeader {
	return m.PeekedTagHeader
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetTimerStateChangeValue) GetPeekedTagNumber() uint8 {
	ctx := context.Background()
	_ = ctx
	return uint8(m.GetPeekedTagHeader().GetActualTagNumber())
}

func (m *_BACnetTimerStateChangeValue) GetPeekedIsContextTag() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetPeekedTagHeader().GetTagClass()) == (TagClass_CONTEXT_SPECIFIC_TAGS)))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetTimerStateChangeValue factory function for _BACnetTimerStateChangeValue
func NewBACnetTimerStateChangeValue(peekedTagHeader BACnetTagHeader, objectTypeArgument BACnetObjectType) *_BACnetTimerStateChangeValue {
	return &_BACnetTimerStateChangeValue{PeekedTagHeader: peekedTagHeader, ObjectTypeArgument: objectTypeArgument}
}

// Deprecated: use the interface for direct cast
func CastBACnetTimerStateChangeValue(structType any) BACnetTimerStateChangeValue {
	if casted, ok := structType.(BACnetTimerStateChangeValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetTimerStateChangeValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetTimerStateChangeValue) GetTypeName() string {
	return "BACnetTimerStateChangeValue"
}

func (m *_BACnetTimerStateChangeValue) GetParentLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetTimerStateChangeValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetTimerStateChangeValueParse(theBytes []byte, objectTypeArgument BACnetObjectType) (BACnetTimerStateChangeValue, error) {
	return BACnetTimerStateChangeValueParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes), objectTypeArgument)
}

func BACnetTimerStateChangeValueParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, objectTypeArgument BACnetObjectType) (BACnetTimerStateChangeValue, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetTimerStateChangeValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetTimerStateChangeValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Peek Field (peekedTagHeader)
	currentPos = positionAware.GetPos()
	if pullErr := readBuffer.PullContext("peekedTagHeader"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for peekedTagHeader")
	}
	peekedTagHeader, _ := BACnetTagHeaderParseWithBuffer(ctx, readBuffer)
	readBuffer.Reset(currentPos)

	// Virtual field
	_peekedTagNumber := peekedTagHeader.GetActualTagNumber()
	peekedTagNumber := uint8(_peekedTagNumber)
	_ = peekedTagNumber

	// Virtual field
	_peekedIsContextTag := bool((peekedTagHeader.GetTagClass()) == (TagClass_CONTEXT_SPECIFIC_TAGS))
	peekedIsContextTag := bool(_peekedIsContextTag)
	_ = peekedIsContextTag

	// Validation
	if !(bool((!(peekedIsContextTag))) || bool((bool(bool(peekedIsContextTag) && bool(bool((peekedTagHeader.GetLengthValueType()) != (0x6)))) && bool(bool((peekedTagHeader.GetLengthValueType()) != (0x7)))))) {
		return nil, errors.WithStack(utils.ParseValidationError{"unexpected opening or closing tag"})
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type BACnetTimerStateChangeValueChildSerializeRequirement interface {
		BACnetTimerStateChangeValue
		InitializeParent(BACnetTimerStateChangeValue, BACnetTagHeader)
		GetParent() BACnetTimerStateChangeValue
	}
	var _childTemp any
	var _child BACnetTimerStateChangeValueChildSerializeRequirement
	var typeSwitchError error
	switch {
	case peekedTagNumber == 0x0 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueNull
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueNullParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x1 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueBoolean
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueBooleanParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x2 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueUnsigned
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueUnsignedParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x3 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueInteger
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueIntegerParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x4 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueReal
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueRealParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x5 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueDouble
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueDoubleParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x6 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueOctetString
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueOctetStringParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x7 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueCharacterString
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueCharacterStringParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x8 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueBitString
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueBitStringParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0x9 && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueEnumerated
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueEnumeratedParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0xA && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueDate
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueDateParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0xB && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueTime
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueTimeParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == 0xC && peekedIsContextTag == bool(false): // BACnetTimerStateChangeValueObjectidentifier
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueObjectidentifierParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == uint8(0) && peekedIsContextTag == bool(true): // BACnetTimerStateChangeValueNoValue
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueNoValueParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == uint8(1) && peekedIsContextTag == bool(true): // BACnetTimerStateChangeValueConstructedValue
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueConstructedValueParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == uint8(2) && peekedIsContextTag == bool(true): // BACnetTimerStateChangeValueDateTime
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueDateTimeParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	case peekedTagNumber == uint8(3) && peekedIsContextTag == bool(true): // BACnetTimerStateChangeValueLightingCommand
		_childTemp, typeSwitchError = BACnetTimerStateChangeValueLightingCommandParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	default:
		typeSwitchError = errors.Errorf("Unmapped type for parameters [peekedTagNumber=%v, peekedIsContextTag=%v]", peekedTagNumber, peekedIsContextTag)
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch of BACnetTimerStateChangeValue")
	}
	_child = _childTemp.(BACnetTimerStateChangeValueChildSerializeRequirement)

	if closeErr := readBuffer.CloseContext("BACnetTimerStateChangeValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetTimerStateChangeValue")
	}

	// Finish initializing
	_child.InitializeParent(_child, peekedTagHeader)
	return _child, nil
}

func (pm *_BACnetTimerStateChangeValue) SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetTimerStateChangeValue, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetTimerStateChangeValue"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetTimerStateChangeValue")
	}
	// Virtual field
	if _peekedTagNumberErr := writeBuffer.WriteVirtual(ctx, "peekedTagNumber", m.GetPeekedTagNumber()); _peekedTagNumberErr != nil {
		return errors.Wrap(_peekedTagNumberErr, "Error serializing 'peekedTagNumber' field")
	}
	// Virtual field
	if _peekedIsContextTagErr := writeBuffer.WriteVirtual(ctx, "peekedIsContextTag", m.GetPeekedIsContextTag()); _peekedIsContextTagErr != nil {
		return errors.Wrap(_peekedIsContextTagErr, "Error serializing 'peekedIsContextTag' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("BACnetTimerStateChangeValue"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetTimerStateChangeValue")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetTimerStateChangeValue) GetObjectTypeArgument() BACnetObjectType {
	return m.ObjectTypeArgument
}

//
////

func (m *_BACnetTimerStateChangeValue) isBACnetTimerStateChangeValue() bool {
	return true
}

func (m *_BACnetTimerStateChangeValue) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
