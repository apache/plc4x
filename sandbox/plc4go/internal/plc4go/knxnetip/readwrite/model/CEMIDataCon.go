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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "reflect"
)

// The data-structure of this message
type CEMIDataCon struct {
    AdditionalInformationLength uint8
    AdditionalInformation []ICEMIAdditionalInformation
    CemiDataFrame ICEMIDataFrame
    CEMI
}

// The corresponding interface
type ICEMIDataCon interface {
    ICEMI
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIDataCon) MessageCode() uint8 {
    return 0x2E
}

func (m CEMIDataCon) initialize() spi.Message {
    return m
}

func NewCEMIDataCon(additionalInformationLength uint8, additionalInformation []ICEMIAdditionalInformation, cemiDataFrame ICEMIDataFrame) CEMIInitializer {
    return &CEMIDataCon{AdditionalInformationLength: additionalInformationLength, AdditionalInformation: additionalInformation, CemiDataFrame: cemiDataFrame}
}

func CastICEMIDataCon(structType interface{}) ICEMIDataCon {
    castFunc := func(typ interface{}) ICEMIDataCon {
        if iCEMIDataCon, ok := typ.(ICEMIDataCon); ok {
            return iCEMIDataCon
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIDataCon(structType interface{}) CEMIDataCon {
    castFunc := func(typ interface{}) CEMIDataCon {
        if sCEMIDataCon, ok := typ.(CEMIDataCon); ok {
            return sCEMIDataCon
        }
        if sCEMIDataCon, ok := typ.(*CEMIDataCon); ok {
            return *sCEMIDataCon
        }
        return CEMIDataCon{}
    }
    return castFunc(structType)
}

func (m CEMIDataCon) LengthInBits() uint16 {
    var lengthInBits uint16 = m.CEMI.LengthInBits()

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

func (m CEMIDataCon) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIDataConParse(io *utils.ReadBuffer) (CEMIInitializer, error) {

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
    for ;io.GetPos() < _additionalInformationEndPos; {
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
    return NewCEMIDataCon(additionalInformationLength, additionalInformation, cemiDataFrame), nil
}

func (m CEMIDataCon) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (additionalInformationLength)
    additionalInformationLength := uint8(m.AdditionalInformationLength)
    _additionalInformationLengthErr := io.WriteUint8(8, (additionalInformationLength))
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
