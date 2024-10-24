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

// BACnetConstructedDataProcessIdentifier is the corresponding interface of BACnetConstructedDataProcessIdentifier
type BACnetConstructedDataProcessIdentifier interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetProcessIdentifier returns ProcessIdentifier (property field)
	GetProcessIdentifier() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// IsBACnetConstructedDataProcessIdentifier is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataProcessIdentifier()
	// CreateBuilder creates a BACnetConstructedDataProcessIdentifierBuilder
	CreateBACnetConstructedDataProcessIdentifierBuilder() BACnetConstructedDataProcessIdentifierBuilder
}

// _BACnetConstructedDataProcessIdentifier is the data-structure of this message
type _BACnetConstructedDataProcessIdentifier struct {
	BACnetConstructedDataContract
	ProcessIdentifier BACnetApplicationTagUnsignedInteger
}

var _ BACnetConstructedDataProcessIdentifier = (*_BACnetConstructedDataProcessIdentifier)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataProcessIdentifier)(nil)

// NewBACnetConstructedDataProcessIdentifier factory function for _BACnetConstructedDataProcessIdentifier
func NewBACnetConstructedDataProcessIdentifier(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, processIdentifier BACnetApplicationTagUnsignedInteger, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataProcessIdentifier {
	if processIdentifier == nil {
		panic("processIdentifier of type BACnetApplicationTagUnsignedInteger for BACnetConstructedDataProcessIdentifier must not be nil")
	}
	_result := &_BACnetConstructedDataProcessIdentifier{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		ProcessIdentifier:             processIdentifier,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataProcessIdentifierBuilder is a builder for BACnetConstructedDataProcessIdentifier
type BACnetConstructedDataProcessIdentifierBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(processIdentifier BACnetApplicationTagUnsignedInteger) BACnetConstructedDataProcessIdentifierBuilder
	// WithProcessIdentifier adds ProcessIdentifier (property field)
	WithProcessIdentifier(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataProcessIdentifierBuilder
	// WithProcessIdentifierBuilder adds ProcessIdentifier (property field) which is build by the builder
	WithProcessIdentifierBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataProcessIdentifierBuilder
	// Build builds the BACnetConstructedDataProcessIdentifier or returns an error if something is wrong
	Build() (BACnetConstructedDataProcessIdentifier, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataProcessIdentifier
}

// NewBACnetConstructedDataProcessIdentifierBuilder() creates a BACnetConstructedDataProcessIdentifierBuilder
func NewBACnetConstructedDataProcessIdentifierBuilder() BACnetConstructedDataProcessIdentifierBuilder {
	return &_BACnetConstructedDataProcessIdentifierBuilder{_BACnetConstructedDataProcessIdentifier: new(_BACnetConstructedDataProcessIdentifier)}
}

type _BACnetConstructedDataProcessIdentifierBuilder struct {
	*_BACnetConstructedDataProcessIdentifier

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataProcessIdentifierBuilder) = (*_BACnetConstructedDataProcessIdentifierBuilder)(nil)

func (b *_BACnetConstructedDataProcessIdentifierBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) WithMandatoryFields(processIdentifier BACnetApplicationTagUnsignedInteger) BACnetConstructedDataProcessIdentifierBuilder {
	return b.WithProcessIdentifier(processIdentifier)
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) WithProcessIdentifier(processIdentifier BACnetApplicationTagUnsignedInteger) BACnetConstructedDataProcessIdentifierBuilder {
	b.ProcessIdentifier = processIdentifier
	return b
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) WithProcessIdentifierBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataProcessIdentifierBuilder {
	builder := builderSupplier(b.ProcessIdentifier.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.ProcessIdentifier, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) Build() (BACnetConstructedDataProcessIdentifier, error) {
	if b.ProcessIdentifier == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'processIdentifier' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataProcessIdentifier.deepCopy(), nil
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) MustBuild() BACnetConstructedDataProcessIdentifier {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataProcessIdentifierBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataProcessIdentifierBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataProcessIdentifierBuilder().(*_BACnetConstructedDataProcessIdentifierBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataProcessIdentifierBuilder creates a BACnetConstructedDataProcessIdentifierBuilder
func (b *_BACnetConstructedDataProcessIdentifier) CreateBACnetConstructedDataProcessIdentifierBuilder() BACnetConstructedDataProcessIdentifierBuilder {
	if b == nil {
		return NewBACnetConstructedDataProcessIdentifierBuilder()
	}
	return &_BACnetConstructedDataProcessIdentifierBuilder{_BACnetConstructedDataProcessIdentifier: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataProcessIdentifier) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataProcessIdentifier) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_PROCESS_IDENTIFIER
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataProcessIdentifier) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataProcessIdentifier) GetProcessIdentifier() BACnetApplicationTagUnsignedInteger {
	return m.ProcessIdentifier
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataProcessIdentifier) GetActualValue() BACnetApplicationTagUnsignedInteger {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagUnsignedInteger(m.GetProcessIdentifier())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataProcessIdentifier(structType any) BACnetConstructedDataProcessIdentifier {
	if casted, ok := structType.(BACnetConstructedDataProcessIdentifier); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataProcessIdentifier); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataProcessIdentifier) GetTypeName() string {
	return "BACnetConstructedDataProcessIdentifier"
}

func (m *_BACnetConstructedDataProcessIdentifier) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (processIdentifier)
	lengthInBits += m.ProcessIdentifier.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataProcessIdentifier) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataProcessIdentifier) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataProcessIdentifier BACnetConstructedDataProcessIdentifier, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataProcessIdentifier"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataProcessIdentifier")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	processIdentifier, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "processIdentifier", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'processIdentifier' field"))
	}
	m.ProcessIdentifier = processIdentifier

	actualValue, err := ReadVirtualField[BACnetApplicationTagUnsignedInteger](ctx, "actualValue", (*BACnetApplicationTagUnsignedInteger)(nil), processIdentifier)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataProcessIdentifier"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataProcessIdentifier")
	}

	return m, nil
}

func (m *_BACnetConstructedDataProcessIdentifier) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataProcessIdentifier) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataProcessIdentifier"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataProcessIdentifier")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "processIdentifier", m.GetProcessIdentifier(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'processIdentifier' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataProcessIdentifier"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataProcessIdentifier")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataProcessIdentifier) IsBACnetConstructedDataProcessIdentifier() {}

func (m *_BACnetConstructedDataProcessIdentifier) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataProcessIdentifier) deepCopy() *_BACnetConstructedDataProcessIdentifier {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataProcessIdentifierCopy := &_BACnetConstructedDataProcessIdentifier{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.ProcessIdentifier.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataProcessIdentifierCopy
}

func (m *_BACnetConstructedDataProcessIdentifier) String() string {
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
