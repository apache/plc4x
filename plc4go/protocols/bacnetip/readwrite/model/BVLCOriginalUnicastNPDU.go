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

// BVLCOriginalUnicastNPDU is the corresponding interface of BVLCOriginalUnicastNPDU
type BVLCOriginalUnicastNPDU interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BVLC
	// GetNpdu returns Npdu (property field)
	GetNpdu() NPDU
}

// BVLCOriginalUnicastNPDUExactly can be used when we want exactly this type and not a type which fulfills BVLCOriginalUnicastNPDU.
// This is useful for switch cases.
type BVLCOriginalUnicastNPDUExactly interface {
	BVLCOriginalUnicastNPDU
	isBVLCOriginalUnicastNPDU() bool
}

// _BVLCOriginalUnicastNPDU is the data-structure of this message
type _BVLCOriginalUnicastNPDU struct {
	*_BVLC
	Npdu NPDU

	// Arguments.
	BvlcPayloadLength uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BVLCOriginalUnicastNPDU) GetBvlcFunction() uint8 {
	return 0x0A
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BVLCOriginalUnicastNPDU) InitializeParent(parent BVLC) {}

func (m *_BVLCOriginalUnicastNPDU) GetParent() BVLC {
	return m._BVLC
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BVLCOriginalUnicastNPDU) GetNpdu() NPDU {
	return m.Npdu
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBVLCOriginalUnicastNPDU factory function for _BVLCOriginalUnicastNPDU
func NewBVLCOriginalUnicastNPDU(npdu NPDU, bvlcPayloadLength uint16) *_BVLCOriginalUnicastNPDU {
	_result := &_BVLCOriginalUnicastNPDU{
		Npdu:  npdu,
		_BVLC: NewBVLC(),
	}
	_result._BVLC._BVLCChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBVLCOriginalUnicastNPDU(structType any) BVLCOriginalUnicastNPDU {
	if casted, ok := structType.(BVLCOriginalUnicastNPDU); ok {
		return casted
	}
	if casted, ok := structType.(*BVLCOriginalUnicastNPDU); ok {
		return *casted
	}
	return nil
}

func (m *_BVLCOriginalUnicastNPDU) GetTypeName() string {
	return "BVLCOriginalUnicastNPDU"
}

func (m *_BVLCOriginalUnicastNPDU) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (npdu)
	lengthInBits += m.Npdu.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BVLCOriginalUnicastNPDU) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BVLCOriginalUnicastNPDUParse(ctx context.Context, theBytes []byte, bvlcPayloadLength uint16) (BVLCOriginalUnicastNPDU, error) {
	return BVLCOriginalUnicastNPDUParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes, utils.WithByteOrderForReadBufferByteBased(binary.BigEndian)), bvlcPayloadLength)
}

func BVLCOriginalUnicastNPDUParseWithBufferProducer(bvlcPayloadLength uint16) func(ctx context.Context, readBuffer utils.ReadBuffer) (BVLCOriginalUnicastNPDU, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BVLCOriginalUnicastNPDU, error) {
		return BVLCOriginalUnicastNPDUParseWithBuffer(ctx, readBuffer, bvlcPayloadLength)
	}
}

func BVLCOriginalUnicastNPDUParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, bvlcPayloadLength uint16) (BVLCOriginalUnicastNPDU, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BVLCOriginalUnicastNPDU"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BVLCOriginalUnicastNPDU")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	npdu, err := ReadSimpleField[NPDU](ctx, "npdu", ReadComplex[NPDU](NPDUParseWithBufferProducer((uint16)(bvlcPayloadLength)), readBuffer), codegen.WithByteOrder(binary.BigEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'npdu' field"))
	}

	if closeErr := readBuffer.CloseContext("BVLCOriginalUnicastNPDU"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BVLCOriginalUnicastNPDU")
	}

	// Create a partially initialized instance
	_child := &_BVLCOriginalUnicastNPDU{
		_BVLC: &_BVLC{},
		Npdu:  npdu,
	}
	_child._BVLC._BVLCChildRequirements = _child
	return _child, nil
}

func (m *_BVLCOriginalUnicastNPDU) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))), utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BVLCOriginalUnicastNPDU) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BVLCOriginalUnicastNPDU"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BVLCOriginalUnicastNPDU")
		}

		if err := WriteSimpleField[NPDU](ctx, "npdu", m.GetNpdu(), WriteComplex[NPDU](writeBuffer), codegen.WithByteOrder(binary.BigEndian)); err != nil {
			return errors.Wrap(err, "Error serializing 'npdu' field")
		}

		if popErr := writeBuffer.PopContext("BVLCOriginalUnicastNPDU"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BVLCOriginalUnicastNPDU")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

////
// Arguments Getter

func (m *_BVLCOriginalUnicastNPDU) GetBvlcPayloadLength() uint16 {
	return m.BvlcPayloadLength
}

//
////

func (m *_BVLCOriginalUnicastNPDU) isBVLCOriginalUnicastNPDU() bool {
	return true
}

func (m *_BVLCOriginalUnicastNPDU) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
