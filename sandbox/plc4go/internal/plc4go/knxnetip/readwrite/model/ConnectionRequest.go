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
	hpaiDiscoveryEndpoint        HPAIDiscoveryEndpoint
	hpaiDataEndpoint             HPAIDataEndpoint
	connectionRequestInformation ConnectionRequestInformation
	KNXNetIPMessage
}

// The corresponding interface
type IConnectionRequest interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ConnectionRequest) MsgType() uint16 {
	return 0x0205
}

func (m ConnectionRequest) initialize() spi.Message {
	return m
}

func NewConnectionRequest(hpaiDiscoveryEndpoint HPAIDiscoveryEndpoint, hpaiDataEndpoint HPAIDataEndpoint, connectionRequestInformation ConnectionRequestInformation) KNXNetIPMessageInitializer {
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

func ConnectionRequestParse(io spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (hpaiDiscoveryEndpoint)
	_hpaiDiscoveryEndpointMessage, _err := HPAIDiscoveryEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiDiscoveryEndpoint'. " + _err.Error())
	}
	var hpaiDiscoveryEndpoint HPAIDiscoveryEndpoint
	hpaiDiscoveryEndpoint, _hpaiDiscoveryEndpointOk := _hpaiDiscoveryEndpointMessage.(HPAIDiscoveryEndpoint)
	if !_hpaiDiscoveryEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiDiscoveryEndpointMessage).Name() + " to HPAIDiscoveryEndpoint")
	}

	// Simple Field (hpaiDataEndpoint)
	_hpaiDataEndpointMessage, _err := HPAIDataEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiDataEndpoint'. " + _err.Error())
	}
	var hpaiDataEndpoint HPAIDataEndpoint
	hpaiDataEndpoint, _hpaiDataEndpointOk := _hpaiDataEndpointMessage.(HPAIDataEndpoint)
	if !_hpaiDataEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiDataEndpointMessage).Name() + " to HPAIDataEndpoint")
	}

	// Simple Field (connectionRequestInformation)
	_connectionRequestInformationMessage, _err := ConnectionRequestInformationParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'connectionRequestInformation'. " + _err.Error())
	}
	var connectionRequestInformation ConnectionRequestInformation
	connectionRequestInformation, _connectionRequestInformationOk := _connectionRequestInformationMessage.(ConnectionRequestInformation)
	if !_connectionRequestInformationOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_connectionRequestInformationMessage).Name() + " to ConnectionRequestInformation")
	}

	// Create the instance
	return NewConnectionRequest(hpaiDiscoveryEndpoint, hpaiDataEndpoint, connectionRequestInformation), nil
}

func (m ConnectionRequest) Serialize(io spi.WriteBuffer) {

	// Simple Field (hpaiDiscoveryEndpoint)
	hpaiDiscoveryEndpoint := HPAIDiscoveryEndpoint(m.hpaiDiscoveryEndpoint)
	hpaiDiscoveryEndpoint.Serialize(io)

	// Simple Field (hpaiDataEndpoint)
	hpaiDataEndpoint := HPAIDataEndpoint(m.hpaiDataEndpoint)
	hpaiDataEndpoint.Serialize(io)

	// Simple Field (connectionRequestInformation)
	connectionRequestInformation := ConnectionRequestInformation(m.connectionRequestInformation)
	connectionRequestInformation.Serialize(io)
}
