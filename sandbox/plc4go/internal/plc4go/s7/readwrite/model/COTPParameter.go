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
type COTPParameter struct {
}

// The corresponding interface
type ICOTPParameter interface {
	spi.Message
	ParameterType() uint8
	Serialize(io spi.WriteBuffer)
}

type COTPParameterInitializer interface {
	initialize() spi.Message
}

func COTPParameterParameterType(m ICOTPParameter) uint8 {
	return m.ParameterType()
}

func CastICOTPParameter(structType interface{}) ICOTPParameter {
	castFunc := func(typ interface{}) ICOTPParameter {
		if iCOTPParameter, ok := typ.(ICOTPParameter); ok {
			return iCOTPParameter
		}
		return nil
	}
	return castFunc(structType)
}

func CastCOTPParameter(structType interface{}) COTPParameter {
	castFunc := func(typ interface{}) COTPParameter {
		if sCOTPParameter, ok := typ.(COTPParameter); ok {
			return sCOTPParameter
		}
		return COTPParameter{}
	}
	return castFunc(structType)
}

func (m COTPParameter) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (parameterType)
	lengthInBits += 8

	// Implicit Field (parameterLength)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m COTPParameter) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPParameterParse(io spi.ReadBuffer, rest uint8) (spi.Message, error) {

	// Discriminator Field (parameterType) (Used as input to a switch field)
	parameterType, _parameterTypeErr := io.ReadUint8(8)
	if _parameterTypeErr != nil {
		return nil, errors.New("Error parsing 'parameterType' field " + _parameterTypeErr.Error())
	}

	// Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _parameterLengthErr := io.ReadUint8(8)
	if _parameterLengthErr != nil {
		return nil, errors.New("Error parsing 'parameterLength' field " + _parameterLengthErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer COTPParameterInitializer
	var typeSwitchError error
	switch {
	case parameterType == 0xC0:
		initializer, typeSwitchError = COTPParameterTpduSizeParse(io)
	case parameterType == 0xC1:
		initializer, typeSwitchError = COTPParameterCallingTsapParse(io)
	case parameterType == 0xC2:
		initializer, typeSwitchError = COTPParameterCalledTsapParse(io)
	case parameterType == 0xC3:
		initializer, typeSwitchError = COTPParameterChecksumParse(io)
	case parameterType == 0xE0:
		initializer, typeSwitchError = COTPParameterDisconnectAdditionalInformationParse(io, rest)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m COTPParameter) Serialize(io spi.WriteBuffer) {
	iCOTPParameter := CastICOTPParameter(m)

	// Discriminator Field (parameterType) (Used as input to a switch field)
	parameterType := uint8(COTPParameterParameterType(iCOTPParameter))
	io.WriteUint8(8, (parameterType))

	// Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	parameterLength := uint8(uint8(uint8(m.LengthInBytes())) - uint8(uint8(2)))
	io.WriteUint8(8, (parameterLength))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iCOTPParameter.Serialize(io)
}
