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

// ProgramDiagnosticDataType is the corresponding interface of ProgramDiagnosticDataType
type ProgramDiagnosticDataType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetCreateSessionId returns CreateSessionId (property field)
	GetCreateSessionId() NodeId
	// GetCreateClientName returns CreateClientName (property field)
	GetCreateClientName() PascalString
	// GetInvocationCreationTime returns InvocationCreationTime (property field)
	GetInvocationCreationTime() int64
	// GetLastTransitionTime returns LastTransitionTime (property field)
	GetLastTransitionTime() int64
	// GetLastMethodCall returns LastMethodCall (property field)
	GetLastMethodCall() PascalString
	// GetLastMethodSessionId returns LastMethodSessionId (property field)
	GetLastMethodSessionId() NodeId
	// GetNoOfLastMethodInputArguments returns NoOfLastMethodInputArguments (property field)
	GetNoOfLastMethodInputArguments() int32
	// GetLastMethodInputArguments returns LastMethodInputArguments (property field)
	GetLastMethodInputArguments() []ExtensionObjectDefinition
	// GetNoOfLastMethodOutputArguments returns NoOfLastMethodOutputArguments (property field)
	GetNoOfLastMethodOutputArguments() int32
	// GetLastMethodOutputArguments returns LastMethodOutputArguments (property field)
	GetLastMethodOutputArguments() []ExtensionObjectDefinition
	// GetLastMethodCallTime returns LastMethodCallTime (property field)
	GetLastMethodCallTime() int64
	// GetLastMethodReturnStatus returns LastMethodReturnStatus (property field)
	GetLastMethodReturnStatus() ExtensionObjectDefinition
}

// ProgramDiagnosticDataTypeExactly can be used when we want exactly this type and not a type which fulfills ProgramDiagnosticDataType.
// This is useful for switch cases.
type ProgramDiagnosticDataTypeExactly interface {
	ProgramDiagnosticDataType
	isProgramDiagnosticDataType() bool
}

// _ProgramDiagnosticDataType is the data-structure of this message
type _ProgramDiagnosticDataType struct {
	*_ExtensionObjectDefinition
	CreateSessionId               NodeId
	CreateClientName              PascalString
	InvocationCreationTime        int64
	LastTransitionTime            int64
	LastMethodCall                PascalString
	LastMethodSessionId           NodeId
	NoOfLastMethodInputArguments  int32
	LastMethodInputArguments      []ExtensionObjectDefinition
	NoOfLastMethodOutputArguments int32
	LastMethodOutputArguments     []ExtensionObjectDefinition
	LastMethodCallTime            int64
	LastMethodReturnStatus        ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ProgramDiagnosticDataType) GetIdentifier() string {
	return "896"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ProgramDiagnosticDataType) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_ProgramDiagnosticDataType) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ProgramDiagnosticDataType) GetCreateSessionId() NodeId {
	return m.CreateSessionId
}

func (m *_ProgramDiagnosticDataType) GetCreateClientName() PascalString {
	return m.CreateClientName
}

func (m *_ProgramDiagnosticDataType) GetInvocationCreationTime() int64 {
	return m.InvocationCreationTime
}

func (m *_ProgramDiagnosticDataType) GetLastTransitionTime() int64 {
	return m.LastTransitionTime
}

func (m *_ProgramDiagnosticDataType) GetLastMethodCall() PascalString {
	return m.LastMethodCall
}

func (m *_ProgramDiagnosticDataType) GetLastMethodSessionId() NodeId {
	return m.LastMethodSessionId
}

func (m *_ProgramDiagnosticDataType) GetNoOfLastMethodInputArguments() int32 {
	return m.NoOfLastMethodInputArguments
}

func (m *_ProgramDiagnosticDataType) GetLastMethodInputArguments() []ExtensionObjectDefinition {
	return m.LastMethodInputArguments
}

func (m *_ProgramDiagnosticDataType) GetNoOfLastMethodOutputArguments() int32 {
	return m.NoOfLastMethodOutputArguments
}

func (m *_ProgramDiagnosticDataType) GetLastMethodOutputArguments() []ExtensionObjectDefinition {
	return m.LastMethodOutputArguments
}

func (m *_ProgramDiagnosticDataType) GetLastMethodCallTime() int64 {
	return m.LastMethodCallTime
}

func (m *_ProgramDiagnosticDataType) GetLastMethodReturnStatus() ExtensionObjectDefinition {
	return m.LastMethodReturnStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewProgramDiagnosticDataType factory function for _ProgramDiagnosticDataType
func NewProgramDiagnosticDataType(createSessionId NodeId, createClientName PascalString, invocationCreationTime int64, lastTransitionTime int64, lastMethodCall PascalString, lastMethodSessionId NodeId, noOfLastMethodInputArguments int32, lastMethodInputArguments []ExtensionObjectDefinition, noOfLastMethodOutputArguments int32, lastMethodOutputArguments []ExtensionObjectDefinition, lastMethodCallTime int64, lastMethodReturnStatus ExtensionObjectDefinition) *_ProgramDiagnosticDataType {
	_result := &_ProgramDiagnosticDataType{
		CreateSessionId:               createSessionId,
		CreateClientName:              createClientName,
		InvocationCreationTime:        invocationCreationTime,
		LastTransitionTime:            lastTransitionTime,
		LastMethodCall:                lastMethodCall,
		LastMethodSessionId:           lastMethodSessionId,
		NoOfLastMethodInputArguments:  noOfLastMethodInputArguments,
		LastMethodInputArguments:      lastMethodInputArguments,
		NoOfLastMethodOutputArguments: noOfLastMethodOutputArguments,
		LastMethodOutputArguments:     lastMethodOutputArguments,
		LastMethodCallTime:            lastMethodCallTime,
		LastMethodReturnStatus:        lastMethodReturnStatus,
		_ExtensionObjectDefinition:    NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastProgramDiagnosticDataType(structType any) ProgramDiagnosticDataType {
	if casted, ok := structType.(ProgramDiagnosticDataType); ok {
		return casted
	}
	if casted, ok := structType.(*ProgramDiagnosticDataType); ok {
		return *casted
	}
	return nil
}

func (m *_ProgramDiagnosticDataType) GetTypeName() string {
	return "ProgramDiagnosticDataType"
}

func (m *_ProgramDiagnosticDataType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (createSessionId)
	lengthInBits += m.CreateSessionId.GetLengthInBits(ctx)

	// Simple field (createClientName)
	lengthInBits += m.CreateClientName.GetLengthInBits(ctx)

	// Simple field (invocationCreationTime)
	lengthInBits += 64

	// Simple field (lastTransitionTime)
	lengthInBits += 64

	// Simple field (lastMethodCall)
	lengthInBits += m.LastMethodCall.GetLengthInBits(ctx)

	// Simple field (lastMethodSessionId)
	lengthInBits += m.LastMethodSessionId.GetLengthInBits(ctx)

	// Simple field (noOfLastMethodInputArguments)
	lengthInBits += 32

	// Array field
	if len(m.LastMethodInputArguments) > 0 {
		for _curItem, element := range m.LastMethodInputArguments {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.LastMethodInputArguments), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (noOfLastMethodOutputArguments)
	lengthInBits += 32

	// Array field
	if len(m.LastMethodOutputArguments) > 0 {
		for _curItem, element := range m.LastMethodOutputArguments {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.LastMethodOutputArguments), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (lastMethodCallTime)
	lengthInBits += 64

	// Simple field (lastMethodReturnStatus)
	lengthInBits += m.LastMethodReturnStatus.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_ProgramDiagnosticDataType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func ProgramDiagnosticDataTypeParse(ctx context.Context, theBytes []byte, identifier string) (ProgramDiagnosticDataType, error) {
	return ProgramDiagnosticDataTypeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func ProgramDiagnosticDataTypeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (ProgramDiagnosticDataType, error) {
	positionAware := readBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pullErr := readBuffer.PullContext("ProgramDiagnosticDataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ProgramDiagnosticDataType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (createSessionId)
	if pullErr := readBuffer.PullContext("createSessionId"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for createSessionId")
	}
	_createSessionId, _createSessionIdErr := NodeIdParseWithBuffer(ctx, readBuffer)
	if _createSessionIdErr != nil {
		return nil, errors.Wrap(_createSessionIdErr, "Error parsing 'createSessionId' field of ProgramDiagnosticDataType")
	}
	createSessionId := _createSessionId.(NodeId)
	if closeErr := readBuffer.CloseContext("createSessionId"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for createSessionId")
	}

	// Simple Field (createClientName)
	if pullErr := readBuffer.PullContext("createClientName"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for createClientName")
	}
	_createClientName, _createClientNameErr := PascalStringParseWithBuffer(ctx, readBuffer)
	if _createClientNameErr != nil {
		return nil, errors.Wrap(_createClientNameErr, "Error parsing 'createClientName' field of ProgramDiagnosticDataType")
	}
	createClientName := _createClientName.(PascalString)
	if closeErr := readBuffer.CloseContext("createClientName"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for createClientName")
	}

	// Simple Field (invocationCreationTime)
	_invocationCreationTime, _invocationCreationTimeErr := readBuffer.ReadInt64("invocationCreationTime", 64)
	if _invocationCreationTimeErr != nil {
		return nil, errors.Wrap(_invocationCreationTimeErr, "Error parsing 'invocationCreationTime' field of ProgramDiagnosticDataType")
	}
	invocationCreationTime := _invocationCreationTime

	// Simple Field (lastTransitionTime)
	_lastTransitionTime, _lastTransitionTimeErr := readBuffer.ReadInt64("lastTransitionTime", 64)
	if _lastTransitionTimeErr != nil {
		return nil, errors.Wrap(_lastTransitionTimeErr, "Error parsing 'lastTransitionTime' field of ProgramDiagnosticDataType")
	}
	lastTransitionTime := _lastTransitionTime

	// Simple Field (lastMethodCall)
	if pullErr := readBuffer.PullContext("lastMethodCall"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lastMethodCall")
	}
	_lastMethodCall, _lastMethodCallErr := PascalStringParseWithBuffer(ctx, readBuffer)
	if _lastMethodCallErr != nil {
		return nil, errors.Wrap(_lastMethodCallErr, "Error parsing 'lastMethodCall' field of ProgramDiagnosticDataType")
	}
	lastMethodCall := _lastMethodCall.(PascalString)
	if closeErr := readBuffer.CloseContext("lastMethodCall"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lastMethodCall")
	}

	// Simple Field (lastMethodSessionId)
	if pullErr := readBuffer.PullContext("lastMethodSessionId"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lastMethodSessionId")
	}
	_lastMethodSessionId, _lastMethodSessionIdErr := NodeIdParseWithBuffer(ctx, readBuffer)
	if _lastMethodSessionIdErr != nil {
		return nil, errors.Wrap(_lastMethodSessionIdErr, "Error parsing 'lastMethodSessionId' field of ProgramDiagnosticDataType")
	}
	lastMethodSessionId := _lastMethodSessionId.(NodeId)
	if closeErr := readBuffer.CloseContext("lastMethodSessionId"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lastMethodSessionId")
	}

	// Simple Field (noOfLastMethodInputArguments)
	_noOfLastMethodInputArguments, _noOfLastMethodInputArgumentsErr := readBuffer.ReadInt32("noOfLastMethodInputArguments", 32)
	if _noOfLastMethodInputArgumentsErr != nil {
		return nil, errors.Wrap(_noOfLastMethodInputArgumentsErr, "Error parsing 'noOfLastMethodInputArguments' field of ProgramDiagnosticDataType")
	}
	noOfLastMethodInputArguments := _noOfLastMethodInputArguments

	// Array field (lastMethodInputArguments)
	if pullErr := readBuffer.PullContext("lastMethodInputArguments", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lastMethodInputArguments")
	}
	// Count array
	lastMethodInputArguments := make([]ExtensionObjectDefinition, max(noOfLastMethodInputArguments, 0))
	// This happens when the size is set conditional to 0
	if len(lastMethodInputArguments) == 0 {
		lastMethodInputArguments = nil
	}
	{
		_numItems := uint16(max(noOfLastMethodInputArguments, 0))
		for _curItem := uint16(0); _curItem < _numItems; _curItem++ {
			arrayCtx := utils.CreateArrayContext(ctx, int(_numItems), int(_curItem))
			_ = arrayCtx
			_ = _curItem
			_item, _err := ExtensionObjectDefinitionParseWithBuffer(arrayCtx, readBuffer, "298")
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'lastMethodInputArguments' field of ProgramDiagnosticDataType")
			}
			lastMethodInputArguments[_curItem] = _item.(ExtensionObjectDefinition)
		}
	}
	if closeErr := readBuffer.CloseContext("lastMethodInputArguments", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lastMethodInputArguments")
	}

	// Simple Field (noOfLastMethodOutputArguments)
	_noOfLastMethodOutputArguments, _noOfLastMethodOutputArgumentsErr := readBuffer.ReadInt32("noOfLastMethodOutputArguments", 32)
	if _noOfLastMethodOutputArgumentsErr != nil {
		return nil, errors.Wrap(_noOfLastMethodOutputArgumentsErr, "Error parsing 'noOfLastMethodOutputArguments' field of ProgramDiagnosticDataType")
	}
	noOfLastMethodOutputArguments := _noOfLastMethodOutputArguments

	// Array field (lastMethodOutputArguments)
	if pullErr := readBuffer.PullContext("lastMethodOutputArguments", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lastMethodOutputArguments")
	}
	// Count array
	lastMethodOutputArguments := make([]ExtensionObjectDefinition, max(noOfLastMethodOutputArguments, 0))
	// This happens when the size is set conditional to 0
	if len(lastMethodOutputArguments) == 0 {
		lastMethodOutputArguments = nil
	}
	{
		_numItems := uint16(max(noOfLastMethodOutputArguments, 0))
		for _curItem := uint16(0); _curItem < _numItems; _curItem++ {
			arrayCtx := utils.CreateArrayContext(ctx, int(_numItems), int(_curItem))
			_ = arrayCtx
			_ = _curItem
			_item, _err := ExtensionObjectDefinitionParseWithBuffer(arrayCtx, readBuffer, "298")
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'lastMethodOutputArguments' field of ProgramDiagnosticDataType")
			}
			lastMethodOutputArguments[_curItem] = _item.(ExtensionObjectDefinition)
		}
	}
	if closeErr := readBuffer.CloseContext("lastMethodOutputArguments", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lastMethodOutputArguments")
	}

	// Simple Field (lastMethodCallTime)
	_lastMethodCallTime, _lastMethodCallTimeErr := readBuffer.ReadInt64("lastMethodCallTime", 64)
	if _lastMethodCallTimeErr != nil {
		return nil, errors.Wrap(_lastMethodCallTimeErr, "Error parsing 'lastMethodCallTime' field of ProgramDiagnosticDataType")
	}
	lastMethodCallTime := _lastMethodCallTime

	// Simple Field (lastMethodReturnStatus)
	if pullErr := readBuffer.PullContext("lastMethodReturnStatus"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lastMethodReturnStatus")
	}
	_lastMethodReturnStatus, _lastMethodReturnStatusErr := ExtensionObjectDefinitionParseWithBuffer(ctx, readBuffer, string("301"))
	if _lastMethodReturnStatusErr != nil {
		return nil, errors.Wrap(_lastMethodReturnStatusErr, "Error parsing 'lastMethodReturnStatus' field of ProgramDiagnosticDataType")
	}
	lastMethodReturnStatus := _lastMethodReturnStatus.(ExtensionObjectDefinition)
	if closeErr := readBuffer.CloseContext("lastMethodReturnStatus"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lastMethodReturnStatus")
	}

	if closeErr := readBuffer.CloseContext("ProgramDiagnosticDataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ProgramDiagnosticDataType")
	}

	// Create a partially initialized instance
	_child := &_ProgramDiagnosticDataType{
		_ExtensionObjectDefinition:    &_ExtensionObjectDefinition{},
		CreateSessionId:               createSessionId,
		CreateClientName:              createClientName,
		InvocationCreationTime:        invocationCreationTime,
		LastTransitionTime:            lastTransitionTime,
		LastMethodCall:                lastMethodCall,
		LastMethodSessionId:           lastMethodSessionId,
		NoOfLastMethodInputArguments:  noOfLastMethodInputArguments,
		LastMethodInputArguments:      lastMethodInputArguments,
		NoOfLastMethodOutputArguments: noOfLastMethodOutputArguments,
		LastMethodOutputArguments:     lastMethodOutputArguments,
		LastMethodCallTime:            lastMethodCallTime,
		LastMethodReturnStatus:        lastMethodReturnStatus,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_ProgramDiagnosticDataType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ProgramDiagnosticDataType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ProgramDiagnosticDataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ProgramDiagnosticDataType")
		}

		// Simple Field (createSessionId)
		if pushErr := writeBuffer.PushContext("createSessionId"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for createSessionId")
		}
		_createSessionIdErr := writeBuffer.WriteSerializable(ctx, m.GetCreateSessionId())
		if popErr := writeBuffer.PopContext("createSessionId"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for createSessionId")
		}
		if _createSessionIdErr != nil {
			return errors.Wrap(_createSessionIdErr, "Error serializing 'createSessionId' field")
		}

		// Simple Field (createClientName)
		if pushErr := writeBuffer.PushContext("createClientName"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for createClientName")
		}
		_createClientNameErr := writeBuffer.WriteSerializable(ctx, m.GetCreateClientName())
		if popErr := writeBuffer.PopContext("createClientName"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for createClientName")
		}
		if _createClientNameErr != nil {
			return errors.Wrap(_createClientNameErr, "Error serializing 'createClientName' field")
		}

		// Simple Field (invocationCreationTime)
		invocationCreationTime := int64(m.GetInvocationCreationTime())
		_invocationCreationTimeErr := writeBuffer.WriteInt64("invocationCreationTime", 64, int64((invocationCreationTime)))
		if _invocationCreationTimeErr != nil {
			return errors.Wrap(_invocationCreationTimeErr, "Error serializing 'invocationCreationTime' field")
		}

		// Simple Field (lastTransitionTime)
		lastTransitionTime := int64(m.GetLastTransitionTime())
		_lastTransitionTimeErr := writeBuffer.WriteInt64("lastTransitionTime", 64, int64((lastTransitionTime)))
		if _lastTransitionTimeErr != nil {
			return errors.Wrap(_lastTransitionTimeErr, "Error serializing 'lastTransitionTime' field")
		}

		// Simple Field (lastMethodCall)
		if pushErr := writeBuffer.PushContext("lastMethodCall"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lastMethodCall")
		}
		_lastMethodCallErr := writeBuffer.WriteSerializable(ctx, m.GetLastMethodCall())
		if popErr := writeBuffer.PopContext("lastMethodCall"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lastMethodCall")
		}
		if _lastMethodCallErr != nil {
			return errors.Wrap(_lastMethodCallErr, "Error serializing 'lastMethodCall' field")
		}

		// Simple Field (lastMethodSessionId)
		if pushErr := writeBuffer.PushContext("lastMethodSessionId"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lastMethodSessionId")
		}
		_lastMethodSessionIdErr := writeBuffer.WriteSerializable(ctx, m.GetLastMethodSessionId())
		if popErr := writeBuffer.PopContext("lastMethodSessionId"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lastMethodSessionId")
		}
		if _lastMethodSessionIdErr != nil {
			return errors.Wrap(_lastMethodSessionIdErr, "Error serializing 'lastMethodSessionId' field")
		}

		// Simple Field (noOfLastMethodInputArguments)
		noOfLastMethodInputArguments := int32(m.GetNoOfLastMethodInputArguments())
		_noOfLastMethodInputArgumentsErr := writeBuffer.WriteInt32("noOfLastMethodInputArguments", 32, int32((noOfLastMethodInputArguments)))
		if _noOfLastMethodInputArgumentsErr != nil {
			return errors.Wrap(_noOfLastMethodInputArgumentsErr, "Error serializing 'noOfLastMethodInputArguments' field")
		}

		// Array Field (lastMethodInputArguments)
		if pushErr := writeBuffer.PushContext("lastMethodInputArguments", utils.WithRenderAsList(true)); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lastMethodInputArguments")
		}
		for _curItem, _element := range m.GetLastMethodInputArguments() {
			_ = _curItem
			arrayCtx := utils.CreateArrayContext(ctx, len(m.GetLastMethodInputArguments()), _curItem)
			_ = arrayCtx
			_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'lastMethodInputArguments' field")
			}
		}
		if popErr := writeBuffer.PopContext("lastMethodInputArguments", utils.WithRenderAsList(true)); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lastMethodInputArguments")
		}

		// Simple Field (noOfLastMethodOutputArguments)
		noOfLastMethodOutputArguments := int32(m.GetNoOfLastMethodOutputArguments())
		_noOfLastMethodOutputArgumentsErr := writeBuffer.WriteInt32("noOfLastMethodOutputArguments", 32, int32((noOfLastMethodOutputArguments)))
		if _noOfLastMethodOutputArgumentsErr != nil {
			return errors.Wrap(_noOfLastMethodOutputArgumentsErr, "Error serializing 'noOfLastMethodOutputArguments' field")
		}

		// Array Field (lastMethodOutputArguments)
		if pushErr := writeBuffer.PushContext("lastMethodOutputArguments", utils.WithRenderAsList(true)); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lastMethodOutputArguments")
		}
		for _curItem, _element := range m.GetLastMethodOutputArguments() {
			_ = _curItem
			arrayCtx := utils.CreateArrayContext(ctx, len(m.GetLastMethodOutputArguments()), _curItem)
			_ = arrayCtx
			_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'lastMethodOutputArguments' field")
			}
		}
		if popErr := writeBuffer.PopContext("lastMethodOutputArguments", utils.WithRenderAsList(true)); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lastMethodOutputArguments")
		}

		// Simple Field (lastMethodCallTime)
		lastMethodCallTime := int64(m.GetLastMethodCallTime())
		_lastMethodCallTimeErr := writeBuffer.WriteInt64("lastMethodCallTime", 64, int64((lastMethodCallTime)))
		if _lastMethodCallTimeErr != nil {
			return errors.Wrap(_lastMethodCallTimeErr, "Error serializing 'lastMethodCallTime' field")
		}

		// Simple Field (lastMethodReturnStatus)
		if pushErr := writeBuffer.PushContext("lastMethodReturnStatus"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lastMethodReturnStatus")
		}
		_lastMethodReturnStatusErr := writeBuffer.WriteSerializable(ctx, m.GetLastMethodReturnStatus())
		if popErr := writeBuffer.PopContext("lastMethodReturnStatus"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lastMethodReturnStatus")
		}
		if _lastMethodReturnStatusErr != nil {
			return errors.Wrap(_lastMethodReturnStatusErr, "Error serializing 'lastMethodReturnStatus' field")
		}

		if popErr := writeBuffer.PopContext("ProgramDiagnosticDataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ProgramDiagnosticDataType")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ProgramDiagnosticDataType) isProgramDiagnosticDataType() bool {
	return true
}

func (m *_ProgramDiagnosticDataType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
