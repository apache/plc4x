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

// NLMInitializeRoutingTableAck is the corresponding interface of NLMInitializeRoutingTableAck
type NLMInitializeRoutingTableAck interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	NLM
	// GetNumberOfPorts returns NumberOfPorts (property field)
	GetNumberOfPorts() uint8
	// GetPortMappings returns PortMappings (property field)
	GetPortMappings() []NLMInitializeRoutingTablePortMapping
	// IsNLMInitializeRoutingTableAck is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsNLMInitializeRoutingTableAck()
	// CreateBuilder creates a NLMInitializeRoutingTableAckBuilder
	CreateNLMInitializeRoutingTableAckBuilder() NLMInitializeRoutingTableAckBuilder
}

// _NLMInitializeRoutingTableAck is the data-structure of this message
type _NLMInitializeRoutingTableAck struct {
	NLMContract
	NumberOfPorts uint8
	PortMappings  []NLMInitializeRoutingTablePortMapping
}

var _ NLMInitializeRoutingTableAck = (*_NLMInitializeRoutingTableAck)(nil)
var _ NLMRequirements = (*_NLMInitializeRoutingTableAck)(nil)

// NewNLMInitializeRoutingTableAck factory function for _NLMInitializeRoutingTableAck
func NewNLMInitializeRoutingTableAck(numberOfPorts uint8, portMappings []NLMInitializeRoutingTablePortMapping, apduLength uint16) *_NLMInitializeRoutingTableAck {
	_result := &_NLMInitializeRoutingTableAck{
		NLMContract:   NewNLM(apduLength),
		NumberOfPorts: numberOfPorts,
		PortMappings:  portMappings,
	}
	_result.NLMContract.(*_NLM)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// NLMInitializeRoutingTableAckBuilder is a builder for NLMInitializeRoutingTableAck
type NLMInitializeRoutingTableAckBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(numberOfPorts uint8, portMappings []NLMInitializeRoutingTablePortMapping) NLMInitializeRoutingTableAckBuilder
	// WithNumberOfPorts adds NumberOfPorts (property field)
	WithNumberOfPorts(uint8) NLMInitializeRoutingTableAckBuilder
	// WithPortMappings adds PortMappings (property field)
	WithPortMappings(...NLMInitializeRoutingTablePortMapping) NLMInitializeRoutingTableAckBuilder
	// Build builds the NLMInitializeRoutingTableAck or returns an error if something is wrong
	Build() (NLMInitializeRoutingTableAck, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() NLMInitializeRoutingTableAck
}

// NewNLMInitializeRoutingTableAckBuilder() creates a NLMInitializeRoutingTableAckBuilder
func NewNLMInitializeRoutingTableAckBuilder() NLMInitializeRoutingTableAckBuilder {
	return &_NLMInitializeRoutingTableAckBuilder{_NLMInitializeRoutingTableAck: new(_NLMInitializeRoutingTableAck)}
}

type _NLMInitializeRoutingTableAckBuilder struct {
	*_NLMInitializeRoutingTableAck

	parentBuilder *_NLMBuilder

	err *utils.MultiError
}

var _ (NLMInitializeRoutingTableAckBuilder) = (*_NLMInitializeRoutingTableAckBuilder)(nil)

func (b *_NLMInitializeRoutingTableAckBuilder) setParent(contract NLMContract) {
	b.NLMContract = contract
}

func (b *_NLMInitializeRoutingTableAckBuilder) WithMandatoryFields(numberOfPorts uint8, portMappings []NLMInitializeRoutingTablePortMapping) NLMInitializeRoutingTableAckBuilder {
	return b.WithNumberOfPorts(numberOfPorts).WithPortMappings(portMappings...)
}

func (b *_NLMInitializeRoutingTableAckBuilder) WithNumberOfPorts(numberOfPorts uint8) NLMInitializeRoutingTableAckBuilder {
	b.NumberOfPorts = numberOfPorts
	return b
}

func (b *_NLMInitializeRoutingTableAckBuilder) WithPortMappings(portMappings ...NLMInitializeRoutingTablePortMapping) NLMInitializeRoutingTableAckBuilder {
	b.PortMappings = portMappings
	return b
}

func (b *_NLMInitializeRoutingTableAckBuilder) Build() (NLMInitializeRoutingTableAck, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._NLMInitializeRoutingTableAck.deepCopy(), nil
}

func (b *_NLMInitializeRoutingTableAckBuilder) MustBuild() NLMInitializeRoutingTableAck {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_NLMInitializeRoutingTableAckBuilder) Done() NLMBuilder {
	return b.parentBuilder
}

func (b *_NLMInitializeRoutingTableAckBuilder) buildForNLM() (NLM, error) {
	return b.Build()
}

func (b *_NLMInitializeRoutingTableAckBuilder) DeepCopy() any {
	_copy := b.CreateNLMInitializeRoutingTableAckBuilder().(*_NLMInitializeRoutingTableAckBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateNLMInitializeRoutingTableAckBuilder creates a NLMInitializeRoutingTableAckBuilder
func (b *_NLMInitializeRoutingTableAck) CreateNLMInitializeRoutingTableAckBuilder() NLMInitializeRoutingTableAckBuilder {
	if b == nil {
		return NewNLMInitializeRoutingTableAckBuilder()
	}
	return &_NLMInitializeRoutingTableAckBuilder{_NLMInitializeRoutingTableAck: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NLMInitializeRoutingTableAck) GetMessageType() uint8 {
	return 0x07
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NLMInitializeRoutingTableAck) GetParent() NLMContract {
	return m.NLMContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_NLMInitializeRoutingTableAck) GetNumberOfPorts() uint8 {
	return m.NumberOfPorts
}

func (m *_NLMInitializeRoutingTableAck) GetPortMappings() []NLMInitializeRoutingTablePortMapping {
	return m.PortMappings
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastNLMInitializeRoutingTableAck(structType any) NLMInitializeRoutingTableAck {
	if casted, ok := structType.(NLMInitializeRoutingTableAck); ok {
		return casted
	}
	if casted, ok := structType.(*NLMInitializeRoutingTableAck); ok {
		return *casted
	}
	return nil
}

func (m *_NLMInitializeRoutingTableAck) GetTypeName() string {
	return "NLMInitializeRoutingTableAck"
}

func (m *_NLMInitializeRoutingTableAck) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.NLMContract.(*_NLM).GetLengthInBits(ctx))

	// Simple field (numberOfPorts)
	lengthInBits += 8

	// Array field
	if len(m.PortMappings) > 0 {
		for _curItem, element := range m.PortMappings {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.PortMappings), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_NLMInitializeRoutingTableAck) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_NLMInitializeRoutingTableAck) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_NLM, apduLength uint16) (__nLMInitializeRoutingTableAck NLMInitializeRoutingTableAck, err error) {
	m.NLMContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NLMInitializeRoutingTableAck"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NLMInitializeRoutingTableAck")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	numberOfPorts, err := ReadSimpleField(ctx, "numberOfPorts", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfPorts' field"))
	}
	m.NumberOfPorts = numberOfPorts

	portMappings, err := ReadCountArrayField[NLMInitializeRoutingTablePortMapping](ctx, "portMappings", ReadComplex[NLMInitializeRoutingTablePortMapping](NLMInitializeRoutingTablePortMappingParseWithBuffer, readBuffer), uint64(numberOfPorts))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'portMappings' field"))
	}
	m.PortMappings = portMappings

	if closeErr := readBuffer.CloseContext("NLMInitializeRoutingTableAck"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NLMInitializeRoutingTableAck")
	}

	return m, nil
}

func (m *_NLMInitializeRoutingTableAck) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NLMInitializeRoutingTableAck) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NLMInitializeRoutingTableAck"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NLMInitializeRoutingTableAck")
		}

		if err := WriteSimpleField[uint8](ctx, "numberOfPorts", m.GetNumberOfPorts(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'numberOfPorts' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "portMappings", m.GetPortMappings(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'portMappings' field")
		}

		if popErr := writeBuffer.PopContext("NLMInitializeRoutingTableAck"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NLMInitializeRoutingTableAck")
		}
		return nil
	}
	return m.NLMContract.(*_NLM).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_NLMInitializeRoutingTableAck) IsNLMInitializeRoutingTableAck() {}

func (m *_NLMInitializeRoutingTableAck) DeepCopy() any {
	return m.deepCopy()
}

func (m *_NLMInitializeRoutingTableAck) deepCopy() *_NLMInitializeRoutingTableAck {
	if m == nil {
		return nil
	}
	_NLMInitializeRoutingTableAckCopy := &_NLMInitializeRoutingTableAck{
		m.NLMContract.(*_NLM).deepCopy(),
		m.NumberOfPorts,
		utils.DeepCopySlice[NLMInitializeRoutingTablePortMapping, NLMInitializeRoutingTablePortMapping](m.PortMappings),
	}
	m.NLMContract.(*_NLM)._SubType = m
	return _NLMInitializeRoutingTableAckCopy
}

func (m *_NLMInitializeRoutingTableAck) String() string {
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
