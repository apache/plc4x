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

// ModbusPDUWriteSingleRegisterResponse is the corresponding interface of ModbusPDUWriteSingleRegisterResponse
type ModbusPDUWriteSingleRegisterResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ModbusPDU
	// GetAddress returns Address (property field)
	GetAddress() uint16
	// GetValue returns Value (property field)
	GetValue() uint16
}

// ModbusPDUWriteSingleRegisterResponseExactly can be used when we want exactly this type and not a type which fulfills ModbusPDUWriteSingleRegisterResponse.
// This is useful for switch cases.
type ModbusPDUWriteSingleRegisterResponseExactly interface {
	ModbusPDUWriteSingleRegisterResponse
	isModbusPDUWriteSingleRegisterResponse() bool
}

// _ModbusPDUWriteSingleRegisterResponse is the data-structure of this message
type _ModbusPDUWriteSingleRegisterResponse struct {
	*_ModbusPDU
	Address uint16
	Value   uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModbusPDUWriteSingleRegisterResponse) GetErrorFlag() bool {
	return bool(false)
}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetFunctionFlag() uint8 {
	return 0x06
}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModbusPDUWriteSingleRegisterResponse) InitializeParent(parent ModbusPDU) {}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetParent() ModbusPDU {
	return m._ModbusPDU
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModbusPDUWriteSingleRegisterResponse) GetAddress() uint16 {
	return m.Address
}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetValue() uint16 {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewModbusPDUWriteSingleRegisterResponse factory function for _ModbusPDUWriteSingleRegisterResponse
func NewModbusPDUWriteSingleRegisterResponse(address uint16, value uint16) *_ModbusPDUWriteSingleRegisterResponse {
	_result := &_ModbusPDUWriteSingleRegisterResponse{
		Address:    address,
		Value:      value,
		_ModbusPDU: NewModbusPDU(),
	}
	_result._ModbusPDU._ModbusPDUChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastModbusPDUWriteSingleRegisterResponse(structType any) ModbusPDUWriteSingleRegisterResponse {
	if casted, ok := structType.(ModbusPDUWriteSingleRegisterResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ModbusPDUWriteSingleRegisterResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetTypeName() string {
	return "ModbusPDUWriteSingleRegisterResponse"
}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (address)
	lengthInBits += 16

	// Simple field (value)
	lengthInBits += 16

	return lengthInBits
}

func (m *_ModbusPDUWriteSingleRegisterResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func ModbusPDUWriteSingleRegisterResponseParse(ctx context.Context, theBytes []byte, response bool) (ModbusPDUWriteSingleRegisterResponse, error) {
	return ModbusPDUWriteSingleRegisterResponseParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), response)
}

func ModbusPDUWriteSingleRegisterResponseParseWithBufferProducer(response bool) func(ctx context.Context, readBuffer utils.ReadBuffer) (ModbusPDUWriteSingleRegisterResponse, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (ModbusPDUWriteSingleRegisterResponse, error) {
		return ModbusPDUWriteSingleRegisterResponseParseWithBuffer(ctx, readBuffer, response)
	}
}

func ModbusPDUWriteSingleRegisterResponseParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, response bool) (ModbusPDUWriteSingleRegisterResponse, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModbusPDUWriteSingleRegisterResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModbusPDUWriteSingleRegisterResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	address, err := ReadSimpleField(ctx, "address", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'address' field"))
	}

	value, err := ReadSimpleField(ctx, "value", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'value' field"))
	}

	if closeErr := readBuffer.CloseContext("ModbusPDUWriteSingleRegisterResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModbusPDUWriteSingleRegisterResponse")
	}

	// Create a partially initialized instance
	_child := &_ModbusPDUWriteSingleRegisterResponse{
		_ModbusPDU: &_ModbusPDU{},
		Address:    address,
		Value:      value,
	}
	_child._ModbusPDU._ModbusPDUChildRequirements = _child
	return _child, nil
}

func (m *_ModbusPDUWriteSingleRegisterResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModbusPDUWriteSingleRegisterResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModbusPDUWriteSingleRegisterResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModbusPDUWriteSingleRegisterResponse")
		}

		if err := WriteSimpleField[uint16](ctx, "address", m.GetAddress(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'address' field")
		}

		if err := WriteSimpleField[uint16](ctx, "value", m.GetValue(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'value' field")
		}

		if popErr := writeBuffer.PopContext("ModbusPDUWriteSingleRegisterResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModbusPDUWriteSingleRegisterResponse")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModbusPDUWriteSingleRegisterResponse) isModbusPDUWriteSingleRegisterResponse() bool {
	return true
}

func (m *_ModbusPDUWriteSingleRegisterResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
