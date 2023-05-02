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

// BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord is the corresponding interface of BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
type BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetPeekedTagHeader returns PeekedTagHeader (property field)
	GetPeekedTagHeader() BACnetTagHeader
	// GetOpeningTag returns OpeningTag (property field)
	GetOpeningTag() BACnetOpeningTag
	// GetClosingTag returns ClosingTag (property field)
	GetClosingTag() BACnetClosingTag
	// GetPeekedTagNumber returns PeekedTagNumber (virtual field)
	GetPeekedTagNumber() uint8
}

// BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordExactly can be used when we want exactly this type and not a type which fulfills BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord.
// This is useful for switch cases.
type BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordExactly interface {
	BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
	isBACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord() bool
}

// _BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord is the data-structure of this message
type _BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord struct {
	_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordChildRequirements
	PeekedTagHeader BACnetTagHeader
	OpeningTag      BACnetOpeningTag
	ClosingTag      BACnetClosingTag
}

type _BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordChildRequirements interface {
	utils.Serializable
	GetLengthInBits(ctx context.Context) uint16
	GetPeekedTagNumber() uint8
}

type BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordParent interface {
	SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord, serializeChildFunction func() error) error
	GetTypeName() string
}

type BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordChild interface {
	utils.Serializable
	InitializeParent(parent BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord, peekedTagHeader BACnetTagHeader, openingTag BACnetOpeningTag, closingTag BACnetClosingTag)
	GetParent() *BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord

	GetTypeName() string
	BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetPeekedTagHeader() BACnetTagHeader {
	return m.PeekedTagHeader
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetOpeningTag() BACnetOpeningTag {
	return m.OpeningTag
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetClosingTag() BACnetClosingTag {
	return m.ClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetPeekedTagNumber() uint8 {
	ctx := context.Background()
	_ = ctx
	return uint8(m.GetPeekedTagHeader().GetActualTagNumber())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord factory function for _BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
func NewBACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord(peekedTagHeader BACnetTagHeader, openingTag BACnetOpeningTag, closingTag BACnetClosingTag) *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord {
	return &_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord{PeekedTagHeader: peekedTagHeader, OpeningTag: openingTag, ClosingTag: closingTag}
}

// Deprecated: use the interface for direct cast
func CastBACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord(structType any) BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord {
	if casted, ok := structType.(BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetTypeName() string {
	return "BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord"
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetParentLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (openingTag)
	lengthInBits += m.OpeningTag.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	// Simple field (closingTag)
	lengthInBits += m.ClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordParse(theBytes []byte) (BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord, error) {
	return BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
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

	// Simple Field (openingTag)
	if pullErr := readBuffer.PullContext("openingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for openingTag")
	}
	_openingTag, _openingTagErr := BACnetOpeningTagParseWithBuffer(ctx, readBuffer, uint8(peekedTagHeader.GetActualTagNumber()))
	if _openingTagErr != nil {
		return nil, errors.Wrap(_openingTagErr, "Error parsing 'openingTag' field of BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
	}
	openingTag := _openingTag.(BACnetOpeningTag)
	if closeErr := readBuffer.CloseContext("openingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for openingTag")
	}

	// Virtual field
	_peekedTagNumber := peekedTagHeader.GetActualTagNumber()
	peekedTagNumber := uint8(_peekedTagNumber)
	_ = peekedTagNumber

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordChildSerializeRequirement interface {
		BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
		InitializeParent(BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord, BACnetTagHeader, BACnetOpeningTag, BACnetClosingTag)
		GetParent() BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
	}
	var _childTemp any
	var _child BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordChildSerializeRequirement
	var typeSwitchError error
	switch {
	case peekedTagNumber == 0x0: // BACnetConfirmedServiceRequestAtomicReadFileStream
		_childTemp, typeSwitchError = BACnetConfirmedServiceRequestAtomicReadFileStreamParseWithBuffer(ctx, readBuffer)
	case peekedTagNumber == 0x1: // BACnetConfirmedServiceRequestAtomicReadFileRecord
		_childTemp, typeSwitchError = BACnetConfirmedServiceRequestAtomicReadFileRecordParseWithBuffer(ctx, readBuffer)
	default:
		typeSwitchError = errors.Errorf("Unmapped type for parameters [peekedTagNumber=%v]", peekedTagNumber)
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch of BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
	}
	_child = _childTemp.(BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecordChildSerializeRequirement)

	// Simple Field (closingTag)
	if pullErr := readBuffer.PullContext("closingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for closingTag")
	}
	_closingTag, _closingTagErr := BACnetClosingTagParseWithBuffer(ctx, readBuffer, uint8(peekedTagHeader.GetActualTagNumber()))
	if _closingTagErr != nil {
		return nil, errors.Wrap(_closingTagErr, "Error parsing 'closingTag' field of BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
	}
	closingTag := _closingTag.(BACnetClosingTag)
	if closeErr := readBuffer.CloseContext("closingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for closingTag")
	}

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
	}

	// Finish initializing
	_child.InitializeParent(_child, peekedTagHeader, openingTag, closingTag)
	return _child, nil
}

func (pm *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
	}

	// Simple Field (openingTag)
	if pushErr := writeBuffer.PushContext("openingTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for openingTag")
	}
	_openingTagErr := writeBuffer.WriteSerializable(ctx, m.GetOpeningTag())
	if popErr := writeBuffer.PopContext("openingTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for openingTag")
	}
	if _openingTagErr != nil {
		return errors.Wrap(_openingTagErr, "Error serializing 'openingTag' field")
	}
	// Virtual field
	if _peekedTagNumberErr := writeBuffer.WriteVirtual(ctx, "peekedTagNumber", m.GetPeekedTagNumber()); _peekedTagNumberErr != nil {
		return errors.Wrap(_peekedTagNumberErr, "Error serializing 'peekedTagNumber' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	// Simple Field (closingTag)
	if pushErr := writeBuffer.PushContext("closingTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for closingTag")
	}
	_closingTagErr := writeBuffer.WriteSerializable(ctx, m.GetClosingTag())
	if popErr := writeBuffer.PopContext("closingTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for closingTag")
	}
	if _closingTagErr != nil {
		return errors.Wrap(_closingTagErr, "Error serializing 'closingTag' field")
	}

	if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord")
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) isBACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord() bool {
	return true
}

func (m *_BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
