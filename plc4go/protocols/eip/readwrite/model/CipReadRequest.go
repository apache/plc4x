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

// CipReadRequest is the corresponding interface of CipReadRequest
type CipReadRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	CipService
	// GetTag returns Tag (property field)
	GetTag() []byte
	// GetElementNb returns ElementNb (property field)
	GetElementNb() uint16
}

// CipReadRequestExactly can be used when we want exactly this type and not a type which fulfills CipReadRequest.
// This is useful for switch cases.
type CipReadRequestExactly interface {
	CipReadRequest
	isCipReadRequest() bool
}

// _CipReadRequest is the data-structure of this message
type _CipReadRequest struct {
	*_CipService
	Tag       []byte
	ElementNb uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_CipReadRequest) GetService() uint8 {
	return 0x4C
}

func (m *_CipReadRequest) GetResponse() bool {
	return bool(false)
}

func (m *_CipReadRequest) GetConnected() bool {
	return false
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CipReadRequest) InitializeParent(parent CipService) {}

func (m *_CipReadRequest) GetParent() CipService {
	return m._CipService
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CipReadRequest) GetTag() []byte {
	return m.Tag
}

func (m *_CipReadRequest) GetElementNb() uint16 {
	return m.ElementNb
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCipReadRequest factory function for _CipReadRequest
func NewCipReadRequest(tag []byte, elementNb uint16, serviceLen uint16) *_CipReadRequest {
	_result := &_CipReadRequest{
		Tag:         tag,
		ElementNb:   elementNb,
		_CipService: NewCipService(serviceLen),
	}
	_result._CipService._CipServiceChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastCipReadRequest(structType any) CipReadRequest {
	if casted, ok := structType.(CipReadRequest); ok {
		return casted
	}
	if casted, ok := structType.(*CipReadRequest); ok {
		return *casted
	}
	return nil
}

func (m *_CipReadRequest) GetTypeName() string {
	return "CipReadRequest"
}

func (m *_CipReadRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Implicit Field (requestPathSize)
	lengthInBits += 8

	// Array field
	if len(m.Tag) > 0 {
		lengthInBits += 8 * uint16(len(m.Tag))
	}

	// Simple field (elementNb)
	lengthInBits += 16

	return lengthInBits
}

func (m *_CipReadRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CipReadRequestParse(ctx context.Context, theBytes []byte, connected bool, serviceLen uint16) (CipReadRequest, error) {
	return CipReadRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), connected, serviceLen)
}

func CipReadRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, connected bool, serviceLen uint16) (CipReadRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pullErr := readBuffer.PullContext("CipReadRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CipReadRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestPathSize, err := ReadImplicitField[uint8](ctx, "requestPathSize", ReadUnsignedByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestPathSize' field"))
	}

	tag, err := readBuffer.ReadByteArray("tag", int((int32(requestPathSize) * int32(int32(2)))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'tag' field"))
	}

	// Simple Field (elementNb)
	_elementNb, _elementNbErr := /*TODO: migrate me*/ readBuffer.ReadUint16("elementNb", 16)
	if _elementNbErr != nil {
		return nil, errors.Wrap(_elementNbErr, "Error parsing 'elementNb' field of CipReadRequest")
	}
	elementNb := _elementNb

	if closeErr := readBuffer.CloseContext("CipReadRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CipReadRequest")
	}

	// Create a partially initialized instance
	_child := &_CipReadRequest{
		_CipService: &_CipService{
			ServiceLen: serviceLen,
		},
		Tag:       tag,
		ElementNb: elementNb,
	}
	_child._CipService._CipServiceChildRequirements = _child
	return _child, nil
}

func (m *_CipReadRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CipReadRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CipReadRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CipReadRequest")
		}

		// Implicit Field (requestPathSize) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		requestPathSize := uint8(uint8(uint8(len(m.GetTag()))) / uint8(uint8(2)))
		_requestPathSizeErr := /*TODO: migrate me*/ writeBuffer.WriteUint8("requestPathSize", 8, uint8((requestPathSize)))
		if _requestPathSizeErr != nil {
			return errors.Wrap(_requestPathSizeErr, "Error serializing 'requestPathSize' field")
		}

		// Array Field (tag)
		// Byte Array field (tag)
		if err := writeBuffer.WriteByteArray("tag", m.GetTag()); err != nil {
			return errors.Wrap(err, "Error serializing 'tag' field")
		}

		// Simple Field (elementNb)
		elementNb := uint16(m.GetElementNb())
		_elementNbErr := /*TODO: migrate me*/ writeBuffer.WriteUint16("elementNb", 16, uint16((elementNb)))
		if _elementNbErr != nil {
			return errors.Wrap(_elementNbErr, "Error serializing 'elementNb' field")
		}

		if popErr := writeBuffer.PopContext("CipReadRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CipReadRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CipReadRequest) isCipReadRequest() bool {
	return true
}

func (m *_CipReadRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
