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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// DescriptionRequest is the corresponding interface of DescriptionRequest
type DescriptionRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	KnxNetIpMessage
	// GetHpaiControlEndpoint returns HpaiControlEndpoint (property field)
	GetHpaiControlEndpoint() HPAIControlEndpoint
}

// DescriptionRequestExactly can be used when we want exactly this type and not a type which fulfills DescriptionRequest.
// This is useful for switch cases.
type DescriptionRequestExactly interface {
	DescriptionRequest
	isDescriptionRequest() bool
}

// _DescriptionRequest is the data-structure of this message
type _DescriptionRequest struct {
	*_KnxNetIpMessage
	HpaiControlEndpoint HPAIControlEndpoint
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_DescriptionRequest) GetMsgType() uint16 {
	return 0x0203
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_DescriptionRequest) InitializeParent(parent KnxNetIpMessage) {}

func (m *_DescriptionRequest) GetParent() KnxNetIpMessage {
	return m._KnxNetIpMessage
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_DescriptionRequest) GetHpaiControlEndpoint() HPAIControlEndpoint {
	return m.HpaiControlEndpoint
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewDescriptionRequest factory function for _DescriptionRequest
func NewDescriptionRequest(hpaiControlEndpoint HPAIControlEndpoint) *_DescriptionRequest {
	_result := &_DescriptionRequest{
		HpaiControlEndpoint: hpaiControlEndpoint,
		_KnxNetIpMessage:    NewKnxNetIpMessage(),
	}
	_result._KnxNetIpMessage._KnxNetIpMessageChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastDescriptionRequest(structType any) DescriptionRequest {
	if casted, ok := structType.(DescriptionRequest); ok {
		return casted
	}
	if casted, ok := structType.(*DescriptionRequest); ok {
		return *casted
	}
	return nil
}

func (m *_DescriptionRequest) GetTypeName() string {
	return "DescriptionRequest"
}

func (m *_DescriptionRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (hpaiControlEndpoint)
	lengthInBits += m.HpaiControlEndpoint.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_DescriptionRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func DescriptionRequestParse(ctx context.Context, theBytes []byte) (DescriptionRequest, error) {
	return DescriptionRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes, utils.WithByteOrderForReadBufferByteBased(binary.BigEndian)))
}

func DescriptionRequestParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (DescriptionRequest, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (DescriptionRequest, error) {
		return DescriptionRequestParseWithBuffer(ctx, readBuffer)
	}
}

func DescriptionRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (DescriptionRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("DescriptionRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for DescriptionRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	hpaiControlEndpoint, err := ReadSimpleField[HPAIControlEndpoint](ctx, "hpaiControlEndpoint", ReadComplex[HPAIControlEndpoint](HPAIControlEndpointParseWithBuffer, readBuffer), codegen.WithByteOrder(binary.BigEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'hpaiControlEndpoint' field"))
	}

	if closeErr := readBuffer.CloseContext("DescriptionRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for DescriptionRequest")
	}

	// Create a partially initialized instance
	_child := &_DescriptionRequest{
		_KnxNetIpMessage:    &_KnxNetIpMessage{},
		HpaiControlEndpoint: hpaiControlEndpoint,
	}
	_child._KnxNetIpMessage._KnxNetIpMessageChildRequirements = _child
	return _child, nil
}

func (m *_DescriptionRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))), utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_DescriptionRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("DescriptionRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for DescriptionRequest")
		}

		if err := WriteSimpleField[HPAIControlEndpoint](ctx, "hpaiControlEndpoint", m.GetHpaiControlEndpoint(), WriteComplex[HPAIControlEndpoint](writeBuffer), codegen.WithByteOrder(binary.BigEndian)); err != nil {
			return errors.Wrap(err, "Error serializing 'hpaiControlEndpoint' field")
		}

		if popErr := writeBuffer.PopContext("DescriptionRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for DescriptionRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_DescriptionRequest) isDescriptionRequest() bool {
	return true
}

func (m *_DescriptionRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
