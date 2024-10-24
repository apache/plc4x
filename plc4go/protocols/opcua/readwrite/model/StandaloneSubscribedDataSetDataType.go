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

// StandaloneSubscribedDataSetDataType is the corresponding interface of StandaloneSubscribedDataSetDataType
type StandaloneSubscribedDataSetDataType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetName returns Name (property field)
	GetName() PascalString
	// GetDataSetFolder returns DataSetFolder (property field)
	GetDataSetFolder() []PascalString
	// GetDataSetMetaData returns DataSetMetaData (property field)
	GetDataSetMetaData() DataSetMetaDataType
	// GetSubscribedDataSet returns SubscribedDataSet (property field)
	GetSubscribedDataSet() ExtensionObject
	// IsStandaloneSubscribedDataSetDataType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsStandaloneSubscribedDataSetDataType()
	// CreateBuilder creates a StandaloneSubscribedDataSetDataTypeBuilder
	CreateStandaloneSubscribedDataSetDataTypeBuilder() StandaloneSubscribedDataSetDataTypeBuilder
}

// _StandaloneSubscribedDataSetDataType is the data-structure of this message
type _StandaloneSubscribedDataSetDataType struct {
	ExtensionObjectDefinitionContract
	Name              PascalString
	DataSetFolder     []PascalString
	DataSetMetaData   DataSetMetaDataType
	SubscribedDataSet ExtensionObject
}

var _ StandaloneSubscribedDataSetDataType = (*_StandaloneSubscribedDataSetDataType)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_StandaloneSubscribedDataSetDataType)(nil)

// NewStandaloneSubscribedDataSetDataType factory function for _StandaloneSubscribedDataSetDataType
func NewStandaloneSubscribedDataSetDataType(name PascalString, dataSetFolder []PascalString, dataSetMetaData DataSetMetaDataType, subscribedDataSet ExtensionObject) *_StandaloneSubscribedDataSetDataType {
	if name == nil {
		panic("name of type PascalString for StandaloneSubscribedDataSetDataType must not be nil")
	}
	if dataSetMetaData == nil {
		panic("dataSetMetaData of type DataSetMetaDataType for StandaloneSubscribedDataSetDataType must not be nil")
	}
	if subscribedDataSet == nil {
		panic("subscribedDataSet of type ExtensionObject for StandaloneSubscribedDataSetDataType must not be nil")
	}
	_result := &_StandaloneSubscribedDataSetDataType{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		Name:                              name,
		DataSetFolder:                     dataSetFolder,
		DataSetMetaData:                   dataSetMetaData,
		SubscribedDataSet:                 subscribedDataSet,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// StandaloneSubscribedDataSetDataTypeBuilder is a builder for StandaloneSubscribedDataSetDataType
type StandaloneSubscribedDataSetDataTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(name PascalString, dataSetFolder []PascalString, dataSetMetaData DataSetMetaDataType, subscribedDataSet ExtensionObject) StandaloneSubscribedDataSetDataTypeBuilder
	// WithName adds Name (property field)
	WithName(PascalString) StandaloneSubscribedDataSetDataTypeBuilder
	// WithNameBuilder adds Name (property field) which is build by the builder
	WithNameBuilder(func(PascalStringBuilder) PascalStringBuilder) StandaloneSubscribedDataSetDataTypeBuilder
	// WithDataSetFolder adds DataSetFolder (property field)
	WithDataSetFolder(...PascalString) StandaloneSubscribedDataSetDataTypeBuilder
	// WithDataSetMetaData adds DataSetMetaData (property field)
	WithDataSetMetaData(DataSetMetaDataType) StandaloneSubscribedDataSetDataTypeBuilder
	// WithDataSetMetaDataBuilder adds DataSetMetaData (property field) which is build by the builder
	WithDataSetMetaDataBuilder(func(DataSetMetaDataTypeBuilder) DataSetMetaDataTypeBuilder) StandaloneSubscribedDataSetDataTypeBuilder
	// WithSubscribedDataSet adds SubscribedDataSet (property field)
	WithSubscribedDataSet(ExtensionObject) StandaloneSubscribedDataSetDataTypeBuilder
	// WithSubscribedDataSetBuilder adds SubscribedDataSet (property field) which is build by the builder
	WithSubscribedDataSetBuilder(func(ExtensionObjectBuilder) ExtensionObjectBuilder) StandaloneSubscribedDataSetDataTypeBuilder
	// Build builds the StandaloneSubscribedDataSetDataType or returns an error if something is wrong
	Build() (StandaloneSubscribedDataSetDataType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() StandaloneSubscribedDataSetDataType
}

// NewStandaloneSubscribedDataSetDataTypeBuilder() creates a StandaloneSubscribedDataSetDataTypeBuilder
func NewStandaloneSubscribedDataSetDataTypeBuilder() StandaloneSubscribedDataSetDataTypeBuilder {
	return &_StandaloneSubscribedDataSetDataTypeBuilder{_StandaloneSubscribedDataSetDataType: new(_StandaloneSubscribedDataSetDataType)}
}

type _StandaloneSubscribedDataSetDataTypeBuilder struct {
	*_StandaloneSubscribedDataSetDataType

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (StandaloneSubscribedDataSetDataTypeBuilder) = (*_StandaloneSubscribedDataSetDataTypeBuilder)(nil)

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithMandatoryFields(name PascalString, dataSetFolder []PascalString, dataSetMetaData DataSetMetaDataType, subscribedDataSet ExtensionObject) StandaloneSubscribedDataSetDataTypeBuilder {
	return b.WithName(name).WithDataSetFolder(dataSetFolder...).WithDataSetMetaData(dataSetMetaData).WithSubscribedDataSet(subscribedDataSet)
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithName(name PascalString) StandaloneSubscribedDataSetDataTypeBuilder {
	b.Name = name
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithNameBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) StandaloneSubscribedDataSetDataTypeBuilder {
	builder := builderSupplier(b.Name.CreatePascalStringBuilder())
	var err error
	b.Name, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithDataSetFolder(dataSetFolder ...PascalString) StandaloneSubscribedDataSetDataTypeBuilder {
	b.DataSetFolder = dataSetFolder
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithDataSetMetaData(dataSetMetaData DataSetMetaDataType) StandaloneSubscribedDataSetDataTypeBuilder {
	b.DataSetMetaData = dataSetMetaData
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithDataSetMetaDataBuilder(builderSupplier func(DataSetMetaDataTypeBuilder) DataSetMetaDataTypeBuilder) StandaloneSubscribedDataSetDataTypeBuilder {
	builder := builderSupplier(b.DataSetMetaData.CreateDataSetMetaDataTypeBuilder())
	var err error
	b.DataSetMetaData, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "DataSetMetaDataTypeBuilder failed"))
	}
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithSubscribedDataSet(subscribedDataSet ExtensionObject) StandaloneSubscribedDataSetDataTypeBuilder {
	b.SubscribedDataSet = subscribedDataSet
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) WithSubscribedDataSetBuilder(builderSupplier func(ExtensionObjectBuilder) ExtensionObjectBuilder) StandaloneSubscribedDataSetDataTypeBuilder {
	builder := builderSupplier(b.SubscribedDataSet.CreateExtensionObjectBuilder())
	var err error
	b.SubscribedDataSet, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "ExtensionObjectBuilder failed"))
	}
	return b
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) Build() (StandaloneSubscribedDataSetDataType, error) {
	if b.Name == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'name' not set"))
	}
	if b.DataSetMetaData == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'dataSetMetaData' not set"))
	}
	if b.SubscribedDataSet == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'subscribedDataSet' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._StandaloneSubscribedDataSetDataType.deepCopy(), nil
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) MustBuild() StandaloneSubscribedDataSetDataType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_StandaloneSubscribedDataSetDataTypeBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_StandaloneSubscribedDataSetDataTypeBuilder) DeepCopy() any {
	_copy := b.CreateStandaloneSubscribedDataSetDataTypeBuilder().(*_StandaloneSubscribedDataSetDataTypeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateStandaloneSubscribedDataSetDataTypeBuilder creates a StandaloneSubscribedDataSetDataTypeBuilder
func (b *_StandaloneSubscribedDataSetDataType) CreateStandaloneSubscribedDataSetDataTypeBuilder() StandaloneSubscribedDataSetDataTypeBuilder {
	if b == nil {
		return NewStandaloneSubscribedDataSetDataTypeBuilder()
	}
	return &_StandaloneSubscribedDataSetDataTypeBuilder{_StandaloneSubscribedDataSetDataType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_StandaloneSubscribedDataSetDataType) GetExtensionId() int32 {
	return int32(23602)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_StandaloneSubscribedDataSetDataType) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_StandaloneSubscribedDataSetDataType) GetName() PascalString {
	return m.Name
}

func (m *_StandaloneSubscribedDataSetDataType) GetDataSetFolder() []PascalString {
	return m.DataSetFolder
}

func (m *_StandaloneSubscribedDataSetDataType) GetDataSetMetaData() DataSetMetaDataType {
	return m.DataSetMetaData
}

func (m *_StandaloneSubscribedDataSetDataType) GetSubscribedDataSet() ExtensionObject {
	return m.SubscribedDataSet
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastStandaloneSubscribedDataSetDataType(structType any) StandaloneSubscribedDataSetDataType {
	if casted, ok := structType.(StandaloneSubscribedDataSetDataType); ok {
		return casted
	}
	if casted, ok := structType.(*StandaloneSubscribedDataSetDataType); ok {
		return *casted
	}
	return nil
}

func (m *_StandaloneSubscribedDataSetDataType) GetTypeName() string {
	return "StandaloneSubscribedDataSetDataType"
}

func (m *_StandaloneSubscribedDataSetDataType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (name)
	lengthInBits += m.Name.GetLengthInBits(ctx)

	// Implicit Field (noOfDataSetFolder)
	lengthInBits += 32

	// Array field
	if len(m.DataSetFolder) > 0 {
		for _curItem, element := range m.DataSetFolder {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.DataSetFolder), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (dataSetMetaData)
	lengthInBits += m.DataSetMetaData.GetLengthInBits(ctx)

	// Simple field (subscribedDataSet)
	lengthInBits += m.SubscribedDataSet.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_StandaloneSubscribedDataSetDataType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_StandaloneSubscribedDataSetDataType) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__standaloneSubscribedDataSetDataType StandaloneSubscribedDataSetDataType, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("StandaloneSubscribedDataSetDataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for StandaloneSubscribedDataSetDataType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	name, err := ReadSimpleField[PascalString](ctx, "name", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'name' field"))
	}
	m.Name = name

	noOfDataSetFolder, err := ReadImplicitField[int32](ctx, "noOfDataSetFolder", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfDataSetFolder' field"))
	}
	_ = noOfDataSetFolder

	dataSetFolder, err := ReadCountArrayField[PascalString](ctx, "dataSetFolder", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer), uint64(noOfDataSetFolder))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dataSetFolder' field"))
	}
	m.DataSetFolder = dataSetFolder

	dataSetMetaData, err := ReadSimpleField[DataSetMetaDataType](ctx, "dataSetMetaData", ReadComplex[DataSetMetaDataType](ExtensionObjectDefinitionParseWithBufferProducer[DataSetMetaDataType]((int32)(int32(14525))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dataSetMetaData' field"))
	}
	m.DataSetMetaData = dataSetMetaData

	subscribedDataSet, err := ReadSimpleField[ExtensionObject](ctx, "subscribedDataSet", ReadComplex[ExtensionObject](ExtensionObjectParseWithBufferProducer[ExtensionObject]((bool)(bool(true))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'subscribedDataSet' field"))
	}
	m.SubscribedDataSet = subscribedDataSet

	if closeErr := readBuffer.CloseContext("StandaloneSubscribedDataSetDataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for StandaloneSubscribedDataSetDataType")
	}

	return m, nil
}

func (m *_StandaloneSubscribedDataSetDataType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_StandaloneSubscribedDataSetDataType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("StandaloneSubscribedDataSetDataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for StandaloneSubscribedDataSetDataType")
		}

		if err := WriteSimpleField[PascalString](ctx, "name", m.GetName(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'name' field")
		}
		noOfDataSetFolder := int32(utils.InlineIf(bool((m.GetDataSetFolder()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetDataSetFolder()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfDataSetFolder", noOfDataSetFolder, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfDataSetFolder' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "dataSetFolder", m.GetDataSetFolder(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'dataSetFolder' field")
		}

		if err := WriteSimpleField[DataSetMetaDataType](ctx, "dataSetMetaData", m.GetDataSetMetaData(), WriteComplex[DataSetMetaDataType](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'dataSetMetaData' field")
		}

		if err := WriteSimpleField[ExtensionObject](ctx, "subscribedDataSet", m.GetSubscribedDataSet(), WriteComplex[ExtensionObject](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'subscribedDataSet' field")
		}

		if popErr := writeBuffer.PopContext("StandaloneSubscribedDataSetDataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for StandaloneSubscribedDataSetDataType")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_StandaloneSubscribedDataSetDataType) IsStandaloneSubscribedDataSetDataType() {}

func (m *_StandaloneSubscribedDataSetDataType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_StandaloneSubscribedDataSetDataType) deepCopy() *_StandaloneSubscribedDataSetDataType {
	if m == nil {
		return nil
	}
	_StandaloneSubscribedDataSetDataTypeCopy := &_StandaloneSubscribedDataSetDataType{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.Name.DeepCopy().(PascalString),
		utils.DeepCopySlice[PascalString, PascalString](m.DataSetFolder),
		m.DataSetMetaData.DeepCopy().(DataSetMetaDataType),
		m.SubscribedDataSet.DeepCopy().(ExtensionObject),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _StandaloneSubscribedDataSetDataTypeCopy
}

func (m *_StandaloneSubscribedDataSetDataType) String() string {
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
