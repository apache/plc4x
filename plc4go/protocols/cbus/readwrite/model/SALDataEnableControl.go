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

// SALDataEnableControl is the corresponding interface of SALDataEnableControl
type SALDataEnableControl interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	SALData
	// GetEnableControlData returns EnableControlData (property field)
	GetEnableControlData() EnableControlData
	// IsSALDataEnableControl is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSALDataEnableControl()
	// CreateBuilder creates a SALDataEnableControlBuilder
	CreateSALDataEnableControlBuilder() SALDataEnableControlBuilder
}

// _SALDataEnableControl is the data-structure of this message
type _SALDataEnableControl struct {
	SALDataContract
	EnableControlData EnableControlData
}

var _ SALDataEnableControl = (*_SALDataEnableControl)(nil)
var _ SALDataRequirements = (*_SALDataEnableControl)(nil)

// NewSALDataEnableControl factory function for _SALDataEnableControl
func NewSALDataEnableControl(salData SALData, enableControlData EnableControlData) *_SALDataEnableControl {
	if enableControlData == nil {
		panic("enableControlData of type EnableControlData for SALDataEnableControl must not be nil")
	}
	_result := &_SALDataEnableControl{
		SALDataContract:   NewSALData(salData),
		EnableControlData: enableControlData,
	}
	_result.SALDataContract.(*_SALData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// SALDataEnableControlBuilder is a builder for SALDataEnableControl
type SALDataEnableControlBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(enableControlData EnableControlData) SALDataEnableControlBuilder
	// WithEnableControlData adds EnableControlData (property field)
	WithEnableControlData(EnableControlData) SALDataEnableControlBuilder
	// WithEnableControlDataBuilder adds EnableControlData (property field) which is build by the builder
	WithEnableControlDataBuilder(func(EnableControlDataBuilder) EnableControlDataBuilder) SALDataEnableControlBuilder
	// Build builds the SALDataEnableControl or returns an error if something is wrong
	Build() (SALDataEnableControl, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() SALDataEnableControl
}

// NewSALDataEnableControlBuilder() creates a SALDataEnableControlBuilder
func NewSALDataEnableControlBuilder() SALDataEnableControlBuilder {
	return &_SALDataEnableControlBuilder{_SALDataEnableControl: new(_SALDataEnableControl)}
}

type _SALDataEnableControlBuilder struct {
	*_SALDataEnableControl

	parentBuilder *_SALDataBuilder

	err *utils.MultiError
}

var _ (SALDataEnableControlBuilder) = (*_SALDataEnableControlBuilder)(nil)

func (b *_SALDataEnableControlBuilder) setParent(contract SALDataContract) {
	b.SALDataContract = contract
}

func (b *_SALDataEnableControlBuilder) WithMandatoryFields(enableControlData EnableControlData) SALDataEnableControlBuilder {
	return b.WithEnableControlData(enableControlData)
}

func (b *_SALDataEnableControlBuilder) WithEnableControlData(enableControlData EnableControlData) SALDataEnableControlBuilder {
	b.EnableControlData = enableControlData
	return b
}

func (b *_SALDataEnableControlBuilder) WithEnableControlDataBuilder(builderSupplier func(EnableControlDataBuilder) EnableControlDataBuilder) SALDataEnableControlBuilder {
	builder := builderSupplier(b.EnableControlData.CreateEnableControlDataBuilder())
	var err error
	b.EnableControlData, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "EnableControlDataBuilder failed"))
	}
	return b
}

func (b *_SALDataEnableControlBuilder) Build() (SALDataEnableControl, error) {
	if b.EnableControlData == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'enableControlData' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._SALDataEnableControl.deepCopy(), nil
}

func (b *_SALDataEnableControlBuilder) MustBuild() SALDataEnableControl {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_SALDataEnableControlBuilder) Done() SALDataBuilder {
	return b.parentBuilder
}

func (b *_SALDataEnableControlBuilder) buildForSALData() (SALData, error) {
	return b.Build()
}

func (b *_SALDataEnableControlBuilder) DeepCopy() any {
	_copy := b.CreateSALDataEnableControlBuilder().(*_SALDataEnableControlBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateSALDataEnableControlBuilder creates a SALDataEnableControlBuilder
func (b *_SALDataEnableControl) CreateSALDataEnableControlBuilder() SALDataEnableControlBuilder {
	if b == nil {
		return NewSALDataEnableControlBuilder()
	}
	return &_SALDataEnableControlBuilder{_SALDataEnableControl: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SALDataEnableControl) GetApplicationId() ApplicationId {
	return ApplicationId_ENABLE_CONTROL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SALDataEnableControl) GetParent() SALDataContract {
	return m.SALDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SALDataEnableControl) GetEnableControlData() EnableControlData {
	return m.EnableControlData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastSALDataEnableControl(structType any) SALDataEnableControl {
	if casted, ok := structType.(SALDataEnableControl); ok {
		return casted
	}
	if casted, ok := structType.(*SALDataEnableControl); ok {
		return *casted
	}
	return nil
}

func (m *_SALDataEnableControl) GetTypeName() string {
	return "SALDataEnableControl"
}

func (m *_SALDataEnableControl) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.SALDataContract.(*_SALData).GetLengthInBits(ctx))

	// Simple field (enableControlData)
	lengthInBits += m.EnableControlData.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_SALDataEnableControl) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_SALDataEnableControl) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_SALData, applicationId ApplicationId) (__sALDataEnableControl SALDataEnableControl, err error) {
	m.SALDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SALDataEnableControl"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SALDataEnableControl")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	enableControlData, err := ReadSimpleField[EnableControlData](ctx, "enableControlData", ReadComplex[EnableControlData](EnableControlDataParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'enableControlData' field"))
	}
	m.EnableControlData = enableControlData

	if closeErr := readBuffer.CloseContext("SALDataEnableControl"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SALDataEnableControl")
	}

	return m, nil
}

func (m *_SALDataEnableControl) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SALDataEnableControl) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SALDataEnableControl"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SALDataEnableControl")
		}

		if err := WriteSimpleField[EnableControlData](ctx, "enableControlData", m.GetEnableControlData(), WriteComplex[EnableControlData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'enableControlData' field")
		}

		if popErr := writeBuffer.PopContext("SALDataEnableControl"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SALDataEnableControl")
		}
		return nil
	}
	return m.SALDataContract.(*_SALData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SALDataEnableControl) IsSALDataEnableControl() {}

func (m *_SALDataEnableControl) DeepCopy() any {
	return m.deepCopy()
}

func (m *_SALDataEnableControl) deepCopy() *_SALDataEnableControl {
	if m == nil {
		return nil
	}
	_SALDataEnableControlCopy := &_SALDataEnableControl{
		m.SALDataContract.(*_SALData).deepCopy(),
		m.EnableControlData.DeepCopy().(EnableControlData),
	}
	m.SALDataContract.(*_SALData)._SubType = m
	return _SALDataEnableControlCopy
}

func (m *_SALDataEnableControl) String() string {
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
