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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// RoutingIndication is the corresponding interface of RoutingIndication
type RoutingIndication interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	KnxNetIpMessage
	// IsRoutingIndication is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsRoutingIndication()
}

// _RoutingIndication is the data-structure of this message
type _RoutingIndication struct {
	KnxNetIpMessageContract
}

var _ RoutingIndication = (*_RoutingIndication)(nil)
var _ KnxNetIpMessageRequirements = (*_RoutingIndication)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_RoutingIndication) GetMsgType() uint16 {
	return 0x0530
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_RoutingIndication) GetParent() KnxNetIpMessageContract {
	return m.KnxNetIpMessageContract
}

// NewRoutingIndication factory function for _RoutingIndication
func NewRoutingIndication() *_RoutingIndication {
	_result := &_RoutingIndication{
		KnxNetIpMessageContract: NewKnxNetIpMessage(),
	}
	_result.KnxNetIpMessageContract.(*_KnxNetIpMessage)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastRoutingIndication(structType any) RoutingIndication {
	if casted, ok := structType.(RoutingIndication); ok {
		return casted
	}
	if casted, ok := structType.(*RoutingIndication); ok {
		return *casted
	}
	return nil
}

func (m *_RoutingIndication) GetTypeName() string {
	return "RoutingIndication"
}

func (m *_RoutingIndication) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.KnxNetIpMessageContract.(*_KnxNetIpMessage).getLengthInBits(ctx))

	return lengthInBits
}

func (m *_RoutingIndication) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_RoutingIndication) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_KnxNetIpMessage) (__routingIndication RoutingIndication, err error) {
	m.KnxNetIpMessageContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("RoutingIndication"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for RoutingIndication")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("RoutingIndication"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for RoutingIndication")
	}

	return m, nil
}

func (m *_RoutingIndication) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))), utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_RoutingIndication) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("RoutingIndication"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for RoutingIndication")
		}

		if popErr := writeBuffer.PopContext("RoutingIndication"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for RoutingIndication")
		}
		return nil
	}
	return m.KnxNetIpMessageContract.(*_KnxNetIpMessage).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_RoutingIndication) IsRoutingIndication() {}

func (m *_RoutingIndication) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
