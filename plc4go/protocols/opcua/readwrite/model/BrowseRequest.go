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

// BrowseRequest is the corresponding interface of BrowseRequest
type BrowseRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() ExtensionObjectDefinition
	// GetView returns View (property field)
	GetView() ExtensionObjectDefinition
	// GetRequestedMaxReferencesPerNode returns RequestedMaxReferencesPerNode (property field)
	GetRequestedMaxReferencesPerNode() uint32
	// GetNoOfNodesToBrowse returns NoOfNodesToBrowse (property field)
	GetNoOfNodesToBrowse() int32
	// GetNodesToBrowse returns NodesToBrowse (property field)
	GetNodesToBrowse() []ExtensionObjectDefinition
}

// BrowseRequestExactly can be used when we want exactly this type and not a type which fulfills BrowseRequest.
// This is useful for switch cases.
type BrowseRequestExactly interface {
	BrowseRequest
	isBrowseRequest() bool
}

// _BrowseRequest is the data-structure of this message
type _BrowseRequest struct {
	*_ExtensionObjectDefinition
	RequestHeader                 ExtensionObjectDefinition
	View                          ExtensionObjectDefinition
	RequestedMaxReferencesPerNode uint32
	NoOfNodesToBrowse             int32
	NodesToBrowse                 []ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BrowseRequest) GetIdentifier() string {
	return "527"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BrowseRequest) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_BrowseRequest) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BrowseRequest) GetRequestHeader() ExtensionObjectDefinition {
	return m.RequestHeader
}

func (m *_BrowseRequest) GetView() ExtensionObjectDefinition {
	return m.View
}

func (m *_BrowseRequest) GetRequestedMaxReferencesPerNode() uint32 {
	return m.RequestedMaxReferencesPerNode
}

func (m *_BrowseRequest) GetNoOfNodesToBrowse() int32 {
	return m.NoOfNodesToBrowse
}

func (m *_BrowseRequest) GetNodesToBrowse() []ExtensionObjectDefinition {
	return m.NodesToBrowse
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBrowseRequest factory function for _BrowseRequest
func NewBrowseRequest(requestHeader ExtensionObjectDefinition, view ExtensionObjectDefinition, requestedMaxReferencesPerNode uint32, noOfNodesToBrowse int32, nodesToBrowse []ExtensionObjectDefinition) *_BrowseRequest {
	_result := &_BrowseRequest{
		RequestHeader:                 requestHeader,
		View:                          view,
		RequestedMaxReferencesPerNode: requestedMaxReferencesPerNode,
		NoOfNodesToBrowse:             noOfNodesToBrowse,
		NodesToBrowse:                 nodesToBrowse,
		_ExtensionObjectDefinition:    NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBrowseRequest(structType any) BrowseRequest {
	if casted, ok := structType.(BrowseRequest); ok {
		return casted
	}
	if casted, ok := structType.(*BrowseRequest); ok {
		return *casted
	}
	return nil
}

func (m *_BrowseRequest) GetTypeName() string {
	return "BrowseRequest"
}

func (m *_BrowseRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Simple field (view)
	lengthInBits += m.View.GetLengthInBits(ctx)

	// Simple field (requestedMaxReferencesPerNode)
	lengthInBits += 32

	// Simple field (noOfNodesToBrowse)
	lengthInBits += 32

	// Array field
	if len(m.NodesToBrowse) > 0 {
		for _curItem, element := range m.NodesToBrowse {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.NodesToBrowse), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_BrowseRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BrowseRequestParse(ctx context.Context, theBytes []byte, identifier string) (BrowseRequest, error) {
	return BrowseRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func BrowseRequestParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (BrowseRequest, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BrowseRequest, error) {
		return BrowseRequestParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func BrowseRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (BrowseRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BrowseRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BrowseRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestHeader, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "requestHeader", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("391")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestHeader' field"))
	}

	view, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "view", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("513")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'view' field"))
	}

	requestedMaxReferencesPerNode, err := ReadSimpleField(ctx, "requestedMaxReferencesPerNode", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestedMaxReferencesPerNode' field"))
	}

	noOfNodesToBrowse, err := ReadSimpleField(ctx, "noOfNodesToBrowse", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfNodesToBrowse' field"))
	}

	nodesToBrowse, err := ReadCountArrayField[ExtensionObjectDefinition](ctx, "nodesToBrowse", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("516")), readBuffer), uint64(noOfNodesToBrowse))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'nodesToBrowse' field"))
	}

	if closeErr := readBuffer.CloseContext("BrowseRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BrowseRequest")
	}

	// Create a partially initialized instance
	_child := &_BrowseRequest{
		_ExtensionObjectDefinition:    &_ExtensionObjectDefinition{},
		RequestHeader:                 requestHeader,
		View:                          view,
		RequestedMaxReferencesPerNode: requestedMaxReferencesPerNode,
		NoOfNodesToBrowse:             noOfNodesToBrowse,
		NodesToBrowse:                 nodesToBrowse,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_BrowseRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BrowseRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BrowseRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BrowseRequest")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "requestHeader", m.GetRequestHeader(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestHeader' field")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "view", m.GetView(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'view' field")
		}

		if err := WriteSimpleField[uint32](ctx, "requestedMaxReferencesPerNode", m.GetRequestedMaxReferencesPerNode(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestedMaxReferencesPerNode' field")
		}

		if err := WriteSimpleField[int32](ctx, "noOfNodesToBrowse", m.GetNoOfNodesToBrowse(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfNodesToBrowse' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "nodesToBrowse", m.GetNodesToBrowse(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'nodesToBrowse' field")
		}

		if popErr := writeBuffer.PopContext("BrowseRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BrowseRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BrowseRequest) isBrowseRequest() bool {
	return true
}

func (m *_BrowseRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
