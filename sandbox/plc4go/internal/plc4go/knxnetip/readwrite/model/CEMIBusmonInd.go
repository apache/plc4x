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
type CEMIBusmonInd struct {
	additionalInformationLength uint8
	additionalInformation       []CEMIAdditionalInformation
	cemiFrame                   CEMIFrame
	CEMI
}

// The corresponding interface
type ICEMIBusmonInd interface {
	ICEMI
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m CEMIBusmonInd) MessageCode() uint8 {
	return 0x2B
}

func (m CEMIBusmonInd) initialize() spi.Message {
	return m
}

func NewCEMIBusmonInd(additionalInformationLength uint8, additionalInformation []CEMIAdditionalInformation, cemiFrame CEMIFrame) CEMIInitializer {
	return &CEMIBusmonInd{additionalInformationLength: additionalInformationLength, additionalInformation: additionalInformation, cemiFrame: cemiFrame}
}

func CastICEMIBusmonInd(structType interface{}) ICEMIBusmonInd {
	castFunc := func(typ interface{}) ICEMIBusmonInd {
		if iCEMIBusmonInd, ok := typ.(ICEMIBusmonInd); ok {
			return iCEMIBusmonInd
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIBusmonInd(structType interface{}) CEMIBusmonInd {
	castFunc := func(typ interface{}) CEMIBusmonInd {
		if sCEMIBusmonInd, ok := typ.(CEMIBusmonInd); ok {
			return sCEMIBusmonInd
		}
		return CEMIBusmonInd{}
	}
	return castFunc(structType)
}

func (m CEMIBusmonInd) LengthInBits() uint16 {
	var lengthInBits uint16 = m.CEMI.LengthInBits()

	// Simple field (additionalInformationLength)
	lengthInBits += 8

	// Array field
	if len(m.additionalInformation) > 0 {
		for _, element := range m.additionalInformation {
			lengthInBits += element.LengthInBits()
		}
	}

	// Simple field (cemiFrame)
	lengthInBits += m.cemiFrame.LengthInBits()

	return lengthInBits
}

func (m CEMIBusmonInd) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIBusmonIndParse(io spi.ReadBuffer) (CEMIInitializer, error) {

	// Simple Field (additionalInformationLength)
	var additionalInformationLength uint8 = io.ReadUint8(8)

	// Array field (additionalInformation)
	var additionalInformation []CEMIAdditionalInformation
	// Length array
	_additionalInformationLength := additionalInformationLength
	_additionalInformationEndPos := io.GetPos() + uint16(_additionalInformationLength)
	for io.GetPos() < _additionalInformationEndPos {
		_message, _err := CEMIAdditionalInformationParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'additionalInformation' field " + _err.Error())
		}
		var _item CEMIAdditionalInformation
		_item, _ok := _message.(CEMIAdditionalInformation)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to CEMIAdditionalInformation")
		}
		additionalInformation = append(additionalInformation, _item)
	}

	// Simple Field (cemiFrame)
	_cemiFrameMessage, _err := CEMIFrameParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'cemiFrame'. " + _err.Error())
	}
	var cemiFrame CEMIFrame
	cemiFrame, _cemiFrameOk := _cemiFrameMessage.(CEMIFrame)
	if !_cemiFrameOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_cemiFrameMessage).Name() + " to CEMIFrame")
	}

	// Create the instance
	return NewCEMIBusmonInd(additionalInformationLength, additionalInformation, cemiFrame), nil
}

func (m CEMIBusmonInd) Serialize(io spi.WriteBuffer) {

	// Simple Field (additionalInformationLength)
	additionalInformationLength := uint8(m.additionalInformationLength)
	io.WriteUint8(8, (additionalInformationLength))

	// Array Field (additionalInformation)
	if m.additionalInformation != nil {
		for _, _element := range m.additionalInformation {
			_element.Serialize(io)
		}
	}

	// Simple Field (cemiFrame)
	cemiFrame := CEMIFrame(m.cemiFrame)
	cemiFrame.Serialize(io)
}
