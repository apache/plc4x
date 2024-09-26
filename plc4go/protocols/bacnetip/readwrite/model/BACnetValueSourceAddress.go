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

// BACnetValueSourceAddress is the corresponding interface of BACnetValueSourceAddress
type BACnetValueSourceAddress interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetValueSource
	// GetAddress returns Address (property field)
	GetAddress() BACnetAddressEnclosed
	// IsBACnetValueSourceAddress is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetValueSourceAddress()
	// CreateBuilder creates a BACnetValueSourceAddressBuilder
	CreateBACnetValueSourceAddressBuilder() BACnetValueSourceAddressBuilder
}

// _BACnetValueSourceAddress is the data-structure of this message
type _BACnetValueSourceAddress struct {
	BACnetValueSourceContract
	Address BACnetAddressEnclosed
}

var _ BACnetValueSourceAddress = (*_BACnetValueSourceAddress)(nil)
var _ BACnetValueSourceRequirements = (*_BACnetValueSourceAddress)(nil)

// NewBACnetValueSourceAddress factory function for _BACnetValueSourceAddress
func NewBACnetValueSourceAddress(peekedTagHeader BACnetTagHeader, address BACnetAddressEnclosed) *_BACnetValueSourceAddress {
	if address == nil {
		panic("address of type BACnetAddressEnclosed for BACnetValueSourceAddress must not be nil")
	}
	_result := &_BACnetValueSourceAddress{
		BACnetValueSourceContract: NewBACnetValueSource(peekedTagHeader),
		Address:                   address,
	}
	_result.BACnetValueSourceContract.(*_BACnetValueSource)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetValueSourceAddressBuilder is a builder for BACnetValueSourceAddress
type BACnetValueSourceAddressBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(address BACnetAddressEnclosed) BACnetValueSourceAddressBuilder
	// WithAddress adds Address (property field)
	WithAddress(BACnetAddressEnclosed) BACnetValueSourceAddressBuilder
	// WithAddressBuilder adds Address (property field) which is build by the builder
	WithAddressBuilder(func(BACnetAddressEnclosedBuilder) BACnetAddressEnclosedBuilder) BACnetValueSourceAddressBuilder
	// Build builds the BACnetValueSourceAddress or returns an error if something is wrong
	Build() (BACnetValueSourceAddress, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetValueSourceAddress
}

// NewBACnetValueSourceAddressBuilder() creates a BACnetValueSourceAddressBuilder
func NewBACnetValueSourceAddressBuilder() BACnetValueSourceAddressBuilder {
	return &_BACnetValueSourceAddressBuilder{_BACnetValueSourceAddress: new(_BACnetValueSourceAddress)}
}

type _BACnetValueSourceAddressBuilder struct {
	*_BACnetValueSourceAddress

	parentBuilder *_BACnetValueSourceBuilder

	err *utils.MultiError
}

var _ (BACnetValueSourceAddressBuilder) = (*_BACnetValueSourceAddressBuilder)(nil)

func (b *_BACnetValueSourceAddressBuilder) setParent(contract BACnetValueSourceContract) {
	b.BACnetValueSourceContract = contract
}

func (b *_BACnetValueSourceAddressBuilder) WithMandatoryFields(address BACnetAddressEnclosed) BACnetValueSourceAddressBuilder {
	return b.WithAddress(address)
}

func (b *_BACnetValueSourceAddressBuilder) WithAddress(address BACnetAddressEnclosed) BACnetValueSourceAddressBuilder {
	b.Address = address
	return b
}

func (b *_BACnetValueSourceAddressBuilder) WithAddressBuilder(builderSupplier func(BACnetAddressEnclosedBuilder) BACnetAddressEnclosedBuilder) BACnetValueSourceAddressBuilder {
	builder := builderSupplier(b.Address.CreateBACnetAddressEnclosedBuilder())
	var err error
	b.Address, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetAddressEnclosedBuilder failed"))
	}
	return b
}

func (b *_BACnetValueSourceAddressBuilder) Build() (BACnetValueSourceAddress, error) {
	if b.Address == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'address' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetValueSourceAddress.deepCopy(), nil
}

func (b *_BACnetValueSourceAddressBuilder) MustBuild() BACnetValueSourceAddress {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetValueSourceAddressBuilder) Done() BACnetValueSourceBuilder {
	return b.parentBuilder
}

func (b *_BACnetValueSourceAddressBuilder) buildForBACnetValueSource() (BACnetValueSource, error) {
	return b.Build()
}

func (b *_BACnetValueSourceAddressBuilder) DeepCopy() any {
	_copy := b.CreateBACnetValueSourceAddressBuilder().(*_BACnetValueSourceAddressBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetValueSourceAddressBuilder creates a BACnetValueSourceAddressBuilder
func (b *_BACnetValueSourceAddress) CreateBACnetValueSourceAddressBuilder() BACnetValueSourceAddressBuilder {
	if b == nil {
		return NewBACnetValueSourceAddressBuilder()
	}
	return &_BACnetValueSourceAddressBuilder{_BACnetValueSourceAddress: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetValueSourceAddress) GetParent() BACnetValueSourceContract {
	return m.BACnetValueSourceContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetValueSourceAddress) GetAddress() BACnetAddressEnclosed {
	return m.Address
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetValueSourceAddress(structType any) BACnetValueSourceAddress {
	if casted, ok := structType.(BACnetValueSourceAddress); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetValueSourceAddress); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetValueSourceAddress) GetTypeName() string {
	return "BACnetValueSourceAddress"
}

func (m *_BACnetValueSourceAddress) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetValueSourceContract.(*_BACnetValueSource).GetLengthInBits(ctx))

	// Simple field (address)
	lengthInBits += m.Address.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetValueSourceAddress) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetValueSourceAddress) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetValueSource) (__bACnetValueSourceAddress BACnetValueSourceAddress, err error) {
	m.BACnetValueSourceContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetValueSourceAddress"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetValueSourceAddress")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	address, err := ReadSimpleField[BACnetAddressEnclosed](ctx, "address", ReadComplex[BACnetAddressEnclosed](BACnetAddressEnclosedParseWithBufferProducer((uint8)(uint8(2))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'address' field"))
	}
	m.Address = address

	if closeErr := readBuffer.CloseContext("BACnetValueSourceAddress"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetValueSourceAddress")
	}

	return m, nil
}

func (m *_BACnetValueSourceAddress) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetValueSourceAddress) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetValueSourceAddress"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetValueSourceAddress")
		}

		if err := WriteSimpleField[BACnetAddressEnclosed](ctx, "address", m.GetAddress(), WriteComplex[BACnetAddressEnclosed](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'address' field")
		}

		if popErr := writeBuffer.PopContext("BACnetValueSourceAddress"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetValueSourceAddress")
		}
		return nil
	}
	return m.BACnetValueSourceContract.(*_BACnetValueSource).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetValueSourceAddress) IsBACnetValueSourceAddress() {}

func (m *_BACnetValueSourceAddress) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetValueSourceAddress) deepCopy() *_BACnetValueSourceAddress {
	if m == nil {
		return nil
	}
	_BACnetValueSourceAddressCopy := &_BACnetValueSourceAddress{
		m.BACnetValueSourceContract.(*_BACnetValueSource).deepCopy(),
		m.Address.DeepCopy().(BACnetAddressEnclosed),
	}
	m.BACnetValueSourceContract.(*_BACnetValueSource)._SubType = m
	return _BACnetValueSourceAddressCopy
}

func (m *_BACnetValueSourceAddress) String() string {
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
