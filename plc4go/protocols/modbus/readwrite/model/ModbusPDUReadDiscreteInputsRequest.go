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

// ModbusPDUReadDiscreteInputsRequest is the corresponding interface of ModbusPDUReadDiscreteInputsRequest
type ModbusPDUReadDiscreteInputsRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ModbusPDU
	// GetStartingAddress returns StartingAddress (property field)
	GetStartingAddress() uint16
	// GetQuantity returns Quantity (property field)
	GetQuantity() uint16
	// IsModbusPDUReadDiscreteInputsRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsModbusPDUReadDiscreteInputsRequest()
	// CreateBuilder creates a ModbusPDUReadDiscreteInputsRequestBuilder
	CreateModbusPDUReadDiscreteInputsRequestBuilder() ModbusPDUReadDiscreteInputsRequestBuilder
}

// _ModbusPDUReadDiscreteInputsRequest is the data-structure of this message
type _ModbusPDUReadDiscreteInputsRequest struct {
	ModbusPDUContract
	StartingAddress uint16
	Quantity        uint16
}

var _ ModbusPDUReadDiscreteInputsRequest = (*_ModbusPDUReadDiscreteInputsRequest)(nil)
var _ ModbusPDURequirements = (*_ModbusPDUReadDiscreteInputsRequest)(nil)

// NewModbusPDUReadDiscreteInputsRequest factory function for _ModbusPDUReadDiscreteInputsRequest
func NewModbusPDUReadDiscreteInputsRequest(startingAddress uint16, quantity uint16) *_ModbusPDUReadDiscreteInputsRequest {
	_result := &_ModbusPDUReadDiscreteInputsRequest{
		ModbusPDUContract: NewModbusPDU(),
		StartingAddress:   startingAddress,
		Quantity:          quantity,
	}
	_result.ModbusPDUContract.(*_ModbusPDU)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ModbusPDUReadDiscreteInputsRequestBuilder is a builder for ModbusPDUReadDiscreteInputsRequest
type ModbusPDUReadDiscreteInputsRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(startingAddress uint16, quantity uint16) ModbusPDUReadDiscreteInputsRequestBuilder
	// WithStartingAddress adds StartingAddress (property field)
	WithStartingAddress(uint16) ModbusPDUReadDiscreteInputsRequestBuilder
	// WithQuantity adds Quantity (property field)
	WithQuantity(uint16) ModbusPDUReadDiscreteInputsRequestBuilder
	// Build builds the ModbusPDUReadDiscreteInputsRequest or returns an error if something is wrong
	Build() (ModbusPDUReadDiscreteInputsRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ModbusPDUReadDiscreteInputsRequest
}

// NewModbusPDUReadDiscreteInputsRequestBuilder() creates a ModbusPDUReadDiscreteInputsRequestBuilder
func NewModbusPDUReadDiscreteInputsRequestBuilder() ModbusPDUReadDiscreteInputsRequestBuilder {
	return &_ModbusPDUReadDiscreteInputsRequestBuilder{_ModbusPDUReadDiscreteInputsRequest: new(_ModbusPDUReadDiscreteInputsRequest)}
}

type _ModbusPDUReadDiscreteInputsRequestBuilder struct {
	*_ModbusPDUReadDiscreteInputsRequest

	parentBuilder *_ModbusPDUBuilder

	err *utils.MultiError
}

var _ (ModbusPDUReadDiscreteInputsRequestBuilder) = (*_ModbusPDUReadDiscreteInputsRequestBuilder)(nil)

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) setParent(contract ModbusPDUContract) {
	b.ModbusPDUContract = contract
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) WithMandatoryFields(startingAddress uint16, quantity uint16) ModbusPDUReadDiscreteInputsRequestBuilder {
	return b.WithStartingAddress(startingAddress).WithQuantity(quantity)
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) WithStartingAddress(startingAddress uint16) ModbusPDUReadDiscreteInputsRequestBuilder {
	b.StartingAddress = startingAddress
	return b
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) WithQuantity(quantity uint16) ModbusPDUReadDiscreteInputsRequestBuilder {
	b.Quantity = quantity
	return b
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) Build() (ModbusPDUReadDiscreteInputsRequest, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ModbusPDUReadDiscreteInputsRequest.deepCopy(), nil
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) MustBuild() ModbusPDUReadDiscreteInputsRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) Done() ModbusPDUBuilder {
	return b.parentBuilder
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) buildForModbusPDU() (ModbusPDU, error) {
	return b.Build()
}

func (b *_ModbusPDUReadDiscreteInputsRequestBuilder) DeepCopy() any {
	_copy := b.CreateModbusPDUReadDiscreteInputsRequestBuilder().(*_ModbusPDUReadDiscreteInputsRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateModbusPDUReadDiscreteInputsRequestBuilder creates a ModbusPDUReadDiscreteInputsRequestBuilder
func (b *_ModbusPDUReadDiscreteInputsRequest) CreateModbusPDUReadDiscreteInputsRequestBuilder() ModbusPDUReadDiscreteInputsRequestBuilder {
	if b == nil {
		return NewModbusPDUReadDiscreteInputsRequestBuilder()
	}
	return &_ModbusPDUReadDiscreteInputsRequestBuilder{_ModbusPDUReadDiscreteInputsRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModbusPDUReadDiscreteInputsRequest) GetErrorFlag() bool {
	return bool(false)
}

func (m *_ModbusPDUReadDiscreteInputsRequest) GetFunctionFlag() uint8 {
	return 0x02
}

func (m *_ModbusPDUReadDiscreteInputsRequest) GetResponse() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModbusPDUReadDiscreteInputsRequest) GetParent() ModbusPDUContract {
	return m.ModbusPDUContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModbusPDUReadDiscreteInputsRequest) GetStartingAddress() uint16 {
	return m.StartingAddress
}

func (m *_ModbusPDUReadDiscreteInputsRequest) GetQuantity() uint16 {
	return m.Quantity
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastModbusPDUReadDiscreteInputsRequest(structType any) ModbusPDUReadDiscreteInputsRequest {
	if casted, ok := structType.(ModbusPDUReadDiscreteInputsRequest); ok {
		return casted
	}
	if casted, ok := structType.(*ModbusPDUReadDiscreteInputsRequest); ok {
		return *casted
	}
	return nil
}

func (m *_ModbusPDUReadDiscreteInputsRequest) GetTypeName() string {
	return "ModbusPDUReadDiscreteInputsRequest"
}

func (m *_ModbusPDUReadDiscreteInputsRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ModbusPDUContract.(*_ModbusPDU).GetLengthInBits(ctx))

	// Simple field (startingAddress)
	lengthInBits += 16

	// Simple field (quantity)
	lengthInBits += 16

	return lengthInBits
}

func (m *_ModbusPDUReadDiscreteInputsRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ModbusPDUReadDiscreteInputsRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ModbusPDU, response bool) (__modbusPDUReadDiscreteInputsRequest ModbusPDUReadDiscreteInputsRequest, err error) {
	m.ModbusPDUContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModbusPDUReadDiscreteInputsRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModbusPDUReadDiscreteInputsRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	startingAddress, err := ReadSimpleField(ctx, "startingAddress", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'startingAddress' field"))
	}
	m.StartingAddress = startingAddress

	quantity, err := ReadSimpleField(ctx, "quantity", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'quantity' field"))
	}
	m.Quantity = quantity

	if closeErr := readBuffer.CloseContext("ModbusPDUReadDiscreteInputsRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModbusPDUReadDiscreteInputsRequest")
	}

	return m, nil
}

func (m *_ModbusPDUReadDiscreteInputsRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModbusPDUReadDiscreteInputsRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModbusPDUReadDiscreteInputsRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModbusPDUReadDiscreteInputsRequest")
		}

		if err := WriteSimpleField[uint16](ctx, "startingAddress", m.GetStartingAddress(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'startingAddress' field")
		}

		if err := WriteSimpleField[uint16](ctx, "quantity", m.GetQuantity(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'quantity' field")
		}

		if popErr := writeBuffer.PopContext("ModbusPDUReadDiscreteInputsRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModbusPDUReadDiscreteInputsRequest")
		}
		return nil
	}
	return m.ModbusPDUContract.(*_ModbusPDU).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModbusPDUReadDiscreteInputsRequest) IsModbusPDUReadDiscreteInputsRequest() {}

func (m *_ModbusPDUReadDiscreteInputsRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ModbusPDUReadDiscreteInputsRequest) deepCopy() *_ModbusPDUReadDiscreteInputsRequest {
	if m == nil {
		return nil
	}
	_ModbusPDUReadDiscreteInputsRequestCopy := &_ModbusPDUReadDiscreteInputsRequest{
		m.ModbusPDUContract.(*_ModbusPDU).deepCopy(),
		m.StartingAddress,
		m.Quantity,
	}
	m.ModbusPDUContract.(*_ModbusPDU)._SubType = m
	return _ModbusPDUReadDiscreteInputsRequestCopy
}

func (m *_ModbusPDUReadDiscreteInputsRequest) String() string {
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
