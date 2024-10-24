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

// BACnetTimerStateChangeValueOctetString is the corresponding interface of BACnetTimerStateChangeValueOctetString
type BACnetTimerStateChangeValueOctetString interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetTimerStateChangeValue
	// GetOctetStringValue returns OctetStringValue (property field)
	GetOctetStringValue() BACnetApplicationTagOctetString
	// IsBACnetTimerStateChangeValueOctetString is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetTimerStateChangeValueOctetString()
	// CreateBuilder creates a BACnetTimerStateChangeValueOctetStringBuilder
	CreateBACnetTimerStateChangeValueOctetStringBuilder() BACnetTimerStateChangeValueOctetStringBuilder
}

// _BACnetTimerStateChangeValueOctetString is the data-structure of this message
type _BACnetTimerStateChangeValueOctetString struct {
	BACnetTimerStateChangeValueContract
	OctetStringValue BACnetApplicationTagOctetString
}

var _ BACnetTimerStateChangeValueOctetString = (*_BACnetTimerStateChangeValueOctetString)(nil)
var _ BACnetTimerStateChangeValueRequirements = (*_BACnetTimerStateChangeValueOctetString)(nil)

// NewBACnetTimerStateChangeValueOctetString factory function for _BACnetTimerStateChangeValueOctetString
func NewBACnetTimerStateChangeValueOctetString(peekedTagHeader BACnetTagHeader, octetStringValue BACnetApplicationTagOctetString, objectTypeArgument BACnetObjectType) *_BACnetTimerStateChangeValueOctetString {
	if octetStringValue == nil {
		panic("octetStringValue of type BACnetApplicationTagOctetString for BACnetTimerStateChangeValueOctetString must not be nil")
	}
	_result := &_BACnetTimerStateChangeValueOctetString{
		BACnetTimerStateChangeValueContract: NewBACnetTimerStateChangeValue(peekedTagHeader, objectTypeArgument),
		OctetStringValue:                    octetStringValue,
	}
	_result.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetTimerStateChangeValueOctetStringBuilder is a builder for BACnetTimerStateChangeValueOctetString
type BACnetTimerStateChangeValueOctetStringBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(octetStringValue BACnetApplicationTagOctetString) BACnetTimerStateChangeValueOctetStringBuilder
	// WithOctetStringValue adds OctetStringValue (property field)
	WithOctetStringValue(BACnetApplicationTagOctetString) BACnetTimerStateChangeValueOctetStringBuilder
	// WithOctetStringValueBuilder adds OctetStringValue (property field) which is build by the builder
	WithOctetStringValueBuilder(func(BACnetApplicationTagOctetStringBuilder) BACnetApplicationTagOctetStringBuilder) BACnetTimerStateChangeValueOctetStringBuilder
	// Build builds the BACnetTimerStateChangeValueOctetString or returns an error if something is wrong
	Build() (BACnetTimerStateChangeValueOctetString, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetTimerStateChangeValueOctetString
}

// NewBACnetTimerStateChangeValueOctetStringBuilder() creates a BACnetTimerStateChangeValueOctetStringBuilder
func NewBACnetTimerStateChangeValueOctetStringBuilder() BACnetTimerStateChangeValueOctetStringBuilder {
	return &_BACnetTimerStateChangeValueOctetStringBuilder{_BACnetTimerStateChangeValueOctetString: new(_BACnetTimerStateChangeValueOctetString)}
}

type _BACnetTimerStateChangeValueOctetStringBuilder struct {
	*_BACnetTimerStateChangeValueOctetString

	parentBuilder *_BACnetTimerStateChangeValueBuilder

	err *utils.MultiError
}

var _ (BACnetTimerStateChangeValueOctetStringBuilder) = (*_BACnetTimerStateChangeValueOctetStringBuilder)(nil)

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) setParent(contract BACnetTimerStateChangeValueContract) {
	b.BACnetTimerStateChangeValueContract = contract
}

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) WithMandatoryFields(octetStringValue BACnetApplicationTagOctetString) BACnetTimerStateChangeValueOctetStringBuilder {
	return b.WithOctetStringValue(octetStringValue)
}

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) WithOctetStringValue(octetStringValue BACnetApplicationTagOctetString) BACnetTimerStateChangeValueOctetStringBuilder {
	b.OctetStringValue = octetStringValue
	return b
}

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) WithOctetStringValueBuilder(builderSupplier func(BACnetApplicationTagOctetStringBuilder) BACnetApplicationTagOctetStringBuilder) BACnetTimerStateChangeValueOctetStringBuilder {
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

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) Build() (BACnetTimerStateChangeValueOctetString, error) {
	if b.OctetStringValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'octetStringValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetTimerStateChangeValueOctetString.deepCopy(), nil
}

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) MustBuild() BACnetTimerStateChangeValueOctetString {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetTimerStateChangeValueOctetStringBuilder) Done() BACnetTimerStateChangeValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) buildForBACnetTimerStateChangeValue() (BACnetTimerStateChangeValue, error) {
	return b.Build()
}

func (b *_BACnetTimerStateChangeValueOctetStringBuilder) DeepCopy() any {
	_copy := b.CreateBACnetTimerStateChangeValueOctetStringBuilder().(*_BACnetTimerStateChangeValueOctetStringBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetTimerStateChangeValueOctetStringBuilder creates a BACnetTimerStateChangeValueOctetStringBuilder
func (b *_BACnetTimerStateChangeValueOctetString) CreateBACnetTimerStateChangeValueOctetStringBuilder() BACnetTimerStateChangeValueOctetStringBuilder {
	if b == nil {
		return NewBACnetTimerStateChangeValueOctetStringBuilder()
	}
	return &_BACnetTimerStateChangeValueOctetStringBuilder{_BACnetTimerStateChangeValueOctetString: b.deepCopy()}
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

func (m *_BACnetTimerStateChangeValueOctetString) GetParent() BACnetTimerStateChangeValueContract {
	return m.BACnetTimerStateChangeValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetTimerStateChangeValueOctetString) GetOctetStringValue() BACnetApplicationTagOctetString {
	return m.OctetStringValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetTimerStateChangeValueOctetString(structType any) BACnetTimerStateChangeValueOctetString {
	if casted, ok := structType.(BACnetTimerStateChangeValueOctetString); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetTimerStateChangeValueOctetString); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetTimerStateChangeValueOctetString) GetTypeName() string {
	return "BACnetTimerStateChangeValueOctetString"
}

func (m *_BACnetTimerStateChangeValueOctetString) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue).GetLengthInBits(ctx))

	// Simple field (octetStringValue)
	lengthInBits += m.OctetStringValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetTimerStateChangeValueOctetString) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetTimerStateChangeValueOctetString) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetTimerStateChangeValue, objectTypeArgument BACnetObjectType) (__bACnetTimerStateChangeValueOctetString BACnetTimerStateChangeValueOctetString, err error) {
	m.BACnetTimerStateChangeValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetTimerStateChangeValueOctetString"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetTimerStateChangeValueOctetString")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	octetStringValue, err := ReadSimpleField[BACnetApplicationTagOctetString](ctx, "octetStringValue", ReadComplex[BACnetApplicationTagOctetString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagOctetString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'octetStringValue' field"))
	}
	m.OctetStringValue = octetStringValue

	if closeErr := readBuffer.CloseContext("BACnetTimerStateChangeValueOctetString"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetTimerStateChangeValueOctetString")
	}

	return m, nil
}

func (m *_BACnetTimerStateChangeValueOctetString) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetTimerStateChangeValueOctetString) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetTimerStateChangeValueOctetString"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetTimerStateChangeValueOctetString")
		}

		if err := WriteSimpleField[BACnetApplicationTagOctetString](ctx, "octetStringValue", m.GetOctetStringValue(), WriteComplex[BACnetApplicationTagOctetString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'octetStringValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetTimerStateChangeValueOctetString"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetTimerStateChangeValueOctetString")
		}
		return nil
	}
	return m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetTimerStateChangeValueOctetString) IsBACnetTimerStateChangeValueOctetString() {}

func (m *_BACnetTimerStateChangeValueOctetString) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetTimerStateChangeValueOctetString) deepCopy() *_BACnetTimerStateChangeValueOctetString {
	if m == nil {
		return nil
	}
	_BACnetTimerStateChangeValueOctetStringCopy := &_BACnetTimerStateChangeValueOctetString{
		m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue).deepCopy(),
		m.OctetStringValue.DeepCopy().(BACnetApplicationTagOctetString),
	}
	m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue)._SubType = m
	return _BACnetTimerStateChangeValueOctetStringCopy
}

func (m *_BACnetTimerStateChangeValueOctetString) String() string {
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
