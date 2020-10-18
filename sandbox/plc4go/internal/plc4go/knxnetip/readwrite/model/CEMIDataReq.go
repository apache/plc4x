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
type CEMIDataReq struct {
	AdditionalInformationLength uint8
	AdditionalInformation       []ICEMIAdditionalInformation
	CemiDataFrame               ICEMIDataFrame
	CEMI
}

// The corresponding interface
type ICEMIDataReq interface {
	ICEMI
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIDataReq) MessageCode() uint8 {
	return 0x11
}

func (m CEMIDataReq) initialize() spi.Message {
	return m
}

func NewCEMIDataReq(additionalInformationLength uint8, additionalInformation []ICEMIAdditionalInformation, cemiDataFrame ICEMIDataFrame) CEMIInitializer {
	return &CEMIDataReq{AdditionalInformationLength: additionalInformationLength, AdditionalInformation: additionalInformation, CemiDataFrame: cemiDataFrame}
}

func CastICEMIDataReq(structType interface{}) ICEMIDataReq {
	castFunc := func(typ interface{}) ICEMIDataReq {
		if iCEMIDataReq, ok := typ.(ICEMIDataReq); ok {
			return iCEMIDataReq
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIDataReq(structType interface{}) CEMIDataReq {
	castFunc := func(typ interface{}) CEMIDataReq {
		if sCEMIDataReq, ok := typ.(CEMIDataReq); ok {
			return sCEMIDataReq
		}
		return CEMIDataReq{}
	}
	return castFunc(structType)
}

func (m CEMIDataReq) LengthInBits() uint16 {
	var lengthInBits = m.CEMI.LengthInBits()

	// Simple field (additionalInformationLength)
	lengthInBits += 8

	// Array field
	if len(m.AdditionalInformation) > 0 {
		for _, element := range m.AdditionalInformation {
			lengthInBits += element.LengthInBits()
		}
	}

	// Simple field (cemiDataFrame)
	lengthInBits += m.CemiDataFrame.LengthInBits()

	return lengthInBits
}

func (m CEMIDataReq) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIDataReqParse(io *spi.ReadBuffer) (CEMIInitializer, error) {

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

	// Simple Field (cemiDataFrame)
	_cemiDataFrameMessage, _err := CEMIDataFrameParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'cemiDataFrame'. " + _err.Error())
	}
	var cemiDataFrame ICEMIDataFrame
	cemiDataFrame, _cemiDataFrameOk := _cemiDataFrameMessage.(ICEMIDataFrame)
	if !_cemiDataFrameOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_cemiDataFrameMessage).Name() + " to ICEMIDataFrame")
	}

	// Create the instance
	return NewCEMIDataReq(additionalInformationLength, additionalInformation, cemiDataFrame), nil
}

func (m CEMIDataReq) Serialize(io spi.WriteBuffer) error {
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

		// Simple Field (cemiDataFrame)
		cemiDataFrame := CastICEMIDataFrame(m.CemiDataFrame)
		_cemiDataFrameErr := cemiDataFrame.Serialize(io)
		if _cemiDataFrameErr != nil {
			return errors.New("Error serializing 'cemiDataFrame' field " + _cemiDataFrameErr.Error())
		}

		return nil
	}
	return CEMISerialize(io, m.CEMI, CastICEMI(m), ser)
}
