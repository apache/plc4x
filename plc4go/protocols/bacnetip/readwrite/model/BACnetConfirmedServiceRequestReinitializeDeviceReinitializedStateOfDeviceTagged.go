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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged is the corresponding interface of BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged
type BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetHeader returns Header (property field)
	GetHeader() BACnetTagHeader
	// GetValue returns Value (property field)
	GetValue() BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice
}

// BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedExactly can be used when we want exactly this type and not a type which fulfills BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged.
// This is useful for switch cases.
type BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedExactly interface {
	BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged
	isBACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged() bool
}

// _BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged is the data-structure of this message
type _BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged struct {
	Header BACnetTagHeader
	Value  BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice

	// Arguments.
	TagNumber uint8
	TagClass  TagClass
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetHeader() BACnetTagHeader {
	return m.Header
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetValue() BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged factory function for _BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged
func NewBACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged(header BACnetTagHeader, value BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice, tagNumber uint8, tagClass TagClass) *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged {
	return &_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged{Header: header, Value: value, TagNumber: tagNumber, TagClass: tagClass}
}

// Deprecated: use the interface for direct cast
func CastBACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged(structType any) BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged {
	if casted, ok := structType.(BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetTypeName() string {
	return "BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged"
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits(ctx)

	// Manual Field (value)
	lengthInBits += uint16(int32(m.GetHeader().GetActualLength()) * int32(int32(8)))

	return lengthInBits
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedParse(ctx context.Context, theBytes []byte, tagNumber uint8, tagClass TagClass) (BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged, error) {
	return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber, tagClass)
}

func BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedParseWithBufferProducer(tagNumber uint8, tagClass TagClass) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged, error) {
		return BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedParseWithBuffer(ctx, readBuffer, tagNumber, tagClass)
	}
}

func BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, tagClass TagClass) (BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	header, err := ReadSimpleField[BACnetTagHeader](ctx, "header", ReadComplex[BACnetTagHeader](BACnetTagHeaderParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'header' field"))
	}

	// Validation
	if !(bool((header.GetTagClass()) == (tagClass))) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "tag class doesn't match"})
	}

	// Validation
	if !(bool((bool((header.GetTagClass()) == (TagClass_APPLICATION_TAGS)))) || bool((bool((header.GetActualTagNumber()) == (tagNumber))))) {
		return nil, errors.WithStack(utils.ParseAssertError{Message: "tagnumber doesn't match"})
	}

	value, err := ReadManualField[BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice](ctx, "value", readBuffer, EnsureType[BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice](ReadEnumGenericFailing(ctx, readBuffer, header.GetActualLength(), BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice_COLDSTART)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'value' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged")
	}

	// Create the instance
	return &_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged{
		TagNumber: tagNumber,
		TagClass:  tagClass,
		Header:    header,
		Value:     value,
	}, nil
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged")
	}

	if err := WriteSimpleField[BACnetTagHeader](ctx, "header", m.GetHeader(), WriteComplex[BACnetTagHeader](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'header' field")
	}

	if err := WriteManualField[BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice](ctx, "value", func(ctx context.Context) error { return WriteEnumGeneric(ctx, writeBuffer, m.GetValue()) }, writeBuffer); err != nil {
		return errors.Wrap(err, "Error serializing 'value' field")
	}

	if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetTagNumber() uint8 {
	return m.TagNumber
}
func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) GetTagClass() TagClass {
	return m.TagClass
}

//
////

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) isBACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged() bool {
	return true
}

func (m *_BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
