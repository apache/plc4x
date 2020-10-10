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
type COTPParameterCallingTsap struct {
	tsapId uint16
	COTPParameter
}

// The corresponding interface
type ICOTPParameterCallingTsap interface {
	ICOTPParameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m COTPParameterCallingTsap) ParameterType() uint8 {
	return 0xC1
}

func (m COTPParameterCallingTsap) initialize() spi.Message {
	return m
}

func NewCOTPParameterCallingTsap(tsapId uint16) COTPParameterInitializer {
	return &COTPParameterCallingTsap{tsapId: tsapId}
}

func (m COTPParameterCallingTsap) LengthInBits() uint16 {
	var lengthInBits uint16 = m.COTPParameter.LengthInBits()

	// Simple field (tsapId)
	lengthInBits += 16

	return lengthInBits
}

func (m COTPParameterCallingTsap) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPParameterCallingTsapParse(io spi.ReadBuffer) (COTPParameterInitializer, error) {

	// Simple Field (tsapId)
	var tsapId uint16 = io.ReadUint16(16)

	// Create the instance
	return NewCOTPParameterCallingTsap(tsapId), nil
}

func (m COTPParameterCallingTsap) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ICOTPParameterCallingTsap); ok {

			// Simple Field (tsapId)
			var tsapId uint16 = m.tsapId
			io.WriteUint16(16, (tsapId))
		}
	}
	serializeFunc(m)
}
