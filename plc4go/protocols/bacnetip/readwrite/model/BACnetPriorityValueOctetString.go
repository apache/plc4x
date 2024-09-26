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

// BACnetPriorityValueOctetString is the corresponding interface of BACnetPriorityValueOctetString
type BACnetPriorityValueOctetString interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetPriorityValue
	// GetOctetStringValue returns OctetStringValue (property field)
	GetOctetStringValue() BACnetApplicationTagOctetString
	// IsBACnetPriorityValueOctetString is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetPriorityValueOctetString()
	// CreateBuilder creates a BACnetPriorityValueOctetStringBuilder
	CreateBACnetPriorityValueOctetStringBuilder() BACnetPriorityValueOctetStringBuilder
}

// _BACnetPriorityValueOctetString is the data-structure of this message
type _BACnetPriorityValueOctetString struct {
	BACnetPriorityValueContract
	OctetStringValue BACnetApplicationTagOctetString
}

var _ BACnetPriorityValueOctetString = (*_BACnetPriorityValueOctetString)(nil)
var _ BACnetPriorityValueRequirements = (*_BACnetPriorityValueOctetString)(nil)

// NewBACnetPriorityValueOctetString factory function for _BACnetPriorityValueOctetString
func NewBACnetPriorityValueOctetString(peekedTagHeader BACnetTagHeader, octetStringValue BACnetApplicationTagOctetString, objectTypeArgument BACnetObjectType) *_BACnetPriorityValueOctetString {
	if octetStringValue == nil {
		panic("octetStringValue of type BACnetApplicationTagOctetString for BACnetPriorityValueOctetString must not be nil")
	}
	_result := &_BACnetPriorityValueOctetString{
		BACnetPriorityValueContract: NewBACnetPriorityValue(peekedTagHeader, objectTypeArgument),
		OctetStringValue:            octetStringValue,
	}
	_result.BACnetPriorityValueContract.(*_BACnetPriorityValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetPriorityValueOctetStringBuilder is a builder for BACnetPriorityValueOctetString
type BACnetPriorityValueOctetStringBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(octetStringValue BACnetApplicationTagOctetString) BACnetPriorityValueOctetStringBuilder
	// WithOctetStringValue adds OctetStringValue (property field)
	WithOctetStringValue(BACnetApplicationTagOctetString) BACnetPriorityValueOctetStringBuilder
	// WithOctetStringValueBuilder adds OctetStringValue (property field) which is build by the builder
	WithOctetStringValueBuilder(func(BACnetApplicationTagOctetStringBuilder) BACnetApplicationTagOctetStringBuilder) BACnetPriorityValueOctetStringBuilder
	// Build builds the BACnetPriorityValueOctetString or returns an error if something is wrong
	Build() (BACnetPriorityValueOctetString, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetPriorityValueOctetString
}

// NewBACnetPriorityValueOctetStringBuilder() creates a BACnetPriorityValueOctetStringBuilder
func NewBACnetPriorityValueOctetStringBuilder() BACnetPriorityValueOctetStringBuilder {
	return &_BACnetPriorityValueOctetStringBuilder{_BACnetPriorityValueOctetString: new(_BACnetPriorityValueOctetString)}
}

type _BACnetPriorityValueOctetStringBuilder struct {
	*_BACnetPriorityValueOctetString

	parentBuilder *_BACnetPriorityValueBuilder

	err *utils.MultiError
}

var _ (BACnetPriorityValueOctetStringBuilder) = (*_BACnetPriorityValueOctetStringBuilder)(nil)

func (b *_BACnetPriorityValueOctetStringBuilder) setParent(contract BACnetPriorityValueContract) {
	b.BACnetPriorityValueContract = contract
}

func (b *_BACnetPriorityValueOctetStringBuilder) WithMandatoryFields(octetStringValue BACnetApplicationTagOctetString) BACnetPriorityValueOctetStringBuilder {
	return b.WithOctetStringValue(octetStringValue)
}

func (b *_BACnetPriorityValueOctetStringBuilder) WithOctetStringValue(octetStringValue BACnetApplicationTagOctetString) BACnetPriorityValueOctetStringBuilder {
	b.OctetStringValue = octetStringValue
	return b
}

func (b *_BACnetPriorityValueOctetStringBuilder) WithOctetStringValueBuilder(builderSupplier func(BACnetApplicationTagOctetStringBuilder) BACnetApplicationTagOctetStringBuilder) BACnetPriorityValueOctetStringBuilder {
	builder := builderSupplier(b.OctetStringValue.CreateBACnetApplicationTagOctetStringBuilder())
	var err error
	b.OctetStringValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagOctetStringBuilder failed"))
	}
	return b
}

func (b *_BACnetPriorityValueOctetStringBuilder) Build() (BACnetPriorityValueOctetString, error) {
	if b.OctetStringValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'octetStringValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetPriorityValueOctetString.deepCopy(), nil
}

func (b *_BACnetPriorityValueOctetStringBuilder) MustBuild() BACnetPriorityValueOctetString {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetPriorityValueOctetStringBuilder) Done() BACnetPriorityValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetPriorityValueOctetStringBuilder) buildForBACnetPriorityValue() (BACnetPriorityValue, error) {
	return b.Build()
}

func (b *_BACnetPriorityValueOctetStringBuilder) DeepCopy() any {
	_copy := b.CreateBACnetPriorityValueOctetStringBuilder().(*_BACnetPriorityValueOctetStringBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetPriorityValueOctetStringBuilder creates a BACnetPriorityValueOctetStringBuilder
func (b *_BACnetPriorityValueOctetString) CreateBACnetPriorityValueOctetStringBuilder() BACnetPriorityValueOctetStringBuilder {
	if b == nil {
		return NewBACnetPriorityValueOctetStringBuilder()
	}
	return &_BACnetPriorityValueOctetStringBuilder{_BACnetPriorityValueOctetString: b.deepCopy()}
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

func (m *_BACnetPriorityValueOctetString) GetParent() BACnetPriorityValueContract {
	return m.BACnetPriorityValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPriorityValueOctetString) GetOctetStringValue() BACnetApplicationTagOctetString {
	return m.OctetStringValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetPriorityValueOctetString(structType any) BACnetPriorityValueOctetString {
	if casted, ok := structType.(BACnetPriorityValueOctetString); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPriorityValueOctetString); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPriorityValueOctetString) GetTypeName() string {
	return "BACnetPriorityValueOctetString"
}

func (m *_BACnetPriorityValueOctetString) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetPriorityValueContract.(*_BACnetPriorityValue).GetLengthInBits(ctx))

	// Simple field (octetStringValue)
	lengthInBits += m.OctetStringValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPriorityValueOctetString) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetPriorityValueOctetString) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetPriorityValue, objectTypeArgument BACnetObjectType) (__bACnetPriorityValueOctetString BACnetPriorityValueOctetString, err error) {
	m.BACnetPriorityValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPriorityValueOctetString"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPriorityValueOctetString")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	octetStringValue, err := ReadSimpleField[BACnetApplicationTagOctetString](ctx, "octetStringValue", ReadComplex[BACnetApplicationTagOctetString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagOctetString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'octetStringValue' field"))
	}
	m.OctetStringValue = octetStringValue

	if closeErr := readBuffer.CloseContext("BACnetPriorityValueOctetString"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPriorityValueOctetString")
	}

	return m, nil
}

func (m *_BACnetPriorityValueOctetString) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPriorityValueOctetString) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPriorityValueOctetString"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPriorityValueOctetString")
		}

		if err := WriteSimpleField[BACnetApplicationTagOctetString](ctx, "octetStringValue", m.GetOctetStringValue(), WriteComplex[BACnetApplicationTagOctetString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'octetStringValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPriorityValueOctetString"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPriorityValueOctetString")
		}
		return nil
	}
	return m.BACnetPriorityValueContract.(*_BACnetPriorityValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPriorityValueOctetString) IsBACnetPriorityValueOctetString() {}

func (m *_BACnetPriorityValueOctetString) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetPriorityValueOctetString) deepCopy() *_BACnetPriorityValueOctetString {
	if m == nil {
		return nil
	}
	_BACnetPriorityValueOctetStringCopy := &_BACnetPriorityValueOctetString{
		m.BACnetPriorityValueContract.(*_BACnetPriorityValue).deepCopy(),
		m.OctetStringValue.DeepCopy().(BACnetApplicationTagOctetString),
	}
	m.BACnetPriorityValueContract.(*_BACnetPriorityValue)._SubType = m
	return _BACnetPriorityValueOctetStringCopy
}

func (m *_BACnetPriorityValueOctetString) String() string {
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
