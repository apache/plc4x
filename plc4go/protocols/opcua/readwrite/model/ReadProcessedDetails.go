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

// ReadProcessedDetails is the corresponding interface of ReadProcessedDetails
type ReadProcessedDetails interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetStartTime returns StartTime (property field)
	GetStartTime() int64
	// GetEndTime returns EndTime (property field)
	GetEndTime() int64
	// GetProcessingInterval returns ProcessingInterval (property field)
	GetProcessingInterval() float64
	// GetAggregateType returns AggregateType (property field)
	GetAggregateType() []NodeId
	// GetAggregateConfiguration returns AggregateConfiguration (property field)
	GetAggregateConfiguration() AggregateConfiguration
	// IsReadProcessedDetails is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsReadProcessedDetails()
	// CreateBuilder creates a ReadProcessedDetailsBuilder
	CreateReadProcessedDetailsBuilder() ReadProcessedDetailsBuilder
}

// _ReadProcessedDetails is the data-structure of this message
type _ReadProcessedDetails struct {
	ExtensionObjectDefinitionContract
	StartTime              int64
	EndTime                int64
	ProcessingInterval     float64
	AggregateType          []NodeId
	AggregateConfiguration AggregateConfiguration
}

var _ ReadProcessedDetails = (*_ReadProcessedDetails)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_ReadProcessedDetails)(nil)

// NewReadProcessedDetails factory function for _ReadProcessedDetails
func NewReadProcessedDetails(startTime int64, endTime int64, processingInterval float64, aggregateType []NodeId, aggregateConfiguration AggregateConfiguration) *_ReadProcessedDetails {
	if aggregateConfiguration == nil {
		panic("aggregateConfiguration of type AggregateConfiguration for ReadProcessedDetails must not be nil")
	}
	_result := &_ReadProcessedDetails{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		StartTime:                         startTime,
		EndTime:                           endTime,
		ProcessingInterval:                processingInterval,
		AggregateType:                     aggregateType,
		AggregateConfiguration:            aggregateConfiguration,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ReadProcessedDetailsBuilder is a builder for ReadProcessedDetails
type ReadProcessedDetailsBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(startTime int64, endTime int64, processingInterval float64, aggregateType []NodeId, aggregateConfiguration AggregateConfiguration) ReadProcessedDetailsBuilder
	// WithStartTime adds StartTime (property field)
	WithStartTime(int64) ReadProcessedDetailsBuilder
	// WithEndTime adds EndTime (property field)
	WithEndTime(int64) ReadProcessedDetailsBuilder
	// WithProcessingInterval adds ProcessingInterval (property field)
	WithProcessingInterval(float64) ReadProcessedDetailsBuilder
	// WithAggregateType adds AggregateType (property field)
	WithAggregateType(...NodeId) ReadProcessedDetailsBuilder
	// WithAggregateConfiguration adds AggregateConfiguration (property field)
	WithAggregateConfiguration(AggregateConfiguration) ReadProcessedDetailsBuilder
	// WithAggregateConfigurationBuilder adds AggregateConfiguration (property field) which is build by the builder
	WithAggregateConfigurationBuilder(func(AggregateConfigurationBuilder) AggregateConfigurationBuilder) ReadProcessedDetailsBuilder
	// Build builds the ReadProcessedDetails or returns an error if something is wrong
	Build() (ReadProcessedDetails, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ReadProcessedDetails
}

// NewReadProcessedDetailsBuilder() creates a ReadProcessedDetailsBuilder
func NewReadProcessedDetailsBuilder() ReadProcessedDetailsBuilder {
	return &_ReadProcessedDetailsBuilder{_ReadProcessedDetails: new(_ReadProcessedDetails)}
}

type _ReadProcessedDetailsBuilder struct {
	*_ReadProcessedDetails

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (ReadProcessedDetailsBuilder) = (*_ReadProcessedDetailsBuilder)(nil)

func (b *_ReadProcessedDetailsBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_ReadProcessedDetailsBuilder) WithMandatoryFields(startTime int64, endTime int64, processingInterval float64, aggregateType []NodeId, aggregateConfiguration AggregateConfiguration) ReadProcessedDetailsBuilder {
	return b.WithStartTime(startTime).WithEndTime(endTime).WithProcessingInterval(processingInterval).WithAggregateType(aggregateType...).WithAggregateConfiguration(aggregateConfiguration)
}

func (b *_ReadProcessedDetailsBuilder) WithStartTime(startTime int64) ReadProcessedDetailsBuilder {
	b.StartTime = startTime
	return b
}

func (b *_ReadProcessedDetailsBuilder) WithEndTime(endTime int64) ReadProcessedDetailsBuilder {
	b.EndTime = endTime
	return b
}

func (b *_ReadProcessedDetailsBuilder) WithProcessingInterval(processingInterval float64) ReadProcessedDetailsBuilder {
	b.ProcessingInterval = processingInterval
	return b
}

func (b *_ReadProcessedDetailsBuilder) WithAggregateType(aggregateType ...NodeId) ReadProcessedDetailsBuilder {
	b.AggregateType = aggregateType
	return b
}

func (b *_ReadProcessedDetailsBuilder) WithAggregateConfiguration(aggregateConfiguration AggregateConfiguration) ReadProcessedDetailsBuilder {
	b.AggregateConfiguration = aggregateConfiguration
	return b
}

func (b *_ReadProcessedDetailsBuilder) WithAggregateConfigurationBuilder(builderSupplier func(AggregateConfigurationBuilder) AggregateConfigurationBuilder) ReadProcessedDetailsBuilder {
	builder := builderSupplier(b.AggregateConfiguration.CreateAggregateConfigurationBuilder())
	var err error
	b.AggregateConfiguration, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "AggregateConfigurationBuilder failed"))
	}
	return b
}

func (b *_ReadProcessedDetailsBuilder) Build() (ReadProcessedDetails, error) {
	if b.AggregateConfiguration == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'aggregateConfiguration' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ReadProcessedDetails.deepCopy(), nil
}

func (b *_ReadProcessedDetailsBuilder) MustBuild() ReadProcessedDetails {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ReadProcessedDetailsBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_ReadProcessedDetailsBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_ReadProcessedDetailsBuilder) DeepCopy() any {
	_copy := b.CreateReadProcessedDetailsBuilder().(*_ReadProcessedDetailsBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateReadProcessedDetailsBuilder creates a ReadProcessedDetailsBuilder
func (b *_ReadProcessedDetails) CreateReadProcessedDetailsBuilder() ReadProcessedDetailsBuilder {
	if b == nil {
		return NewReadProcessedDetailsBuilder()
	}
	return &_ReadProcessedDetailsBuilder{_ReadProcessedDetails: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ReadProcessedDetails) GetExtensionId() int32 {
	return int32(652)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ReadProcessedDetails) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ReadProcessedDetails) GetStartTime() int64 {
	return m.StartTime
}

func (m *_ReadProcessedDetails) GetEndTime() int64 {
	return m.EndTime
}

func (m *_ReadProcessedDetails) GetProcessingInterval() float64 {
	return m.ProcessingInterval
}

func (m *_ReadProcessedDetails) GetAggregateType() []NodeId {
	return m.AggregateType
}

func (m *_ReadProcessedDetails) GetAggregateConfiguration() AggregateConfiguration {
	return m.AggregateConfiguration
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastReadProcessedDetails(structType any) ReadProcessedDetails {
	if casted, ok := structType.(ReadProcessedDetails); ok {
		return casted
	}
	if casted, ok := structType.(*ReadProcessedDetails); ok {
		return *casted
	}
	return nil
}

func (m *_ReadProcessedDetails) GetTypeName() string {
	return "ReadProcessedDetails"
}

func (m *_ReadProcessedDetails) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (startTime)
	lengthInBits += 64

	// Simple field (endTime)
	lengthInBits += 64

	// Simple field (processingInterval)
	lengthInBits += 64

	// Implicit Field (noOfAggregateType)
	lengthInBits += 32

	// Array field
	if len(m.AggregateType) > 0 {
		for _curItem, element := range m.AggregateType {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.AggregateType), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (aggregateConfiguration)
	lengthInBits += m.AggregateConfiguration.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_ReadProcessedDetails) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ReadProcessedDetails) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__readProcessedDetails ReadProcessedDetails, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ReadProcessedDetails"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ReadProcessedDetails")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	startTime, err := ReadSimpleField(ctx, "startTime", ReadSignedLong(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'startTime' field"))
	}
	m.StartTime = startTime

	endTime, err := ReadSimpleField(ctx, "endTime", ReadSignedLong(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'endTime' field"))
	}
	m.EndTime = endTime

	processingInterval, err := ReadSimpleField(ctx, "processingInterval", ReadDouble(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'processingInterval' field"))
	}
	m.ProcessingInterval = processingInterval

	noOfAggregateType, err := ReadImplicitField[int32](ctx, "noOfAggregateType", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfAggregateType' field"))
	}
	_ = noOfAggregateType

	aggregateType, err := ReadCountArrayField[NodeId](ctx, "aggregateType", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer), uint64(noOfAggregateType))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'aggregateType' field"))
	}
	m.AggregateType = aggregateType

	aggregateConfiguration, err := ReadSimpleField[AggregateConfiguration](ctx, "aggregateConfiguration", ReadComplex[AggregateConfiguration](ExtensionObjectDefinitionParseWithBufferProducer[AggregateConfiguration]((int32)(int32(950))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'aggregateConfiguration' field"))
	}
	m.AggregateConfiguration = aggregateConfiguration

	if closeErr := readBuffer.CloseContext("ReadProcessedDetails"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ReadProcessedDetails")
	}

	return m, nil
}

func (m *_ReadProcessedDetails) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ReadProcessedDetails) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ReadProcessedDetails"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ReadProcessedDetails")
		}

		if err := WriteSimpleField[int64](ctx, "startTime", m.GetStartTime(), WriteSignedLong(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'startTime' field")
		}

		if err := WriteSimpleField[int64](ctx, "endTime", m.GetEndTime(), WriteSignedLong(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'endTime' field")
		}

		if err := WriteSimpleField[float64](ctx, "processingInterval", m.GetProcessingInterval(), WriteDouble(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'processingInterval' field")
		}
		noOfAggregateType := int32(utils.InlineIf(bool((m.GetAggregateType()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetAggregateType()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfAggregateType", noOfAggregateType, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfAggregateType' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "aggregateType", m.GetAggregateType(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'aggregateType' field")
		}

		if err := WriteSimpleField[AggregateConfiguration](ctx, "aggregateConfiguration", m.GetAggregateConfiguration(), WriteComplex[AggregateConfiguration](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'aggregateConfiguration' field")
		}

		if popErr := writeBuffer.PopContext("ReadProcessedDetails"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ReadProcessedDetails")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ReadProcessedDetails) IsReadProcessedDetails() {}

func (m *_ReadProcessedDetails) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ReadProcessedDetails) deepCopy() *_ReadProcessedDetails {
	if m == nil {
		return nil
	}
	_ReadProcessedDetailsCopy := &_ReadProcessedDetails{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.StartTime,
		m.EndTime,
		m.ProcessingInterval,
		utils.DeepCopySlice[NodeId, NodeId](m.AggregateType),
		m.AggregateConfiguration.DeepCopy().(AggregateConfiguration),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _ReadProcessedDetailsCopy
}

func (m *_ReadProcessedDetails) String() string {
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
