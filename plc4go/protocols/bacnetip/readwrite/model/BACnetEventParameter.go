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

// BACnetEventParameter is the corresponding interface of BACnetEventParameter
type BACnetEventParameter interface {
	BACnetEventParameterContract
	BACnetEventParameterRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// IsBACnetEventParameter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetEventParameter()
	// CreateBuilder creates a BACnetEventParameterBuilder
	CreateBACnetEventParameterBuilder() BACnetEventParameterBuilder
}

// BACnetEventParameterContract provides a set of functions which can be overwritten by a sub struct
type BACnetEventParameterContract interface {
	// GetPeekedTagHeader returns PeekedTagHeader (property field)
	GetPeekedTagHeader() BACnetTagHeader
	// GetPeekedTagNumber returns PeekedTagNumber (virtual field)
	GetPeekedTagNumber() uint8
	// IsBACnetEventParameter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetEventParameter()
	// CreateBuilder creates a BACnetEventParameterBuilder
	CreateBACnetEventParameterBuilder() BACnetEventParameterBuilder
}

// BACnetEventParameterRequirements provides a set of functions which need to be implemented by a sub struct
type BACnetEventParameterRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetPeekedTagNumber returns PeekedTagNumber (discriminator field)
	GetPeekedTagNumber() uint8
}

// _BACnetEventParameter is the data-structure of this message
type _BACnetEventParameter struct {
	_SubType        BACnetEventParameter
	PeekedTagHeader BACnetTagHeader
}

var _ BACnetEventParameterContract = (*_BACnetEventParameter)(nil)

// NewBACnetEventParameter factory function for _BACnetEventParameter
func NewBACnetEventParameter(peekedTagHeader BACnetTagHeader) *_BACnetEventParameter {
	if peekedTagHeader == nil {
		panic("peekedTagHeader of type BACnetTagHeader for BACnetEventParameter must not be nil")
	}
	return &_BACnetEventParameter{PeekedTagHeader: peekedTagHeader}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetEventParameterBuilder is a builder for BACnetEventParameter
type BACnetEventParameterBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(peekedTagHeader BACnetTagHeader) BACnetEventParameterBuilder
	// WithPeekedTagHeader adds PeekedTagHeader (property field)
	WithPeekedTagHeader(BACnetTagHeader) BACnetEventParameterBuilder
	// WithPeekedTagHeaderBuilder adds PeekedTagHeader (property field) which is build by the builder
	WithPeekedTagHeaderBuilder(func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetEventParameterBuilder
	// AsBACnetEventParameterChangeOfBitstring converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfBitstring() interface {
		BACnetEventParameterChangeOfBitstringBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfState converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfState() interface {
		BACnetEventParameterChangeOfStateBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfValue converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfValue() interface {
		BACnetEventParameterChangeOfValueBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterCommandFailure converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterCommandFailure() interface {
		BACnetEventParameterCommandFailureBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterFloatingLimit converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterFloatingLimit() interface {
		BACnetEventParameterFloatingLimitBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterOutOfRange converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterOutOfRange() interface {
		BACnetEventParameterOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfLifeSavety converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfLifeSavety() interface {
		BACnetEventParameterChangeOfLifeSavetyBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterExtended converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterExtended() interface {
		BACnetEventParameterExtendedBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterBufferReady converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterBufferReady() interface {
		BACnetEventParameterBufferReadyBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterUnsignedRange converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterUnsignedRange() interface {
		BACnetEventParameterUnsignedRangeBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterAccessEvent converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterAccessEvent() interface {
		BACnetEventParameterAccessEventBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterDoubleOutOfRange converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterDoubleOutOfRange() interface {
		BACnetEventParameterDoubleOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterSignedOutOfRange converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterSignedOutOfRange() interface {
		BACnetEventParameterSignedOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterUnsignedOutOfRange converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterUnsignedOutOfRange() interface {
		BACnetEventParameterUnsignedOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfCharacterString converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfCharacterString() interface {
		BACnetEventParameterChangeOfCharacterStringBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfStatusFlags converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfStatusFlags() interface {
		BACnetEventParameterChangeOfStatusFlagsBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterNone converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterNone() interface {
		BACnetEventParameterNoneBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfDiscreteValue converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfDiscreteValue() interface {
		BACnetEventParameterChangeOfDiscreteValueBuilder
		Done() BACnetEventParameterBuilder
	}
	// AsBACnetEventParameterChangeOfTimer converts this build to a subType of BACnetEventParameter. It is always possible to return to current builder using Done()
	AsBACnetEventParameterChangeOfTimer() interface {
		BACnetEventParameterChangeOfTimerBuilder
		Done() BACnetEventParameterBuilder
	}
	// Build builds the BACnetEventParameter or returns an error if something is wrong
	PartialBuild() (BACnetEventParameterContract, error)
	// MustBuild does the same as Build but panics on error
	PartialMustBuild() BACnetEventParameterContract
	// Build builds the BACnetEventParameter or returns an error if something is wrong
	Build() (BACnetEventParameter, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetEventParameter
}

// NewBACnetEventParameterBuilder() creates a BACnetEventParameterBuilder
func NewBACnetEventParameterBuilder() BACnetEventParameterBuilder {
	return &_BACnetEventParameterBuilder{_BACnetEventParameter: new(_BACnetEventParameter)}
}

type _BACnetEventParameterChildBuilder interface {
	utils.Copyable
	setParent(BACnetEventParameterContract)
	buildForBACnetEventParameter() (BACnetEventParameter, error)
}

type _BACnetEventParameterBuilder struct {
	*_BACnetEventParameter

	childBuilder _BACnetEventParameterChildBuilder

	err *utils.MultiError
}

var _ (BACnetEventParameterBuilder) = (*_BACnetEventParameterBuilder)(nil)

func (b *_BACnetEventParameterBuilder) WithMandatoryFields(peekedTagHeader BACnetTagHeader) BACnetEventParameterBuilder {
	return b.WithPeekedTagHeader(peekedTagHeader)
}

func (b *_BACnetEventParameterBuilder) WithPeekedTagHeader(peekedTagHeader BACnetTagHeader) BACnetEventParameterBuilder {
	b.PeekedTagHeader = peekedTagHeader
	return b
}

func (b *_BACnetEventParameterBuilder) WithPeekedTagHeaderBuilder(builderSupplier func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetEventParameterBuilder {
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

func (b *_BACnetEventParameterBuilder) PartialBuild() (BACnetEventParameterContract, error) {
	if b.PeekedTagHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'peekedTagHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetEventParameter.deepCopy(), nil
}

func (b *_BACnetEventParameterBuilder) PartialMustBuild() BACnetEventParameterContract {
	build, err := b.PartialBuild()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfBitstring() interface {
	BACnetEventParameterChangeOfBitstringBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfBitstringBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfBitstringBuilder().(*_BACnetEventParameterChangeOfBitstringBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfState() interface {
	BACnetEventParameterChangeOfStateBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfStateBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfStateBuilder().(*_BACnetEventParameterChangeOfStateBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfValue() interface {
	BACnetEventParameterChangeOfValueBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfValueBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfValueBuilder().(*_BACnetEventParameterChangeOfValueBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterCommandFailure() interface {
	BACnetEventParameterCommandFailureBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterCommandFailureBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterCommandFailureBuilder().(*_BACnetEventParameterCommandFailureBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterFloatingLimit() interface {
	BACnetEventParameterFloatingLimitBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterFloatingLimitBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterFloatingLimitBuilder().(*_BACnetEventParameterFloatingLimitBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterOutOfRange() interface {
	BACnetEventParameterOutOfRangeBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterOutOfRangeBuilder().(*_BACnetEventParameterOutOfRangeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfLifeSavety() interface {
	BACnetEventParameterChangeOfLifeSavetyBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfLifeSavetyBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfLifeSavetyBuilder().(*_BACnetEventParameterChangeOfLifeSavetyBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterExtended() interface {
	BACnetEventParameterExtendedBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterExtendedBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterExtendedBuilder().(*_BACnetEventParameterExtendedBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterBufferReady() interface {
	BACnetEventParameterBufferReadyBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterBufferReadyBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterBufferReadyBuilder().(*_BACnetEventParameterBufferReadyBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterUnsignedRange() interface {
	BACnetEventParameterUnsignedRangeBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterUnsignedRangeBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterUnsignedRangeBuilder().(*_BACnetEventParameterUnsignedRangeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterAccessEvent() interface {
	BACnetEventParameterAccessEventBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterAccessEventBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterAccessEventBuilder().(*_BACnetEventParameterAccessEventBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterDoubleOutOfRange() interface {
	BACnetEventParameterDoubleOutOfRangeBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterDoubleOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterDoubleOutOfRangeBuilder().(*_BACnetEventParameterDoubleOutOfRangeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterSignedOutOfRange() interface {
	BACnetEventParameterSignedOutOfRangeBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterSignedOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterSignedOutOfRangeBuilder().(*_BACnetEventParameterSignedOutOfRangeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterUnsignedOutOfRange() interface {
	BACnetEventParameterUnsignedOutOfRangeBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterUnsignedOutOfRangeBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterUnsignedOutOfRangeBuilder().(*_BACnetEventParameterUnsignedOutOfRangeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfCharacterString() interface {
	BACnetEventParameterChangeOfCharacterStringBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfCharacterStringBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfCharacterStringBuilder().(*_BACnetEventParameterChangeOfCharacterStringBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfStatusFlags() interface {
	BACnetEventParameterChangeOfStatusFlagsBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfStatusFlagsBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfStatusFlagsBuilder().(*_BACnetEventParameterChangeOfStatusFlagsBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterNone() interface {
	BACnetEventParameterNoneBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterNoneBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterNoneBuilder().(*_BACnetEventParameterNoneBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfDiscreteValue() interface {
	BACnetEventParameterChangeOfDiscreteValueBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfDiscreteValueBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfDiscreteValueBuilder().(*_BACnetEventParameterChangeOfDiscreteValueBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) AsBACnetEventParameterChangeOfTimer() interface {
	BACnetEventParameterChangeOfTimerBuilder
	Done() BACnetEventParameterBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetEventParameterChangeOfTimerBuilder
		Done() BACnetEventParameterBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetEventParameterChangeOfTimerBuilder().(*_BACnetEventParameterChangeOfTimerBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetEventParameterBuilder) Build() (BACnetEventParameter, error) {
	v, err := b.PartialBuild()
	if err != nil {
		return nil, errors.Wrap(err, "error occurred during partial build")
	}
	if b.childBuilder == nil {
		return nil, errors.New("no child builder present")
	}
	b.childBuilder.setParent(v)
	return b.childBuilder.buildForBACnetEventParameter()
}

func (b *_BACnetEventParameterBuilder) MustBuild() BACnetEventParameter {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetEventParameterBuilder) DeepCopy() any {
	_copy := b.CreateBACnetEventParameterBuilder().(*_BACnetEventParameterBuilder)
	_copy.childBuilder = b.childBuilder.DeepCopy().(_BACnetEventParameterChildBuilder)
	_copy.childBuilder.setParent(_copy)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetEventParameterBuilder creates a BACnetEventParameterBuilder
func (b *_BACnetEventParameter) CreateBACnetEventParameterBuilder() BACnetEventParameterBuilder {
	if b == nil {
		return NewBACnetEventParameterBuilder()
	}
	return &_BACnetEventParameterBuilder{_BACnetEventParameter: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetEventParameter) GetPeekedTagHeader() BACnetTagHeader {
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

func (pm *_BACnetEventParameter) GetPeekedTagNumber() uint8 {
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
func CastBACnetEventParameter(structType any) BACnetEventParameter {
	if casted, ok := structType.(BACnetEventParameter); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetEventParameter); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetEventParameter) GetTypeName() string {
	return "BACnetEventParameter"
}

func (m *_BACnetEventParameter) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetEventParameter) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func BACnetEventParameterParse[T BACnetEventParameter](ctx context.Context, theBytes []byte) (T, error) {
	return BACnetEventParameterParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetEventParameterParseWithBufferProducer[T BACnetEventParameter]() func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := BACnetEventParameterParseWithBuffer[T](ctx, readBuffer)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, nil
	}
}

func BACnetEventParameterParseWithBuffer[T BACnetEventParameter](ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	v, err := (&_BACnetEventParameter{}).parse(ctx, readBuffer)
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

func (m *_BACnetEventParameter) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__bACnetEventParameter BACnetEventParameter, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetEventParameter"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetEventParameter")
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
	var _child BACnetEventParameter
	switch {
	case peekedTagNumber == uint8(0): // BACnetEventParameterChangeOfBitstring
		if _child, err = new(_BACnetEventParameterChangeOfBitstring).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfBitstring for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(1): // BACnetEventParameterChangeOfState
		if _child, err = new(_BACnetEventParameterChangeOfState).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfState for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(2): // BACnetEventParameterChangeOfValue
		if _child, err = new(_BACnetEventParameterChangeOfValue).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfValue for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(3): // BACnetEventParameterCommandFailure
		if _child, err = new(_BACnetEventParameterCommandFailure).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterCommandFailure for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(4): // BACnetEventParameterFloatingLimit
		if _child, err = new(_BACnetEventParameterFloatingLimit).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterFloatingLimit for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(5): // BACnetEventParameterOutOfRange
		if _child, err = new(_BACnetEventParameterOutOfRange).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterOutOfRange for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(8): // BACnetEventParameterChangeOfLifeSavety
		if _child, err = new(_BACnetEventParameterChangeOfLifeSavety).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfLifeSavety for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(9): // BACnetEventParameterExtended
		if _child, err = new(_BACnetEventParameterExtended).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterExtended for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(10): // BACnetEventParameterBufferReady
		if _child, err = new(_BACnetEventParameterBufferReady).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterBufferReady for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(11): // BACnetEventParameterUnsignedRange
		if _child, err = new(_BACnetEventParameterUnsignedRange).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterUnsignedRange for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(13): // BACnetEventParameterAccessEvent
		if _child, err = new(_BACnetEventParameterAccessEvent).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterAccessEvent for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(14): // BACnetEventParameterDoubleOutOfRange
		if _child, err = new(_BACnetEventParameterDoubleOutOfRange).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterDoubleOutOfRange for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(15): // BACnetEventParameterSignedOutOfRange
		if _child, err = new(_BACnetEventParameterSignedOutOfRange).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterSignedOutOfRange for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(16): // BACnetEventParameterUnsignedOutOfRange
		if _child, err = new(_BACnetEventParameterUnsignedOutOfRange).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterUnsignedOutOfRange for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(17): // BACnetEventParameterChangeOfCharacterString
		if _child, err = new(_BACnetEventParameterChangeOfCharacterString).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfCharacterString for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(18): // BACnetEventParameterChangeOfStatusFlags
		if _child, err = new(_BACnetEventParameterChangeOfStatusFlags).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfStatusFlags for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(20): // BACnetEventParameterNone
		if _child, err = new(_BACnetEventParameterNone).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterNone for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(21): // BACnetEventParameterChangeOfDiscreteValue
		if _child, err = new(_BACnetEventParameterChangeOfDiscreteValue).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfDiscreteValue for type-switch of BACnetEventParameter")
		}
	case peekedTagNumber == uint8(22): // BACnetEventParameterChangeOfTimer
		if _child, err = new(_BACnetEventParameterChangeOfTimer).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetEventParameterChangeOfTimer for type-switch of BACnetEventParameter")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [peekedTagNumber=%v]", peekedTagNumber)
	}

	if closeErr := readBuffer.CloseContext("BACnetEventParameter"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetEventParameter")
	}

	return _child, nil
}

func (pm *_BACnetEventParameter) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetEventParameter, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetEventParameter"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetEventParameter")
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

	if popErr := writeBuffer.PopContext("BACnetEventParameter"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetEventParameter")
	}
	return nil
}

func (m *_BACnetEventParameter) IsBACnetEventParameter() {}

func (m *_BACnetEventParameter) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetEventParameter) deepCopy() *_BACnetEventParameter {
	if m == nil {
		return nil
	}
	_BACnetEventParameterCopy := &_BACnetEventParameter{
		nil, // will be set by child
		m.PeekedTagHeader.DeepCopy().(BACnetTagHeader),
	}
	return _BACnetEventParameterCopy
}
