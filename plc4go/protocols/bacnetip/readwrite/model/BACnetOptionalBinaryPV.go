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

// BACnetOptionalBinaryPV is the corresponding interface of BACnetOptionalBinaryPV
type BACnetOptionalBinaryPV interface {
	BACnetOptionalBinaryPVContract
	BACnetOptionalBinaryPVRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// IsBACnetOptionalBinaryPV is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetOptionalBinaryPV()
	// CreateBuilder creates a BACnetOptionalBinaryPVBuilder
	CreateBACnetOptionalBinaryPVBuilder() BACnetOptionalBinaryPVBuilder
}

// BACnetOptionalBinaryPVContract provides a set of functions which can be overwritten by a sub struct
type BACnetOptionalBinaryPVContract interface {
	// GetPeekedTagHeader returns PeekedTagHeader (property field)
	GetPeekedTagHeader() BACnetTagHeader
	// GetPeekedTagNumber returns PeekedTagNumber (virtual field)
	GetPeekedTagNumber() uint8
	// IsBACnetOptionalBinaryPV is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetOptionalBinaryPV()
	// CreateBuilder creates a BACnetOptionalBinaryPVBuilder
	CreateBACnetOptionalBinaryPVBuilder() BACnetOptionalBinaryPVBuilder
}

// BACnetOptionalBinaryPVRequirements provides a set of functions which need to be implemented by a sub struct
type BACnetOptionalBinaryPVRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetPeekedTagNumber returns PeekedTagNumber (discriminator field)
	GetPeekedTagNumber() uint8
}

// _BACnetOptionalBinaryPV is the data-structure of this message
type _BACnetOptionalBinaryPV struct {
	_SubType        BACnetOptionalBinaryPV
	PeekedTagHeader BACnetTagHeader
}

var _ BACnetOptionalBinaryPVContract = (*_BACnetOptionalBinaryPV)(nil)

// NewBACnetOptionalBinaryPV factory function for _BACnetOptionalBinaryPV
func NewBACnetOptionalBinaryPV(peekedTagHeader BACnetTagHeader) *_BACnetOptionalBinaryPV {
	if peekedTagHeader == nil {
		panic("peekedTagHeader of type BACnetTagHeader for BACnetOptionalBinaryPV must not be nil")
	}
	return &_BACnetOptionalBinaryPV{PeekedTagHeader: peekedTagHeader}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetOptionalBinaryPVBuilder is a builder for BACnetOptionalBinaryPV
type BACnetOptionalBinaryPVBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(peekedTagHeader BACnetTagHeader) BACnetOptionalBinaryPVBuilder
	// WithPeekedTagHeader adds PeekedTagHeader (property field)
	WithPeekedTagHeader(BACnetTagHeader) BACnetOptionalBinaryPVBuilder
	// WithPeekedTagHeaderBuilder adds PeekedTagHeader (property field) which is build by the builder
	WithPeekedTagHeaderBuilder(func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetOptionalBinaryPVBuilder
	// AsBACnetOptionalBinaryPVNull converts this build to a subType of BACnetOptionalBinaryPV. It is always possible to return to current builder using Done()
	AsBACnetOptionalBinaryPVNull() interface {
		BACnetOptionalBinaryPVNullBuilder
		Done() BACnetOptionalBinaryPVBuilder
	}
	// AsBACnetOptionalBinaryPVValue converts this build to a subType of BACnetOptionalBinaryPV. It is always possible to return to current builder using Done()
	AsBACnetOptionalBinaryPVValue() interface {
		BACnetOptionalBinaryPVValueBuilder
		Done() BACnetOptionalBinaryPVBuilder
	}
	// Build builds the BACnetOptionalBinaryPV or returns an error if something is wrong
	PartialBuild() (BACnetOptionalBinaryPVContract, error)
	// MustBuild does the same as Build but panics on error
	PartialMustBuild() BACnetOptionalBinaryPVContract
	// Build builds the BACnetOptionalBinaryPV or returns an error if something is wrong
	Build() (BACnetOptionalBinaryPV, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetOptionalBinaryPV
}

// NewBACnetOptionalBinaryPVBuilder() creates a BACnetOptionalBinaryPVBuilder
func NewBACnetOptionalBinaryPVBuilder() BACnetOptionalBinaryPVBuilder {
	return &_BACnetOptionalBinaryPVBuilder{_BACnetOptionalBinaryPV: new(_BACnetOptionalBinaryPV)}
}

type _BACnetOptionalBinaryPVChildBuilder interface {
	utils.Copyable
	setParent(BACnetOptionalBinaryPVContract)
	buildForBACnetOptionalBinaryPV() (BACnetOptionalBinaryPV, error)
}

type _BACnetOptionalBinaryPVBuilder struct {
	*_BACnetOptionalBinaryPV

	childBuilder _BACnetOptionalBinaryPVChildBuilder

	err *utils.MultiError
}

var _ (BACnetOptionalBinaryPVBuilder) = (*_BACnetOptionalBinaryPVBuilder)(nil)

func (b *_BACnetOptionalBinaryPVBuilder) WithMandatoryFields(peekedTagHeader BACnetTagHeader) BACnetOptionalBinaryPVBuilder {
	return b.WithPeekedTagHeader(peekedTagHeader)
}

func (b *_BACnetOptionalBinaryPVBuilder) WithPeekedTagHeader(peekedTagHeader BACnetTagHeader) BACnetOptionalBinaryPVBuilder {
	b.PeekedTagHeader = peekedTagHeader
	return b
}

func (b *_BACnetOptionalBinaryPVBuilder) WithPeekedTagHeaderBuilder(builderSupplier func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetOptionalBinaryPVBuilder {
	builder := builderSupplier(b.PeekedTagHeader.CreateBACnetTagHeaderBuilder())
	var err error
	b.PeekedTagHeader, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetTagHeaderBuilder failed"))
	}
	return b
}

func (b *_BACnetOptionalBinaryPVBuilder) PartialBuild() (BACnetOptionalBinaryPVContract, error) {
	if b.PeekedTagHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'peekedTagHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetOptionalBinaryPV.deepCopy(), nil
}

func (b *_BACnetOptionalBinaryPVBuilder) PartialMustBuild() BACnetOptionalBinaryPVContract {
	build, err := b.PartialBuild()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetOptionalBinaryPVBuilder) AsBACnetOptionalBinaryPVNull() interface {
	BACnetOptionalBinaryPVNullBuilder
	Done() BACnetOptionalBinaryPVBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetOptionalBinaryPVNullBuilder
		Done() BACnetOptionalBinaryPVBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetOptionalBinaryPVNullBuilder().(*_BACnetOptionalBinaryPVNullBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetOptionalBinaryPVBuilder) AsBACnetOptionalBinaryPVValue() interface {
	BACnetOptionalBinaryPVValueBuilder
	Done() BACnetOptionalBinaryPVBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetOptionalBinaryPVValueBuilder
		Done() BACnetOptionalBinaryPVBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetOptionalBinaryPVValueBuilder().(*_BACnetOptionalBinaryPVValueBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetOptionalBinaryPVBuilder) Build() (BACnetOptionalBinaryPV, error) {
	v, err := b.PartialBuild()
	if err != nil {
		return nil, errors.Wrap(err, "error occurred during partial build")
	}
	if b.childBuilder == nil {
		return nil, errors.New("no child builder present")
	}
	b.childBuilder.setParent(v)
	return b.childBuilder.buildForBACnetOptionalBinaryPV()
}

func (b *_BACnetOptionalBinaryPVBuilder) MustBuild() BACnetOptionalBinaryPV {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetOptionalBinaryPVBuilder) DeepCopy() any {
	_copy := b.CreateBACnetOptionalBinaryPVBuilder().(*_BACnetOptionalBinaryPVBuilder)
	_copy.childBuilder = b.childBuilder.DeepCopy().(_BACnetOptionalBinaryPVChildBuilder)
	_copy.childBuilder.setParent(_copy)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetOptionalBinaryPVBuilder creates a BACnetOptionalBinaryPVBuilder
func (b *_BACnetOptionalBinaryPV) CreateBACnetOptionalBinaryPVBuilder() BACnetOptionalBinaryPVBuilder {
	if b == nil {
		return NewBACnetOptionalBinaryPVBuilder()
	}
	return &_BACnetOptionalBinaryPVBuilder{_BACnetOptionalBinaryPV: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetOptionalBinaryPV) GetPeekedTagHeader() BACnetTagHeader {
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

func (pm *_BACnetOptionalBinaryPV) GetPeekedTagNumber() uint8 {
	m := pm._SubType
	ctx := context.Background()
	_ = ctx
	return uint8(m.GetPeekedTagHeader().GetActualTagNumber())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetOptionalBinaryPV(structType any) BACnetOptionalBinaryPV {
	if casted, ok := structType.(BACnetOptionalBinaryPV); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetOptionalBinaryPV); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetOptionalBinaryPV) GetTypeName() string {
	return "BACnetOptionalBinaryPV"
}

func (m *_BACnetOptionalBinaryPV) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetOptionalBinaryPV) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func BACnetOptionalBinaryPVParse[T BACnetOptionalBinaryPV](ctx context.Context, theBytes []byte) (T, error) {
	return BACnetOptionalBinaryPVParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetOptionalBinaryPVParseWithBufferProducer[T BACnetOptionalBinaryPV]() func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := BACnetOptionalBinaryPVParseWithBuffer[T](ctx, readBuffer)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, nil
	}
}

func BACnetOptionalBinaryPVParseWithBuffer[T BACnetOptionalBinaryPV](ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	v, err := (&_BACnetOptionalBinaryPV{}).parse(ctx, readBuffer)
	if err != nil {
		var zero T
		return zero, err
	}
	vc, ok := v.(T)
	if !ok {
		var zero T
		return zero, errors.Errorf("Unexpected type %T. Expected type %T", v, *new(T))
	}
	return vc, nil
}

func (m *_BACnetOptionalBinaryPV) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__bACnetOptionalBinaryPV BACnetOptionalBinaryPV, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetOptionalBinaryPV"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetOptionalBinaryPV")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	peekedTagHeader, err := ReadPeekField[BACnetTagHeader](ctx, "peekedTagHeader", ReadComplex[BACnetTagHeader](BACnetTagHeaderParseWithBuffer, readBuffer), 0)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'peekedTagHeader' field"))
	}
	m.PeekedTagHeader = peekedTagHeader

	peekedTagNumber, err := ReadVirtualField[uint8](ctx, "peekedTagNumber", (*uint8)(nil), peekedTagHeader.GetActualTagNumber())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'peekedTagNumber' field"))
	}
	_ = peekedTagNumber

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _child BACnetOptionalBinaryPV
	switch {
	case peekedTagNumber == uint8(0): // BACnetOptionalBinaryPVNull
		if _child, err = new(_BACnetOptionalBinaryPVNull).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetOptionalBinaryPVNull for type-switch of BACnetOptionalBinaryPV")
		}
	case 0 == 0: // BACnetOptionalBinaryPVValue
		if _child, err = new(_BACnetOptionalBinaryPVValue).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetOptionalBinaryPVValue for type-switch of BACnetOptionalBinaryPV")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [peekedTagNumber=%v]", peekedTagNumber)
	}

	if closeErr := readBuffer.CloseContext("BACnetOptionalBinaryPV"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetOptionalBinaryPV")
	}

	return _child, nil
}

func (pm *_BACnetOptionalBinaryPV) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetOptionalBinaryPV, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetOptionalBinaryPV"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetOptionalBinaryPV")
	}
	// Virtual field
	peekedTagNumber := m.GetPeekedTagNumber()
	_ = peekedTagNumber
	if _peekedTagNumberErr := writeBuffer.WriteVirtual(ctx, "peekedTagNumber", m.GetPeekedTagNumber()); _peekedTagNumberErr != nil {
		return errors.Wrap(_peekedTagNumberErr, "Error serializing 'peekedTagNumber' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("BACnetOptionalBinaryPV"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetOptionalBinaryPV")
	}
	return nil
}

func (m *_BACnetOptionalBinaryPV) IsBACnetOptionalBinaryPV() {}

func (m *_BACnetOptionalBinaryPV) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetOptionalBinaryPV) deepCopy() *_BACnetOptionalBinaryPV {
	if m == nil {
		return nil
	}
	_BACnetOptionalBinaryPVCopy := &_BACnetOptionalBinaryPV{
		nil, // will be set by child
		m.PeekedTagHeader.DeepCopy().(BACnetTagHeader),
	}
	return _BACnetOptionalBinaryPVCopy
}
