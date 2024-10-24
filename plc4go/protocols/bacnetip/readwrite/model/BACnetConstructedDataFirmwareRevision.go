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

// BACnetConstructedDataFirmwareRevision is the corresponding interface of BACnetConstructedDataFirmwareRevision
type BACnetConstructedDataFirmwareRevision interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetFirmwareRevision returns FirmwareRevision (property field)
	GetFirmwareRevision() BACnetApplicationTagCharacterString
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagCharacterString
	// IsBACnetConstructedDataFirmwareRevision is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataFirmwareRevision()
	// CreateBuilder creates a BACnetConstructedDataFirmwareRevisionBuilder
	CreateBACnetConstructedDataFirmwareRevisionBuilder() BACnetConstructedDataFirmwareRevisionBuilder
}

// _BACnetConstructedDataFirmwareRevision is the data-structure of this message
type _BACnetConstructedDataFirmwareRevision struct {
	BACnetConstructedDataContract
	FirmwareRevision BACnetApplicationTagCharacterString
}

var _ BACnetConstructedDataFirmwareRevision = (*_BACnetConstructedDataFirmwareRevision)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataFirmwareRevision)(nil)

// NewBACnetConstructedDataFirmwareRevision factory function for _BACnetConstructedDataFirmwareRevision
func NewBACnetConstructedDataFirmwareRevision(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, firmwareRevision BACnetApplicationTagCharacterString, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataFirmwareRevision {
	if firmwareRevision == nil {
		panic("firmwareRevision of type BACnetApplicationTagCharacterString for BACnetConstructedDataFirmwareRevision must not be nil")
	}
	_result := &_BACnetConstructedDataFirmwareRevision{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		FirmwareRevision:              firmwareRevision,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataFirmwareRevisionBuilder is a builder for BACnetConstructedDataFirmwareRevision
type BACnetConstructedDataFirmwareRevisionBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(firmwareRevision BACnetApplicationTagCharacterString) BACnetConstructedDataFirmwareRevisionBuilder
	// WithFirmwareRevision adds FirmwareRevision (property field)
	WithFirmwareRevision(BACnetApplicationTagCharacterString) BACnetConstructedDataFirmwareRevisionBuilder
	// WithFirmwareRevisionBuilder adds FirmwareRevision (property field) which is build by the builder
	WithFirmwareRevisionBuilder(func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetConstructedDataFirmwareRevisionBuilder
	// Build builds the BACnetConstructedDataFirmwareRevision or returns an error if something is wrong
	Build() (BACnetConstructedDataFirmwareRevision, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataFirmwareRevision
}

// NewBACnetConstructedDataFirmwareRevisionBuilder() creates a BACnetConstructedDataFirmwareRevisionBuilder
func NewBACnetConstructedDataFirmwareRevisionBuilder() BACnetConstructedDataFirmwareRevisionBuilder {
	return &_BACnetConstructedDataFirmwareRevisionBuilder{_BACnetConstructedDataFirmwareRevision: new(_BACnetConstructedDataFirmwareRevision)}
}

type _BACnetConstructedDataFirmwareRevisionBuilder struct {
	*_BACnetConstructedDataFirmwareRevision

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataFirmwareRevisionBuilder) = (*_BACnetConstructedDataFirmwareRevisionBuilder)(nil)

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) WithMandatoryFields(firmwareRevision BACnetApplicationTagCharacterString) BACnetConstructedDataFirmwareRevisionBuilder {
	return b.WithFirmwareRevision(firmwareRevision)
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) WithFirmwareRevision(firmwareRevision BACnetApplicationTagCharacterString) BACnetConstructedDataFirmwareRevisionBuilder {
	b.FirmwareRevision = firmwareRevision
	return b
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) WithFirmwareRevisionBuilder(builderSupplier func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetConstructedDataFirmwareRevisionBuilder {
	builder := builderSupplier(b.FirmwareRevision.CreateBACnetApplicationTagCharacterStringBuilder())
	var err error
	b.FirmwareRevision, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagCharacterStringBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) Build() (BACnetConstructedDataFirmwareRevision, error) {
	if b.FirmwareRevision == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'firmwareRevision' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataFirmwareRevision.deepCopy(), nil
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) MustBuild() BACnetConstructedDataFirmwareRevision {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataFirmwareRevisionBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataFirmwareRevisionBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataFirmwareRevisionBuilder().(*_BACnetConstructedDataFirmwareRevisionBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataFirmwareRevisionBuilder creates a BACnetConstructedDataFirmwareRevisionBuilder
func (b *_BACnetConstructedDataFirmwareRevision) CreateBACnetConstructedDataFirmwareRevisionBuilder() BACnetConstructedDataFirmwareRevisionBuilder {
	if b == nil {
		return NewBACnetConstructedDataFirmwareRevisionBuilder()
	}
	return &_BACnetConstructedDataFirmwareRevisionBuilder{_BACnetConstructedDataFirmwareRevision: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataFirmwareRevision) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataFirmwareRevision) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_FIRMWARE_REVISION
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataFirmwareRevision) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataFirmwareRevision) GetFirmwareRevision() BACnetApplicationTagCharacterString {
	return m.FirmwareRevision
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataFirmwareRevision) GetActualValue() BACnetApplicationTagCharacterString {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagCharacterString(m.GetFirmwareRevision())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataFirmwareRevision(structType any) BACnetConstructedDataFirmwareRevision {
	if casted, ok := structType.(BACnetConstructedDataFirmwareRevision); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataFirmwareRevision); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataFirmwareRevision) GetTypeName() string {
	return "BACnetConstructedDataFirmwareRevision"
}

func (m *_BACnetConstructedDataFirmwareRevision) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (firmwareRevision)
	lengthInBits += m.FirmwareRevision.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataFirmwareRevision) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataFirmwareRevision) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataFirmwareRevision BACnetConstructedDataFirmwareRevision, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataFirmwareRevision"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataFirmwareRevision")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	firmwareRevision, err := ReadSimpleField[BACnetApplicationTagCharacterString](ctx, "firmwareRevision", ReadComplex[BACnetApplicationTagCharacterString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagCharacterString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'firmwareRevision' field"))
	}
	m.FirmwareRevision = firmwareRevision

	actualValue, err := ReadVirtualField[BACnetApplicationTagCharacterString](ctx, "actualValue", (*BACnetApplicationTagCharacterString)(nil), firmwareRevision)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataFirmwareRevision"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataFirmwareRevision")
	}

	return m, nil
}

func (m *_BACnetConstructedDataFirmwareRevision) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataFirmwareRevision) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataFirmwareRevision"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataFirmwareRevision")
		}

		if err := WriteSimpleField[BACnetApplicationTagCharacterString](ctx, "firmwareRevision", m.GetFirmwareRevision(), WriteComplex[BACnetApplicationTagCharacterString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'firmwareRevision' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataFirmwareRevision"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataFirmwareRevision")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataFirmwareRevision) IsBACnetConstructedDataFirmwareRevision() {}

func (m *_BACnetConstructedDataFirmwareRevision) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataFirmwareRevision) deepCopy() *_BACnetConstructedDataFirmwareRevision {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataFirmwareRevisionCopy := &_BACnetConstructedDataFirmwareRevision{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.FirmwareRevision.DeepCopy().(BACnetApplicationTagCharacterString),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataFirmwareRevisionCopy
}

func (m *_BACnetConstructedDataFirmwareRevision) String() string {
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
