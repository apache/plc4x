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

// BACnetFaultParameter is the corresponding interface of BACnetFaultParameter
type BACnetFaultParameter interface {
	BACnetFaultParameterContract
	BACnetFaultParameterRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// IsBACnetFaultParameter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetFaultParameter()
	// CreateBuilder creates a BACnetFaultParameterBuilder
	CreateBACnetFaultParameterBuilder() BACnetFaultParameterBuilder
}

// BACnetFaultParameterContract provides a set of functions which can be overwritten by a sub struct
type BACnetFaultParameterContract interface {
	// GetPeekedTagHeader returns PeekedTagHeader (property field)
	GetPeekedTagHeader() BACnetTagHeader
	// GetPeekedTagNumber returns PeekedTagNumber (virtual field)
	GetPeekedTagNumber() uint8
	// IsBACnetFaultParameter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetFaultParameter()
	// CreateBuilder creates a BACnetFaultParameterBuilder
	CreateBACnetFaultParameterBuilder() BACnetFaultParameterBuilder
}

// BACnetFaultParameterRequirements provides a set of functions which need to be implemented by a sub struct
type BACnetFaultParameterRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetPeekedTagNumber returns PeekedTagNumber (discriminator field)
	GetPeekedTagNumber() uint8
}

// _BACnetFaultParameter is the data-structure of this message
type _BACnetFaultParameter struct {
	_SubType        BACnetFaultParameter
	PeekedTagHeader BACnetTagHeader
}

var _ BACnetFaultParameterContract = (*_BACnetFaultParameter)(nil)

// NewBACnetFaultParameter factory function for _BACnetFaultParameter
func NewBACnetFaultParameter(peekedTagHeader BACnetTagHeader) *_BACnetFaultParameter {
	if peekedTagHeader == nil {
		panic("peekedTagHeader of type BACnetTagHeader for BACnetFaultParameter must not be nil")
	}
	return &_BACnetFaultParameter{PeekedTagHeader: peekedTagHeader}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetFaultParameterBuilder is a builder for BACnetFaultParameter
type BACnetFaultParameterBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(peekedTagHeader BACnetTagHeader) BACnetFaultParameterBuilder
	// WithPeekedTagHeader adds PeekedTagHeader (property field)
	WithPeekedTagHeader(BACnetTagHeader) BACnetFaultParameterBuilder
	// WithPeekedTagHeaderBuilder adds PeekedTagHeader (property field) which is build by the builder
	WithPeekedTagHeaderBuilder(func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetFaultParameterBuilder
	// AsBACnetFaultParameterNone converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterNone() interface {
		BACnetFaultParameterNoneBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultCharacterString converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultCharacterString() interface {
		BACnetFaultParameterFaultCharacterStringBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultExtended converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultExtended() interface {
		BACnetFaultParameterFaultExtendedBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultLifeSafety converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultLifeSafety() interface {
		BACnetFaultParameterFaultLifeSafetyBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultState converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultState() interface {
		BACnetFaultParameterFaultStateBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultStatusFlags converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultStatusFlags() interface {
		BACnetFaultParameterFaultStatusFlagsBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultOutOfRange converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultOutOfRange() interface {
		BACnetFaultParameterFaultOutOfRangeBuilder
		Done() BACnetFaultParameterBuilder
	}
	// AsBACnetFaultParameterFaultListed converts this build to a subType of BACnetFaultParameter. It is always possible to return to current builder using Done()
	AsBACnetFaultParameterFaultListed() interface {
		BACnetFaultParameterFaultListedBuilder
		Done() BACnetFaultParameterBuilder
	}
	// Build builds the BACnetFaultParameter or returns an error if something is wrong
	PartialBuild() (BACnetFaultParameterContract, error)
	// MustBuild does the same as Build but panics on error
	PartialMustBuild() BACnetFaultParameterContract
	// Build builds the BACnetFaultParameter or returns an error if something is wrong
	Build() (BACnetFaultParameter, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetFaultParameter
}

// NewBACnetFaultParameterBuilder() creates a BACnetFaultParameterBuilder
func NewBACnetFaultParameterBuilder() BACnetFaultParameterBuilder {
	return &_BACnetFaultParameterBuilder{_BACnetFaultParameter: new(_BACnetFaultParameter)}
}

type _BACnetFaultParameterChildBuilder interface {
	utils.Copyable
	setParent(BACnetFaultParameterContract)
	buildForBACnetFaultParameter() (BACnetFaultParameter, error)
}

type _BACnetFaultParameterBuilder struct {
	*_BACnetFaultParameter

	childBuilder _BACnetFaultParameterChildBuilder

	err *utils.MultiError
}

var _ (BACnetFaultParameterBuilder) = (*_BACnetFaultParameterBuilder)(nil)

func (b *_BACnetFaultParameterBuilder) WithMandatoryFields(peekedTagHeader BACnetTagHeader) BACnetFaultParameterBuilder {
	return b.WithPeekedTagHeader(peekedTagHeader)
}

func (b *_BACnetFaultParameterBuilder) WithPeekedTagHeader(peekedTagHeader BACnetTagHeader) BACnetFaultParameterBuilder {
	b.PeekedTagHeader = peekedTagHeader
	return b
}

func (b *_BACnetFaultParameterBuilder) WithPeekedTagHeaderBuilder(builderSupplier func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetFaultParameterBuilder {
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

func (b *_BACnetFaultParameterBuilder) PartialBuild() (BACnetFaultParameterContract, error) {
	if b.PeekedTagHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'peekedTagHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetFaultParameter.deepCopy(), nil
}

func (b *_BACnetFaultParameterBuilder) PartialMustBuild() BACnetFaultParameterContract {
	build, err := b.PartialBuild()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterNone() interface {
	BACnetFaultParameterNoneBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterNoneBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterNoneBuilder().(*_BACnetFaultParameterNoneBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultCharacterString() interface {
	BACnetFaultParameterFaultCharacterStringBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultCharacterStringBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultCharacterStringBuilder().(*_BACnetFaultParameterFaultCharacterStringBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultExtended() interface {
	BACnetFaultParameterFaultExtendedBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultExtendedBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultExtendedBuilder().(*_BACnetFaultParameterFaultExtendedBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultLifeSafety() interface {
	BACnetFaultParameterFaultLifeSafetyBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultLifeSafetyBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultLifeSafetyBuilder().(*_BACnetFaultParameterFaultLifeSafetyBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultState() interface {
	BACnetFaultParameterFaultStateBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultStateBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultStateBuilder().(*_BACnetFaultParameterFaultStateBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultStatusFlags() interface {
	BACnetFaultParameterFaultStatusFlagsBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultStatusFlagsBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultStatusFlagsBuilder().(*_BACnetFaultParameterFaultStatusFlagsBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultOutOfRange() interface {
	BACnetFaultParameterFaultOutOfRangeBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultOutOfRangeBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultOutOfRangeBuilder().(*_BACnetFaultParameterFaultOutOfRangeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) AsBACnetFaultParameterFaultListed() interface {
	BACnetFaultParameterFaultListedBuilder
	Done() BACnetFaultParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetFaultParameterFaultListedBuilder
		Done() BACnetFaultParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetFaultParameterFaultListedBuilder().(*_BACnetFaultParameterFaultListedBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetFaultParameterBuilder) Build() (BACnetFaultParameter, error) {
	v, err := b.PartialBuild()
	if err != nil {
		return nil, errors.Wrap(err, "error occurred during partial build")
	}
	if b.childBuilder == nil {
		return nil, errors.New("no child builder present")
	}
	b.childBuilder.setParent(v)
	return b.childBuilder.buildForBACnetFaultParameter()
}

func (b *_BACnetFaultParameterBuilder) MustBuild() BACnetFaultParameter {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetFaultParameterBuilder) DeepCopy() any {
	_copy := b.CreateBACnetFaultParameterBuilder().(*_BACnetFaultParameterBuilder)
	_copy.childBuilder = b.childBuilder.DeepCopy().(_BACnetFaultParameterChildBuilder)
	_copy.childBuilder.setParent(_copy)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetFaultParameterBuilder creates a BACnetFaultParameterBuilder
func (b *_BACnetFaultParameter) CreateBACnetFaultParameterBuilder() BACnetFaultParameterBuilder {
	if b == nil {
		return NewBACnetFaultParameterBuilder()
	}
	return &_BACnetFaultParameterBuilder{_BACnetFaultParameter: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetFaultParameter) GetPeekedTagHeader() BACnetTagHeader {
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

func (pm *_BACnetFaultParameter) GetPeekedTagNumber() uint8 {
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
func CastBACnetFaultParameter(structType any) BACnetFaultParameter {
	if casted, ok := structType.(BACnetFaultParameter); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetFaultParameter); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetFaultParameter) GetTypeName() string {
	return "BACnetFaultParameter"
}

func (m *_BACnetFaultParameter) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetFaultParameter) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func BACnetFaultParameterParse[T BACnetFaultParameter](ctx context.Context, theBytes []byte) (T, error) {
	return BACnetFaultParameterParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetFaultParameterParseWithBufferProducer[T BACnetFaultParameter]() func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := BACnetFaultParameterParseWithBuffer[T](ctx, readBuffer)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, nil
	}
}

func BACnetFaultParameterParseWithBuffer[T BACnetFaultParameter](ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	v, err := (&_BACnetFaultParameter{}).parse(ctx, readBuffer)
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

func (m *_BACnetFaultParameter) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__bACnetFaultParameter BACnetFaultParameter, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetFaultParameter"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetFaultParameter")
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
	var _child BACnetFaultParameter
	switch {
	case peekedTagNumber == uint8(0): // BACnetFaultParameterNone
		if _child, err = new(_BACnetFaultParameterNone).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterNone for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(1): // BACnetFaultParameterFaultCharacterString
		if _child, err = new(_BACnetFaultParameterFaultCharacterString).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultCharacterString for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(2): // BACnetFaultParameterFaultExtended
		if _child, err = new(_BACnetFaultParameterFaultExtended).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultExtended for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(3): // BACnetFaultParameterFaultLifeSafety
		if _child, err = new(_BACnetFaultParameterFaultLifeSafety).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultLifeSafety for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(4): // BACnetFaultParameterFaultState
		if _child, err = new(_BACnetFaultParameterFaultState).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultState for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(5): // BACnetFaultParameterFaultStatusFlags
		if _child, err = new(_BACnetFaultParameterFaultStatusFlags).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultStatusFlags for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(6): // BACnetFaultParameterFaultOutOfRange
		if _child, err = new(_BACnetFaultParameterFaultOutOfRange).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultOutOfRange for type-switch of BACnetFaultParameter")
		}
	case peekedTagNumber == uint8(7): // BACnetFaultParameterFaultListed
		if _child, err = new(_BACnetFaultParameterFaultListed).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetFaultParameterFaultListed for type-switch of BACnetFaultParameter")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [peekedTagNumber=%v]", peekedTagNumber)
	}

	if closeErr := readBuffer.CloseContext("BACnetFaultParameter"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetFaultParameter")
	}

	return _child, nil
}

func (pm *_BACnetFaultParameter) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetFaultParameter, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetFaultParameter"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetFaultParameter")
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

	if popErr := writeBuffer.PopContext("BACnetFaultParameter"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetFaultParameter")
	}
	return nil
}

func (m *_BACnetFaultParameter) IsBACnetFaultParameter() {}

func (m *_BACnetFaultParameter) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetFaultParameter) deepCopy() *_BACnetFaultParameter {
	if m == nil {
		return nil
	}
	_BACnetFaultParameterCopy := &_BACnetFaultParameter{
		nil, // will be set by child
		m.PeekedTagHeader.DeepCopy().(BACnetTagHeader),
	}
	return _BACnetFaultParameterCopy
}
