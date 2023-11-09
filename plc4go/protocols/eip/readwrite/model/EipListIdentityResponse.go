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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Code generated by code-generation. DO NOT EDIT.

// EipListIdentityResponse is the corresponding interface of EipListIdentityResponse
type EipListIdentityResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	EipPacket
	// GetItems returns Items (property field)
	GetItems() []CommandSpecificDataItem
}

// EipListIdentityResponseExactly can be used when we want exactly this type and not a type which fulfills EipListIdentityResponse.
// This is useful for switch cases.
type EipListIdentityResponseExactly interface {
	EipListIdentityResponse
	isEipListIdentityResponse() bool
}

// _EipListIdentityResponse is the data-structure of this message
type _EipListIdentityResponse struct {
	*_EipPacket
	Items []CommandSpecificDataItem
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_EipListIdentityResponse) GetCommand() uint16 {
	return 0x0063
}

func (m *_EipListIdentityResponse) GetResponse() bool {
	return bool(true)
}

func (m *_EipListIdentityResponse) GetPacketLength() uint16 {
	return 0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_EipListIdentityResponse) InitializeParent(parent EipPacket, sessionHandle uint32, status uint32, senderContext []byte, options uint32) {
	m.SessionHandle = sessionHandle
	m.Status = status
	m.SenderContext = senderContext
	m.Options = options
}

func (m *_EipListIdentityResponse) GetParent() EipPacket {
	return m._EipPacket
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_EipListIdentityResponse) GetItems() []CommandSpecificDataItem {
	return m.Items
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewEipListIdentityResponse factory function for _EipListIdentityResponse
func NewEipListIdentityResponse(items []CommandSpecificDataItem, sessionHandle uint32, status uint32, senderContext []byte, options uint32) *_EipListIdentityResponse {
	_result := &_EipListIdentityResponse{
		Items:      items,
		_EipPacket: NewEipPacket(sessionHandle, status, senderContext, options),
	}
	_result._EipPacket._EipPacketChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastEipListIdentityResponse(structType any) EipListIdentityResponse {
	if casted, ok := structType.(EipListIdentityResponse); ok {
		return casted
	}
	if casted, ok := structType.(*EipListIdentityResponse); ok {
		return *casted
	}
	return nil
}

func (m *_EipListIdentityResponse) GetTypeName() string {
	return "EipListIdentityResponse"
}

func (m *_EipListIdentityResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Implicit Field (itemCount)
	lengthInBits += 16

	// Array field
	if len(m.Items) > 0 {
		for _curItem, element := range m.Items {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.Items), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_EipListIdentityResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func EipListIdentityResponseParse(ctx context.Context, theBytes []byte, response bool) (EipListIdentityResponse, error) {
	return EipListIdentityResponseParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), response)
}

func EipListIdentityResponseParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, response bool) (EipListIdentityResponse, error) {
	positionAware := readBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pullErr := readBuffer.PullContext("EipListIdentityResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for EipListIdentityResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Implicit Field (itemCount) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
	itemCount, _itemCountErr := readBuffer.ReadUint16("itemCount", 16)
	_ = itemCount
	if _itemCountErr != nil {
		return nil, errors.Wrap(_itemCountErr, "Error parsing 'itemCount' field of EipListIdentityResponse")
	}

	// Array field (items)
	if pullErr := readBuffer.PullContext("items", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for items")
	}
	// Count array
	items := make([]CommandSpecificDataItem, utils.Max(itemCount, 0))
	// This happens when the size is set conditional to 0
	if len(items) == 0 {
		items = nil
	}
	{
		_numItems := uint16(utils.Max(itemCount, 0))
		for _curItem := uint16(0); _curItem < _numItems; _curItem++ {
			arrayCtx := utils.CreateArrayContext(ctx, int(_numItems), int(_curItem))
			_ = arrayCtx
			_ = _curItem
			_item, _err := CommandSpecificDataItemParseWithBuffer(arrayCtx, readBuffer)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'items' field of EipListIdentityResponse")
			}
			items[_curItem] = _item.(CommandSpecificDataItem)
		}
	}
	if closeErr := readBuffer.CloseContext("items", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for items")
	}

	if closeErr := readBuffer.CloseContext("EipListIdentityResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for EipListIdentityResponse")
	}

	// Create a partially initialized instance
	_child := &_EipListIdentityResponse{
		_EipPacket: &_EipPacket{},
		Items:      items,
	}
	_child._EipPacket._EipPacketChildRequirements = _child
	return _child, nil
}

func (m *_EipListIdentityResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_EipListIdentityResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("EipListIdentityResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for EipListIdentityResponse")
		}

		// Implicit Field (itemCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		itemCount := uint16(uint16(len(m.GetItems())))
		_itemCountErr := writeBuffer.WriteUint16("itemCount", 16, (itemCount))
		if _itemCountErr != nil {
			return errors.Wrap(_itemCountErr, "Error serializing 'itemCount' field")
		}

		// Array Field (items)
		if pushErr := writeBuffer.PushContext("items", utils.WithRenderAsList(true)); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for items")
		}
		for _curItem, _element := range m.GetItems() {
			_ = _curItem
			arrayCtx := utils.CreateArrayContext(ctx, len(m.GetItems()), _curItem)
			_ = arrayCtx
			_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'items' field")
			}
		}
		if popErr := writeBuffer.PopContext("items", utils.WithRenderAsList(true)); popErr != nil {
			return errors.Wrap(popErr, "Error popping for items")
		}

		if popErr := writeBuffer.PopContext("EipListIdentityResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for EipListIdentityResponse")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_EipListIdentityResponse) isEipListIdentityResponse() bool {
	return true
}

func (m *_EipListIdentityResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}