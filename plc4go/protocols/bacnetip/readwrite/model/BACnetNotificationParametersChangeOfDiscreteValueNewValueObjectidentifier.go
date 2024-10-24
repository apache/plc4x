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

// BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier is the corresponding interface of BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier
type BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetNotificationParametersChangeOfDiscreteValueNewValue
	// GetObjectidentifierValue returns ObjectidentifierValue (property field)
	GetObjectidentifierValue() BACnetApplicationTagObjectIdentifier
	// IsBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier()
	// CreateBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
	CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
}

// _BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier is the data-structure of this message
type _BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier struct {
	BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
	ObjectidentifierValue BACnetApplicationTagObjectIdentifier
}

var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier)(nil)
var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueRequirements = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier)(nil)

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier factory function for _BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, objectidentifierValue BACnetApplicationTagObjectIdentifier, tagNumber uint8) *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier {
	if objectidentifierValue == nil {
		panic("objectidentifierValue of type BACnetApplicationTagObjectIdentifier for BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier must not be nil")
	}
	_result := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier{
		BACnetNotificationParametersChangeOfDiscreteValueNewValueContract: NewBACnetNotificationParametersChangeOfDiscreteValueNewValue(openingTag, peekedTagHeader, closingTag, tagNumber),
		ObjectidentifierValue: objectidentifierValue,
	}
	_result.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder is a builder for BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier
type BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(objectidentifierValue BACnetApplicationTagObjectIdentifier) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
	// WithObjectidentifierValue adds ObjectidentifierValue (property field)
	WithObjectidentifierValue(BACnetApplicationTagObjectIdentifier) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
	// WithObjectidentifierValueBuilder adds ObjectidentifierValue (property field) which is build by the builder
	WithObjectidentifierValueBuilder(func(BACnetApplicationTagObjectIdentifierBuilder) BACnetApplicationTagObjectIdentifierBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
	// Build builds the BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier or returns an error if something is wrong
	Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier
}

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder() creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder {
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier: new(_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier)}
}

type _BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder struct {
	*_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier

	parentBuilder *_BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder

	err *utils.MultiError
}

var _ (BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder)(nil)

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) setParent(contract BACnetNotificationParametersChangeOfDiscreteValueNewValueContract) {
	b.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = contract
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) WithMandatoryFields(objectidentifierValue BACnetApplicationTagObjectIdentifier) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder {
	return b.WithObjectidentifierValue(objectidentifierValue)
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) WithObjectidentifierValue(objectidentifierValue BACnetApplicationTagObjectIdentifier) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder {
	b.ObjectidentifierValue = objectidentifierValue
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) WithObjectidentifierValueBuilder(builderSupplier func(BACnetApplicationTagObjectIdentifierBuilder) BACnetApplicationTagObjectIdentifierBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder {
	builder := builderSupplier(b.ObjectidentifierValue.CreateBACnetApplicationTagObjectIdentifierBuilder())
	var err error
	b.ObjectidentifierValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagObjectIdentifierBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier, error) {
	if b.ObjectidentifierValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'objectidentifierValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier.deepCopy(), nil
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) Done() BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) buildForBACnetNotificationParametersChangeOfDiscreteValueNewValue() (BACnetNotificationParametersChangeOfDiscreteValueNewValue, error) {
	return b.Build()
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder) DeepCopy() any {
	_copy := b.CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder().(*_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder {
	if b == nil {
		return NewBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder()
	}
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier: b.deepCopy()}
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

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) GetParent() BACnetNotificationParametersChangeOfDiscreteValueNewValueContract {
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) GetObjectidentifierValue() BACnetApplicationTagObjectIdentifier {
	return m.ObjectidentifierValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier(structType any) BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier {
	if casted, ok := structType.(BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) GetTypeName() string {
	return "BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier"
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).GetLengthInBits(ctx))

	// Simple field (objectidentifierValue)
	lengthInBits += m.ObjectidentifierValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetNotificationParametersChangeOfDiscreteValueNewValue, tagNumber uint8) (__bACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier, err error) {
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	objectidentifierValue, err := ReadSimpleField[BACnetApplicationTagObjectIdentifier](ctx, "objectidentifierValue", ReadComplex[BACnetApplicationTagObjectIdentifier](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagObjectIdentifier](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'objectidentifierValue' field"))
	}
	m.ObjectidentifierValue = objectidentifierValue

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier")
	}

	return m, nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier")
		}

		if err := WriteSimpleField[BACnetApplicationTagObjectIdentifier](ctx, "objectidentifierValue", m.GetObjectidentifierValue(), WriteComplex[BACnetApplicationTagObjectIdentifier](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'objectidentifierValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier")
		}
		return nil
	}
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) IsBACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier() {
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) deepCopy() *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier {
	if m == nil {
		return nil
	}
	_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierCopy := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier{
		m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).deepCopy(),
		m.ObjectidentifierValue.DeepCopy().(BACnetApplicationTagObjectIdentifier),
	}
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = m
	return _BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifierCopy
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier) String() string {
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
