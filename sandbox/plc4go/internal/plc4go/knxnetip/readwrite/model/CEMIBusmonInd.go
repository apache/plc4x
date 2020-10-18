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
	AdditionalInformationLength uint8
	AdditionalInformation       []ICEMIAdditionalInformation
	CemiFrame                   ICEMIFrame
	CEMI
}

// The corresponding interface
type ICEMIBusmonInd interface {
	ICEMI
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIBusmonInd) MessageCode() uint8 {
	return 0x2B
}

func (m CEMIBusmonInd) initialize() spi.Message {
	return m
}

func NewCEMIBusmonInd(additionalInformationLength uint8, additionalInformation []ICEMIAdditionalInformation, cemiFrame ICEMIFrame) CEMIInitializer {
	return &CEMIBusmonInd{AdditionalInformationLength: additionalInformationLength, AdditionalInformation: additionalInformation, CemiFrame: cemiFrame}
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
	var lengthInBits = m.CEMI.LengthInBits()

	// Simple field (additionalInformationLength)
	lengthInBits += 8

	// Array field
	if len(m.AdditionalInformation) > 0 {
		for _, element := range m.AdditionalInformation {
			lengthInBits += element.LengthInBits()
		}
	}

	// Simple field (cemiFrame)
	lengthInBits += m.CemiFrame.LengthInBits()

	return lengthInBits
}

func (m CEMIBusmonInd) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIBusmonIndParse(io *spi.ReadBuffer) (CEMIInitializer, error) {

	// Simple Field (additionalInformationLength)
	additionalInformationLength, _additionalInformationLengthErr := io.ReadUint8(8)
	if _additionalInformationLengthErr != nil {
		return nil, errors.New("Error parsing 'additionalInformationLength' field " + _additionalInformationLengthErr.Error())
	}

	// Array field (additionalInformation)
	// Length array
	additionalInformation := make([]ICEMIAdditionalInformation, 0)
	_additionalInformationLength := additionalInformationLength
	_additionalInformationEndPos := io.GetPos() + uint16(_additionalInformationLength)
	for io.GetPos() < _additionalInformationEndPos {
		_message, _err := CEMIAdditionalInformationParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'additionalInformation' field " + _err.Error())
		}
		var _item ICEMIAdditionalInformation
		_item, _ok := _message.(ICEMIAdditionalInformation)
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
	var cemiFrame ICEMIFrame
	cemiFrame, _cemiFrameOk := _cemiFrameMessage.(ICEMIFrame)
	if !_cemiFrameOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_cemiFrameMessage).Name() + " to ICEMIFrame")
	}

	// Create the instance
	return NewCEMIBusmonInd(additionalInformationLength, additionalInformation, cemiFrame), nil
}

func (m CEMIBusmonInd) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (additionalInformationLength)
		additionalInformationLength := uint8(m.AdditionalInformationLength)
		_additionalInformationLengthErr := io.WriteUint8(8, additionalInformationLength)
		if _additionalInformationLengthErr != nil {
			return errors.New("Error serializing 'additionalInformationLength' field " + _additionalInformationLengthErr.Error())
		}

		// Array Field (additionalInformation)
		if m.AdditionalInformation != nil {
			for _, _element := range m.AdditionalInformation {
				_elementErr := _element.Serialize(io)
				if _elementErr != nil {
					return errors.New("Error serializing 'additionalInformation' field " + _elementErr.Error())
				}
			}
		}

		// Simple Field (cemiFrame)
		cemiFrame := CastICEMIFrame(m.CemiFrame)
		_cemiFrameErr := cemiFrame.Serialize(io)
		if _cemiFrameErr != nil {
			return errors.New("Error serializing 'cemiFrame' field " + _cemiFrameErr.Error())
		}

		return nil
	}
	return CEMISerialize(io, m.CEMI, CastICEMI(m), ser)
}
