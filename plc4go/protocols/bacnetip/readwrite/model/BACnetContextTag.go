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

// BACnetContextTag is the corresponding interface of BACnetContextTag
type BACnetContextTag interface {
	BACnetContextTagContract
	BACnetContextTagRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// IsBACnetContextTag is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetContextTag()
	// CreateBuilder creates a BACnetContextTagBuilder
	CreateBACnetContextTagBuilder() BACnetContextTagBuilder
}

// BACnetContextTagContract provides a set of functions which can be overwritten by a sub struct
type BACnetContextTagContract interface {
	// GetHeader returns Header (property field)
	GetHeader() BACnetTagHeader
	// GetTagNumber returns TagNumber (virtual field)
	GetTagNumber() uint8
	// GetActualLength returns ActualLength (virtual field)
	GetActualLength() uint32
	// GetTagNumberArgument() returns a parser argument
	GetTagNumberArgument() uint8
	// IsBACnetContextTag is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetContextTag()
	// CreateBuilder creates a BACnetContextTagBuilder
	CreateBACnetContextTagBuilder() BACnetContextTagBuilder
}

// BACnetContextTagRequirements provides a set of functions which need to be implemented by a sub struct
type BACnetContextTagRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetDataType returns DataType (discriminator field)
	GetDataType() BACnetDataType
}

// _BACnetContextTag is the data-structure of this message
type _BACnetContextTag struct {
	_SubType BACnetContextTag
	Header   BACnetTagHeader

	// Arguments.
	TagNumberArgument uint8
}

var _ BACnetContextTagContract = (*_BACnetContextTag)(nil)

// NewBACnetContextTag factory function for _BACnetContextTag
func NewBACnetContextTag(header BACnetTagHeader, tagNumberArgument uint8) *_BACnetContextTag {
	if header == nil {
		panic("header of type BACnetTagHeader for BACnetContextTag must not be nil")
	}
	return &_BACnetContextTag{Header: header, TagNumberArgument: tagNumberArgument}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetContextTagBuilder is a builder for BACnetContextTag
type BACnetContextTagBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(header BACnetTagHeader) BACnetContextTagBuilder
	// WithHeader adds Header (property field)
	WithHeader(BACnetTagHeader) BACnetContextTagBuilder
	// WithHeaderBuilder adds Header (property field) which is build by the builder
	WithHeaderBuilder(func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetContextTagBuilder
	// AsBACnetContextTagNull converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagNull() interface {
		BACnetContextTagNullBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagBoolean converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagBoolean() interface {
		BACnetContextTagBooleanBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagUnsignedInteger converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagUnsignedInteger() interface {
		BACnetContextTagUnsignedIntegerBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagSignedInteger converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagSignedInteger() interface {
		BACnetContextTagSignedIntegerBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagReal converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagReal() interface {
		BACnetContextTagRealBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagDouble converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagDouble() interface {
		BACnetContextTagDoubleBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagOctetString converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagOctetString() interface {
		BACnetContextTagOctetStringBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagCharacterString converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagCharacterString() interface {
		BACnetContextTagCharacterStringBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagBitString converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagBitString() interface {
		BACnetContextTagBitStringBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagEnumerated converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagEnumerated() interface {
		BACnetContextTagEnumeratedBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagDate converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagDate() interface {
		BACnetContextTagDateBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagTime converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagTime() interface {
		BACnetContextTagTimeBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagObjectIdentifier converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagObjectIdentifier() interface {
		BACnetContextTagObjectIdentifierBuilder
		Done() BACnetContextTagBuilder
	}
	// AsBACnetContextTagUnknown converts this build to a subType of BACnetContextTag. It is always possible to return to current builder using Done()
	AsBACnetContextTagUnknown() interface {
		BACnetContextTagUnknownBuilder
		Done() BACnetContextTagBuilder
	}
	// Build builds the BACnetContextTag or returns an error if something is wrong
	PartialBuild() (BACnetContextTagContract, error)
	// MustBuild does the same as Build but panics on error
	PartialMustBuild() BACnetContextTagContract
	// Build builds the BACnetContextTag or returns an error if something is wrong
	Build() (BACnetContextTag, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetContextTag
}

// NewBACnetContextTagBuilder() creates a BACnetContextTagBuilder
func NewBACnetContextTagBuilder() BACnetContextTagBuilder {
	return &_BACnetContextTagBuilder{_BACnetContextTag: new(_BACnetContextTag)}
}

type _BACnetContextTagChildBuilder interface {
	utils.Copyable
	setParent(BACnetContextTagContract)
	buildForBACnetContextTag() (BACnetContextTag, error)
}

type _BACnetContextTagBuilder struct {
	*_BACnetContextTag

	childBuilder _BACnetContextTagChildBuilder

	err *utils.MultiError
}

var _ (BACnetContextTagBuilder) = (*_BACnetContextTagBuilder)(nil)

func (b *_BACnetContextTagBuilder) WithMandatoryFields(header BACnetTagHeader) BACnetContextTagBuilder {
	return b.WithHeader(header)
}

func (b *_BACnetContextTagBuilder) WithHeader(header BACnetTagHeader) BACnetContextTagBuilder {
	b.Header = header
	return b
}

func (b *_BACnetContextTagBuilder) WithHeaderBuilder(builderSupplier func(BACnetTagHeaderBuilder) BACnetTagHeaderBuilder) BACnetContextTagBuilder {
	builder := builderSupplier(b.Header.CreateBACnetTagHeaderBuilder())
	var err error
	b.Header, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetTagHeaderBuilder failed"))
	}
	return b
}

func (b *_BACnetContextTagBuilder) PartialBuild() (BACnetContextTagContract, error) {
	if b.Header == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'header' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetContextTag.deepCopy(), nil
}

func (b *_BACnetContextTagBuilder) PartialMustBuild() BACnetContextTagContract {
	build, err := b.PartialBuild()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagNull() interface {
	BACnetContextTagNullBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagNullBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagNullBuilder().(*_BACnetContextTagNullBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagBoolean() interface {
	BACnetContextTagBooleanBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagBooleanBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagBooleanBuilder().(*_BACnetContextTagBooleanBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagUnsignedInteger() interface {
	BACnetContextTagUnsignedIntegerBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagUnsignedIntegerBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagUnsignedIntegerBuilder().(*_BACnetContextTagUnsignedIntegerBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagSignedInteger() interface {
	BACnetContextTagSignedIntegerBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagSignedIntegerBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagSignedIntegerBuilder().(*_BACnetContextTagSignedIntegerBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagReal() interface {
	BACnetContextTagRealBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagRealBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagRealBuilder().(*_BACnetContextTagRealBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagDouble() interface {
	BACnetContextTagDoubleBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagDoubleBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagDoubleBuilder().(*_BACnetContextTagDoubleBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagOctetString() interface {
	BACnetContextTagOctetStringBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagOctetStringBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagOctetStringBuilder().(*_BACnetContextTagOctetStringBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagCharacterString() interface {
	BACnetContextTagCharacterStringBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagCharacterStringBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagCharacterStringBuilder().(*_BACnetContextTagCharacterStringBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagBitString() interface {
	BACnetContextTagBitStringBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagBitStringBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagBitStringBuilder().(*_BACnetContextTagBitStringBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagEnumerated() interface {
	BACnetContextTagEnumeratedBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagEnumeratedBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagEnumeratedBuilder().(*_BACnetContextTagEnumeratedBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagDate() interface {
	BACnetContextTagDateBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagDateBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagDateBuilder().(*_BACnetContextTagDateBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagTime() interface {
	BACnetContextTagTimeBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagTimeBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagTimeBuilder().(*_BACnetContextTagTimeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagObjectIdentifier() interface {
	BACnetContextTagObjectIdentifierBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagObjectIdentifierBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagObjectIdentifierBuilder().(*_BACnetContextTagObjectIdentifierBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) AsBACnetContextTagUnknown() interface {
	BACnetContextTagUnknownBuilder
	Done() BACnetContextTagBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		BACnetContextTagUnknownBuilder
		Done() BACnetContextTagBuilder
	}); ok {
		return cb
	}
	cb := NewBACnetContextTagUnknownBuilder().(*_BACnetContextTagUnknownBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_BACnetContextTagBuilder) Build() (BACnetContextTag, error) {
	v, err := b.PartialBuild()
	if err != nil {
		return nil, errors.Wrap(err, "error occurred during partial build")
	}
	if b.childBuilder == nil {
		return nil, errors.New("no child builder present")
	}
	b.childBuilder.setParent(v)
	return b.childBuilder.buildForBACnetContextTag()
}

func (b *_BACnetContextTagBuilder) MustBuild() BACnetContextTag {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_BACnetContextTagBuilder) DeepCopy() any {
	_copy := b.CreateBACnetContextTagBuilder().(*_BACnetContextTagBuilder)
	_copy.childBuilder = b.childBuilder.DeepCopy().(_BACnetContextTagChildBuilder)
	_copy.childBuilder.setParent(_copy)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetContextTagBuilder creates a BACnetContextTagBuilder
func (b *_BACnetContextTag) CreateBACnetContextTagBuilder() BACnetContextTagBuilder {
	if b == nil {
		return NewBACnetContextTagBuilder()
	}
	return &_BACnetContextTagBuilder{_BACnetContextTag: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetContextTag) GetHeader() BACnetTagHeader {
	return m.Header
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (pm *_BACnetContextTag) GetTagNumber() uint8 {
	m := pm._SubType
	ctx := context.Background()
	_ = ctx
	return uint8(m.GetHeader().GetTagNumber())
}

func (pm *_BACnetContextTag) GetActualLength() uint32 {
	m := pm._SubType
	ctx := context.Background()
	_ = ctx
	return uint32(m.GetHeader().GetActualLength())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetContextTag(structType any) BACnetContextTag {
	if casted, ok := structType.(BACnetContextTag); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetContextTag); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetContextTag) GetTypeName() string {
	return "BACnetContextTag"
}

func (m *_BACnetContextTag) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetContextTag) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func BACnetContextTagParse[T BACnetContextTag](ctx context.Context, theBytes []byte, tagNumberArgument uint8, dataType BACnetDataType) (T, error) {
	return BACnetContextTagParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes), tagNumberArgument, dataType)
}

func BACnetContextTagParseWithBufferProducer[T BACnetContextTag](tagNumberArgument uint8, dataType BACnetDataType) func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := BACnetContextTagParseWithBuffer[T](ctx, readBuffer, tagNumberArgument, dataType)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, nil
	}
}

func BACnetContextTagParseWithBuffer[T BACnetContextTag](ctx context.Context, readBuffer utils.ReadBuffer, tagNumberArgument uint8, dataType BACnetDataType) (T, error) {
	v, err := (&_BACnetContextTag{TagNumberArgument: tagNumberArgument}).parse(ctx, readBuffer, tagNumberArgument, dataType)
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

func (m *_BACnetContextTag) parse(ctx context.Context, readBuffer utils.ReadBuffer, tagNumberArgument uint8, dataType BACnetDataType) (__bACnetContextTag BACnetContextTag, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetContextTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetContextTag")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	header, err := ReadSimpleField[BACnetTagHeader](ctx, "header", ReadComplex[BACnetTagHeader](BACnetTagHeaderParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'header' field"))
	}
	m.Header = header

	// Validation
	if !(bool((header.GetActualTagNumber()) == (tagNumberArgument))) {
		return nil, errors.WithStack(utils.ParseAssertError{Message: "tagnumber doesn't match"})
	}

	// Validation
	if !(bool((header.GetTagClass()) == (TagClass_CONTEXT_SPECIFIC_TAGS))) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "should be a context tag"})
	}

	tagNumber, err := ReadVirtualField[uint8](ctx, "tagNumber", (*uint8)(nil), header.GetTagNumber())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'tagNumber' field"))
	}
	_ = tagNumber

	actualLength, err := ReadVirtualField[uint32](ctx, "actualLength", (*uint32)(nil), header.GetActualLength())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualLength' field"))
	}
	_ = actualLength

	// Validation
	if !(bool(bool((header.GetLengthValueType()) != (6))) && bool(bool((header.GetLengthValueType()) != (7)))) {
		return nil, errors.WithStack(utils.ParseAssertError{Message: "length 6 and 7 reserved for opening and closing tag"})
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _child BACnetContextTag
	switch {
	case dataType == BACnetDataType_NULL: // BACnetContextTagNull
		if _child, err = new(_BACnetContextTagNull).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagNull for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_BOOLEAN: // BACnetContextTagBoolean
		if _child, err = new(_BACnetContextTagBoolean).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagBoolean for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_UNSIGNED_INTEGER: // BACnetContextTagUnsignedInteger
		if _child, err = new(_BACnetContextTagUnsignedInteger).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagUnsignedInteger for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_SIGNED_INTEGER: // BACnetContextTagSignedInteger
		if _child, err = new(_BACnetContextTagSignedInteger).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagSignedInteger for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_REAL: // BACnetContextTagReal
		if _child, err = new(_BACnetContextTagReal).parse(ctx, readBuffer, m, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagReal for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_DOUBLE: // BACnetContextTagDouble
		if _child, err = new(_BACnetContextTagDouble).parse(ctx, readBuffer, m, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagDouble for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_OCTET_STRING: // BACnetContextTagOctetString
		if _child, err = new(_BACnetContextTagOctetString).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagOctetString for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_CHARACTER_STRING: // BACnetContextTagCharacterString
		if _child, err = new(_BACnetContextTagCharacterString).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagCharacterString for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_BIT_STRING: // BACnetContextTagBitString
		if _child, err = new(_BACnetContextTagBitString).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagBitString for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_ENUMERATED: // BACnetContextTagEnumerated
		if _child, err = new(_BACnetContextTagEnumerated).parse(ctx, readBuffer, m, header, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagEnumerated for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_DATE: // BACnetContextTagDate
		if _child, err = new(_BACnetContextTagDate).parse(ctx, readBuffer, m, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagDate for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_TIME: // BACnetContextTagTime
		if _child, err = new(_BACnetContextTagTime).parse(ctx, readBuffer, m, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagTime for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_BACNET_OBJECT_IDENTIFIER: // BACnetContextTagObjectIdentifier
		if _child, err = new(_BACnetContextTagObjectIdentifier).parse(ctx, readBuffer, m, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagObjectIdentifier for type-switch of BACnetContextTag")
		}
	case dataType == BACnetDataType_UNKNOWN: // BACnetContextTagUnknown
		if _child, err = new(_BACnetContextTagUnknown).parse(ctx, readBuffer, m, actualLength, tagNumberArgument, dataType); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type BACnetContextTagUnknown for type-switch of BACnetContextTag")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [dataType=%v]", dataType)
	}

	if closeErr := readBuffer.CloseContext("BACnetContextTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetContextTag")
	}

	return _child, nil
}

func (pm *_BACnetContextTag) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child BACnetContextTag, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetContextTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetContextTag")
	}

	if err := WriteSimpleField[BACnetTagHeader](ctx, "header", m.GetHeader(), WriteComplex[BACnetTagHeader](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'header' field")
	}
	// Virtual field
	tagNumber := m.GetTagNumber()
	_ = tagNumber
	if _tagNumberErr := writeBuffer.WriteVirtual(ctx, "tagNumber", m.GetTagNumber()); _tagNumberErr != nil {
		return errors.Wrap(_tagNumberErr, "Error serializing 'tagNumber' field")
	}
	// Virtual field
	actualLength := m.GetActualLength()
	_ = actualLength
	if _actualLengthErr := writeBuffer.WriteVirtual(ctx, "actualLength", m.GetActualLength()); _actualLengthErr != nil {
		return errors.Wrap(_actualLengthErr, "Error serializing 'actualLength' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("BACnetContextTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetContextTag")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetContextTag) GetTagNumberArgument() uint8 {
	return m.TagNumberArgument
}

//
////

func (m *_BACnetContextTag) IsBACnetContextTag() {}

func (m *_BACnetContextTag) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetContextTag) deepCopy() *_BACnetContextTag {
	if m == nil {
		return nil
	}
	_BACnetContextTagCopy := &_BACnetContextTag{
		nil, // will be set by child
		m.Header.DeepCopy().(BACnetTagHeader),
		m.TagNumberArgument,
	}
	return _BACnetContextTagCopy
}
