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

// AdsWriteControlResponse is the corresponding interface of AdsWriteControlResponse
type AdsWriteControlResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	AmsPacket
	// GetResult returns Result (property field)
	GetResult() ReturnCode
	// IsAdsWriteControlResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsAdsWriteControlResponse()
	// CreateBuilder creates a AdsWriteControlResponseBuilder
	CreateAdsWriteControlResponseBuilder() AdsWriteControlResponseBuilder
}

// _AdsWriteControlResponse is the data-structure of this message
type _AdsWriteControlResponse struct {
	AmsPacketContract
	Result ReturnCode
}

var _ AdsWriteControlResponse = (*_AdsWriteControlResponse)(nil)
var _ AmsPacketRequirements = (*_AdsWriteControlResponse)(nil)

// NewAdsWriteControlResponse factory function for _AdsWriteControlResponse
func NewAdsWriteControlResponse(targetAmsNetId AmsNetId, targetAmsPort uint16, sourceAmsNetId AmsNetId, sourceAmsPort uint16, errorCode uint32, invokeId uint32, result ReturnCode) *_AdsWriteControlResponse {
	_result := &_AdsWriteControlResponse{
		AmsPacketContract: NewAmsPacket(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, errorCode, invokeId),
		Result:            result,
	}
	_result.AmsPacketContract.(*_AmsPacket)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// AdsWriteControlResponseBuilder is a builder for AdsWriteControlResponse
type AdsWriteControlResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(result ReturnCode) AdsWriteControlResponseBuilder
	// WithResult adds Result (property field)
	WithResult(ReturnCode) AdsWriteControlResponseBuilder
	// Build builds the AdsWriteControlResponse or returns an error if something is wrong
	Build() (AdsWriteControlResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() AdsWriteControlResponse
}

// NewAdsWriteControlResponseBuilder() creates a AdsWriteControlResponseBuilder
func NewAdsWriteControlResponseBuilder() AdsWriteControlResponseBuilder {
	return &_AdsWriteControlResponseBuilder{_AdsWriteControlResponse: new(_AdsWriteControlResponse)}
}

type _AdsWriteControlResponseBuilder struct {
	*_AdsWriteControlResponse

	parentBuilder *_AmsPacketBuilder

	err *utils.MultiError
}

var _ (AdsWriteControlResponseBuilder) = (*_AdsWriteControlResponseBuilder)(nil)

func (b *_AdsWriteControlResponseBuilder) setParent(contract AmsPacketContract) {
	b.AmsPacketContract = contract
}

func (b *_AdsWriteControlResponseBuilder) WithMandatoryFields(result ReturnCode) AdsWriteControlResponseBuilder {
	return b.WithResult(result)
}

func (b *_AdsWriteControlResponseBuilder) WithResult(result ReturnCode) AdsWriteControlResponseBuilder {
	b.Result = result
	return b
}

func (b *_AdsWriteControlResponseBuilder) Build() (AdsWriteControlResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._AdsWriteControlResponse.deepCopy(), nil
}

func (b *_AdsWriteControlResponseBuilder) MustBuild() AdsWriteControlResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_AdsWriteControlResponseBuilder) Done() AmsPacketBuilder {
	return b.parentBuilder
}

func (b *_AdsWriteControlResponseBuilder) buildForAmsPacket() (AmsPacket, error) {
	return b.Build()
}

func (b *_AdsWriteControlResponseBuilder) DeepCopy() any {
	_copy := b.CreateAdsWriteControlResponseBuilder().(*_AdsWriteControlResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateAdsWriteControlResponseBuilder creates a AdsWriteControlResponseBuilder
func (b *_AdsWriteControlResponse) CreateAdsWriteControlResponseBuilder() AdsWriteControlResponseBuilder {
	if b == nil {
		return NewAdsWriteControlResponseBuilder()
	}
	return &_AdsWriteControlResponseBuilder{_AdsWriteControlResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_AdsWriteControlResponse) GetCommandId() CommandId {
	return CommandId_ADS_WRITE_CONTROL
}

func (m *_AdsWriteControlResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AdsWriteControlResponse) GetParent() AmsPacketContract {
	return m.AmsPacketContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AdsWriteControlResponse) GetResult() ReturnCode {
	return m.Result
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastAdsWriteControlResponse(structType any) AdsWriteControlResponse {
	if casted, ok := structType.(AdsWriteControlResponse); ok {
		return casted
	}
	if casted, ok := structType.(*AdsWriteControlResponse); ok {
		return *casted
	}
	return nil
}

func (m *_AdsWriteControlResponse) GetTypeName() string {
	return "AdsWriteControlResponse"
}

func (m *_AdsWriteControlResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.AmsPacketContract.(*_AmsPacket).GetLengthInBits(ctx))

	// Simple field (result)
	lengthInBits += 32

	return lengthInBits
}

func (m *_AdsWriteControlResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_AdsWriteControlResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_AmsPacket) (__adsWriteControlResponse AdsWriteControlResponse, err error) {
	m.AmsPacketContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AdsWriteControlResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AdsWriteControlResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	result, err := ReadEnumField[ReturnCode](ctx, "result", "ReturnCode", ReadEnum(ReturnCodeByValue, ReadUnsignedInt(readBuffer, uint8(32))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'result' field"))
	}
	m.Result = result

	if closeErr := readBuffer.CloseContext("AdsWriteControlResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AdsWriteControlResponse")
	}

	return m, nil
}

func (m *_AdsWriteControlResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AdsWriteControlResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AdsWriteControlResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AdsWriteControlResponse")
		}

		if err := WriteSimpleEnumField[ReturnCode](ctx, "result", "ReturnCode", m.GetResult(), WriteEnum[ReturnCode, uint32](ReturnCode.GetValue, ReturnCode.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'result' field")
		}

		if popErr := writeBuffer.PopContext("AdsWriteControlResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AdsWriteControlResponse")
		}
		return nil
	}
	return m.AmsPacketContract.(*_AmsPacket).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_AdsWriteControlResponse) IsAdsWriteControlResponse() {}

func (m *_AdsWriteControlResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_AdsWriteControlResponse) deepCopy() *_AdsWriteControlResponse {
	if m == nil {
		return nil
	}
	_AdsWriteControlResponseCopy := &_AdsWriteControlResponse{
		m.AmsPacketContract.(*_AmsPacket).deepCopy(),
		m.Result,
	}
	m.AmsPacketContract.(*_AmsPacket)._SubType = m
	return _AdsWriteControlResponseCopy
}

func (m *_AdsWriteControlResponse) String() string {
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
