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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// EncodedTicket is the corresponding interface of EncodedTicket
type EncodedTicket interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// IsEncodedTicket is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsEncodedTicket()
}

// _EncodedTicket is the data-structure of this message
type _EncodedTicket struct {
}

var _ EncodedTicket = (*_EncodedTicket)(nil)

// NewEncodedTicket factory function for _EncodedTicket
func NewEncodedTicket() *_EncodedTicket {
	return &_EncodedTicket{}
}

// Deprecated: use the interface for direct cast
func CastEncodedTicket(structType any) EncodedTicket {
	if casted, ok := structType.(EncodedTicket); ok {
		return casted
	}
	if casted, ok := structType.(*EncodedTicket); ok {
		return *casted
	}
	return nil
}

func (m *_EncodedTicket) GetTypeName() string {
	return "EncodedTicket"
}

func (m *_EncodedTicket) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	return lengthInBits
}

func (m *_EncodedTicket) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func EncodedTicketParse(ctx context.Context, theBytes []byte) (EncodedTicket, error) {
	return EncodedTicketParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func EncodedTicketParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (EncodedTicket, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (EncodedTicket, error) {
		return EncodedTicketParseWithBuffer(ctx, readBuffer)
	}
}

func EncodedTicketParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (EncodedTicket, error) {
	v, err := (&_EncodedTicket{}).parse(ctx, readBuffer)
	if err != nil {
		return nil, err
	}
	return v, err
}

func (m *_EncodedTicket) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__encodedTicket EncodedTicket, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("EncodedTicket"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for EncodedTicket")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("EncodedTicket"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for EncodedTicket")
	}

	return m, nil
}

func (m *_EncodedTicket) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_EncodedTicket) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("EncodedTicket"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for EncodedTicket")
	}

	if popErr := writeBuffer.PopContext("EncodedTicket"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for EncodedTicket")
	}
	return nil
}

func (m *_EncodedTicket) IsEncodedTicket() {}

func (m *_EncodedTicket) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
