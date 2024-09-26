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

// BACnetConstructedDataSlaveProxyEnable is the corresponding interface of BACnetConstructedDataSlaveProxyEnable
type BACnetConstructedDataSlaveProxyEnable interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetSlaveProxyEnable returns SlaveProxyEnable (property field)
	GetSlaveProxyEnable() BACnetApplicationTagBoolean
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagBoolean
	// IsBACnetConstructedDataSlaveProxyEnable is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataSlaveProxyEnable()
	// CreateBuilder creates a BACnetConstructedDataSlaveProxyEnableBuilder
	CreateBACnetConstructedDataSlaveProxyEnableBuilder() BACnetConstructedDataSlaveProxyEnableBuilder
}

// _BACnetConstructedDataSlaveProxyEnable is the data-structure of this message
type _BACnetConstructedDataSlaveProxyEnable struct {
	BACnetConstructedDataContract
	SlaveProxyEnable BACnetApplicationTagBoolean
}

var _ BACnetConstructedDataSlaveProxyEnable = (*_BACnetConstructedDataSlaveProxyEnable)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataSlaveProxyEnable)(nil)

// NewBACnetConstructedDataSlaveProxyEnable factory function for _BACnetConstructedDataSlaveProxyEnable
func NewBACnetConstructedDataSlaveProxyEnable(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, slaveProxyEnable BACnetApplicationTagBoolean, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataSlaveProxyEnable {
	if slaveProxyEnable == nil {
		panic("slaveProxyEnable of type BACnetApplicationTagBoolean for BACnetConstructedDataSlaveProxyEnable must not be nil")
	}
	_result := &_BACnetConstructedDataSlaveProxyEnable{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		SlaveProxyEnable:              slaveProxyEnable,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataSlaveProxyEnableBuilder is a builder for BACnetConstructedDataSlaveProxyEnable
type BACnetConstructedDataSlaveProxyEnableBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(slaveProxyEnable BACnetApplicationTagBoolean) BACnetConstructedDataSlaveProxyEnableBuilder
	// WithSlaveProxyEnable adds SlaveProxyEnable (property field)
	WithSlaveProxyEnable(BACnetApplicationTagBoolean) BACnetConstructedDataSlaveProxyEnableBuilder
	// WithSlaveProxyEnableBuilder adds SlaveProxyEnable (property field) which is build by the builder
	WithSlaveProxyEnableBuilder(func(BACnetApplicationTagBooleanBuilder) BACnetApplicationTagBooleanBuilder) BACnetConstructedDataSlaveProxyEnableBuilder
	// Build builds the BACnetConstructedDataSlaveProxyEnable or returns an error if something is wrong
	Build() (BACnetConstructedDataSlaveProxyEnable, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataSlaveProxyEnable
}

// NewBACnetConstructedDataSlaveProxyEnableBuilder() creates a BACnetConstructedDataSlaveProxyEnableBuilder
func NewBACnetConstructedDataSlaveProxyEnableBuilder() BACnetConstructedDataSlaveProxyEnableBuilder {
	return &_BACnetConstructedDataSlaveProxyEnableBuilder{_BACnetConstructedDataSlaveProxyEnable: new(_BACnetConstructedDataSlaveProxyEnable)}
}

type _BACnetConstructedDataSlaveProxyEnableBuilder struct {
	*_BACnetConstructedDataSlaveProxyEnable

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataSlaveProxyEnableBuilder) = (*_BACnetConstructedDataSlaveProxyEnableBuilder)(nil)

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) WithMandatoryFields(slaveProxyEnable BACnetApplicationTagBoolean) BACnetConstructedDataSlaveProxyEnableBuilder {
	return b.WithSlaveProxyEnable(slaveProxyEnable)
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) WithSlaveProxyEnable(slaveProxyEnable BACnetApplicationTagBoolean) BACnetConstructedDataSlaveProxyEnableBuilder {
	b.SlaveProxyEnable = slaveProxyEnable
	return b
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) WithSlaveProxyEnableBuilder(builderSupplier func(BACnetApplicationTagBooleanBuilder) BACnetApplicationTagBooleanBuilder) BACnetConstructedDataSlaveProxyEnableBuilder {
	builder := builderSupplier(b.SlaveProxyEnable.CreateBACnetApplicationTagBooleanBuilder())
	var err error
	b.SlaveProxyEnable, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagBooleanBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) Build() (BACnetConstructedDataSlaveProxyEnable, error) {
	if b.SlaveProxyEnable == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'slaveProxyEnable' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataSlaveProxyEnable.deepCopy(), nil
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) MustBuild() BACnetConstructedDataSlaveProxyEnable {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataSlaveProxyEnableBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataSlaveProxyEnableBuilder().(*_BACnetConstructedDataSlaveProxyEnableBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataSlaveProxyEnableBuilder creates a BACnetConstructedDataSlaveProxyEnableBuilder
func (b *_BACnetConstructedDataSlaveProxyEnable) CreateBACnetConstructedDataSlaveProxyEnableBuilder() BACnetConstructedDataSlaveProxyEnableBuilder {
	if b == nil {
		return NewBACnetConstructedDataSlaveProxyEnableBuilder()
	}
	return &_BACnetConstructedDataSlaveProxyEnableBuilder{_BACnetConstructedDataSlaveProxyEnable: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataSlaveProxyEnable) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataSlaveProxyEnable) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_SLAVE_PROXY_ENABLE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataSlaveProxyEnable) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataSlaveProxyEnable) GetSlaveProxyEnable() BACnetApplicationTagBoolean {
	return m.SlaveProxyEnable
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataSlaveProxyEnable) GetActualValue() BACnetApplicationTagBoolean {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagBoolean(m.GetSlaveProxyEnable())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataSlaveProxyEnable(structType any) BACnetConstructedDataSlaveProxyEnable {
	if casted, ok := structType.(BACnetConstructedDataSlaveProxyEnable); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataSlaveProxyEnable); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataSlaveProxyEnable) GetTypeName() string {
	return "BACnetConstructedDataSlaveProxyEnable"
}

func (m *_BACnetConstructedDataSlaveProxyEnable) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (slaveProxyEnable)
	lengthInBits += m.SlaveProxyEnable.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataSlaveProxyEnable) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataSlaveProxyEnable) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataSlaveProxyEnable BACnetConstructedDataSlaveProxyEnable, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataSlaveProxyEnable"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataSlaveProxyEnable")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	slaveProxyEnable, err := ReadSimpleField[BACnetApplicationTagBoolean](ctx, "slaveProxyEnable", ReadComplex[BACnetApplicationTagBoolean](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagBoolean](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'slaveProxyEnable' field"))
	}
	m.SlaveProxyEnable = slaveProxyEnable

	actualValue, err := ReadVirtualField[BACnetApplicationTagBoolean](ctx, "actualValue", (*BACnetApplicationTagBoolean)(nil), slaveProxyEnable)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataSlaveProxyEnable"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataSlaveProxyEnable")
	}

	return m, nil
}

func (m *_BACnetConstructedDataSlaveProxyEnable) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataSlaveProxyEnable) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataSlaveProxyEnable"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataSlaveProxyEnable")
		}

		if err := WriteSimpleField[BACnetApplicationTagBoolean](ctx, "slaveProxyEnable", m.GetSlaveProxyEnable(), WriteComplex[BACnetApplicationTagBoolean](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'slaveProxyEnable' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataSlaveProxyEnable"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataSlaveProxyEnable")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataSlaveProxyEnable) IsBACnetConstructedDataSlaveProxyEnable() {}

func (m *_BACnetConstructedDataSlaveProxyEnable) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataSlaveProxyEnable) deepCopy() *_BACnetConstructedDataSlaveProxyEnable {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataSlaveProxyEnableCopy := &_BACnetConstructedDataSlaveProxyEnable{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.SlaveProxyEnable.DeepCopy().(BACnetApplicationTagBoolean),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataSlaveProxyEnableCopy
}

func (m *_BACnetConstructedDataSlaveProxyEnable) String() string {
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
