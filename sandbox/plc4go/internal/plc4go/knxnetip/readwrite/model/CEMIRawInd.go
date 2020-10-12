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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type CEMIRawInd struct {
	CEMI
}

// The corresponding interface
type ICEMIRawInd interface {
	ICEMI
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m CEMIRawInd) MessageCode() uint8 {
	return 0x2D
}

func (m CEMIRawInd) initialize() spi.Message {
	return m
}

func NewCEMIRawInd() CEMIInitializer {
	return &CEMIRawInd{}
}

func CastICEMIRawInd(structType interface{}) ICEMIRawInd {
	castFunc := func(typ interface{}) ICEMIRawInd {
		if iCEMIRawInd, ok := typ.(ICEMIRawInd); ok {
			return iCEMIRawInd
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIRawInd(structType interface{}) CEMIRawInd {
	castFunc := func(typ interface{}) CEMIRawInd {
		if sCEMIRawInd, ok := typ.(CEMIRawInd); ok {
			return sCEMIRawInd
		}
		return CEMIRawInd{}
	}
	return castFunc(structType)
}

func (m CEMIRawInd) LengthInBits() uint16 {
	var lengthInBits uint16 = m.CEMI.LengthInBits()

	return lengthInBits
}

func (m CEMIRawInd) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIRawIndParse(io spi.ReadBuffer) (CEMIInitializer, error) {

	// Create the instance
	return NewCEMIRawInd(), nil
}

func (m CEMIRawInd) Serialize(io spi.WriteBuffer) {

}
