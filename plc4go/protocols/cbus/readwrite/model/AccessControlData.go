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

// AccessControlData is the corresponding interface of AccessControlData
type AccessControlData interface {
	AccessControlDataContract
	AccessControlDataRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// IsAccessControlData is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsAccessControlData()
	// CreateBuilder creates a AccessControlDataBuilder
	CreateAccessControlDataBuilder() AccessControlDataBuilder
}

// AccessControlDataContract provides a set of functions which can be overwritten by a sub struct
type AccessControlDataContract interface {
	// GetCommandTypeContainer returns CommandTypeContainer (property field)
	GetCommandTypeContainer() AccessControlCommandTypeContainer
	// GetNetworkId returns NetworkId (property field)
	GetNetworkId() byte
	// GetAccessPointId returns AccessPointId (property field)
	GetAccessPointId() byte
	// GetCommandType returns CommandType (virtual field)
	GetCommandType() AccessControlCommandType
	// IsAccessControlData is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsAccessControlData()
	// CreateBuilder creates a AccessControlDataBuilder
	CreateAccessControlDataBuilder() AccessControlDataBuilder
}

// AccessControlDataRequirements provides a set of functions which need to be implemented by a sub struct
type AccessControlDataRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetCommandType returns CommandType (discriminator field)
	GetCommandType() AccessControlCommandType
}

// _AccessControlData is the data-structure of this message
type _AccessControlData struct {
	_SubType             AccessControlData
	CommandTypeContainer AccessControlCommandTypeContainer
	NetworkId            byte
	AccessPointId        byte
}

var _ AccessControlDataContract = (*_AccessControlData)(nil)

// NewAccessControlData factory function for _AccessControlData
func NewAccessControlData(commandTypeContainer AccessControlCommandTypeContainer, networkId byte, accessPointId byte) *_AccessControlData {
	return &_AccessControlData{CommandTypeContainer: commandTypeContainer, NetworkId: networkId, AccessPointId: accessPointId}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// AccessControlDataBuilder is a builder for AccessControlData
type AccessControlDataBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(commandTypeContainer AccessControlCommandTypeContainer, networkId byte, accessPointId byte) AccessControlDataBuilder
	// WithCommandTypeContainer adds CommandTypeContainer (property field)
	WithCommandTypeContainer(AccessControlCommandTypeContainer) AccessControlDataBuilder
	// WithNetworkId adds NetworkId (property field)
	WithNetworkId(byte) AccessControlDataBuilder
	// WithAccessPointId adds AccessPointId (property field)
	WithAccessPointId(byte) AccessControlDataBuilder
	// AsAccessControlDataValidAccessRequest converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataValidAccessRequest() interface {
		AccessControlDataValidAccessRequestBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataInvalidAccessRequest converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataInvalidAccessRequest() interface {
		AccessControlDataInvalidAccessRequestBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataAccessPointLeftOpen converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataAccessPointLeftOpen() interface {
		AccessControlDataAccessPointLeftOpenBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataAccessPointForcedOpen converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataAccessPointForcedOpen() interface {
		AccessControlDataAccessPointForcedOpenBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataAccessPointClosed converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataAccessPointClosed() interface {
		AccessControlDataAccessPointClosedBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataRequestToExit converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataRequestToExit() interface {
		AccessControlDataRequestToExitBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataCloseAccessPoint converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataCloseAccessPoint() interface {
		AccessControlDataCloseAccessPointBuilder
		Done() AccessControlDataBuilder
	}
	// AsAccessControlDataLockAccessPoint converts this build to a subType of AccessControlData. It is always possible to return to current builder using Done()
	AsAccessControlDataLockAccessPoint() interface {
		AccessControlDataLockAccessPointBuilder
		Done() AccessControlDataBuilder
	}
	// Build builds the AccessControlData or returns an error if something is wrong
	PartialBuild() (AccessControlDataContract, error)
	// MustBuild does the same as Build but panics on error
	PartialMustBuild() AccessControlDataContract
	// Build builds the AccessControlData or returns an error if something is wrong
	Build() (AccessControlData, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() AccessControlData
}

// NewAccessControlDataBuilder() creates a AccessControlDataBuilder
func NewAccessControlDataBuilder() AccessControlDataBuilder {
	return &_AccessControlDataBuilder{_AccessControlData: new(_AccessControlData)}
}

type _AccessControlDataChildBuilder interface {
	utils.Copyable
	setParent(AccessControlDataContract)
	buildForAccessControlData() (AccessControlData, error)
}

type _AccessControlDataBuilder struct {
	*_AccessControlData

	childBuilder _AccessControlDataChildBuilder

	err *utils.MultiError
}

var _ (AccessControlDataBuilder) = (*_AccessControlDataBuilder)(nil)

func (b *_AccessControlDataBuilder) WithMandatoryFields(commandTypeContainer AccessControlCommandTypeContainer, networkId byte, accessPointId byte) AccessControlDataBuilder {
	return b.WithCommandTypeContainer(commandTypeContainer).WithNetworkId(networkId).WithAccessPointId(accessPointId)
}

func (b *_AccessControlDataBuilder) WithCommandTypeContainer(commandTypeContainer AccessControlCommandTypeContainer) AccessControlDataBuilder {
	b.CommandTypeContainer = commandTypeContainer
	return b
}

func (b *_AccessControlDataBuilder) WithNetworkId(networkId byte) AccessControlDataBuilder {
	b.NetworkId = networkId
	return b
}

func (b *_AccessControlDataBuilder) WithAccessPointId(accessPointId byte) AccessControlDataBuilder {
	b.AccessPointId = accessPointId
	return b
}

func (b *_AccessControlDataBuilder) PartialBuild() (AccessControlDataContract, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._AccessControlData.deepCopy(), nil
}

func (b *_AccessControlDataBuilder) PartialMustBuild() AccessControlDataContract {
	build, err := b.PartialBuild()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_AccessControlDataBuilder) AsAccessControlDataValidAccessRequest() interface {
	AccessControlDataValidAccessRequestBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataValidAccessRequestBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataValidAccessRequestBuilder().(*_AccessControlDataValidAccessRequestBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataInvalidAccessRequest() interface {
	AccessControlDataInvalidAccessRequestBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataInvalidAccessRequestBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataInvalidAccessRequestBuilder().(*_AccessControlDataInvalidAccessRequestBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataAccessPointLeftOpen() interface {
	AccessControlDataAccessPointLeftOpenBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataAccessPointLeftOpenBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataAccessPointLeftOpenBuilder().(*_AccessControlDataAccessPointLeftOpenBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataAccessPointForcedOpen() interface {
	AccessControlDataAccessPointForcedOpenBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataAccessPointForcedOpenBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataAccessPointForcedOpenBuilder().(*_AccessControlDataAccessPointForcedOpenBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataAccessPointClosed() interface {
	AccessControlDataAccessPointClosedBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataAccessPointClosedBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataAccessPointClosedBuilder().(*_AccessControlDataAccessPointClosedBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataRequestToExit() interface {
	AccessControlDataRequestToExitBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataRequestToExitBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataRequestToExitBuilder().(*_AccessControlDataRequestToExitBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataCloseAccessPoint() interface {
	AccessControlDataCloseAccessPointBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataCloseAccessPointBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataCloseAccessPointBuilder().(*_AccessControlDataCloseAccessPointBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) AsAccessControlDataLockAccessPoint() interface {
	AccessControlDataLockAccessPointBuilder
	Done() AccessControlDataBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		AccessControlDataLockAccessPointBuilder
		Done() AccessControlDataBuilder
	}); ok {
		return cb
	}
	cb := NewAccessControlDataLockAccessPointBuilder().(*_AccessControlDataLockAccessPointBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_AccessControlDataBuilder) Build() (AccessControlData, error) {
	v, err := b.PartialBuild()
	if err != nil {
		return nil, errors.Wrap(err, "error occurred during partial build")
	}
	if b.childBuilder == nil {
		return nil, errors.New("no child builder present")
	}
	b.childBuilder.setParent(v)
	return b.childBuilder.buildForAccessControlData()
}

func (b *_AccessControlDataBuilder) MustBuild() AccessControlData {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_AccessControlDataBuilder) DeepCopy() any {
	_copy := b.CreateAccessControlDataBuilder().(*_AccessControlDataBuilder)
	_copy.childBuilder = b.childBuilder.DeepCopy().(_AccessControlDataChildBuilder)
	_copy.childBuilder.setParent(_copy)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateAccessControlDataBuilder creates a AccessControlDataBuilder
func (b *_AccessControlData) CreateAccessControlDataBuilder() AccessControlDataBuilder {
	if b == nil {
		return NewAccessControlDataBuilder()
	}
	return &_AccessControlDataBuilder{_AccessControlData: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AccessControlData) GetCommandTypeContainer() AccessControlCommandTypeContainer {
	return m.CommandTypeContainer
}

func (m *_AccessControlData) GetNetworkId() byte {
	return m.NetworkId
}

func (m *_AccessControlData) GetAccessPointId() byte {
	return m.AccessPointId
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (pm *_AccessControlData) GetCommandType() AccessControlCommandType {
	m := pm._SubType
	ctx := context.Background()
	_ = ctx
	return CastAccessControlCommandType(m.GetCommandTypeContainer().CommandType())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastAccessControlData(structType any) AccessControlData {
	if casted, ok := structType.(AccessControlData); ok {
		return casted
	}
	if casted, ok := structType.(*AccessControlData); ok {
		return *casted
	}
	return nil
}

func (m *_AccessControlData) GetTypeName() string {
	return "AccessControlData"
}

func (m *_AccessControlData) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (commandTypeContainer)
	lengthInBits += 8

	// A virtual field doesn't have any in- or output.

	// Simple field (networkId)
	lengthInBits += 8

	// Simple field (accessPointId)
	lengthInBits += 8

	return lengthInBits
}

func (m *_AccessControlData) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func AccessControlDataParse[T AccessControlData](ctx context.Context, theBytes []byte) (T, error) {
	return AccessControlDataParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes))
}

func AccessControlDataParseWithBufferProducer[T AccessControlData]() func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := AccessControlDataParseWithBuffer[T](ctx, readBuffer)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, nil
	}
}

func AccessControlDataParseWithBuffer[T AccessControlData](ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	v, err := (&_AccessControlData{}).parse(ctx, readBuffer)
	if err != nil {
		var zero T
		return zero, err
	}
	vc, ok := v.(T)
	if !ok {
		var zero T
		return zero, errors.Errorf("Unexpected type %T. Expected type %T", v, *new(T))
	}
	return vc, nil
}

func (m *_AccessControlData) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__accessControlData AccessControlData, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AccessControlData"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AccessControlData")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Validation
	if !(KnowsAccessControlCommandTypeContainer(ctx, readBuffer)) {
		return nil, errors.WithStack(utils.ParseAssertError{Message: "no command type could be found"})
	}

	commandTypeContainer, err := ReadEnumField[AccessControlCommandTypeContainer](ctx, "commandTypeContainer", "AccessControlCommandTypeContainer", ReadEnum(AccessControlCommandTypeContainerByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'commandTypeContainer' field"))
	}
	m.CommandTypeContainer = commandTypeContainer

	commandType, err := ReadVirtualField[AccessControlCommandType](ctx, "commandType", (*AccessControlCommandType)(nil), commandTypeContainer.CommandType())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'commandType' field"))
	}
	_ = commandType

	networkId, err := ReadSimpleField(ctx, "networkId", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'networkId' field"))
	}
	m.NetworkId = networkId

	accessPointId, err := ReadSimpleField(ctx, "accessPointId", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'accessPointId' field"))
	}
	m.AccessPointId = accessPointId

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _child AccessControlData
	switch {
	case commandType == AccessControlCommandType_VALID_ACCESS: // AccessControlDataValidAccessRequest
		if _child, err = new(_AccessControlDataValidAccessRequest).parse(ctx, readBuffer, m, commandTypeContainer); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataValidAccessRequest for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_INVALID_ACCESS: // AccessControlDataInvalidAccessRequest
		if _child, err = new(_AccessControlDataInvalidAccessRequest).parse(ctx, readBuffer, m, commandTypeContainer); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataInvalidAccessRequest for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_ACCESS_POINT_LEFT_OPEN: // AccessControlDataAccessPointLeftOpen
		if _child, err = new(_AccessControlDataAccessPointLeftOpen).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataAccessPointLeftOpen for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_ACCESS_POINT_FORCED_OPEN: // AccessControlDataAccessPointForcedOpen
		if _child, err = new(_AccessControlDataAccessPointForcedOpen).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataAccessPointForcedOpen for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_ACCESS_POINT_CLOSED: // AccessControlDataAccessPointClosed
		if _child, err = new(_AccessControlDataAccessPointClosed).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataAccessPointClosed for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_REQUEST_TO_EXIT: // AccessControlDataRequestToExit
		if _child, err = new(_AccessControlDataRequestToExit).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataRequestToExit for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_CLOSE_ACCESS_POINT: // AccessControlDataCloseAccessPoint
		if _child, err = new(_AccessControlDataCloseAccessPoint).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataCloseAccessPoint for type-switch of AccessControlData")
		}
	case commandType == AccessControlCommandType_LOCK_ACCESS_POINT: // AccessControlDataLockAccessPoint
		if _child, err = new(_AccessControlDataLockAccessPoint).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type AccessControlDataLockAccessPoint for type-switch of AccessControlData")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [commandType=%v]", commandType)
	}

	if closeErr := readBuffer.CloseContext("AccessControlData"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AccessControlData")
	}

	return _child, nil
}

func (pm *_AccessControlData) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child AccessControlData, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("AccessControlData"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for AccessControlData")
	}

	if err := WriteSimpleEnumField[AccessControlCommandTypeContainer](ctx, "commandTypeContainer", "AccessControlCommandTypeContainer", m.GetCommandTypeContainer(), WriteEnum[AccessControlCommandTypeContainer, uint8](AccessControlCommandTypeContainer.GetValue, AccessControlCommandTypeContainer.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
		return errors.Wrap(err, "Error serializing 'commandTypeContainer' field")
	}
	// Virtual field
	commandType := m.GetCommandType()
	_ = commandType
	if _commandTypeErr := writeBuffer.WriteVirtual(ctx, "commandType", m.GetCommandType()); _commandTypeErr != nil {
		return errors.Wrap(_commandTypeErr, "Error serializing 'commandType' field")
	}

	if err := WriteSimpleField[byte](ctx, "networkId", m.GetNetworkId(), WriteByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'networkId' field")
	}

	if err := WriteSimpleField[byte](ctx, "accessPointId", m.GetAccessPointId(), WriteByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'accessPointId' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("AccessControlData"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for AccessControlData")
	}
	return nil
}

func (m *_AccessControlData) IsAccessControlData() {}

func (m *_AccessControlData) DeepCopy() any {
	return m.deepCopy()
}

func (m *_AccessControlData) deepCopy() *_AccessControlData {
	if m == nil {
		return nil
	}
	_AccessControlDataCopy := &_AccessControlData{
		nil, // will be set by child
		m.CommandTypeContainer,
		m.NetworkId,
		m.AccessPointId,
	}
	return _AccessControlDataCopy
}
