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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type COTPPacketDisconnectResponse struct {
    DestinationReference uint16
    SourceReference uint16
    COTPPacket
}

// The corresponding interface
type ICOTPPacketDisconnectResponse interface {
    ICOTPPacket
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m COTPPacketDisconnectResponse) TpduCode() uint8 {
    return 0xC0
}

func (m COTPPacketDisconnectResponse) initialize(parameters []ICOTPParameter, payload *IS7Message) spi.Message {
    m.Parameters = parameters
    m.Payload = payload
    return m
}

func NewCOTPPacketDisconnectResponse(destinationReference uint16, sourceReference uint16) COTPPacketInitializer {
    return &COTPPacketDisconnectResponse{DestinationReference: destinationReference, SourceReference: sourceReference}
}

func CastICOTPPacketDisconnectResponse(structType interface{}) ICOTPPacketDisconnectResponse {
    castFunc := func(typ interface{}) ICOTPPacketDisconnectResponse {
        if iCOTPPacketDisconnectResponse, ok := typ.(ICOTPPacketDisconnectResponse); ok {
            return iCOTPPacketDisconnectResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastCOTPPacketDisconnectResponse(structType interface{}) COTPPacketDisconnectResponse {
    castFunc := func(typ interface{}) COTPPacketDisconnectResponse {
        if sCOTPPacketDisconnectResponse, ok := typ.(COTPPacketDisconnectResponse); ok {
            return sCOTPPacketDisconnectResponse
        }
        if sCOTPPacketDisconnectResponse, ok := typ.(*COTPPacketDisconnectResponse); ok {
            return *sCOTPPacketDisconnectResponse
        }
        return COTPPacketDisconnectResponse{}
    }
    return castFunc(structType)
}

func (m COTPPacketDisconnectResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.COTPPacket.LengthInBits()

    // Simple field (destinationReference)
    lengthInBits += 16

    // Simple field (sourceReference)
    lengthInBits += 16

    return lengthInBits
}

func (m COTPPacketDisconnectResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPPacketDisconnectResponseParse(io *utils.ReadBuffer) (COTPPacketInitializer, error) {

    // Simple Field (destinationReference)
    destinationReference, _destinationReferenceErr := io.ReadUint16(16)
    if _destinationReferenceErr != nil {
        return nil, errors.New("Error parsing 'destinationReference' field " + _destinationReferenceErr.Error())
    }

    // Simple Field (sourceReference)
    sourceReference, _sourceReferenceErr := io.ReadUint16(16)
    if _sourceReferenceErr != nil {
        return nil, errors.New("Error parsing 'sourceReference' field " + _sourceReferenceErr.Error())
    }

    // Create the instance
    return NewCOTPPacketDisconnectResponse(destinationReference, sourceReference), nil
}

func (m COTPPacketDisconnectResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (destinationReference)
    destinationReference := uint16(m.DestinationReference)
    _destinationReferenceErr := io.WriteUint16(16, (destinationReference))
    if _destinationReferenceErr != nil {
        return errors.New("Error serializing 'destinationReference' field " + _destinationReferenceErr.Error())
    }

    // Simple Field (sourceReference)
    sourceReference := uint16(m.SourceReference)
    _sourceReferenceErr := io.WriteUint16(16, (sourceReference))
    if _sourceReferenceErr != nil {
        return errors.New("Error serializing 'sourceReference' field " + _sourceReferenceErr.Error())
    }

        return nil
    }
    return COTPPacketSerialize(io, m.COTPPacket, CastICOTPPacket(m), ser)
}
