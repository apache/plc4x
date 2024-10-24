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

// BACnetConstructedDataLoggingObject is the corresponding interface of BACnetConstructedDataLoggingObject
type BACnetConstructedDataLoggingObject interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetLoggingObject returns LoggingObject (property field)
	GetLoggingObject() BACnetApplicationTagObjectIdentifier
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagObjectIdentifier
	// IsBACnetConstructedDataLoggingObject is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataLoggingObject()
	// CreateBuilder creates a BACnetConstructedDataLoggingObjectBuilder
	CreateBACnetConstructedDataLoggingObjectBuilder() BACnetConstructedDataLoggingObjectBuilder
}

// _BACnetConstructedDataLoggingObject is the data-structure of this message
type _BACnetConstructedDataLoggingObject struct {
	BACnetConstructedDataContract
	LoggingObject BACnetApplicationTagObjectIdentifier
}

var _ BACnetConstructedDataLoggingObject = (*_BACnetConstructedDataLoggingObject)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataLoggingObject)(nil)

// NewBACnetConstructedDataLoggingObject factory function for _BACnetConstructedDataLoggingObject
func NewBACnetConstructedDataLoggingObject(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, loggingObject BACnetApplicationTagObjectIdentifier, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLoggingObject {
	if loggingObject == nil {
		panic("loggingObject of type BACnetApplicationTagObjectIdentifier for BACnetConstructedDataLoggingObject must not be nil")
	}
	_result := &_BACnetConstructedDataLoggingObject{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		LoggingObject:                 loggingObject,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataLoggingObjectBuilder is a builder for BACnetConstructedDataLoggingObject
type BACnetConstructedDataLoggingObjectBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(loggingObject BACnetApplicationTagObjectIdentifier) BACnetConstructedDataLoggingObjectBuilder
	// WithLoggingObject adds LoggingObject (property field)
	WithLoggingObject(BACnetApplicationTagObjectIdentifier) BACnetConstructedDataLoggingObjectBuilder
	// WithLoggingObjectBuilder adds LoggingObject (property field) which is build by the builder
	WithLoggingObjectBuilder(func(BACnetApplicationTagObjectIdentifierBuilder) BACnetApplicationTagObjectIdentifierBuilder) BACnetConstructedDataLoggingObjectBuilder
	// Build builds the BACnetConstructedDataLoggingObject or returns an error if something is wrong
	Build() (BACnetConstructedDataLoggingObject, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataLoggingObject
}

// NewBACnetConstructedDataLoggingObjectBuilder() creates a BACnetConstructedDataLoggingObjectBuilder
func NewBACnetConstructedDataLoggingObjectBuilder() BACnetConstructedDataLoggingObjectBuilder {
	return &_BACnetConstructedDataLoggingObjectBuilder{_BACnetConstructedDataLoggingObject: new(_BACnetConstructedDataLoggingObject)}
}

type _BACnetConstructedDataLoggingObjectBuilder struct {
	*_BACnetConstructedDataLoggingObject

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataLoggingObjectBuilder) = (*_BACnetConstructedDataLoggingObjectBuilder)(nil)

func (b *_BACnetConstructedDataLoggingObjectBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) WithMandatoryFields(loggingObject BACnetApplicationTagObjectIdentifier) BACnetConstructedDataLoggingObjectBuilder {
	return b.WithLoggingObject(loggingObject)
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) WithLoggingObject(loggingObject BACnetApplicationTagObjectIdentifier) BACnetConstructedDataLoggingObjectBuilder {
	b.LoggingObject = loggingObject
	return b
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) WithLoggingObjectBuilder(builderSupplier func(BACnetApplicationTagObjectIdentifierBuilder) BACnetApplicationTagObjectIdentifierBuilder) BACnetConstructedDataLoggingObjectBuilder {
	builder := builderSupplier(b.LoggingObject.CreateBACnetApplicationTagObjectIdentifierBuilder())
	var err error
	b.LoggingObject, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagObjectIdentifierBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) Build() (BACnetConstructedDataLoggingObject, error) {
	if b.LoggingObject == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'loggingObject' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataLoggingObject.deepCopy(), nil
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) MustBuild() BACnetConstructedDataLoggingObject {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataLoggingObjectBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataLoggingObjectBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataLoggingObjectBuilder().(*_BACnetConstructedDataLoggingObjectBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataLoggingObjectBuilder creates a BACnetConstructedDataLoggingObjectBuilder
func (b *_BACnetConstructedDataLoggingObject) CreateBACnetConstructedDataLoggingObjectBuilder() BACnetConstructedDataLoggingObjectBuilder {
	if b == nil {
		return NewBACnetConstructedDataLoggingObjectBuilder()
	}
	return &_BACnetConstructedDataLoggingObjectBuilder{_BACnetConstructedDataLoggingObject: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLoggingObject) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataLoggingObject) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_LOGGING_OBJECT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLoggingObject) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLoggingObject) GetLoggingObject() BACnetApplicationTagObjectIdentifier {
	return m.LoggingObject
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataLoggingObject) GetActualValue() BACnetApplicationTagObjectIdentifier {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagObjectIdentifier(m.GetLoggingObject())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLoggingObject(structType any) BACnetConstructedDataLoggingObject {
	if casted, ok := structType.(BACnetConstructedDataLoggingObject); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLoggingObject); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLoggingObject) GetTypeName() string {
	return "BACnetConstructedDataLoggingObject"
}

func (m *_BACnetConstructedDataLoggingObject) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (loggingObject)
	lengthInBits += m.LoggingObject.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataLoggingObject) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataLoggingObject) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataLoggingObject BACnetConstructedDataLoggingObject, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLoggingObject"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLoggingObject")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	loggingObject, err := ReadSimpleField[BACnetApplicationTagObjectIdentifier](ctx, "loggingObject", ReadComplex[BACnetApplicationTagObjectIdentifier](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagObjectIdentifier](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'loggingObject' field"))
	}
	m.LoggingObject = loggingObject

	actualValue, err := ReadVirtualField[BACnetApplicationTagObjectIdentifier](ctx, "actualValue", (*BACnetApplicationTagObjectIdentifier)(nil), loggingObject)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLoggingObject"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLoggingObject")
	}

	return m, nil
}

func (m *_BACnetConstructedDataLoggingObject) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataLoggingObject) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLoggingObject"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLoggingObject")
		}

		if err := WriteSimpleField[BACnetApplicationTagObjectIdentifier](ctx, "loggingObject", m.GetLoggingObject(), WriteComplex[BACnetApplicationTagObjectIdentifier](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'loggingObject' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLoggingObject"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLoggingObject")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLoggingObject) IsBACnetConstructedDataLoggingObject() {}

func (m *_BACnetConstructedDataLoggingObject) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataLoggingObject) deepCopy() *_BACnetConstructedDataLoggingObject {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataLoggingObjectCopy := &_BACnetConstructedDataLoggingObject{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.LoggingObject.DeepCopy().(BACnetApplicationTagObjectIdentifier),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataLoggingObjectCopy
}

func (m *_BACnetConstructedDataLoggingObject) String() string {
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
