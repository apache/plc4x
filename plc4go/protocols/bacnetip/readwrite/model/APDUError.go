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

// APDUError is the corresponding interface of APDUError
type APDUError interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	APDU
	// GetOriginalInvokeId returns OriginalInvokeId (property field)
	GetOriginalInvokeId() uint8
	// GetErrorChoice returns ErrorChoice (property field)
	GetErrorChoice() BACnetConfirmedServiceChoice
	// GetError returns Error (property field)
	GetError() BACnetError
}

// APDUErrorExactly can be used when we want exactly this type and not a type which fulfills APDUError.
// This is useful for switch cases.
type APDUErrorExactly interface {
	APDUError
	isAPDUError() bool
}

// _APDUError is the data-structure of this message
type _APDUError struct {
	*_APDU
	OriginalInvokeId uint8
	ErrorChoice      BACnetConfirmedServiceChoice
	Error            BACnetError
	// Reserved Fields
	reservedField0 *uint8
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_APDUError) GetApduType() ApduType {
	return ApduType_ERROR_PDU
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_APDUError) InitializeParent(parent APDU) {}

func (m *_APDUError) GetParent() APDU {
	return m._APDU
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_APDUError) GetOriginalInvokeId() uint8 {
	return m.OriginalInvokeId
}

func (m *_APDUError) GetErrorChoice() BACnetConfirmedServiceChoice {
	return m.ErrorChoice
}

func (m *_APDUError) GetError() BACnetError {
	return m.Error
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAPDUError factory function for _APDUError
func NewAPDUError(originalInvokeId uint8, errorChoice BACnetConfirmedServiceChoice, error BACnetError, apduLength uint16) *_APDUError {
	_result := &_APDUError{
		OriginalInvokeId: originalInvokeId,
		ErrorChoice:      errorChoice,
		Error:            error,
		_APDU:            NewAPDU(apduLength),
	}
	_result._APDU._APDUChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastAPDUError(structType any) APDUError {
	if casted, ok := structType.(APDUError); ok {
		return casted
	}
	if casted, ok := structType.(*APDUError); ok {
		return *casted
	}
	return nil
}

func (m *_APDUError) GetTypeName() string {
	return "APDUError"
}

func (m *_APDUError) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Reserved Field (reserved)
	lengthInBits += 4

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (errorChoice)
	lengthInBits += 8

	// Simple field (error)
	lengthInBits += m.Error.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_APDUError) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func APDUErrorParse(ctx context.Context, theBytes []byte, apduLength uint16) (APDUError, error) {
	return APDUErrorParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), apduLength)
}

func APDUErrorParseWithBufferProducer(apduLength uint16) func(ctx context.Context, readBuffer utils.ReadBuffer) (APDUError, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (APDUError, error) {
		return APDUErrorParseWithBuffer(ctx, readBuffer, apduLength)
	}
}

func APDUErrorParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, apduLength uint16) (APDUError, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("APDUError"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for APDUError")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(4)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	originalInvokeId, err := ReadSimpleField(ctx, "originalInvokeId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'originalInvokeId' field"))
	}

	errorChoice, err := ReadEnumField[BACnetConfirmedServiceChoice](ctx, "errorChoice", "BACnetConfirmedServiceChoice", ReadEnum(BACnetConfirmedServiceChoiceByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'errorChoice' field"))
	}

	error, err := ReadSimpleField[BACnetError](ctx, "error", ReadComplex[BACnetError](BACnetErrorParseWithBufferProducer[BACnetError]((BACnetConfirmedServiceChoice)(errorChoice)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'error' field"))
	}

	if closeErr := readBuffer.CloseContext("APDUError"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for APDUError")
	}

	// Create a partially initialized instance
	_child := &_APDUError{
		_APDU: &_APDU{
			ApduLength: apduLength,
		},
		OriginalInvokeId: originalInvokeId,
		ErrorChoice:      errorChoice,
		Error:            error,
		reservedField0:   reservedField0,
	}
	_child._APDU._APDUChildRequirements = _child
	return _child, nil
}

func (m *_APDUError) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_APDUError) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("APDUError"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for APDUError")
		}

		if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 4)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteSimpleField[uint8](ctx, "originalInvokeId", m.GetOriginalInvokeId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'originalInvokeId' field")
		}

		if err := WriteSimpleEnumField[BACnetConfirmedServiceChoice](ctx, "errorChoice", "BACnetConfirmedServiceChoice", m.GetErrorChoice(), WriteEnum[BACnetConfirmedServiceChoice, uint8](BACnetConfirmedServiceChoice.GetValue, BACnetConfirmedServiceChoice.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'errorChoice' field")
		}

		if err := WriteSimpleField[BACnetError](ctx, "error", m.GetError(), WriteComplex[BACnetError](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'error' field")
		}

		if popErr := writeBuffer.PopContext("APDUError"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for APDUError")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_APDUError) isAPDUError() bool {
	return true
}

func (m *_APDUError) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
