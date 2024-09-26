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

// DatagramConnectionTransportDataType is the corresponding interface of DatagramConnectionTransportDataType
type DatagramConnectionTransportDataType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetDiscoveryAddress returns DiscoveryAddress (property field)
	GetDiscoveryAddress() ExtensionObject
	// IsDatagramConnectionTransportDataType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsDatagramConnectionTransportDataType()
	// CreateBuilder creates a DatagramConnectionTransportDataTypeBuilder
	CreateDatagramConnectionTransportDataTypeBuilder() DatagramConnectionTransportDataTypeBuilder
}

// _DatagramConnectionTransportDataType is the data-structure of this message
type _DatagramConnectionTransportDataType struct {
	ExtensionObjectDefinitionContract
	DiscoveryAddress ExtensionObject
}

var _ DatagramConnectionTransportDataType = (*_DatagramConnectionTransportDataType)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_DatagramConnectionTransportDataType)(nil)

// NewDatagramConnectionTransportDataType factory function for _DatagramConnectionTransportDataType
func NewDatagramConnectionTransportDataType(discoveryAddress ExtensionObject) *_DatagramConnectionTransportDataType {
	if discoveryAddress == nil {
		panic("discoveryAddress of type ExtensionObject for DatagramConnectionTransportDataType must not be nil")
	}
	_result := &_DatagramConnectionTransportDataType{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		DiscoveryAddress:                  discoveryAddress,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// DatagramConnectionTransportDataTypeBuilder is a builder for DatagramConnectionTransportDataType
type DatagramConnectionTransportDataTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(discoveryAddress ExtensionObject) DatagramConnectionTransportDataTypeBuilder
	// WithDiscoveryAddress adds DiscoveryAddress (property field)
	WithDiscoveryAddress(ExtensionObject) DatagramConnectionTransportDataTypeBuilder
	// WithDiscoveryAddressBuilder adds DiscoveryAddress (property field) which is build by the builder
	WithDiscoveryAddressBuilder(func(ExtensionObjectBuilder) ExtensionObjectBuilder) DatagramConnectionTransportDataTypeBuilder
	// Build builds the DatagramConnectionTransportDataType or returns an error if something is wrong
	Build() (DatagramConnectionTransportDataType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() DatagramConnectionTransportDataType
}

// NewDatagramConnectionTransportDataTypeBuilder() creates a DatagramConnectionTransportDataTypeBuilder
func NewDatagramConnectionTransportDataTypeBuilder() DatagramConnectionTransportDataTypeBuilder {
	return &_DatagramConnectionTransportDataTypeBuilder{_DatagramConnectionTransportDataType: new(_DatagramConnectionTransportDataType)}
}

type _DatagramConnectionTransportDataTypeBuilder struct {
	*_DatagramConnectionTransportDataType

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (DatagramConnectionTransportDataTypeBuilder) = (*_DatagramConnectionTransportDataTypeBuilder)(nil)

func (b *_DatagramConnectionTransportDataTypeBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_DatagramConnectionTransportDataTypeBuilder) WithMandatoryFields(discoveryAddress ExtensionObject) DatagramConnectionTransportDataTypeBuilder {
	return b.WithDiscoveryAddress(discoveryAddress)
}

func (b *_DatagramConnectionTransportDataTypeBuilder) WithDiscoveryAddress(discoveryAddress ExtensionObject) DatagramConnectionTransportDataTypeBuilder {
	b.DiscoveryAddress = discoveryAddress
	return b
}

func (b *_DatagramConnectionTransportDataTypeBuilder) WithDiscoveryAddressBuilder(builderSupplier func(ExtensionObjectBuilder) ExtensionObjectBuilder) DatagramConnectionTransportDataTypeBuilder {
	builder := builderSupplier(b.DiscoveryAddress.CreateExtensionObjectBuilder())
	var err error
	b.DiscoveryAddress, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "ExtensionObjectBuilder failed"))
	}
	return b
}

func (b *_DatagramConnectionTransportDataTypeBuilder) Build() (DatagramConnectionTransportDataType, error) {
	if b.DiscoveryAddress == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'discoveryAddress' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._DatagramConnectionTransportDataType.deepCopy(), nil
}

func (b *_DatagramConnectionTransportDataTypeBuilder) MustBuild() DatagramConnectionTransportDataType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_DatagramConnectionTransportDataTypeBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_DatagramConnectionTransportDataTypeBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_DatagramConnectionTransportDataTypeBuilder) DeepCopy() any {
	_copy := b.CreateDatagramConnectionTransportDataTypeBuilder().(*_DatagramConnectionTransportDataTypeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateDatagramConnectionTransportDataTypeBuilder creates a DatagramConnectionTransportDataTypeBuilder
func (b *_DatagramConnectionTransportDataType) CreateDatagramConnectionTransportDataTypeBuilder() DatagramConnectionTransportDataTypeBuilder {
	if b == nil {
		return NewDatagramConnectionTransportDataTypeBuilder()
	}
	return &_DatagramConnectionTransportDataTypeBuilder{_DatagramConnectionTransportDataType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_DatagramConnectionTransportDataType) GetExtensionId() int32 {
	return int32(17469)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_DatagramConnectionTransportDataType) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_DatagramConnectionTransportDataType) GetDiscoveryAddress() ExtensionObject {
	return m.DiscoveryAddress
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastDatagramConnectionTransportDataType(structType any) DatagramConnectionTransportDataType {
	if casted, ok := structType.(DatagramConnectionTransportDataType); ok {
		return casted
	}
	if casted, ok := structType.(*DatagramConnectionTransportDataType); ok {
		return *casted
	}
	return nil
}

func (m *_DatagramConnectionTransportDataType) GetTypeName() string {
	return "DatagramConnectionTransportDataType"
}

func (m *_DatagramConnectionTransportDataType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (discoveryAddress)
	lengthInBits += m.DiscoveryAddress.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_DatagramConnectionTransportDataType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_DatagramConnectionTransportDataType) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__datagramConnectionTransportDataType DatagramConnectionTransportDataType, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("DatagramConnectionTransportDataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for DatagramConnectionTransportDataType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	discoveryAddress, err := ReadSimpleField[ExtensionObject](ctx, "discoveryAddress", ReadComplex[ExtensionObject](ExtensionObjectParseWithBufferProducer[ExtensionObject]((bool)(bool(true))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'discoveryAddress' field"))
	}
	m.DiscoveryAddress = discoveryAddress

	if closeErr := readBuffer.CloseContext("DatagramConnectionTransportDataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for DatagramConnectionTransportDataType")
	}

	return m, nil
}

func (m *_DatagramConnectionTransportDataType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_DatagramConnectionTransportDataType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("DatagramConnectionTransportDataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for DatagramConnectionTransportDataType")
		}

		if err := WriteSimpleField[ExtensionObject](ctx, "discoveryAddress", m.GetDiscoveryAddress(), WriteComplex[ExtensionObject](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'discoveryAddress' field")
		}

		if popErr := writeBuffer.PopContext("DatagramConnectionTransportDataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for DatagramConnectionTransportDataType")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_DatagramConnectionTransportDataType) IsDatagramConnectionTransportDataType() {}

func (m *_DatagramConnectionTransportDataType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_DatagramConnectionTransportDataType) deepCopy() *_DatagramConnectionTransportDataType {
	if m == nil {
		return nil
	}
	_DatagramConnectionTransportDataTypeCopy := &_DatagramConnectionTransportDataType{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.DiscoveryAddress.DeepCopy().(ExtensionObject),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _DatagramConnectionTransportDataTypeCopy
}

func (m *_DatagramConnectionTransportDataType) String() string {
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
