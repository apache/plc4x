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
type CEMI struct {
}

// The corresponding interface
type ICEMI interface {
	spi.Message
	MessageCode() uint8
	Serialize(io spi.WriteBuffer)
}

type CEMIInitializer interface {
	initialize() spi.Message
}

func CEMIMessageCode(m ICEMI) uint8 {
	return m.MessageCode()
}

func (m CEMI) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (messageCode)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m CEMI) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIParse(io spi.ReadBuffer, size uint8) (spi.Message, error) {

	// Discriminator Field (messageCode) (Used as input to a switch field)
	var messageCode uint8 = io.ReadUint8(8)

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer CEMIInitializer
	var typeSwitchError error
	switch {
	case messageCode == 0x11:
		initializer, typeSwitchError = CEMIDataReqParse(io)
	case messageCode == 0x2E:
		initializer, typeSwitchError = CEMIDataConParse(io)
	case messageCode == 0x29:
		initializer, typeSwitchError = CEMIDataIndParse(io)
	case messageCode == 0x10:
		initializer, typeSwitchError = CEMIRawReqParse(io)
	case messageCode == 0x2F:
		initializer, typeSwitchError = CEMIRawConParse(io)
	case messageCode == 0x2D:
		initializer, typeSwitchError = CEMIRawIndParse(io)
	case messageCode == 0x13:
		initializer, typeSwitchError = CEMIPollDataReqParse(io)
	case messageCode == 0x25:
		initializer, typeSwitchError = CEMIPollDataConParse(io)
	case messageCode == 0x2B:
		initializer, typeSwitchError = CEMIBusmonIndParse(io)
	case messageCode == 0xFC:
		initializer, typeSwitchError = CEMIMPropReadReqParse(io)
	case messageCode == 0xFB:
		initializer, typeSwitchError = CEMIMPropReadConParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m CEMI) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if iCEMI, ok := typ.(ICEMI); ok {

			// Discriminator Field (messageCode) (Used as input to a switch field)
			messageCode := CEMIMessageCode(iCEMI)
			io.WriteUint8(8, (messageCode))

			// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
			iCEMI.Serialize(io)
		}
	}
	serializeFunc(m)
}
