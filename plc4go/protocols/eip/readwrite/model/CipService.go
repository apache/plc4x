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

// CipService is the corresponding interface of CipService
type CipService interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetConnected returns Connected (discriminator field)
	GetConnected() bool
	// GetResponse returns Response (discriminator field)
	GetResponse() bool
	// GetService returns Service (discriminator field)
	GetService() uint8
}

// CipServiceExactly can be used when we want exactly this type and not a type which fulfills CipService.
// This is useful for switch cases.
type CipServiceExactly interface {
	CipService
	isCipService() bool
}

// _CipService is the data-structure of this message
type _CipService struct {
	_CipServiceChildRequirements

	// Arguments.
	ServiceLen uint16
}

type _CipServiceChildRequirements interface {
	utils.Serializable
	GetLengthInBits(ctx context.Context) uint16
	GetService() uint8
	GetResponse() bool
	GetConnected() bool
}

type CipServiceParent interface {
	SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child CipService, serializeChildFunction func() error) error
	GetTypeName() string
}

type CipServiceChild interface {
	utils.Serializable
	InitializeParent(parent CipService)
	GetParent() *CipService

	GetTypeName() string
	CipService
}

// NewCipService factory function for _CipService
func NewCipService(serviceLen uint16) *_CipService {
	return &_CipService{ServiceLen: serviceLen}
}

// Deprecated: use the interface for direct cast
func CastCipService(structType any) CipService {
	if casted, ok := structType.(CipService); ok {
		return casted
	}
	if casted, ok := structType.(*CipService); ok {
		return *casted
	}
	return nil
}

func (m *_CipService) GetTypeName() string {
	return "CipService"
}

func (m *_CipService) GetParentLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)
	// Discriminator Field (response)
	lengthInBits += 1
	// Discriminator Field (service)
	lengthInBits += 7

	return lengthInBits
}

func (m *_CipService) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CipServiceParse(theBytes []byte, connected bool, serviceLen uint16) (CipService, error) {
	return CipServiceParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes), connected, serviceLen)
}

func CipServiceParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, connected bool, serviceLen uint16) (CipService, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CipService"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CipService")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Discriminator Field (response) (Used as input to a switch field)
	response, _responseErr := readBuffer.ReadBit("response")
	if _responseErr != nil {
		return nil, errors.Wrap(_responseErr, "Error parsing 'response' field of CipService")
	}

	// Discriminator Field (service) (Used as input to a switch field)
	service, _serviceErr := readBuffer.ReadUint8("service", 7)
	if _serviceErr != nil {
		return nil, errors.Wrap(_serviceErr, "Error parsing 'service' field of CipService")
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type CipServiceChildSerializeRequirement interface {
		CipService
		InitializeParent(CipService)
		GetParent() CipService
	}
	var _childTemp any
	var _child CipServiceChildSerializeRequirement
	var typeSwitchError error
	switch {
	case service == 0x01 && response == bool(false): // GetAttributeAllRequest
		_childTemp, typeSwitchError = GetAttributeAllRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x01 && response == bool(true): // GetAttributeAllResponse
		_childTemp, typeSwitchError = GetAttributeAllResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x02 && response == bool(false): // SetAttributeAllRequest
		_childTemp, typeSwitchError = SetAttributeAllRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x02 && response == bool(true): // SetAttributeAllResponse
		_childTemp, typeSwitchError = SetAttributeAllResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x03 && response == bool(false): // GetAttributeListRequest
		_childTemp, typeSwitchError = GetAttributeListRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x03 && response == bool(true): // GetAttributeListResponse
		_childTemp, typeSwitchError = GetAttributeListResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x04 && response == bool(false): // SetAttributeListRequest
		_childTemp, typeSwitchError = SetAttributeListRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x04 && response == bool(true): // SetAttributeListResponse
		_childTemp, typeSwitchError = SetAttributeListResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x0A && response == bool(false): // MultipleServiceRequest
		_childTemp, typeSwitchError = MultipleServiceRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x0A && response == bool(true): // MultipleServiceResponse
		_childTemp, typeSwitchError = MultipleServiceResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x0E && response == bool(false): // GetAttributeSingleRequest
		_childTemp, typeSwitchError = GetAttributeSingleRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x0E && response == bool(true): // GetAttributeSingleResponse
		_childTemp, typeSwitchError = GetAttributeSingleResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x10 && response == bool(false): // SetAttributeSingleRequest
		_childTemp, typeSwitchError = SetAttributeSingleRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x10 && response == bool(true): // SetAttributeSingleResponse
		_childTemp, typeSwitchError = SetAttributeSingleResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x4C && response == bool(false): // CipReadRequest
		_childTemp, typeSwitchError = CipReadRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x4C && response == bool(true): // CipReadResponse
		_childTemp, typeSwitchError = CipReadResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x4D && response == bool(false): // CipWriteRequest
		_childTemp, typeSwitchError = CipWriteRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x4D && response == bool(true): // CipWriteResponse
		_childTemp, typeSwitchError = CipWriteResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x4E && response == bool(false): // CipConnectionManagerCloseRequest
		_childTemp, typeSwitchError = CipConnectionManagerCloseRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x4E && response == bool(true): // CipConnectionManagerCloseResponse
		_childTemp, typeSwitchError = CipConnectionManagerCloseResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x52 && response == bool(false) && connected == bool(false): // CipUnconnectedRequest
		_childTemp, typeSwitchError = CipUnconnectedRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x52 && response == bool(false) && connected == bool(true): // CipConnectedRequest
		_childTemp, typeSwitchError = CipConnectedRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x52 && response == bool(true): // CipConnectedResponse
		_childTemp, typeSwitchError = CipConnectedResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x5B && response == bool(false): // CipConnectionManagerRequest
		_childTemp, typeSwitchError = CipConnectionManagerRequestParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	case service == 0x5B && response == bool(true): // CipConnectionManagerResponse
		_childTemp, typeSwitchError = CipConnectionManagerResponseParseWithBuffer(ctx, readBuffer, connected, serviceLen)
	default:
		typeSwitchError = errors.Errorf("Unmapped type for parameters [service=%v, response=%v, connected=%v]", service, response, connected)
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch of CipService")
	}
	_child = _childTemp.(CipServiceChildSerializeRequirement)

	if closeErr := readBuffer.CloseContext("CipService"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CipService")
	}

	// Finish initializing
	_child.InitializeParent(_child)
	return _child, nil
}

func (pm *_CipService) SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child CipService, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("CipService"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for CipService")
	}

	// Discriminator Field (response) (Used as input to a switch field)
	response := bool(child.GetResponse())
	_responseErr := writeBuffer.WriteBit("response", (response))

	if _responseErr != nil {
		return errors.Wrap(_responseErr, "Error serializing 'response' field")
	}

	// Discriminator Field (service) (Used as input to a switch field)
	service := uint8(child.GetService())
	_serviceErr := writeBuffer.WriteUint8("service", 7, (service))

	if _serviceErr != nil {
		return errors.Wrap(_serviceErr, "Error serializing 'service' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("CipService"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for CipService")
	}
	return nil
}

////
// Arguments Getter

func (m *_CipService) GetServiceLen() uint16 {
	return m.ServiceLen
}

//
////

func (m *_CipService) isCipService() bool {
	return true
}

func (m *_CipService) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
