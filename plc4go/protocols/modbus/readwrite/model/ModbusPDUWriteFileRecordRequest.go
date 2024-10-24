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

// ModbusPDUWriteFileRecordRequest is the corresponding interface of ModbusPDUWriteFileRecordRequest
type ModbusPDUWriteFileRecordRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ModbusPDU
	// GetItems returns Items (property field)
	GetItems() []ModbusPDUWriteFileRecordRequestItem
	// IsModbusPDUWriteFileRecordRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsModbusPDUWriteFileRecordRequest()
	// CreateBuilder creates a ModbusPDUWriteFileRecordRequestBuilder
	CreateModbusPDUWriteFileRecordRequestBuilder() ModbusPDUWriteFileRecordRequestBuilder
}

// _ModbusPDUWriteFileRecordRequest is the data-structure of this message
type _ModbusPDUWriteFileRecordRequest struct {
	ModbusPDUContract
	Items []ModbusPDUWriteFileRecordRequestItem
}

var _ ModbusPDUWriteFileRecordRequest = (*_ModbusPDUWriteFileRecordRequest)(nil)
var _ ModbusPDURequirements = (*_ModbusPDUWriteFileRecordRequest)(nil)

// NewModbusPDUWriteFileRecordRequest factory function for _ModbusPDUWriteFileRecordRequest
func NewModbusPDUWriteFileRecordRequest(items []ModbusPDUWriteFileRecordRequestItem) *_ModbusPDUWriteFileRecordRequest {
	_result := &_ModbusPDUWriteFileRecordRequest{
		ModbusPDUContract: NewModbusPDU(),
		Items:             items,
	}
	_result.ModbusPDUContract.(*_ModbusPDU)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ModbusPDUWriteFileRecordRequestBuilder is a builder for ModbusPDUWriteFileRecordRequest
type ModbusPDUWriteFileRecordRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(items []ModbusPDUWriteFileRecordRequestItem) ModbusPDUWriteFileRecordRequestBuilder
	// WithItems adds Items (property field)
	WithItems(...ModbusPDUWriteFileRecordRequestItem) ModbusPDUWriteFileRecordRequestBuilder
	// Build builds the ModbusPDUWriteFileRecordRequest or returns an error if something is wrong
	Build() (ModbusPDUWriteFileRecordRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ModbusPDUWriteFileRecordRequest
}

// NewModbusPDUWriteFileRecordRequestBuilder() creates a ModbusPDUWriteFileRecordRequestBuilder
func NewModbusPDUWriteFileRecordRequestBuilder() ModbusPDUWriteFileRecordRequestBuilder {
	return &_ModbusPDUWriteFileRecordRequestBuilder{_ModbusPDUWriteFileRecordRequest: new(_ModbusPDUWriteFileRecordRequest)}
}

type _ModbusPDUWriteFileRecordRequestBuilder struct {
	*_ModbusPDUWriteFileRecordRequest

	parentBuilder *_ModbusPDUBuilder

	err *utils.MultiError
}

var _ (ModbusPDUWriteFileRecordRequestBuilder) = (*_ModbusPDUWriteFileRecordRequestBuilder)(nil)

func (b *_ModbusPDUWriteFileRecordRequestBuilder) setParent(contract ModbusPDUContract) {
	b.ModbusPDUContract = contract
}

func (b *_ModbusPDUWriteFileRecordRequestBuilder) WithMandatoryFields(items []ModbusPDUWriteFileRecordRequestItem) ModbusPDUWriteFileRecordRequestBuilder {
	return b.WithItems(items...)
}

func (b *_ModbusPDUWriteFileRecordRequestBuilder) WithItems(items ...ModbusPDUWriteFileRecordRequestItem) ModbusPDUWriteFileRecordRequestBuilder {
	b.Items = items
	return b
}

func (b *_ModbusPDUWriteFileRecordRequestBuilder) Build() (ModbusPDUWriteFileRecordRequest, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ModbusPDUWriteFileRecordRequest.deepCopy(), nil
}

func (b *_ModbusPDUWriteFileRecordRequestBuilder) MustBuild() ModbusPDUWriteFileRecordRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ModbusPDUWriteFileRecordRequestBuilder) Done() ModbusPDUBuilder {
	return b.parentBuilder
}

func (b *_ModbusPDUWriteFileRecordRequestBuilder) buildForModbusPDU() (ModbusPDU, error) {
	return b.Build()
}

func (b *_ModbusPDUWriteFileRecordRequestBuilder) DeepCopy() any {
	_copy := b.CreateModbusPDUWriteFileRecordRequestBuilder().(*_ModbusPDUWriteFileRecordRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateModbusPDUWriteFileRecordRequestBuilder creates a ModbusPDUWriteFileRecordRequestBuilder
func (b *_ModbusPDUWriteFileRecordRequest) CreateModbusPDUWriteFileRecordRequestBuilder() ModbusPDUWriteFileRecordRequestBuilder {
	if b == nil {
		return NewModbusPDUWriteFileRecordRequestBuilder()
	}
	return &_ModbusPDUWriteFileRecordRequestBuilder{_ModbusPDUWriteFileRecordRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModbusPDUWriteFileRecordRequest) GetErrorFlag() bool {
	return bool(false)
}

func (m *_ModbusPDUWriteFileRecordRequest) GetFunctionFlag() uint8 {
	return 0x15
}

func (m *_ModbusPDUWriteFileRecordRequest) GetResponse() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModbusPDUWriteFileRecordRequest) GetParent() ModbusPDUContract {
	return m.ModbusPDUContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModbusPDUWriteFileRecordRequest) GetItems() []ModbusPDUWriteFileRecordRequestItem {
	return m.Items
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastModbusPDUWriteFileRecordRequest(structType any) ModbusPDUWriteFileRecordRequest {
	if casted, ok := structType.(ModbusPDUWriteFileRecordRequest); ok {
		return casted
	}
	if casted, ok := structType.(*ModbusPDUWriteFileRecordRequest); ok {
		return *casted
	}
	return nil
}

func (m *_ModbusPDUWriteFileRecordRequest) GetTypeName() string {
	return "ModbusPDUWriteFileRecordRequest"
}

func (m *_ModbusPDUWriteFileRecordRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ModbusPDUContract.(*_ModbusPDU).GetLengthInBits(ctx))

	// Implicit Field (byteCount)
	lengthInBits += 8

	// Array field
	if len(m.Items) > 0 {
		for _, element := range m.Items {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_ModbusPDUWriteFileRecordRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ModbusPDUWriteFileRecordRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ModbusPDU, response bool) (__modbusPDUWriteFileRecordRequest ModbusPDUWriteFileRecordRequest, err error) {
	m.ModbusPDUContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModbusPDUWriteFileRecordRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModbusPDUWriteFileRecordRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	byteCount, err := ReadImplicitField[uint8](ctx, "byteCount", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'byteCount' field"))
	}
	_ = byteCount

	items, err := ReadLengthArrayField[ModbusPDUWriteFileRecordRequestItem](ctx, "items", ReadComplex[ModbusPDUWriteFileRecordRequestItem](ModbusPDUWriteFileRecordRequestItemParseWithBuffer, readBuffer), int(byteCount))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'items' field"))
	}
	m.Items = items

	if closeErr := readBuffer.CloseContext("ModbusPDUWriteFileRecordRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModbusPDUWriteFileRecordRequest")
	}

	return m, nil
}

func (m *_ModbusPDUWriteFileRecordRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModbusPDUWriteFileRecordRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	itemsArraySizeInBytes := func(items []ModbusPDUWriteFileRecordRequestItem) uint32 {
		var sizeInBytes uint32 = 0
		for _, v := range items {
			sizeInBytes += uint32(v.GetLengthInBytes(ctx))
		}
		return sizeInBytes
	}
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModbusPDUWriteFileRecordRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModbusPDUWriteFileRecordRequest")
		}
		byteCount := uint8(uint8(itemsArraySizeInBytes(m.GetItems())))
		if err := WriteImplicitField(ctx, "byteCount", byteCount, WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'byteCount' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "items", m.GetItems(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'items' field")
		}

		if popErr := writeBuffer.PopContext("ModbusPDUWriteFileRecordRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModbusPDUWriteFileRecordRequest")
		}
		return nil
	}
	return m.ModbusPDUContract.(*_ModbusPDU).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModbusPDUWriteFileRecordRequest) IsModbusPDUWriteFileRecordRequest() {}

func (m *_ModbusPDUWriteFileRecordRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ModbusPDUWriteFileRecordRequest) deepCopy() *_ModbusPDUWriteFileRecordRequest {
	if m == nil {
		return nil
	}
	_ModbusPDUWriteFileRecordRequestCopy := &_ModbusPDUWriteFileRecordRequest{
		m.ModbusPDUContract.(*_ModbusPDU).deepCopy(),
		utils.DeepCopySlice[ModbusPDUWriteFileRecordRequestItem, ModbusPDUWriteFileRecordRequestItem](m.Items),
	}
	m.ModbusPDUContract.(*_ModbusPDU)._SubType = m
	return _ModbusPDUWriteFileRecordRequestCopy
}

func (m *_ModbusPDUWriteFileRecordRequest) String() string {
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
