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
type COTPParameterDisconnectAdditionalInformation struct {
	data []uint8
	COTPParameter
}

// The corresponding interface
type ICOTPParameterDisconnectAdditionalInformation interface {
	ICOTPParameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m COTPParameterDisconnectAdditionalInformation) ParameterType() uint8 {
	return 0xE0
}

func (m COTPParameterDisconnectAdditionalInformation) initialize() spi.Message {
	return m
}

func NewCOTPParameterDisconnectAdditionalInformation(data []uint8) COTPParameterInitializer {
	return &COTPParameterDisconnectAdditionalInformation{data: data}
}

func (m COTPParameterDisconnectAdditionalInformation) LengthInBits() uint16 {
	var lengthInBits uint16 = m.COTPParameter.LengthInBits()

	// Array field
	if len(m.data) > 0 {
		lengthInBits += 8 * uint16(len(m.data))
	}

	return lengthInBits
}

func (m COTPParameterDisconnectAdditionalInformation) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPParameterDisconnectAdditionalInformationParse(io spi.ReadBuffer, rest uint8) (COTPParameterInitializer, error) {

	// Array field (data)
	var data []uint8
	// Count array
	{
		data := make([]uint8, rest)
		for curItem := uint16(0); curItem < uint16(rest); curItem++ {

			data = append(data, io.ReadUint8(8))
		}
	}

	// Create the instance
	return NewCOTPParameterDisconnectAdditionalInformation(data), nil
}

func (m COTPParameterDisconnectAdditionalInformation) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ICOTPParameterDisconnectAdditionalInformation); ok {

			// Array Field (data)
			if m.data != nil {
				for _, _element := range m.data {
					io.WriteUint8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
