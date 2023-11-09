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

// UnregisterNodesRequest is the corresponding interface of UnregisterNodesRequest
type UnregisterNodesRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() ExtensionObjectDefinition
	// GetNoOfNodesToUnregister returns NoOfNodesToUnregister (property field)
	GetNoOfNodesToUnregister() int32
	// GetNodesToUnregister returns NodesToUnregister (property field)
	GetNodesToUnregister() []NodeId
}

// UnregisterNodesRequestExactly can be used when we want exactly this type and not a type which fulfills UnregisterNodesRequest.
// This is useful for switch cases.
type UnregisterNodesRequestExactly interface {
	UnregisterNodesRequest
	isUnregisterNodesRequest() bool
}

// _UnregisterNodesRequest is the data-structure of this message
type _UnregisterNodesRequest struct {
	*_ExtensionObjectDefinition
	RequestHeader         ExtensionObjectDefinition
	NoOfNodesToUnregister int32
	NodesToUnregister     []NodeId
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_UnregisterNodesRequest) GetIdentifier() string {
	return "566"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_UnregisterNodesRequest) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_UnregisterNodesRequest) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_UnregisterNodesRequest) GetRequestHeader() ExtensionObjectDefinition {
	return m.RequestHeader
}

func (m *_UnregisterNodesRequest) GetNoOfNodesToUnregister() int32 {
	return m.NoOfNodesToUnregister
}

func (m *_UnregisterNodesRequest) GetNodesToUnregister() []NodeId {
	return m.NodesToUnregister
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewUnregisterNodesRequest factory function for _UnregisterNodesRequest
func NewUnregisterNodesRequest(requestHeader ExtensionObjectDefinition, noOfNodesToUnregister int32, nodesToUnregister []NodeId) *_UnregisterNodesRequest {
	_result := &_UnregisterNodesRequest{
		RequestHeader:              requestHeader,
		NoOfNodesToUnregister:      noOfNodesToUnregister,
		NodesToUnregister:          nodesToUnregister,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastUnregisterNodesRequest(structType any) UnregisterNodesRequest {
	if casted, ok := structType.(UnregisterNodesRequest); ok {
		return casted
	}
	if casted, ok := structType.(*UnregisterNodesRequest); ok {
		return *casted
	}
	return nil
}

func (m *_UnregisterNodesRequest) GetTypeName() string {
	return "UnregisterNodesRequest"
}

func (m *_UnregisterNodesRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Simple field (noOfNodesToUnregister)
	lengthInBits += 32

	// Array field
	if len(m.NodesToUnregister) > 0 {
		for _curItem, element := range m.NodesToUnregister {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.NodesToUnregister), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_UnregisterNodesRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func UnregisterNodesRequestParse(ctx context.Context, theBytes []byte, identifier string) (UnregisterNodesRequest, error) {
	return UnregisterNodesRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func UnregisterNodesRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (UnregisterNodesRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pullErr := readBuffer.PullContext("UnregisterNodesRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for UnregisterNodesRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (requestHeader)
	if pullErr := readBuffer.PullContext("requestHeader"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for requestHeader")
	}
	_requestHeader, _requestHeaderErr := ExtensionObjectDefinitionParseWithBuffer(ctx, readBuffer, string("391"))
	if _requestHeaderErr != nil {
		return nil, errors.Wrap(_requestHeaderErr, "Error parsing 'requestHeader' field of UnregisterNodesRequest")
	}
	requestHeader := _requestHeader.(ExtensionObjectDefinition)
	if closeErr := readBuffer.CloseContext("requestHeader"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for requestHeader")
	}

	// Simple Field (noOfNodesToUnregister)
	_noOfNodesToUnregister, _noOfNodesToUnregisterErr := readBuffer.ReadInt32("noOfNodesToUnregister", 32)
	if _noOfNodesToUnregisterErr != nil {
		return nil, errors.Wrap(_noOfNodesToUnregisterErr, "Error parsing 'noOfNodesToUnregister' field of UnregisterNodesRequest")
	}
	noOfNodesToUnregister := _noOfNodesToUnregister

	// Array field (nodesToUnregister)
	if pullErr := readBuffer.PullContext("nodesToUnregister", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for nodesToUnregister")
	}
	// Count array
	nodesToUnregister := make([]NodeId, utils.Max(noOfNodesToUnregister, 0))
	// This happens when the size is set conditional to 0
	if len(nodesToUnregister) == 0 {
		nodesToUnregister = nil
	}
	{
		_numItems := uint16(utils.Max(noOfNodesToUnregister, 0))
		for _curItem := uint16(0); _curItem < _numItems; _curItem++ {
			arrayCtx := utils.CreateArrayContext(ctx, int(_numItems), int(_curItem))
			_ = arrayCtx
			_ = _curItem
			_item, _err := NodeIdParseWithBuffer(arrayCtx, readBuffer)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'nodesToUnregister' field of UnregisterNodesRequest")
			}
			nodesToUnregister[_curItem] = _item.(NodeId)
		}
	}
	if closeErr := readBuffer.CloseContext("nodesToUnregister", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for nodesToUnregister")
	}

	if closeErr := readBuffer.CloseContext("UnregisterNodesRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for UnregisterNodesRequest")
	}

	// Create a partially initialized instance
	_child := &_UnregisterNodesRequest{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		RequestHeader:              requestHeader,
		NoOfNodesToUnregister:      noOfNodesToUnregister,
		NodesToUnregister:          nodesToUnregister,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_UnregisterNodesRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_UnregisterNodesRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("UnregisterNodesRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for UnregisterNodesRequest")
		}

		// Simple Field (requestHeader)
		if pushErr := writeBuffer.PushContext("requestHeader"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for requestHeader")
		}
		_requestHeaderErr := writeBuffer.WriteSerializable(ctx, m.GetRequestHeader())
		if popErr := writeBuffer.PopContext("requestHeader"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for requestHeader")
		}
		if _requestHeaderErr != nil {
			return errors.Wrap(_requestHeaderErr, "Error serializing 'requestHeader' field")
		}

		// Simple Field (noOfNodesToUnregister)
		noOfNodesToUnregister := int32(m.GetNoOfNodesToUnregister())
		_noOfNodesToUnregisterErr := writeBuffer.WriteInt32("noOfNodesToUnregister", 32, (noOfNodesToUnregister))
		if _noOfNodesToUnregisterErr != nil {
			return errors.Wrap(_noOfNodesToUnregisterErr, "Error serializing 'noOfNodesToUnregister' field")
		}

		// Array Field (nodesToUnregister)
		if pushErr := writeBuffer.PushContext("nodesToUnregister", utils.WithRenderAsList(true)); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for nodesToUnregister")
		}
		for _curItem, _element := range m.GetNodesToUnregister() {
			_ = _curItem
			arrayCtx := utils.CreateArrayContext(ctx, len(m.GetNodesToUnregister()), _curItem)
			_ = arrayCtx
			_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'nodesToUnregister' field")
			}
		}
		if popErr := writeBuffer.PopContext("nodesToUnregister", utils.WithRenderAsList(true)); popErr != nil {
			return errors.Wrap(popErr, "Error popping for nodesToUnregister")
		}

		if popErr := writeBuffer.PopContext("UnregisterNodesRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for UnregisterNodesRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_UnregisterNodesRequest) isUnregisterNodesRequest() bool {
	return true
}

func (m *_UnregisterNodesRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}