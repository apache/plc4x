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
type COTPPacketTpduError struct {
	destinationReference uint16
	rejectCause          uint8
	COTPPacket
}

// The corresponding interface
type ICOTPPacketTpduError interface {
	ICOTPPacket
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m COTPPacketTpduError) TpduCode() uint8 {
	return 0x70
}

func (m COTPPacketTpduError) initialize(parameters []ICOTPParameter, payload *IS7Message) spi.Message {
	m.parameters = parameters
	m.payload = payload
	return m
}

func NewCOTPPacketTpduError(destinationReference uint16, rejectCause uint8) COTPPacketInitializer {
	return &COTPPacketTpduError{destinationReference: destinationReference, rejectCause: rejectCause}
}

func CastICOTPPacketTpduError(structType interface{}) ICOTPPacketTpduError {
	castFunc := func(typ interface{}) ICOTPPacketTpduError {
		if iCOTPPacketTpduError, ok := typ.(ICOTPPacketTpduError); ok {
			return iCOTPPacketTpduError
		}
		return nil
	}
	return castFunc(structType)
}

func CastCOTPPacketTpduError(structType interface{}) COTPPacketTpduError {
	castFunc := func(typ interface{}) COTPPacketTpduError {
		if sCOTPPacketTpduError, ok := typ.(COTPPacketTpduError); ok {
			return sCOTPPacketTpduError
		}
		return COTPPacketTpduError{}
	}
	return castFunc(structType)
}

func (m COTPPacketTpduError) LengthInBits() uint16 {
	var lengthInBits uint16 = m.COTPPacket.LengthInBits()

	// Simple field (destinationReference)
	lengthInBits += 16

	// Simple field (rejectCause)
	lengthInBits += 8

	return lengthInBits
}

func (m COTPPacketTpduError) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPPacketTpduErrorParse(io spi.ReadBuffer) (COTPPacketInitializer, error) {

	// Simple Field (destinationReference)
	destinationReference, _destinationReferenceErr := io.ReadUint16(16)
	if _destinationReferenceErr != nil {
		return nil, errors.New("Error parsing 'destinationReference' field " + _destinationReferenceErr.Error())
	}

	// Simple Field (rejectCause)
	rejectCause, _rejectCauseErr := io.ReadUint8(8)
	if _rejectCauseErr != nil {
		return nil, errors.New("Error parsing 'rejectCause' field " + _rejectCauseErr.Error())
	}

	// Create the instance
	return NewCOTPPacketTpduError(destinationReference, rejectCause), nil
}

func (m COTPPacketTpduError) Serialize(io spi.WriteBuffer) {

	// Simple Field (destinationReference)
	destinationReference := uint16(m.destinationReference)
	io.WriteUint16(16, (destinationReference))

	// Simple Field (rejectCause)
	rejectCause := uint8(m.rejectCause)
	io.WriteUint8(8, (rejectCause))
}
