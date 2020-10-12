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
)

// The data-structure of this message
type COTPPacketConnectionRequest struct {
	destinationReference uint16
	sourceReference      uint16
	protocolClass        COTPProtocolClass
	COTPPacket
}

// The corresponding interface
type ICOTPPacketConnectionRequest interface {
	ICOTPPacket
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m COTPPacketConnectionRequest) TpduCode() uint8 {
	return 0xE0
}

func (m COTPPacketConnectionRequest) initialize(parameters []COTPParameter, payload *S7Message) spi.Message {
	m.parameters = parameters
	m.payload = payload
	return m
}

func NewCOTPPacketConnectionRequest(destinationReference uint16, sourceReference uint16, protocolClass COTPProtocolClass) COTPPacketInitializer {
	return &COTPPacketConnectionRequest{destinationReference: destinationReference, sourceReference: sourceReference, protocolClass: protocolClass}
}

func CastICOTPPacketConnectionRequest(structType interface{}) ICOTPPacketConnectionRequest {
	castFunc := func(typ interface{}) ICOTPPacketConnectionRequest {
		if iCOTPPacketConnectionRequest, ok := typ.(ICOTPPacketConnectionRequest); ok {
			return iCOTPPacketConnectionRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastCOTPPacketConnectionRequest(structType interface{}) COTPPacketConnectionRequest {
	castFunc := func(typ interface{}) COTPPacketConnectionRequest {
		if sCOTPPacketConnectionRequest, ok := typ.(COTPPacketConnectionRequest); ok {
			return sCOTPPacketConnectionRequest
		}
		return COTPPacketConnectionRequest{}
	}
	return castFunc(structType)
}

func (m COTPPacketConnectionRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.COTPPacket.LengthInBits()

	// Simple field (destinationReference)
	lengthInBits += 16

	// Simple field (sourceReference)
	lengthInBits += 16

	// Enum Field (protocolClass)
	lengthInBits += 8

	return lengthInBits
}

func (m COTPPacketConnectionRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPPacketConnectionRequestParse(io spi.ReadBuffer) (COTPPacketInitializer, error) {

	// Simple Field (destinationReference)
	var destinationReference uint16 = io.ReadUint16(16)

	// Simple Field (sourceReference)
	var sourceReference uint16 = io.ReadUint16(16)

	// Enum field (protocolClass)
	protocolClass, _protocolClassErr := COTPProtocolClassParse(io)
	if _protocolClassErr != nil {
		return nil, errors.New("Error parsing 'protocolClass' field " + _protocolClassErr.Error())
	}

	// Create the instance
	return NewCOTPPacketConnectionRequest(destinationReference, sourceReference, protocolClass), nil
}

func (m COTPPacketConnectionRequest) Serialize(io spi.WriteBuffer) {

	// Simple Field (destinationReference)
	destinationReference := uint16(m.destinationReference)
	io.WriteUint16(16, (destinationReference))

	// Simple Field (sourceReference)
	sourceReference := uint16(m.sourceReference)
	io.WriteUint16(16, (sourceReference))

	// Enum field (protocolClass)
	protocolClass := COTPProtocolClass(m.protocolClass)
	protocolClass.Serialize(io)
}
