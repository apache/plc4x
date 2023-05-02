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
)

// Code generated by code-generation. DO NOT EDIT.

// ServiceId is the corresponding interface of ServiceId
type ServiceId interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetServiceType returns ServiceType (discriminator field)
	GetServiceType() uint8
}

// ServiceIdExactly can be used when we want exactly this type and not a type which fulfills ServiceId.
// This is useful for switch cases.
type ServiceIdExactly interface {
	ServiceId
	isServiceId() bool
}

// _ServiceId is the data-structure of this message
type _ServiceId struct {
	_ServiceIdChildRequirements
}

type _ServiceIdChildRequirements interface {
	utils.Serializable
	GetLengthInBits(ctx context.Context) uint16
	GetServiceType() uint8
}

type ServiceIdParent interface {
	SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child ServiceId, serializeChildFunction func() error) error
	GetTypeName() string
}

type ServiceIdChild interface {
	utils.Serializable
	InitializeParent(parent ServiceId)
	GetParent() *ServiceId

	GetTypeName() string
	ServiceId
}

// NewServiceId factory function for _ServiceId
func NewServiceId() *_ServiceId {
	return &_ServiceId{}
}

// Deprecated: use the interface for direct cast
func CastServiceId(structType any) ServiceId {
	if casted, ok := structType.(ServiceId); ok {
		return casted
	}
	if casted, ok := structType.(*ServiceId); ok {
		return *casted
	}
	return nil
}

func (m *_ServiceId) GetTypeName() string {
	return "ServiceId"
}

func (m *_ServiceId) GetParentLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)
	// Discriminator Field (serviceType)
	lengthInBits += 8

	return lengthInBits
}

func (m *_ServiceId) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func ServiceIdParse(theBytes []byte) (ServiceId, error) {
	return ServiceIdParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func ServiceIdParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (ServiceId, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ServiceId"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ServiceId")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Discriminator Field (serviceType) (Used as input to a switch field)
	serviceType, _serviceTypeErr := readBuffer.ReadUint8("serviceType", 8)
	if _serviceTypeErr != nil {
		return nil, errors.Wrap(_serviceTypeErr, "Error parsing 'serviceType' field of ServiceId")
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type ServiceIdChildSerializeRequirement interface {
		ServiceId
		InitializeParent(ServiceId)
		GetParent() ServiceId
	}
	var _childTemp any
	var _child ServiceIdChildSerializeRequirement
	var typeSwitchError error
	switch {
	case serviceType == 0x02: // KnxNetIpCore
		_childTemp, typeSwitchError = KnxNetIpCoreParseWithBuffer(ctx, readBuffer)
	case serviceType == 0x03: // KnxNetIpDeviceManagement
		_childTemp, typeSwitchError = KnxNetIpDeviceManagementParseWithBuffer(ctx, readBuffer)
	case serviceType == 0x04: // KnxNetIpTunneling
		_childTemp, typeSwitchError = KnxNetIpTunnelingParseWithBuffer(ctx, readBuffer)
	case serviceType == 0x05: // KnxNetIpRouting
		_childTemp, typeSwitchError = KnxNetIpRoutingParseWithBuffer(ctx, readBuffer)
	case serviceType == 0x06: // KnxNetRemoteLogging
		_childTemp, typeSwitchError = KnxNetRemoteLoggingParseWithBuffer(ctx, readBuffer)
	case serviceType == 0x07: // KnxNetRemoteConfigurationAndDiagnosis
		_childTemp, typeSwitchError = KnxNetRemoteConfigurationAndDiagnosisParseWithBuffer(ctx, readBuffer)
	case serviceType == 0x08: // KnxNetObjectServer
		_childTemp, typeSwitchError = KnxNetObjectServerParseWithBuffer(ctx, readBuffer)
	default:
		typeSwitchError = errors.Errorf("Unmapped type for parameters [serviceType=%v]", serviceType)
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch of ServiceId")
	}
	_child = _childTemp.(ServiceIdChildSerializeRequirement)

	if closeErr := readBuffer.CloseContext("ServiceId"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ServiceId")
	}

	// Finish initializing
	_child.InitializeParent(_child)
	return _child, nil
}

func (pm *_ServiceId) SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child ServiceId, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("ServiceId"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for ServiceId")
	}

	// Discriminator Field (serviceType) (Used as input to a switch field)
	serviceType := uint8(child.GetServiceType())
	_serviceTypeErr := writeBuffer.WriteUint8("serviceType", 8, (serviceType))

	if _serviceTypeErr != nil {
		return errors.Wrap(_serviceTypeErr, "Error serializing 'serviceType' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("ServiceId"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for ServiceId")
	}
	return nil
}

func (m *_ServiceId) isServiceId() bool {
	return true
}

func (m *_ServiceId) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
