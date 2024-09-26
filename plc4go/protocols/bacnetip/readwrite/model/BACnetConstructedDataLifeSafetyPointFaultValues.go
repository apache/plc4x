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

// BACnetConstructedDataLifeSafetyPointFaultValues is the corresponding interface of BACnetConstructedDataLifeSafetyPointFaultValues
type BACnetConstructedDataLifeSafetyPointFaultValues interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetFaultValues returns FaultValues (property field)
	GetFaultValues() []BACnetLifeSafetyStateTagged
	// IsBACnetConstructedDataLifeSafetyPointFaultValues is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataLifeSafetyPointFaultValues()
	// CreateBuilder creates a BACnetConstructedDataLifeSafetyPointFaultValuesBuilder
	CreateBACnetConstructedDataLifeSafetyPointFaultValuesBuilder() BACnetConstructedDataLifeSafetyPointFaultValuesBuilder
}

// _BACnetConstructedDataLifeSafetyPointFaultValues is the data-structure of this message
type _BACnetConstructedDataLifeSafetyPointFaultValues struct {
	BACnetConstructedDataContract
	FaultValues []BACnetLifeSafetyStateTagged
}

var _ BACnetConstructedDataLifeSafetyPointFaultValues = (*_BACnetConstructedDataLifeSafetyPointFaultValues)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataLifeSafetyPointFaultValues)(nil)

// NewBACnetConstructedDataLifeSafetyPointFaultValues factory function for _BACnetConstructedDataLifeSafetyPointFaultValues
func NewBACnetConstructedDataLifeSafetyPointFaultValues(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, faultValues []BACnetLifeSafetyStateTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLifeSafetyPointFaultValues {
	_result := &_BACnetConstructedDataLifeSafetyPointFaultValues{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		FaultValues:                   faultValues,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataLifeSafetyPointFaultValuesBuilder is a builder for BACnetConstructedDataLifeSafetyPointFaultValues
type BACnetConstructedDataLifeSafetyPointFaultValuesBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(faultValues []BACnetLifeSafetyStateTagged) BACnetConstructedDataLifeSafetyPointFaultValuesBuilder
	// WithFaultValues adds FaultValues (property field)
	WithFaultValues(...BACnetLifeSafetyStateTagged) BACnetConstructedDataLifeSafetyPointFaultValuesBuilder
	// Build builds the BACnetConstructedDataLifeSafetyPointFaultValues or returns an error if something is wrong
	Build() (BACnetConstructedDataLifeSafetyPointFaultValues, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataLifeSafetyPointFaultValues
}

// NewBACnetConstructedDataLifeSafetyPointFaultValuesBuilder() creates a BACnetConstructedDataLifeSafetyPointFaultValuesBuilder
func NewBACnetConstructedDataLifeSafetyPointFaultValuesBuilder() BACnetConstructedDataLifeSafetyPointFaultValuesBuilder {
	return &_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder{_BACnetConstructedDataLifeSafetyPointFaultValues: new(_BACnetConstructedDataLifeSafetyPointFaultValues)}
}

type _BACnetConstructedDataLifeSafetyPointFaultValuesBuilder struct {
	*_BACnetConstructedDataLifeSafetyPointFaultValues

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) = (*_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder)(nil)

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) WithMandatoryFields(faultValues []BACnetLifeSafetyStateTagged) BACnetConstructedDataLifeSafetyPointFaultValuesBuilder {
	return b.WithFaultValues(faultValues...)
}

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) WithFaultValues(faultValues ...BACnetLifeSafetyStateTagged) BACnetConstructedDataLifeSafetyPointFaultValuesBuilder {
	b.FaultValues = faultValues
	return b
}

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) Build() (BACnetConstructedDataLifeSafetyPointFaultValues, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataLifeSafetyPointFaultValues.deepCopy(), nil
}

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) MustBuild() BACnetConstructedDataLifeSafetyPointFaultValues {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataLifeSafetyPointFaultValuesBuilder().(*_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataLifeSafetyPointFaultValuesBuilder creates a BACnetConstructedDataLifeSafetyPointFaultValuesBuilder
func (b *_BACnetConstructedDataLifeSafetyPointFaultValues) CreateBACnetConstructedDataLifeSafetyPointFaultValuesBuilder() BACnetConstructedDataLifeSafetyPointFaultValuesBuilder {
	if b == nil {
		return NewBACnetConstructedDataLifeSafetyPointFaultValuesBuilder()
	}
	return &_BACnetConstructedDataLifeSafetyPointFaultValuesBuilder{_BACnetConstructedDataLifeSafetyPointFaultValues: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_LIFE_SAFETY_POINT
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_FAULT_VALUES
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetFaultValues() []BACnetLifeSafetyStateTagged {
	return m.FaultValues
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLifeSafetyPointFaultValues(structType any) BACnetConstructedDataLifeSafetyPointFaultValues {
	if casted, ok := structType.(BACnetConstructedDataLifeSafetyPointFaultValues); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLifeSafetyPointFaultValues); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetTypeName() string {
	return "BACnetConstructedDataLifeSafetyPointFaultValues"
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Array field
	if len(m.FaultValues) > 0 {
		for _, element := range m.FaultValues {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataLifeSafetyPointFaultValues BACnetConstructedDataLifeSafetyPointFaultValues, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLifeSafetyPointFaultValues"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLifeSafetyPointFaultValues")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	faultValues, err := ReadTerminatedArrayField[BACnetLifeSafetyStateTagged](ctx, "faultValues", ReadComplex[BACnetLifeSafetyStateTagged](BACnetLifeSafetyStateTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'faultValues' field"))
	}
	m.FaultValues = faultValues

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLifeSafetyPointFaultValues"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLifeSafetyPointFaultValues")
	}

	return m, nil
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLifeSafetyPointFaultValues"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLifeSafetyPointFaultValues")
		}

		if err := WriteComplexTypeArrayField(ctx, "faultValues", m.GetFaultValues(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'faultValues' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLifeSafetyPointFaultValues"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLifeSafetyPointFaultValues")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) IsBACnetConstructedDataLifeSafetyPointFaultValues() {
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) deepCopy() *_BACnetConstructedDataLifeSafetyPointFaultValues {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataLifeSafetyPointFaultValuesCopy := &_BACnetConstructedDataLifeSafetyPointFaultValues{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		utils.DeepCopySlice[BACnetLifeSafetyStateTagged, BACnetLifeSafetyStateTagged](m.FaultValues),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataLifeSafetyPointFaultValuesCopy
}

func (m *_BACnetConstructedDataLifeSafetyPointFaultValues) String() string {
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
