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

// CBusPointToPointToMultiPointCommandNormal is the corresponding interface of CBusPointToPointToMultiPointCommandNormal
type CBusPointToPointToMultiPointCommandNormal interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	CBusPointToPointToMultiPointCommand
	// GetApplication returns Application (property field)
	GetApplication() ApplicationIdContainer
	// GetSalData returns SalData (property field)
	GetSalData() SALData
	// IsCBusPointToPointToMultiPointCommandNormal is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCBusPointToPointToMultiPointCommandNormal()
	// CreateBuilder creates a CBusPointToPointToMultiPointCommandNormalBuilder
	CreateCBusPointToPointToMultiPointCommandNormalBuilder() CBusPointToPointToMultiPointCommandNormalBuilder
}

// _CBusPointToPointToMultiPointCommandNormal is the data-structure of this message
type _CBusPointToPointToMultiPointCommandNormal struct {
	CBusPointToPointToMultiPointCommandContract
	Application ApplicationIdContainer
	SalData     SALData
}

var _ CBusPointToPointToMultiPointCommandNormal = (*_CBusPointToPointToMultiPointCommandNormal)(nil)
var _ CBusPointToPointToMultiPointCommandRequirements = (*_CBusPointToPointToMultiPointCommandNormal)(nil)

// NewCBusPointToPointToMultiPointCommandNormal factory function for _CBusPointToPointToMultiPointCommandNormal
func NewCBusPointToPointToMultiPointCommandNormal(bridgeAddress BridgeAddress, networkRoute NetworkRoute, peekedApplication byte, application ApplicationIdContainer, salData SALData, cBusOptions CBusOptions) *_CBusPointToPointToMultiPointCommandNormal {
	if salData == nil {
		panic("salData of type SALData for CBusPointToPointToMultiPointCommandNormal must not be nil")
	}
	_result := &_CBusPointToPointToMultiPointCommandNormal{
		CBusPointToPointToMultiPointCommandContract: NewCBusPointToPointToMultiPointCommand(bridgeAddress, networkRoute, peekedApplication, cBusOptions),
		Application: application,
		SalData:     salData,
	}
	_result.CBusPointToPointToMultiPointCommandContract.(*_CBusPointToPointToMultiPointCommand)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CBusPointToPointToMultiPointCommandNormalBuilder is a builder for CBusPointToPointToMultiPointCommandNormal
type CBusPointToPointToMultiPointCommandNormalBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(application ApplicationIdContainer, salData SALData) CBusPointToPointToMultiPointCommandNormalBuilder
	// WithApplication adds Application (property field)
	WithApplication(ApplicationIdContainer) CBusPointToPointToMultiPointCommandNormalBuilder
	// WithSalData adds SalData (property field)
	WithSalData(SALData) CBusPointToPointToMultiPointCommandNormalBuilder
	// WithSalDataBuilder adds SalData (property field) which is build by the builder
	WithSalDataBuilder(func(SALDataBuilder) SALDataBuilder) CBusPointToPointToMultiPointCommandNormalBuilder
	// Build builds the CBusPointToPointToMultiPointCommandNormal or returns an error if something is wrong
	Build() (CBusPointToPointToMultiPointCommandNormal, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CBusPointToPointToMultiPointCommandNormal
}

// NewCBusPointToPointToMultiPointCommandNormalBuilder() creates a CBusPointToPointToMultiPointCommandNormalBuilder
func NewCBusPointToPointToMultiPointCommandNormalBuilder() CBusPointToPointToMultiPointCommandNormalBuilder {
	return &_CBusPointToPointToMultiPointCommandNormalBuilder{_CBusPointToPointToMultiPointCommandNormal: new(_CBusPointToPointToMultiPointCommandNormal)}
}

type _CBusPointToPointToMultiPointCommandNormalBuilder struct {
	*_CBusPointToPointToMultiPointCommandNormal

	parentBuilder *_CBusPointToPointToMultiPointCommandBuilder

	err *utils.MultiError
}

var _ (CBusPointToPointToMultiPointCommandNormalBuilder) = (*_CBusPointToPointToMultiPointCommandNormalBuilder)(nil)

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) setParent(contract CBusPointToPointToMultiPointCommandContract) {
	b.CBusPointToPointToMultiPointCommandContract = contract
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) WithMandatoryFields(application ApplicationIdContainer, salData SALData) CBusPointToPointToMultiPointCommandNormalBuilder {
	return b.WithApplication(application).WithSalData(salData)
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) WithApplication(application ApplicationIdContainer) CBusPointToPointToMultiPointCommandNormalBuilder {
	b.Application = application
	return b
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) WithSalData(salData SALData) CBusPointToPointToMultiPointCommandNormalBuilder {
	b.SalData = salData
	return b
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) WithSalDataBuilder(builderSupplier func(SALDataBuilder) SALDataBuilder) CBusPointToPointToMultiPointCommandNormalBuilder {
	builder := builderSupplier(b.SalData.CreateSALDataBuilder())
	var err error
	b.SalData, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "SALDataBuilder failed"))
	}
	return b
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) Build() (CBusPointToPointToMultiPointCommandNormal, error) {
	if b.SalData == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'salData' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CBusPointToPointToMultiPointCommandNormal.deepCopy(), nil
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) MustBuild() CBusPointToPointToMultiPointCommandNormal {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) Done() CBusPointToPointToMultiPointCommandBuilder {
	return b.parentBuilder
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) buildForCBusPointToPointToMultiPointCommand() (CBusPointToPointToMultiPointCommand, error) {
	return b.Build()
}

func (b *_CBusPointToPointToMultiPointCommandNormalBuilder) DeepCopy() any {
	_copy := b.CreateCBusPointToPointToMultiPointCommandNormalBuilder().(*_CBusPointToPointToMultiPointCommandNormalBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCBusPointToPointToMultiPointCommandNormalBuilder creates a CBusPointToPointToMultiPointCommandNormalBuilder
func (b *_CBusPointToPointToMultiPointCommandNormal) CreateCBusPointToPointToMultiPointCommandNormalBuilder() CBusPointToPointToMultiPointCommandNormalBuilder {
	if b == nil {
		return NewCBusPointToPointToMultiPointCommandNormalBuilder()
	}
	return &_CBusPointToPointToMultiPointCommandNormalBuilder{_CBusPointToPointToMultiPointCommandNormal: b.deepCopy()}
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

func (m *_CBusPointToPointToMultiPointCommandNormal) GetParent() CBusPointToPointToMultiPointCommandContract {
	return m.CBusPointToPointToMultiPointCommandContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CBusPointToPointToMultiPointCommandNormal) GetApplication() ApplicationIdContainer {
	return m.Application
}

func (m *_CBusPointToPointToMultiPointCommandNormal) GetSalData() SALData {
	return m.SalData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastCBusPointToPointToMultiPointCommandNormal(structType any) CBusPointToPointToMultiPointCommandNormal {
	if casted, ok := structType.(CBusPointToPointToMultiPointCommandNormal); ok {
		return casted
	}
	if casted, ok := structType.(*CBusPointToPointToMultiPointCommandNormal); ok {
		return *casted
	}
	return nil
}

func (m *_CBusPointToPointToMultiPointCommandNormal) GetTypeName() string {
	return "CBusPointToPointToMultiPointCommandNormal"
}

func (m *_CBusPointToPointToMultiPointCommandNormal) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.CBusPointToPointToMultiPointCommandContract.(*_CBusPointToPointToMultiPointCommand).GetLengthInBits(ctx))

	// Simple field (application)
	lengthInBits += 8

	// Simple field (salData)
	lengthInBits += m.SalData.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_CBusPointToPointToMultiPointCommandNormal) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CBusPointToPointToMultiPointCommandNormal) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_CBusPointToPointToMultiPointCommand, cBusOptions CBusOptions) (__cBusPointToPointToMultiPointCommandNormal CBusPointToPointToMultiPointCommandNormal, err error) {
	m.CBusPointToPointToMultiPointCommandContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CBusPointToPointToMultiPointCommandNormal"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CBusPointToPointToMultiPointCommandNormal")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	application, err := ReadEnumField[ApplicationIdContainer](ctx, "application", "ApplicationIdContainer", ReadEnum(ApplicationIdContainerByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'application' field"))
	}
	m.Application = application

	salData, err := ReadSimpleField[SALData](ctx, "salData", ReadComplex[SALData](SALDataParseWithBufferProducer[SALData]((ApplicationId)(application.ApplicationId())), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'salData' field"))
	}
	m.SalData = salData

	if closeErr := readBuffer.CloseContext("CBusPointToPointToMultiPointCommandNormal"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CBusPointToPointToMultiPointCommandNormal")
	}

	return m, nil
}

func (m *_CBusPointToPointToMultiPointCommandNormal) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CBusPointToPointToMultiPointCommandNormal) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CBusPointToPointToMultiPointCommandNormal"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CBusPointToPointToMultiPointCommandNormal")
		}

		if err := WriteSimpleEnumField[ApplicationIdContainer](ctx, "application", "ApplicationIdContainer", m.GetApplication(), WriteEnum[ApplicationIdContainer, uint8](ApplicationIdContainer.GetValue, ApplicationIdContainer.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'application' field")
		}

		if err := WriteSimpleField[SALData](ctx, "salData", m.GetSalData(), WriteComplex[SALData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'salData' field")
		}

		if popErr := writeBuffer.PopContext("CBusPointToPointToMultiPointCommandNormal"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CBusPointToPointToMultiPointCommandNormal")
		}
		return nil
	}
	return m.CBusPointToPointToMultiPointCommandContract.(*_CBusPointToPointToMultiPointCommand).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CBusPointToPointToMultiPointCommandNormal) IsCBusPointToPointToMultiPointCommandNormal() {}

func (m *_CBusPointToPointToMultiPointCommandNormal) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CBusPointToPointToMultiPointCommandNormal) deepCopy() *_CBusPointToPointToMultiPointCommandNormal {
	if m == nil {
		return nil
	}
	_CBusPointToPointToMultiPointCommandNormalCopy := &_CBusPointToPointToMultiPointCommandNormal{
		m.CBusPointToPointToMultiPointCommandContract.(*_CBusPointToPointToMultiPointCommand).deepCopy(),
		m.Application,
		m.SalData.DeepCopy().(SALData),
	}
	m.CBusPointToPointToMultiPointCommandContract.(*_CBusPointToPointToMultiPointCommand)._SubType = m
	return _CBusPointToPointToMultiPointCommandNormalCopy
}

func (m *_CBusPointToPointToMultiPointCommandNormal) String() string {
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
