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

// CipConnectionManagerRequest is the corresponding interface of CipConnectionManagerRequest
type CipConnectionManagerRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	CipService
	// GetClassSegment returns ClassSegment (property field)
	GetClassSegment() PathSegment
	// GetInstanceSegment returns InstanceSegment (property field)
	GetInstanceSegment() PathSegment
	// GetPriority returns Priority (property field)
	GetPriority() uint8
	// GetTickTime returns TickTime (property field)
	GetTickTime() uint8
	// GetTimeoutTicks returns TimeoutTicks (property field)
	GetTimeoutTicks() uint8
	// GetOtConnectionId returns OtConnectionId (property field)
	GetOtConnectionId() uint32
	// GetToConnectionId returns ToConnectionId (property field)
	GetToConnectionId() uint32
	// GetConnectionSerialNumber returns ConnectionSerialNumber (property field)
	GetConnectionSerialNumber() uint16
	// GetOriginatorVendorId returns OriginatorVendorId (property field)
	GetOriginatorVendorId() uint16
	// GetOriginatorSerialNumber returns OriginatorSerialNumber (property field)
	GetOriginatorSerialNumber() uint32
	// GetTimeoutMultiplier returns TimeoutMultiplier (property field)
	GetTimeoutMultiplier() uint8
	// GetOtRpi returns OtRpi (property field)
	GetOtRpi() uint32
	// GetOtConnectionParameters returns OtConnectionParameters (property field)
	GetOtConnectionParameters() NetworkConnectionParameters
	// GetToRpi returns ToRpi (property field)
	GetToRpi() uint32
	// GetToConnectionParameters returns ToConnectionParameters (property field)
	GetToConnectionParameters() NetworkConnectionParameters
	// GetTransportType returns TransportType (property field)
	GetTransportType() TransportType
	// GetConnectionPathSize returns ConnectionPathSize (property field)
	GetConnectionPathSize() uint8
	// GetConnectionPaths returns ConnectionPaths (property field)
	GetConnectionPaths() []PathSegment
	// IsCipConnectionManagerRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCipConnectionManagerRequest()
	// CreateBuilder creates a CipConnectionManagerRequestBuilder
	CreateCipConnectionManagerRequestBuilder() CipConnectionManagerRequestBuilder
}

// _CipConnectionManagerRequest is the data-structure of this message
type _CipConnectionManagerRequest struct {
	CipServiceContract
	ClassSegment           PathSegment
	InstanceSegment        PathSegment
	Priority               uint8
	TickTime               uint8
	TimeoutTicks           uint8
	OtConnectionId         uint32
	ToConnectionId         uint32
	ConnectionSerialNumber uint16
	OriginatorVendorId     uint16
	OriginatorSerialNumber uint32
	TimeoutMultiplier      uint8
	OtRpi                  uint32
	OtConnectionParameters NetworkConnectionParameters
	ToRpi                  uint32
	ToConnectionParameters NetworkConnectionParameters
	TransportType          TransportType
	ConnectionPathSize     uint8
	ConnectionPaths        []PathSegment
	// Reserved Fields
	reservedField0 *uint32
}

var _ CipConnectionManagerRequest = (*_CipConnectionManagerRequest)(nil)
var _ CipServiceRequirements = (*_CipConnectionManagerRequest)(nil)

// NewCipConnectionManagerRequest factory function for _CipConnectionManagerRequest
func NewCipConnectionManagerRequest(classSegment PathSegment, instanceSegment PathSegment, priority uint8, tickTime uint8, timeoutTicks uint8, otConnectionId uint32, toConnectionId uint32, connectionSerialNumber uint16, originatorVendorId uint16, originatorSerialNumber uint32, timeoutMultiplier uint8, otRpi uint32, otConnectionParameters NetworkConnectionParameters, toRpi uint32, toConnectionParameters NetworkConnectionParameters, transportType TransportType, connectionPathSize uint8, connectionPaths []PathSegment, serviceLen uint16) *_CipConnectionManagerRequest {
	if classSegment == nil {
		panic("classSegment of type PathSegment for CipConnectionManagerRequest must not be nil")
	}
	if instanceSegment == nil {
		panic("instanceSegment of type PathSegment for CipConnectionManagerRequest must not be nil")
	}
	if otConnectionParameters == nil {
		panic("otConnectionParameters of type NetworkConnectionParameters for CipConnectionManagerRequest must not be nil")
	}
	if toConnectionParameters == nil {
		panic("toConnectionParameters of type NetworkConnectionParameters for CipConnectionManagerRequest must not be nil")
	}
	if transportType == nil {
		panic("transportType of type TransportType for CipConnectionManagerRequest must not be nil")
	}
	_result := &_CipConnectionManagerRequest{
		CipServiceContract:     NewCipService(serviceLen),
		ClassSegment:           classSegment,
		InstanceSegment:        instanceSegment,
		Priority:               priority,
		TickTime:               tickTime,
		TimeoutTicks:           timeoutTicks,
		OtConnectionId:         otConnectionId,
		ToConnectionId:         toConnectionId,
		ConnectionSerialNumber: connectionSerialNumber,
		OriginatorVendorId:     originatorVendorId,
		OriginatorSerialNumber: originatorSerialNumber,
		TimeoutMultiplier:      timeoutMultiplier,
		OtRpi:                  otRpi,
		OtConnectionParameters: otConnectionParameters,
		ToRpi:                  toRpi,
		ToConnectionParameters: toConnectionParameters,
		TransportType:          transportType,
		ConnectionPathSize:     connectionPathSize,
		ConnectionPaths:        connectionPaths,
	}
	_result.CipServiceContract.(*_CipService)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CipConnectionManagerRequestBuilder is a builder for CipConnectionManagerRequest
type CipConnectionManagerRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(classSegment PathSegment, instanceSegment PathSegment, priority uint8, tickTime uint8, timeoutTicks uint8, otConnectionId uint32, toConnectionId uint32, connectionSerialNumber uint16, originatorVendorId uint16, originatorSerialNumber uint32, timeoutMultiplier uint8, otRpi uint32, otConnectionParameters NetworkConnectionParameters, toRpi uint32, toConnectionParameters NetworkConnectionParameters, transportType TransportType, connectionPathSize uint8, connectionPaths []PathSegment) CipConnectionManagerRequestBuilder
	// WithClassSegment adds ClassSegment (property field)
	WithClassSegment(PathSegment) CipConnectionManagerRequestBuilder
	// WithClassSegmentBuilder adds ClassSegment (property field) which is build by the builder
	WithClassSegmentBuilder(func(PathSegmentBuilder) PathSegmentBuilder) CipConnectionManagerRequestBuilder
	// WithInstanceSegment adds InstanceSegment (property field)
	WithInstanceSegment(PathSegment) CipConnectionManagerRequestBuilder
	// WithInstanceSegmentBuilder adds InstanceSegment (property field) which is build by the builder
	WithInstanceSegmentBuilder(func(PathSegmentBuilder) PathSegmentBuilder) CipConnectionManagerRequestBuilder
	// WithPriority adds Priority (property field)
	WithPriority(uint8) CipConnectionManagerRequestBuilder
	// WithTickTime adds TickTime (property field)
	WithTickTime(uint8) CipConnectionManagerRequestBuilder
	// WithTimeoutTicks adds TimeoutTicks (property field)
	WithTimeoutTicks(uint8) CipConnectionManagerRequestBuilder
	// WithOtConnectionId adds OtConnectionId (property field)
	WithOtConnectionId(uint32) CipConnectionManagerRequestBuilder
	// WithToConnectionId adds ToConnectionId (property field)
	WithToConnectionId(uint32) CipConnectionManagerRequestBuilder
	// WithConnectionSerialNumber adds ConnectionSerialNumber (property field)
	WithConnectionSerialNumber(uint16) CipConnectionManagerRequestBuilder
	// WithOriginatorVendorId adds OriginatorVendorId (property field)
	WithOriginatorVendorId(uint16) CipConnectionManagerRequestBuilder
	// WithOriginatorSerialNumber adds OriginatorSerialNumber (property field)
	WithOriginatorSerialNumber(uint32) CipConnectionManagerRequestBuilder
	// WithTimeoutMultiplier adds TimeoutMultiplier (property field)
	WithTimeoutMultiplier(uint8) CipConnectionManagerRequestBuilder
	// WithOtRpi adds OtRpi (property field)
	WithOtRpi(uint32) CipConnectionManagerRequestBuilder
	// WithOtConnectionParameters adds OtConnectionParameters (property field)
	WithOtConnectionParameters(NetworkConnectionParameters) CipConnectionManagerRequestBuilder
	// WithOtConnectionParametersBuilder adds OtConnectionParameters (property field) which is build by the builder
	WithOtConnectionParametersBuilder(func(NetworkConnectionParametersBuilder) NetworkConnectionParametersBuilder) CipConnectionManagerRequestBuilder
	// WithToRpi adds ToRpi (property field)
	WithToRpi(uint32) CipConnectionManagerRequestBuilder
	// WithToConnectionParameters adds ToConnectionParameters (property field)
	WithToConnectionParameters(NetworkConnectionParameters) CipConnectionManagerRequestBuilder
	// WithToConnectionParametersBuilder adds ToConnectionParameters (property field) which is build by the builder
	WithToConnectionParametersBuilder(func(NetworkConnectionParametersBuilder) NetworkConnectionParametersBuilder) CipConnectionManagerRequestBuilder
	// WithTransportType adds TransportType (property field)
	WithTransportType(TransportType) CipConnectionManagerRequestBuilder
	// WithTransportTypeBuilder adds TransportType (property field) which is build by the builder
	WithTransportTypeBuilder(func(TransportTypeBuilder) TransportTypeBuilder) CipConnectionManagerRequestBuilder
	// WithConnectionPathSize adds ConnectionPathSize (property field)
	WithConnectionPathSize(uint8) CipConnectionManagerRequestBuilder
	// WithConnectionPaths adds ConnectionPaths (property field)
	WithConnectionPaths(...PathSegment) CipConnectionManagerRequestBuilder
	// Build builds the CipConnectionManagerRequest or returns an error if something is wrong
	Build() (CipConnectionManagerRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CipConnectionManagerRequest
}

// NewCipConnectionManagerRequestBuilder() creates a CipConnectionManagerRequestBuilder
func NewCipConnectionManagerRequestBuilder() CipConnectionManagerRequestBuilder {
	return &_CipConnectionManagerRequestBuilder{_CipConnectionManagerRequest: new(_CipConnectionManagerRequest)}
}

type _CipConnectionManagerRequestBuilder struct {
	*_CipConnectionManagerRequest

	parentBuilder *_CipServiceBuilder

	err *utils.MultiError
}

var _ (CipConnectionManagerRequestBuilder) = (*_CipConnectionManagerRequestBuilder)(nil)

func (b *_CipConnectionManagerRequestBuilder) setParent(contract CipServiceContract) {
	b.CipServiceContract = contract
}

func (b *_CipConnectionManagerRequestBuilder) WithMandatoryFields(classSegment PathSegment, instanceSegment PathSegment, priority uint8, tickTime uint8, timeoutTicks uint8, otConnectionId uint32, toConnectionId uint32, connectionSerialNumber uint16, originatorVendorId uint16, originatorSerialNumber uint32, timeoutMultiplier uint8, otRpi uint32, otConnectionParameters NetworkConnectionParameters, toRpi uint32, toConnectionParameters NetworkConnectionParameters, transportType TransportType, connectionPathSize uint8, connectionPaths []PathSegment) CipConnectionManagerRequestBuilder {
	return b.WithClassSegment(classSegment).WithInstanceSegment(instanceSegment).WithPriority(priority).WithTickTime(tickTime).WithTimeoutTicks(timeoutTicks).WithOtConnectionId(otConnectionId).WithToConnectionId(toConnectionId).WithConnectionSerialNumber(connectionSerialNumber).WithOriginatorVendorId(originatorVendorId).WithOriginatorSerialNumber(originatorSerialNumber).WithTimeoutMultiplier(timeoutMultiplier).WithOtRpi(otRpi).WithOtConnectionParameters(otConnectionParameters).WithToRpi(toRpi).WithToConnectionParameters(toConnectionParameters).WithTransportType(transportType).WithConnectionPathSize(connectionPathSize).WithConnectionPaths(connectionPaths...)
}

func (b *_CipConnectionManagerRequestBuilder) WithClassSegment(classSegment PathSegment) CipConnectionManagerRequestBuilder {
	b.ClassSegment = classSegment
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithClassSegmentBuilder(builderSupplier func(PathSegmentBuilder) PathSegmentBuilder) CipConnectionManagerRequestBuilder {
	builder := builderSupplier(b.ClassSegment.CreatePathSegmentBuilder())
	var err error
	b.ClassSegment, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PathSegmentBuilder failed"))
	}
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithInstanceSegment(instanceSegment PathSegment) CipConnectionManagerRequestBuilder {
	b.InstanceSegment = instanceSegment
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithInstanceSegmentBuilder(builderSupplier func(PathSegmentBuilder) PathSegmentBuilder) CipConnectionManagerRequestBuilder {
	builder := builderSupplier(b.InstanceSegment.CreatePathSegmentBuilder())
	var err error
	b.InstanceSegment, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PathSegmentBuilder failed"))
	}
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithPriority(priority uint8) CipConnectionManagerRequestBuilder {
	b.Priority = priority
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithTickTime(tickTime uint8) CipConnectionManagerRequestBuilder {
	b.TickTime = tickTime
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithTimeoutTicks(timeoutTicks uint8) CipConnectionManagerRequestBuilder {
	b.TimeoutTicks = timeoutTicks
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithOtConnectionId(otConnectionId uint32) CipConnectionManagerRequestBuilder {
	b.OtConnectionId = otConnectionId
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithToConnectionId(toConnectionId uint32) CipConnectionManagerRequestBuilder {
	b.ToConnectionId = toConnectionId
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithConnectionSerialNumber(connectionSerialNumber uint16) CipConnectionManagerRequestBuilder {
	b.ConnectionSerialNumber = connectionSerialNumber
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithOriginatorVendorId(originatorVendorId uint16) CipConnectionManagerRequestBuilder {
	b.OriginatorVendorId = originatorVendorId
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithOriginatorSerialNumber(originatorSerialNumber uint32) CipConnectionManagerRequestBuilder {
	b.OriginatorSerialNumber = originatorSerialNumber
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithTimeoutMultiplier(timeoutMultiplier uint8) CipConnectionManagerRequestBuilder {
	b.TimeoutMultiplier = timeoutMultiplier
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithOtRpi(otRpi uint32) CipConnectionManagerRequestBuilder {
	b.OtRpi = otRpi
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithOtConnectionParameters(otConnectionParameters NetworkConnectionParameters) CipConnectionManagerRequestBuilder {
	b.OtConnectionParameters = otConnectionParameters
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithOtConnectionParametersBuilder(builderSupplier func(NetworkConnectionParametersBuilder) NetworkConnectionParametersBuilder) CipConnectionManagerRequestBuilder {
	builder := builderSupplier(b.OtConnectionParameters.CreateNetworkConnectionParametersBuilder())
	var err error
	b.OtConnectionParameters, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "NetworkConnectionParametersBuilder failed"))
	}
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithToRpi(toRpi uint32) CipConnectionManagerRequestBuilder {
	b.ToRpi = toRpi
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithToConnectionParameters(toConnectionParameters NetworkConnectionParameters) CipConnectionManagerRequestBuilder {
	b.ToConnectionParameters = toConnectionParameters
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithToConnectionParametersBuilder(builderSupplier func(NetworkConnectionParametersBuilder) NetworkConnectionParametersBuilder) CipConnectionManagerRequestBuilder {
	builder := builderSupplier(b.ToConnectionParameters.CreateNetworkConnectionParametersBuilder())
	var err error
	b.ToConnectionParameters, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "NetworkConnectionParametersBuilder failed"))
	}
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithTransportType(transportType TransportType) CipConnectionManagerRequestBuilder {
	b.TransportType = transportType
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithTransportTypeBuilder(builderSupplier func(TransportTypeBuilder) TransportTypeBuilder) CipConnectionManagerRequestBuilder {
	builder := builderSupplier(b.TransportType.CreateTransportTypeBuilder())
	var err error
	b.TransportType, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "TransportTypeBuilder failed"))
	}
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithConnectionPathSize(connectionPathSize uint8) CipConnectionManagerRequestBuilder {
	b.ConnectionPathSize = connectionPathSize
	return b
}

func (b *_CipConnectionManagerRequestBuilder) WithConnectionPaths(connectionPaths ...PathSegment) CipConnectionManagerRequestBuilder {
	b.ConnectionPaths = connectionPaths
	return b
}

func (b *_CipConnectionManagerRequestBuilder) Build() (CipConnectionManagerRequest, error) {
	if b.ClassSegment == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'classSegment' not set"))
	}
	if b.InstanceSegment == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'instanceSegment' not set"))
	}
	if b.OtConnectionParameters == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'otConnectionParameters' not set"))
	}
	if b.ToConnectionParameters == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'toConnectionParameters' not set"))
	}
	if b.TransportType == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'transportType' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CipConnectionManagerRequest.deepCopy(), nil
}

func (b *_CipConnectionManagerRequestBuilder) MustBuild() CipConnectionManagerRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_CipConnectionManagerRequestBuilder) Done() CipServiceBuilder {
	return b.parentBuilder
}

func (b *_CipConnectionManagerRequestBuilder) buildForCipService() (CipService, error) {
	return b.Build()
}

func (b *_CipConnectionManagerRequestBuilder) DeepCopy() any {
	_copy := b.CreateCipConnectionManagerRequestBuilder().(*_CipConnectionManagerRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCipConnectionManagerRequestBuilder creates a CipConnectionManagerRequestBuilder
func (b *_CipConnectionManagerRequest) CreateCipConnectionManagerRequestBuilder() CipConnectionManagerRequestBuilder {
	if b == nil {
		return NewCipConnectionManagerRequestBuilder()
	}
	return &_CipConnectionManagerRequestBuilder{_CipConnectionManagerRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_CipConnectionManagerRequest) GetService() uint8 {
	return 0x5B
}

func (m *_CipConnectionManagerRequest) GetResponse() bool {
	return bool(false)
}

func (m *_CipConnectionManagerRequest) GetConnected() bool {
	return false
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CipConnectionManagerRequest) GetParent() CipServiceContract {
	return m.CipServiceContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CipConnectionManagerRequest) GetClassSegment() PathSegment {
	return m.ClassSegment
}

func (m *_CipConnectionManagerRequest) GetInstanceSegment() PathSegment {
	return m.InstanceSegment
}

func (m *_CipConnectionManagerRequest) GetPriority() uint8 {
	return m.Priority
}

func (m *_CipConnectionManagerRequest) GetTickTime() uint8 {
	return m.TickTime
}

func (m *_CipConnectionManagerRequest) GetTimeoutTicks() uint8 {
	return m.TimeoutTicks
}

func (m *_CipConnectionManagerRequest) GetOtConnectionId() uint32 {
	return m.OtConnectionId
}

func (m *_CipConnectionManagerRequest) GetToConnectionId() uint32 {
	return m.ToConnectionId
}

func (m *_CipConnectionManagerRequest) GetConnectionSerialNumber() uint16 {
	return m.ConnectionSerialNumber
}

func (m *_CipConnectionManagerRequest) GetOriginatorVendorId() uint16 {
	return m.OriginatorVendorId
}

func (m *_CipConnectionManagerRequest) GetOriginatorSerialNumber() uint32 {
	return m.OriginatorSerialNumber
}

func (m *_CipConnectionManagerRequest) GetTimeoutMultiplier() uint8 {
	return m.TimeoutMultiplier
}

func (m *_CipConnectionManagerRequest) GetOtRpi() uint32 {
	return m.OtRpi
}

func (m *_CipConnectionManagerRequest) GetOtConnectionParameters() NetworkConnectionParameters {
	return m.OtConnectionParameters
}

func (m *_CipConnectionManagerRequest) GetToRpi() uint32 {
	return m.ToRpi
}

func (m *_CipConnectionManagerRequest) GetToConnectionParameters() NetworkConnectionParameters {
	return m.ToConnectionParameters
}

func (m *_CipConnectionManagerRequest) GetTransportType() TransportType {
	return m.TransportType
}

func (m *_CipConnectionManagerRequest) GetConnectionPathSize() uint8 {
	return m.ConnectionPathSize
}

func (m *_CipConnectionManagerRequest) GetConnectionPaths() []PathSegment {
	return m.ConnectionPaths
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastCipConnectionManagerRequest(structType any) CipConnectionManagerRequest {
	if casted, ok := structType.(CipConnectionManagerRequest); ok {
		return casted
	}
	if casted, ok := structType.(*CipConnectionManagerRequest); ok {
		return *casted
	}
	return nil
}

func (m *_CipConnectionManagerRequest) GetTypeName() string {
	return "CipConnectionManagerRequest"
}

func (m *_CipConnectionManagerRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.CipServiceContract.(*_CipService).GetLengthInBits(ctx))

	// Implicit Field (requestPathSize)
	lengthInBits += 8

	// Simple field (classSegment)
	lengthInBits += m.ClassSegment.GetLengthInBits(ctx)

	// Simple field (instanceSegment)
	lengthInBits += m.InstanceSegment.GetLengthInBits(ctx)

	// Simple field (priority)
	lengthInBits += 4

	// Simple field (tickTime)
	lengthInBits += 4

	// Simple field (timeoutTicks)
	lengthInBits += 8

	// Simple field (otConnectionId)
	lengthInBits += 32

	// Simple field (toConnectionId)
	lengthInBits += 32

	// Simple field (connectionSerialNumber)
	lengthInBits += 16

	// Simple field (originatorVendorId)
	lengthInBits += 16

	// Simple field (originatorSerialNumber)
	lengthInBits += 32

	// Simple field (timeoutMultiplier)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 24

	// Simple field (otRpi)
	lengthInBits += 32

	// Simple field (otConnectionParameters)
	lengthInBits += m.OtConnectionParameters.GetLengthInBits(ctx)

	// Simple field (toRpi)
	lengthInBits += 32

	// Simple field (toConnectionParameters)
	lengthInBits += m.ToConnectionParameters.GetLengthInBits(ctx)

	// Simple field (transportType)
	lengthInBits += m.TransportType.GetLengthInBits(ctx)

	// Simple field (connectionPathSize)
	lengthInBits += 8

	// Array field
	if len(m.ConnectionPaths) > 0 {
		for _, element := range m.ConnectionPaths {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_CipConnectionManagerRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CipConnectionManagerRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_CipService, connected bool, serviceLen uint16) (__cipConnectionManagerRequest CipConnectionManagerRequest, err error) {
	m.CipServiceContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CipConnectionManagerRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CipConnectionManagerRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestPathSize, err := ReadImplicitField[uint8](ctx, "requestPathSize", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestPathSize' field"))
	}
	_ = requestPathSize

	classSegment, err := ReadSimpleField[PathSegment](ctx, "classSegment", ReadComplex[PathSegment](PathSegmentParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'classSegment' field"))
	}
	m.ClassSegment = classSegment

	instanceSegment, err := ReadSimpleField[PathSegment](ctx, "instanceSegment", ReadComplex[PathSegment](PathSegmentParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'instanceSegment' field"))
	}
	m.InstanceSegment = instanceSegment

	priority, err := ReadSimpleField(ctx, "priority", ReadUnsignedByte(readBuffer, uint8(4)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'priority' field"))
	}
	m.Priority = priority

	tickTime, err := ReadSimpleField(ctx, "tickTime", ReadUnsignedByte(readBuffer, uint8(4)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'tickTime' field"))
	}
	m.TickTime = tickTime

	timeoutTicks, err := ReadSimpleField(ctx, "timeoutTicks", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timeoutTicks' field"))
	}
	m.TimeoutTicks = timeoutTicks

	otConnectionId, err := ReadSimpleField(ctx, "otConnectionId", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'otConnectionId' field"))
	}
	m.OtConnectionId = otConnectionId

	toConnectionId, err := ReadSimpleField(ctx, "toConnectionId", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'toConnectionId' field"))
	}
	m.ToConnectionId = toConnectionId

	connectionSerialNumber, err := ReadSimpleField(ctx, "connectionSerialNumber", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'connectionSerialNumber' field"))
	}
	m.ConnectionSerialNumber = connectionSerialNumber

	originatorVendorId, err := ReadSimpleField(ctx, "originatorVendorId", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'originatorVendorId' field"))
	}
	m.OriginatorVendorId = originatorVendorId

	originatorSerialNumber, err := ReadSimpleField(ctx, "originatorSerialNumber", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'originatorSerialNumber' field"))
	}
	m.OriginatorSerialNumber = originatorSerialNumber

	timeoutMultiplier, err := ReadSimpleField(ctx, "timeoutMultiplier", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timeoutMultiplier' field"))
	}
	m.TimeoutMultiplier = timeoutMultiplier

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedInt(readBuffer, uint8(24)), uint32(0x000000))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}
	m.reservedField0 = reservedField0

	otRpi, err := ReadSimpleField(ctx, "otRpi", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'otRpi' field"))
	}
	m.OtRpi = otRpi

	otConnectionParameters, err := ReadSimpleField[NetworkConnectionParameters](ctx, "otConnectionParameters", ReadComplex[NetworkConnectionParameters](NetworkConnectionParametersParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'otConnectionParameters' field"))
	}
	m.OtConnectionParameters = otConnectionParameters

	toRpi, err := ReadSimpleField(ctx, "toRpi", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'toRpi' field"))
	}
	m.ToRpi = toRpi

	toConnectionParameters, err := ReadSimpleField[NetworkConnectionParameters](ctx, "toConnectionParameters", ReadComplex[NetworkConnectionParameters](NetworkConnectionParametersParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'toConnectionParameters' field"))
	}
	m.ToConnectionParameters = toConnectionParameters

	transportType, err := ReadSimpleField[TransportType](ctx, "transportType", ReadComplex[TransportType](TransportTypeParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'transportType' field"))
	}
	m.TransportType = transportType

	connectionPathSize, err := ReadSimpleField(ctx, "connectionPathSize", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'connectionPathSize' field"))
	}
	m.ConnectionPathSize = connectionPathSize

	connectionPaths, err := ReadTerminatedArrayField[PathSegment](ctx, "connectionPaths", ReadComplex[PathSegment](PathSegmentParseWithBuffer, readBuffer), NoMorePathSegments(ctx, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'connectionPaths' field"))
	}
	m.ConnectionPaths = connectionPaths

	if closeErr := readBuffer.CloseContext("CipConnectionManagerRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CipConnectionManagerRequest")
	}

	return m, nil
}

func (m *_CipConnectionManagerRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CipConnectionManagerRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CipConnectionManagerRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CipConnectionManagerRequest")
		}
		requestPathSize := uint8(uint8((uint8(m.GetClassSegment().GetLengthInBytes(ctx)) + uint8(m.GetInstanceSegment().GetLengthInBytes(ctx)))) / uint8(uint8(2)))
		if err := WriteImplicitField(ctx, "requestPathSize", requestPathSize, WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestPathSize' field")
		}

		if err := WriteSimpleField[PathSegment](ctx, "classSegment", m.GetClassSegment(), WriteComplex[PathSegment](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'classSegment' field")
		}

		if err := WriteSimpleField[PathSegment](ctx, "instanceSegment", m.GetInstanceSegment(), WriteComplex[PathSegment](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'instanceSegment' field")
		}

		if err := WriteSimpleField[uint8](ctx, "priority", m.GetPriority(), WriteUnsignedByte(writeBuffer, 4)); err != nil {
			return errors.Wrap(err, "Error serializing 'priority' field")
		}

		if err := WriteSimpleField[uint8](ctx, "tickTime", m.GetTickTime(), WriteUnsignedByte(writeBuffer, 4)); err != nil {
			return errors.Wrap(err, "Error serializing 'tickTime' field")
		}

		if err := WriteSimpleField[uint8](ctx, "timeoutTicks", m.GetTimeoutTicks(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'timeoutTicks' field")
		}

		if err := WriteSimpleField[uint32](ctx, "otConnectionId", m.GetOtConnectionId(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'otConnectionId' field")
		}

		if err := WriteSimpleField[uint32](ctx, "toConnectionId", m.GetToConnectionId(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'toConnectionId' field")
		}

		if err := WriteSimpleField[uint16](ctx, "connectionSerialNumber", m.GetConnectionSerialNumber(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'connectionSerialNumber' field")
		}

		if err := WriteSimpleField[uint16](ctx, "originatorVendorId", m.GetOriginatorVendorId(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'originatorVendorId' field")
		}

		if err := WriteSimpleField[uint32](ctx, "originatorSerialNumber", m.GetOriginatorSerialNumber(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'originatorSerialNumber' field")
		}

		if err := WriteSimpleField[uint8](ctx, "timeoutMultiplier", m.GetTimeoutMultiplier(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'timeoutMultiplier' field")
		}

		if err := WriteReservedField[uint32](ctx, "reserved", uint32(0x000000), WriteUnsignedInt(writeBuffer, 24)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteSimpleField[uint32](ctx, "otRpi", m.GetOtRpi(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'otRpi' field")
		}

		if err := WriteSimpleField[NetworkConnectionParameters](ctx, "otConnectionParameters", m.GetOtConnectionParameters(), WriteComplex[NetworkConnectionParameters](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'otConnectionParameters' field")
		}

		if err := WriteSimpleField[uint32](ctx, "toRpi", m.GetToRpi(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'toRpi' field")
		}

		if err := WriteSimpleField[NetworkConnectionParameters](ctx, "toConnectionParameters", m.GetToConnectionParameters(), WriteComplex[NetworkConnectionParameters](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'toConnectionParameters' field")
		}

		if err := WriteSimpleField[TransportType](ctx, "transportType", m.GetTransportType(), WriteComplex[TransportType](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'transportType' field")
		}

		if err := WriteSimpleField[uint8](ctx, "connectionPathSize", m.GetConnectionPathSize(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'connectionPathSize' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "connectionPaths", m.GetConnectionPaths(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'connectionPaths' field")
		}

		if popErr := writeBuffer.PopContext("CipConnectionManagerRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CipConnectionManagerRequest")
		}
		return nil
	}
	return m.CipServiceContract.(*_CipService).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CipConnectionManagerRequest) IsCipConnectionManagerRequest() {}

func (m *_CipConnectionManagerRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CipConnectionManagerRequest) deepCopy() *_CipConnectionManagerRequest {
	if m == nil {
		return nil
	}
	_CipConnectionManagerRequestCopy := &_CipConnectionManagerRequest{
		m.CipServiceContract.(*_CipService).deepCopy(),
		m.ClassSegment.DeepCopy().(PathSegment),
		m.InstanceSegment.DeepCopy().(PathSegment),
		m.Priority,
		m.TickTime,
		m.TimeoutTicks,
		m.OtConnectionId,
		m.ToConnectionId,
		m.ConnectionSerialNumber,
		m.OriginatorVendorId,
		m.OriginatorSerialNumber,
		m.TimeoutMultiplier,
		m.OtRpi,
		m.OtConnectionParameters.DeepCopy().(NetworkConnectionParameters),
		m.ToRpi,
		m.ToConnectionParameters.DeepCopy().(NetworkConnectionParameters),
		m.TransportType.DeepCopy().(TransportType),
		m.ConnectionPathSize,
		utils.DeepCopySlice[PathSegment, PathSegment](m.ConnectionPaths),
		m.reservedField0,
	}
	m.CipServiceContract.(*_CipService)._SubType = m
	return _CipConnectionManagerRequestCopy
}

func (m *_CipConnectionManagerRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
