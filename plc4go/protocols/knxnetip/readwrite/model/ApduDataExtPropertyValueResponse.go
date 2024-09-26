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

// ApduDataExtPropertyValueResponse is the corresponding interface of ApduDataExtPropertyValueResponse
type ApduDataExtPropertyValueResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduDataExt
	// GetObjectIndex returns ObjectIndex (property field)
	GetObjectIndex() uint8
	// GetPropertyId returns PropertyId (property field)
	GetPropertyId() uint8
	// GetCount returns Count (property field)
	GetCount() uint8
	// GetIndex returns Index (property field)
	GetIndex() uint16
	// GetData returns Data (property field)
	GetData() []byte
	// IsApduDataExtPropertyValueResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataExtPropertyValueResponse()
	// CreateBuilder creates a ApduDataExtPropertyValueResponseBuilder
	CreateApduDataExtPropertyValueResponseBuilder() ApduDataExtPropertyValueResponseBuilder
}

// _ApduDataExtPropertyValueResponse is the data-structure of this message
type _ApduDataExtPropertyValueResponse struct {
	ApduDataExtContract
	ObjectIndex uint8
	PropertyId  uint8
	Count       uint8
	Index       uint16
	Data        []byte
}

var _ ApduDataExtPropertyValueResponse = (*_ApduDataExtPropertyValueResponse)(nil)
var _ ApduDataExtRequirements = (*_ApduDataExtPropertyValueResponse)(nil)

// NewApduDataExtPropertyValueResponse factory function for _ApduDataExtPropertyValueResponse
func NewApduDataExtPropertyValueResponse(objectIndex uint8, propertyId uint8, count uint8, index uint16, data []byte, length uint8) *_ApduDataExtPropertyValueResponse {
	_result := &_ApduDataExtPropertyValueResponse{
		ApduDataExtContract: NewApduDataExt(length),
		ObjectIndex:         objectIndex,
		PropertyId:          propertyId,
		Count:               count,
		Index:               index,
		Data:                data,
	}
	_result.ApduDataExtContract.(*_ApduDataExt)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataExtPropertyValueResponseBuilder is a builder for ApduDataExtPropertyValueResponse
type ApduDataExtPropertyValueResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(objectIndex uint8, propertyId uint8, count uint8, index uint16, data []byte) ApduDataExtPropertyValueResponseBuilder
	// WithObjectIndex adds ObjectIndex (property field)
	WithObjectIndex(uint8) ApduDataExtPropertyValueResponseBuilder
	// WithPropertyId adds PropertyId (property field)
	WithPropertyId(uint8) ApduDataExtPropertyValueResponseBuilder
	// WithCount adds Count (property field)
	WithCount(uint8) ApduDataExtPropertyValueResponseBuilder
	// WithIndex adds Index (property field)
	WithIndex(uint16) ApduDataExtPropertyValueResponseBuilder
	// WithData adds Data (property field)
	WithData(...byte) ApduDataExtPropertyValueResponseBuilder
	// Build builds the ApduDataExtPropertyValueResponse or returns an error if something is wrong
	Build() (ApduDataExtPropertyValueResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataExtPropertyValueResponse
}

// NewApduDataExtPropertyValueResponseBuilder() creates a ApduDataExtPropertyValueResponseBuilder
func NewApduDataExtPropertyValueResponseBuilder() ApduDataExtPropertyValueResponseBuilder {
	return &_ApduDataExtPropertyValueResponseBuilder{_ApduDataExtPropertyValueResponse: new(_ApduDataExtPropertyValueResponse)}
}

type _ApduDataExtPropertyValueResponseBuilder struct {
	*_ApduDataExtPropertyValueResponse

	parentBuilder *_ApduDataExtBuilder

	err *utils.MultiError
}

var _ (ApduDataExtPropertyValueResponseBuilder) = (*_ApduDataExtPropertyValueResponseBuilder)(nil)

func (b *_ApduDataExtPropertyValueResponseBuilder) setParent(contract ApduDataExtContract) {
	b.ApduDataExtContract = contract
}

func (b *_ApduDataExtPropertyValueResponseBuilder) WithMandatoryFields(objectIndex uint8, propertyId uint8, count uint8, index uint16, data []byte) ApduDataExtPropertyValueResponseBuilder {
	return b.WithObjectIndex(objectIndex).WithPropertyId(propertyId).WithCount(count).WithIndex(index).WithData(data...)
}

func (b *_ApduDataExtPropertyValueResponseBuilder) WithObjectIndex(objectIndex uint8) ApduDataExtPropertyValueResponseBuilder {
	b.ObjectIndex = objectIndex
	return b
}

func (b *_ApduDataExtPropertyValueResponseBuilder) WithPropertyId(propertyId uint8) ApduDataExtPropertyValueResponseBuilder {
	b.PropertyId = propertyId
	return b
}

func (b *_ApduDataExtPropertyValueResponseBuilder) WithCount(count uint8) ApduDataExtPropertyValueResponseBuilder {
	b.Count = count
	return b
}

func (b *_ApduDataExtPropertyValueResponseBuilder) WithIndex(index uint16) ApduDataExtPropertyValueResponseBuilder {
	b.Index = index
	return b
}

func (b *_ApduDataExtPropertyValueResponseBuilder) WithData(data ...byte) ApduDataExtPropertyValueResponseBuilder {
	b.Data = data
	return b
}

func (b *_ApduDataExtPropertyValueResponseBuilder) Build() (ApduDataExtPropertyValueResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataExtPropertyValueResponse.deepCopy(), nil
}

func (b *_ApduDataExtPropertyValueResponseBuilder) MustBuild() ApduDataExtPropertyValueResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataExtPropertyValueResponseBuilder) Done() ApduDataExtBuilder {
	return b.parentBuilder
}

func (b *_ApduDataExtPropertyValueResponseBuilder) buildForApduDataExt() (ApduDataExt, error) {
	return b.Build()
}

func (b *_ApduDataExtPropertyValueResponseBuilder) DeepCopy() any {
	_copy := b.CreateApduDataExtPropertyValueResponseBuilder().(*_ApduDataExtPropertyValueResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataExtPropertyValueResponseBuilder creates a ApduDataExtPropertyValueResponseBuilder
func (b *_ApduDataExtPropertyValueResponse) CreateApduDataExtPropertyValueResponseBuilder() ApduDataExtPropertyValueResponseBuilder {
	if b == nil {
		return NewApduDataExtPropertyValueResponseBuilder()
	}
	return &_ApduDataExtPropertyValueResponseBuilder{_ApduDataExtPropertyValueResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataExtPropertyValueResponse) GetExtApciType() uint8 {
	return 0x16
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataExtPropertyValueResponse) GetParent() ApduDataExtContract {
	return m.ApduDataExtContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ApduDataExtPropertyValueResponse) GetObjectIndex() uint8 {
	return m.ObjectIndex
}

func (m *_ApduDataExtPropertyValueResponse) GetPropertyId() uint8 {
	return m.PropertyId
}

func (m *_ApduDataExtPropertyValueResponse) GetCount() uint8 {
	return m.Count
}

func (m *_ApduDataExtPropertyValueResponse) GetIndex() uint16 {
	return m.Index
}

func (m *_ApduDataExtPropertyValueResponse) GetData() []byte {
	return m.Data
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastApduDataExtPropertyValueResponse(structType any) ApduDataExtPropertyValueResponse {
	if casted, ok := structType.(ApduDataExtPropertyValueResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataExtPropertyValueResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataExtPropertyValueResponse) GetTypeName() string {
	return "ApduDataExtPropertyValueResponse"
}

func (m *_ApduDataExtPropertyValueResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataExtContract.(*_ApduDataExt).GetLengthInBits(ctx))

	// Simple field (objectIndex)
	lengthInBits += 8

	// Simple field (propertyId)
	lengthInBits += 8

	// Simple field (count)
	lengthInBits += 4

	// Simple field (index)
	lengthInBits += 12

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *_ApduDataExtPropertyValueResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataExtPropertyValueResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduDataExt, length uint8) (__apduDataExtPropertyValueResponse ApduDataExtPropertyValueResponse, err error) {
	m.ApduDataExtContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataExtPropertyValueResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataExtPropertyValueResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	objectIndex, err := ReadSimpleField(ctx, "objectIndex", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'objectIndex' field"))
	}
	m.ObjectIndex = objectIndex

	propertyId, err := ReadSimpleField(ctx, "propertyId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'propertyId' field"))
	}
	m.PropertyId = propertyId

	count, err := ReadSimpleField(ctx, "count", ReadUnsignedByte(readBuffer, uint8(4)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'count' field"))
	}
	m.Count = count

	index, err := ReadSimpleField(ctx, "index", ReadUnsignedShort(readBuffer, uint8(12)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'index' field"))
	}
	m.Index = index

	data, err := readBuffer.ReadByteArray("data", int(int32(length)-int32(int32(5))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'data' field"))
	}
	m.Data = data

	if closeErr := readBuffer.CloseContext("ApduDataExtPropertyValueResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataExtPropertyValueResponse")
	}

	return m, nil
}

func (m *_ApduDataExtPropertyValueResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataExtPropertyValueResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtPropertyValueResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataExtPropertyValueResponse")
		}

		if err := WriteSimpleField[uint8](ctx, "objectIndex", m.GetObjectIndex(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'objectIndex' field")
		}

		if err := WriteSimpleField[uint8](ctx, "propertyId", m.GetPropertyId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'propertyId' field")
		}

		if err := WriteSimpleField[uint8](ctx, "count", m.GetCount(), WriteUnsignedByte(writeBuffer, 4)); err != nil {
			return errors.Wrap(err, "Error serializing 'count' field")
		}

		if err := WriteSimpleField[uint16](ctx, "index", m.GetIndex(), WriteUnsignedShort(writeBuffer, 12)); err != nil {
			return errors.Wrap(err, "Error serializing 'index' field")
		}

		if err := WriteByteArrayField(ctx, "data", m.GetData(), WriteByteArray(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'data' field")
		}

		if popErr := writeBuffer.PopContext("ApduDataExtPropertyValueResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataExtPropertyValueResponse")
		}
		return nil
	}
	return m.ApduDataExtContract.(*_ApduDataExt).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataExtPropertyValueResponse) IsApduDataExtPropertyValueResponse() {}

func (m *_ApduDataExtPropertyValueResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataExtPropertyValueResponse) deepCopy() *_ApduDataExtPropertyValueResponse {
	if m == nil {
		return nil
	}
	_ApduDataExtPropertyValueResponseCopy := &_ApduDataExtPropertyValueResponse{
		m.ApduDataExtContract.(*_ApduDataExt).deepCopy(),
		m.ObjectIndex,
		m.PropertyId,
		m.Count,
		m.Index,
		utils.DeepCopySlice[byte, byte](m.Data),
	}
	m.ApduDataExtContract.(*_ApduDataExt)._SubType = m
	return _ApduDataExtPropertyValueResponseCopy
}

func (m *_ApduDataExtPropertyValueResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
