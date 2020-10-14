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
	"strconv"
)

// Constant values.
const CEMIAdditionalInformationRelativeTimestamp_LEN uint8 = 2

// The data-structure of this message
type CEMIAdditionalInformationRelativeTimestamp struct {
	RelativeTimestamp IRelativeTimestamp
	CEMIAdditionalInformation
}

// The corresponding interface
type ICEMIAdditionalInformationRelativeTimestamp interface {
	ICEMIAdditionalInformation
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIAdditionalInformationRelativeTimestamp) AdditionalInformationType() uint8 {
	return 0x04
}

func (m CEMIAdditionalInformationRelativeTimestamp) initialize() spi.Message {
	return m
}

func NewCEMIAdditionalInformationRelativeTimestamp(relativeTimestamp IRelativeTimestamp) CEMIAdditionalInformationInitializer {
	return &CEMIAdditionalInformationRelativeTimestamp{RelativeTimestamp: relativeTimestamp}
}

func CastICEMIAdditionalInformationRelativeTimestamp(structType interface{}) ICEMIAdditionalInformationRelativeTimestamp {
	castFunc := func(typ interface{}) ICEMIAdditionalInformationRelativeTimestamp {
		if iCEMIAdditionalInformationRelativeTimestamp, ok := typ.(ICEMIAdditionalInformationRelativeTimestamp); ok {
			return iCEMIAdditionalInformationRelativeTimestamp
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIAdditionalInformationRelativeTimestamp(structType interface{}) CEMIAdditionalInformationRelativeTimestamp {
	castFunc := func(typ interface{}) CEMIAdditionalInformationRelativeTimestamp {
		if sCEMIAdditionalInformationRelativeTimestamp, ok := typ.(CEMIAdditionalInformationRelativeTimestamp); ok {
			return sCEMIAdditionalInformationRelativeTimestamp
		}
		return CEMIAdditionalInformationRelativeTimestamp{}
	}
	return castFunc(structType)
}

func (m CEMIAdditionalInformationRelativeTimestamp) LengthInBits() uint16 {
	var lengthInBits uint16 = m.CEMIAdditionalInformation.LengthInBits()

	// Const Field (len)
	lengthInBits += 8

	// Simple field (relativeTimestamp)
	lengthInBits += m.RelativeTimestamp.LengthInBits()

	return lengthInBits
}

func (m CEMIAdditionalInformationRelativeTimestamp) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIAdditionalInformationRelativeTimestampParse(io *spi.ReadBuffer) (CEMIAdditionalInformationInitializer, error) {

	// Const Field (len)
	len, _lenErr := io.ReadUint8(8)
	if _lenErr != nil {
		return nil, errors.New("Error parsing 'len' field " + _lenErr.Error())
	}
	if len != CEMIAdditionalInformationRelativeTimestamp_LEN {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(CEMIAdditionalInformationRelativeTimestamp_LEN)) + " but got " + strconv.Itoa(int(len)))
	}

	// Simple Field (relativeTimestamp)
	_relativeTimestampMessage, _err := RelativeTimestampParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'relativeTimestamp'. " + _err.Error())
	}
	var relativeTimestamp IRelativeTimestamp
	relativeTimestamp, _relativeTimestampOk := _relativeTimestampMessage.(IRelativeTimestamp)
	if !_relativeTimestampOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_relativeTimestampMessage).Name() + " to IRelativeTimestamp")
	}

	// Create the instance
	return NewCEMIAdditionalInformationRelativeTimestamp(relativeTimestamp), nil
}

func (m CEMIAdditionalInformationRelativeTimestamp) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Const Field (len)
		_lenErr := io.WriteUint8(8, 2)
		if _lenErr != nil {
			return errors.New("Error serializing 'len' field " + _lenErr.Error())
		}

		// Simple Field (relativeTimestamp)
		relativeTimestamp := CastIRelativeTimestamp(m.RelativeTimestamp)
		_relativeTimestampErr := relativeTimestamp.Serialize(io)
		if _relativeTimestampErr != nil {
			return errors.New("Error serializing 'relativeTimestamp' field " + _relativeTimestampErr.Error())
		}

		return nil
	}
	return CEMIAdditionalInformationSerialize(io, m.CEMIAdditionalInformation, CastICEMIAdditionalInformation(m), ser)
}
