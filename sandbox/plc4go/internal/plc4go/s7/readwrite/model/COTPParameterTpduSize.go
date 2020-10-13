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
type COTPParameterTpduSize struct {
	tpduSize ICOTPTpduSize
	COTPParameter
}

// The corresponding interface
type ICOTPParameterTpduSize interface {
	ICOTPParameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m COTPParameterTpduSize) ParameterType() uint8 {
	return 0xC0
}

func (m COTPParameterTpduSize) initialize() spi.Message {
	return m
}

func NewCOTPParameterTpduSize(tpduSize ICOTPTpduSize) COTPParameterInitializer {
	return &COTPParameterTpduSize{tpduSize: tpduSize}
}

func CastICOTPParameterTpduSize(structType interface{}) ICOTPParameterTpduSize {
	castFunc := func(typ interface{}) ICOTPParameterTpduSize {
		if iCOTPParameterTpduSize, ok := typ.(ICOTPParameterTpduSize); ok {
			return iCOTPParameterTpduSize
		}
		return nil
	}
	return castFunc(structType)
}

func CastCOTPParameterTpduSize(structType interface{}) COTPParameterTpduSize {
	castFunc := func(typ interface{}) COTPParameterTpduSize {
		if sCOTPParameterTpduSize, ok := typ.(COTPParameterTpduSize); ok {
			return sCOTPParameterTpduSize
		}
		return COTPParameterTpduSize{}
	}
	return castFunc(structType)
}

func (m COTPParameterTpduSize) LengthInBits() uint16 {
	var lengthInBits uint16 = m.COTPParameter.LengthInBits()

	// Enum Field (tpduSize)
	lengthInBits += 8

	return lengthInBits
}

func (m COTPParameterTpduSize) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPParameterTpduSizeParse(io *spi.ReadBuffer) (COTPParameterInitializer, error) {

	// Enum field (tpduSize)
	tpduSize, _tpduSizeErr := COTPTpduSizeParse(io)
	if _tpduSizeErr != nil {
		return nil, errors.New("Error parsing 'tpduSize' field " + _tpduSizeErr.Error())
	}

	// Create the instance
	return NewCOTPParameterTpduSize(tpduSize), nil
}

func (m COTPParameterTpduSize) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Enum field (tpduSize)
		tpduSize := CastCOTPTpduSize(m.tpduSize)
		tpduSize.Serialize(io)

	}
	COTPParameterSerialize(io, m.COTPParameter, CastICOTPParameter(m), ser)
}
