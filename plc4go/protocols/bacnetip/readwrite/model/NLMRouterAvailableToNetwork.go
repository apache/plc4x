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

// NLMRouterAvailableToNetwork is the corresponding interface of NLMRouterAvailableToNetwork
type NLMRouterAvailableToNetwork interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	NLM
	// GetDestinationNetworkAddresses returns DestinationNetworkAddresses (property field)
	GetDestinationNetworkAddresses() []uint16
	// IsNLMRouterAvailableToNetwork is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsNLMRouterAvailableToNetwork()
	// CreateBuilder creates a NLMRouterAvailableToNetworkBuilder
	CreateNLMRouterAvailableToNetworkBuilder() NLMRouterAvailableToNetworkBuilder
}

// _NLMRouterAvailableToNetwork is the data-structure of this message
type _NLMRouterAvailableToNetwork struct {
	NLMContract
	DestinationNetworkAddresses []uint16
}

var _ NLMRouterAvailableToNetwork = (*_NLMRouterAvailableToNetwork)(nil)
var _ NLMRequirements = (*_NLMRouterAvailableToNetwork)(nil)

// NewNLMRouterAvailableToNetwork factory function for _NLMRouterAvailableToNetwork
func NewNLMRouterAvailableToNetwork(destinationNetworkAddresses []uint16, apduLength uint16) *_NLMRouterAvailableToNetwork {
	_result := &_NLMRouterAvailableToNetwork{
		NLMContract:                 NewNLM(apduLength),
		DestinationNetworkAddresses: destinationNetworkAddresses,
	}
	_result.NLMContract.(*_NLM)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// NLMRouterAvailableToNetworkBuilder is a builder for NLMRouterAvailableToNetwork
type NLMRouterAvailableToNetworkBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(destinationNetworkAddresses []uint16) NLMRouterAvailableToNetworkBuilder
	// WithDestinationNetworkAddresses adds DestinationNetworkAddresses (property field)
	WithDestinationNetworkAddresses(...uint16) NLMRouterAvailableToNetworkBuilder
	// Build builds the NLMRouterAvailableToNetwork or returns an error if something is wrong
	Build() (NLMRouterAvailableToNetwork, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() NLMRouterAvailableToNetwork
}

// NewNLMRouterAvailableToNetworkBuilder() creates a NLMRouterAvailableToNetworkBuilder
func NewNLMRouterAvailableToNetworkBuilder() NLMRouterAvailableToNetworkBuilder {
	return &_NLMRouterAvailableToNetworkBuilder{_NLMRouterAvailableToNetwork: new(_NLMRouterAvailableToNetwork)}
}

type _NLMRouterAvailableToNetworkBuilder struct {
	*_NLMRouterAvailableToNetwork

	parentBuilder *_NLMBuilder

	err *utils.MultiError
}

var _ (NLMRouterAvailableToNetworkBuilder) = (*_NLMRouterAvailableToNetworkBuilder)(nil)

func (b *_NLMRouterAvailableToNetworkBuilder) setParent(contract NLMContract) {
	b.NLMContract = contract
}

func (b *_NLMRouterAvailableToNetworkBuilder) WithMandatoryFields(destinationNetworkAddresses []uint16) NLMRouterAvailableToNetworkBuilder {
	return b.WithDestinationNetworkAddresses(destinationNetworkAddresses...)
}

func (b *_NLMRouterAvailableToNetworkBuilder) WithDestinationNetworkAddresses(destinationNetworkAddresses ...uint16) NLMRouterAvailableToNetworkBuilder {
	b.DestinationNetworkAddresses = destinationNetworkAddresses
	return b
}

func (b *_NLMRouterAvailableToNetworkBuilder) Build() (NLMRouterAvailableToNetwork, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._NLMRouterAvailableToNetwork.deepCopy(), nil
}

func (b *_NLMRouterAvailableToNetworkBuilder) MustBuild() NLMRouterAvailableToNetwork {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_NLMRouterAvailableToNetworkBuilder) Done() NLMBuilder {
	return b.parentBuilder
}

func (b *_NLMRouterAvailableToNetworkBuilder) buildForNLM() (NLM, error) {
	return b.Build()
}

func (b *_NLMRouterAvailableToNetworkBuilder) DeepCopy() any {
	_copy := b.CreateNLMRouterAvailableToNetworkBuilder().(*_NLMRouterAvailableToNetworkBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateNLMRouterAvailableToNetworkBuilder creates a NLMRouterAvailableToNetworkBuilder
func (b *_NLMRouterAvailableToNetwork) CreateNLMRouterAvailableToNetworkBuilder() NLMRouterAvailableToNetworkBuilder {
	if b == nil {
		return NewNLMRouterAvailableToNetworkBuilder()
	}
	return &_NLMRouterAvailableToNetworkBuilder{_NLMRouterAvailableToNetwork: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NLMRouterAvailableToNetwork) GetMessageType() uint8 {
	return 0x05
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NLMRouterAvailableToNetwork) GetParent() NLMContract {
	return m.NLMContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_NLMRouterAvailableToNetwork) GetDestinationNetworkAddresses() []uint16 {
	return m.DestinationNetworkAddresses
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastNLMRouterAvailableToNetwork(structType any) NLMRouterAvailableToNetwork {
	if casted, ok := structType.(NLMRouterAvailableToNetwork); ok {
		return casted
	}
	if casted, ok := structType.(*NLMRouterAvailableToNetwork); ok {
		return *casted
	}
	return nil
}

func (m *_NLMRouterAvailableToNetwork) GetTypeName() string {
	return "NLMRouterAvailableToNetwork"
}

func (m *_NLMRouterAvailableToNetwork) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.NLMContract.(*_NLM).GetLengthInBits(ctx))

	// Array field
	if len(m.DestinationNetworkAddresses) > 0 {
		lengthInBits += 16 * uint16(len(m.DestinationNetworkAddresses))
	}

	return lengthInBits
}

func (m *_NLMRouterAvailableToNetwork) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_NLMRouterAvailableToNetwork) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_NLM, apduLength uint16) (__nLMRouterAvailableToNetwork NLMRouterAvailableToNetwork, err error) {
	m.NLMContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NLMRouterAvailableToNetwork"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NLMRouterAvailableToNetwork")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	destinationNetworkAddresses, err := ReadLengthArrayField[uint16](ctx, "destinationNetworkAddresses", ReadUnsignedShort(readBuffer, uint8(16)), int(int32(apduLength)-int32(int32(1))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'destinationNetworkAddresses' field"))
	}
	m.DestinationNetworkAddresses = destinationNetworkAddresses

	if closeErr := readBuffer.CloseContext("NLMRouterAvailableToNetwork"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NLMRouterAvailableToNetwork")
	}

	return m, nil
}

func (m *_NLMRouterAvailableToNetwork) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NLMRouterAvailableToNetwork) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NLMRouterAvailableToNetwork"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NLMRouterAvailableToNetwork")
		}

		if err := WriteSimpleTypeArrayField(ctx, "destinationNetworkAddresses", m.GetDestinationNetworkAddresses(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'destinationNetworkAddresses' field")
		}

		if popErr := writeBuffer.PopContext("NLMRouterAvailableToNetwork"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NLMRouterAvailableToNetwork")
		}
		return nil
	}
	return m.NLMContract.(*_NLM).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_NLMRouterAvailableToNetwork) IsNLMRouterAvailableToNetwork() {}

func (m *_NLMRouterAvailableToNetwork) DeepCopy() any {
	return m.deepCopy()
}

func (m *_NLMRouterAvailableToNetwork) deepCopy() *_NLMRouterAvailableToNetwork {
	if m == nil {
		return nil
	}
	_NLMRouterAvailableToNetworkCopy := &_NLMRouterAvailableToNetwork{
		m.NLMContract.(*_NLM).deepCopy(),
		utils.DeepCopySlice[uint16, uint16](m.DestinationNetworkAddresses),
	}
	m.NLMContract.(*_NLM)._SubType = m
	return _NLMRouterAvailableToNetworkCopy
}

func (m *_NLMRouterAvailableToNetwork) String() string {
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
