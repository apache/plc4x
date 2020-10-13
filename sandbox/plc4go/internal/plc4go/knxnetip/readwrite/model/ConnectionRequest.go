//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
)

// The data-structure of this message
type ConnectionRequest struct {
	hpaiDiscoveryEndpoint        IHPAIDiscoveryEndpoint
	hpaiDataEndpoint             IHPAIDataEndpoint
	connectionRequestInformation IConnectionRequestInformation
	KNXNetIPMessage
}

// The corresponding interface
type IConnectionRequest interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ConnectionRequest) MsgType() uint16 {
	return 0x0205
}

func (m ConnectionRequest) initialize() spi.Message {
	return m
}

func NewConnectionRequest(hpaiDiscoveryEndpoint IHPAIDiscoveryEndpoint, hpaiDataEndpoint IHPAIDataEndpoint, connectionRequestInformation IConnectionRequestInformation) KNXNetIPMessageInitializer {
	return &ConnectionRequest{hpaiDiscoveryEndpoint: hpaiDiscoveryEndpoint, hpaiDataEndpoint: hpaiDataEndpoint, connectionRequestInformation: connectionRequestInformation}
}

func CastIConnectionRequest(structType interface{}) IConnectionRequest {
	castFunc := func(typ interface{}) IConnectionRequest {
		if iConnectionRequest, ok := typ.(IConnectionRequest); ok {
			return iConnectionRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionRequest(structType interface{}) ConnectionRequest {
	castFunc := func(typ interface{}) ConnectionRequest {
		if sConnectionRequest, ok := typ.(ConnectionRequest); ok {
			return sConnectionRequest
		}
		return ConnectionRequest{}
	}
	return castFunc(structType)
}

func (m ConnectionRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (hpaiDiscoveryEndpoint)
	lengthInBits += m.hpaiDiscoveryEndpoint.LengthInBits()

	// Simple field (hpaiDataEndpoint)
	lengthInBits += m.hpaiDataEndpoint.LengthInBits()

	// Simple field (connectionRequestInformation)
	lengthInBits += m.connectionRequestInformation.LengthInBits()

	return lengthInBits
}

func (m ConnectionRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionRequestParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (hpaiDiscoveryEndpoint)
	_hpaiDiscoveryEndpointMessage, _err := HPAIDiscoveryEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiDiscoveryEndpoint'. " + _err.Error())
	}
	var hpaiDiscoveryEndpoint IHPAIDiscoveryEndpoint
	hpaiDiscoveryEndpoint, _hpaiDiscoveryEndpointOk := _hpaiDiscoveryEndpointMessage.(IHPAIDiscoveryEndpoint)
	if !_hpaiDiscoveryEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiDiscoveryEndpointMessage).Name() + " to IHPAIDiscoveryEndpoint")
	}

	// Simple Field (hpaiDataEndpoint)
	_hpaiDataEndpointMessage, _err := HPAIDataEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiDataEndpoint'. " + _err.Error())
	}
	var hpaiDataEndpoint IHPAIDataEndpoint
	hpaiDataEndpoint, _hpaiDataEndpointOk := _hpaiDataEndpointMessage.(IHPAIDataEndpoint)
	if !_hpaiDataEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiDataEndpointMessage).Name() + " to IHPAIDataEndpoint")
	}

	// Simple Field (connectionRequestInformation)
	_connectionRequestInformationMessage, _err := ConnectionRequestInformationParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'connectionRequestInformation'. " + _err.Error())
	}
	var connectionRequestInformation IConnectionRequestInformation
	connectionRequestInformation, _connectionRequestInformationOk := _connectionRequestInformationMessage.(IConnectionRequestInformation)
	if !_connectionRequestInformationOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_connectionRequestInformationMessage).Name() + " to IConnectionRequestInformation")
	}

	// Create the instance
	return NewConnectionRequest(hpaiDiscoveryEndpoint, hpaiDataEndpoint, connectionRequestInformation), nil
}

func (m ConnectionRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (hpaiDiscoveryEndpoint)
		hpaiDiscoveryEndpoint := CastIHPAIDiscoveryEndpoint(m.hpaiDiscoveryEndpoint)
		_hpaiDiscoveryEndpointErr := hpaiDiscoveryEndpoint.Serialize(io)
		if _hpaiDiscoveryEndpointErr != nil {
			return errors.New("Error serializing 'hpaiDiscoveryEndpoint' field " + _hpaiDiscoveryEndpointErr.Error())
		}

		// Simple Field (hpaiDataEndpoint)
		hpaiDataEndpoint := CastIHPAIDataEndpoint(m.hpaiDataEndpoint)
		_hpaiDataEndpointErr := hpaiDataEndpoint.Serialize(io)
		if _hpaiDataEndpointErr != nil {
			return errors.New("Error serializing 'hpaiDataEndpoint' field " + _hpaiDataEndpointErr.Error())
		}

		// Simple Field (connectionRequestInformation)
		connectionRequestInformation := CastIConnectionRequestInformation(m.connectionRequestInformation)
		_connectionRequestInformationErr := connectionRequestInformation.Serialize(io)
		if _connectionRequestInformationErr != nil {
			return errors.New("Error serializing 'connectionRequestInformation' field " + _connectionRequestInformationErr.Error())
		}

		return nil
	}
	return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
