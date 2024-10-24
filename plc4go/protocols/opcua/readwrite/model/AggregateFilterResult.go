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

// AggregateFilterResult is the corresponding interface of AggregateFilterResult
type AggregateFilterResult interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetRevisedStartTime returns RevisedStartTime (property field)
	GetRevisedStartTime() int64
	// GetRevisedProcessingInterval returns RevisedProcessingInterval (property field)
	GetRevisedProcessingInterval() float64
	// GetRevisedAggregateConfiguration returns RevisedAggregateConfiguration (property field)
	GetRevisedAggregateConfiguration() AggregateConfiguration
	// IsAggregateFilterResult is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsAggregateFilterResult()
	// CreateBuilder creates a AggregateFilterResultBuilder
	CreateAggregateFilterResultBuilder() AggregateFilterResultBuilder
}

// _AggregateFilterResult is the data-structure of this message
type _AggregateFilterResult struct {
	ExtensionObjectDefinitionContract
	RevisedStartTime              int64
	RevisedProcessingInterval     float64
	RevisedAggregateConfiguration AggregateConfiguration
}

var _ AggregateFilterResult = (*_AggregateFilterResult)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_AggregateFilterResult)(nil)

// NewAggregateFilterResult factory function for _AggregateFilterResult
func NewAggregateFilterResult(revisedStartTime int64, revisedProcessingInterval float64, revisedAggregateConfiguration AggregateConfiguration) *_AggregateFilterResult {
	if revisedAggregateConfiguration == nil {
		panic("revisedAggregateConfiguration of type AggregateConfiguration for AggregateFilterResult must not be nil")
	}
	_result := &_AggregateFilterResult{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		RevisedStartTime:                  revisedStartTime,
		RevisedProcessingInterval:         revisedProcessingInterval,
		RevisedAggregateConfiguration:     revisedAggregateConfiguration,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// AggregateFilterResultBuilder is a builder for AggregateFilterResult
type AggregateFilterResultBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(revisedStartTime int64, revisedProcessingInterval float64, revisedAggregateConfiguration AggregateConfiguration) AggregateFilterResultBuilder
	// WithRevisedStartTime adds RevisedStartTime (property field)
	WithRevisedStartTime(int64) AggregateFilterResultBuilder
	// WithRevisedProcessingInterval adds RevisedProcessingInterval (property field)
	WithRevisedProcessingInterval(float64) AggregateFilterResultBuilder
	// WithRevisedAggregateConfiguration adds RevisedAggregateConfiguration (property field)
	WithRevisedAggregateConfiguration(AggregateConfiguration) AggregateFilterResultBuilder
	// WithRevisedAggregateConfigurationBuilder adds RevisedAggregateConfiguration (property field) which is build by the builder
	WithRevisedAggregateConfigurationBuilder(func(AggregateConfigurationBuilder) AggregateConfigurationBuilder) AggregateFilterResultBuilder
	// Build builds the AggregateFilterResult or returns an error if something is wrong
	Build() (AggregateFilterResult, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() AggregateFilterResult
}

// NewAggregateFilterResultBuilder() creates a AggregateFilterResultBuilder
func NewAggregateFilterResultBuilder() AggregateFilterResultBuilder {
	return &_AggregateFilterResultBuilder{_AggregateFilterResult: new(_AggregateFilterResult)}
}

type _AggregateFilterResultBuilder struct {
	*_AggregateFilterResult

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (AggregateFilterResultBuilder) = (*_AggregateFilterResultBuilder)(nil)

func (b *_AggregateFilterResultBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_AggregateFilterResultBuilder) WithMandatoryFields(revisedStartTime int64, revisedProcessingInterval float64, revisedAggregateConfiguration AggregateConfiguration) AggregateFilterResultBuilder {
	return b.WithRevisedStartTime(revisedStartTime).WithRevisedProcessingInterval(revisedProcessingInterval).WithRevisedAggregateConfiguration(revisedAggregateConfiguration)
}

func (b *_AggregateFilterResultBuilder) WithRevisedStartTime(revisedStartTime int64) AggregateFilterResultBuilder {
	b.RevisedStartTime = revisedStartTime
	return b
}

func (b *_AggregateFilterResultBuilder) WithRevisedProcessingInterval(revisedProcessingInterval float64) AggregateFilterResultBuilder {
	b.RevisedProcessingInterval = revisedProcessingInterval
	return b
}

func (b *_AggregateFilterResultBuilder) WithRevisedAggregateConfiguration(revisedAggregateConfiguration AggregateConfiguration) AggregateFilterResultBuilder {
	b.RevisedAggregateConfiguration = revisedAggregateConfiguration
	return b
}

func (b *_AggregateFilterResultBuilder) WithRevisedAggregateConfigurationBuilder(builderSupplier func(AggregateConfigurationBuilder) AggregateConfigurationBuilder) AggregateFilterResultBuilder {
	builder := builderSupplier(b.RevisedAggregateConfiguration.CreateAggregateConfigurationBuilder())
	var err error
	b.RevisedAggregateConfiguration, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "AggregateConfigurationBuilder failed"))
	}
	return b
}

func (b *_AggregateFilterResultBuilder) Build() (AggregateFilterResult, error) {
	if b.RevisedAggregateConfiguration == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'revisedAggregateConfiguration' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._AggregateFilterResult.deepCopy(), nil
}

func (b *_AggregateFilterResultBuilder) MustBuild() AggregateFilterResult {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_AggregateFilterResultBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_AggregateFilterResultBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_AggregateFilterResultBuilder) DeepCopy() any {
	_copy := b.CreateAggregateFilterResultBuilder().(*_AggregateFilterResultBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateAggregateFilterResultBuilder creates a AggregateFilterResultBuilder
func (b *_AggregateFilterResult) CreateAggregateFilterResultBuilder() AggregateFilterResultBuilder {
	if b == nil {
		return NewAggregateFilterResultBuilder()
	}
	return &_AggregateFilterResultBuilder{_AggregateFilterResult: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_AggregateFilterResult) GetExtensionId() int32 {
	return int32(739)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AggregateFilterResult) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AggregateFilterResult) GetRevisedStartTime() int64 {
	return m.RevisedStartTime
}

func (m *_AggregateFilterResult) GetRevisedProcessingInterval() float64 {
	return m.RevisedProcessingInterval
}

func (m *_AggregateFilterResult) GetRevisedAggregateConfiguration() AggregateConfiguration {
	return m.RevisedAggregateConfiguration
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastAggregateFilterResult(structType any) AggregateFilterResult {
	if casted, ok := structType.(AggregateFilterResult); ok {
		return casted
	}
	if casted, ok := structType.(*AggregateFilterResult); ok {
		return *casted
	}
	return nil
}

func (m *_AggregateFilterResult) GetTypeName() string {
	return "AggregateFilterResult"
}

func (m *_AggregateFilterResult) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (revisedStartTime)
	lengthInBits += 64

	// Simple field (revisedProcessingInterval)
	lengthInBits += 64

	// Simple field (revisedAggregateConfiguration)
	lengthInBits += m.RevisedAggregateConfiguration.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_AggregateFilterResult) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_AggregateFilterResult) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__aggregateFilterResult AggregateFilterResult, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AggregateFilterResult"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AggregateFilterResult")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	revisedStartTime, err := ReadSimpleField(ctx, "revisedStartTime", ReadSignedLong(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedStartTime' field"))
	}
	m.RevisedStartTime = revisedStartTime

	revisedProcessingInterval, err := ReadSimpleField(ctx, "revisedProcessingInterval", ReadDouble(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedProcessingInterval' field"))
	}
	m.RevisedProcessingInterval = revisedProcessingInterval

	revisedAggregateConfiguration, err := ReadSimpleField[AggregateConfiguration](ctx, "revisedAggregateConfiguration", ReadComplex[AggregateConfiguration](ExtensionObjectDefinitionParseWithBufferProducer[AggregateConfiguration]((int32)(int32(950))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedAggregateConfiguration' field"))
	}
	m.RevisedAggregateConfiguration = revisedAggregateConfiguration

	if closeErr := readBuffer.CloseContext("AggregateFilterResult"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AggregateFilterResult")
	}

	return m, nil
}

func (m *_AggregateFilterResult) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AggregateFilterResult) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AggregateFilterResult"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AggregateFilterResult")
		}

		if err := WriteSimpleField[int64](ctx, "revisedStartTime", m.GetRevisedStartTime(), WriteSignedLong(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedStartTime' field")
		}

		if err := WriteSimpleField[float64](ctx, "revisedProcessingInterval", m.GetRevisedProcessingInterval(), WriteDouble(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedProcessingInterval' field")
		}

		if err := WriteSimpleField[AggregateConfiguration](ctx, "revisedAggregateConfiguration", m.GetRevisedAggregateConfiguration(), WriteComplex[AggregateConfiguration](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedAggregateConfiguration' field")
		}

		if popErr := writeBuffer.PopContext("AggregateFilterResult"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AggregateFilterResult")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_AggregateFilterResult) IsAggregateFilterResult() {}

func (m *_AggregateFilterResult) DeepCopy() any {
	return m.deepCopy()
}

func (m *_AggregateFilterResult) deepCopy() *_AggregateFilterResult {
	if m == nil {
		return nil
	}
	_AggregateFilterResultCopy := &_AggregateFilterResult{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.RevisedStartTime,
		m.RevisedProcessingInterval,
		m.RevisedAggregateConfiguration.DeepCopy().(AggregateConfiguration),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _AggregateFilterResultCopy
}

func (m *_AggregateFilterResult) String() string {
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
