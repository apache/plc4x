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
    "strconv"
)

// Constant values.
const CEMIAdditionalInformationBusmonitorInfo_LEN uint8 = 1

// The data-structure of this message
type CEMIAdditionalInformationBusmonitorInfo struct {
    FrameErrorFlag bool
    BitErrorFlag bool
    ParityErrorFlag bool
    UnknownFlag bool
    LostFlag bool
    SequenceNumber uint8
    CEMIAdditionalInformation
}

// The corresponding interface
type ICEMIAdditionalInformationBusmonitorInfo interface {
    ICEMIAdditionalInformation
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIAdditionalInformationBusmonitorInfo) AdditionalInformationType() uint8 {
    return 0x03
}

func (m CEMIAdditionalInformationBusmonitorInfo) initialize() spi.Message {
    return m
}

func NewCEMIAdditionalInformationBusmonitorInfo(frameErrorFlag bool, bitErrorFlag bool, parityErrorFlag bool, unknownFlag bool, lostFlag bool, sequenceNumber uint8) CEMIAdditionalInformationInitializer {
    return &CEMIAdditionalInformationBusmonitorInfo{FrameErrorFlag: frameErrorFlag, BitErrorFlag: bitErrorFlag, ParityErrorFlag: parityErrorFlag, UnknownFlag: unknownFlag, LostFlag: lostFlag, SequenceNumber: sequenceNumber}
}

func CastICEMIAdditionalInformationBusmonitorInfo(structType interface{}) ICEMIAdditionalInformationBusmonitorInfo {
    castFunc := func(typ interface{}) ICEMIAdditionalInformationBusmonitorInfo {
        if iCEMIAdditionalInformationBusmonitorInfo, ok := typ.(ICEMIAdditionalInformationBusmonitorInfo); ok {
            return iCEMIAdditionalInformationBusmonitorInfo
        }
        return nil
    }
    return castFunc(structType)
}

func CastCEMIAdditionalInformationBusmonitorInfo(structType interface{}) CEMIAdditionalInformationBusmonitorInfo {
    castFunc := func(typ interface{}) CEMIAdditionalInformationBusmonitorInfo {
        if sCEMIAdditionalInformationBusmonitorInfo, ok := typ.(CEMIAdditionalInformationBusmonitorInfo); ok {
            return sCEMIAdditionalInformationBusmonitorInfo
        }
        if sCEMIAdditionalInformationBusmonitorInfo, ok := typ.(*CEMIAdditionalInformationBusmonitorInfo); ok {
            return *sCEMIAdditionalInformationBusmonitorInfo
        }
        return CEMIAdditionalInformationBusmonitorInfo{}
    }
    return castFunc(structType)
}

func (m CEMIAdditionalInformationBusmonitorInfo) LengthInBits() uint16 {
    var lengthInBits uint16 = m.CEMIAdditionalInformation.LengthInBits()

    // Const Field (len)
    lengthInBits += 8

    // Simple field (frameErrorFlag)
    lengthInBits += 1

    // Simple field (bitErrorFlag)
    lengthInBits += 1

    // Simple field (parityErrorFlag)
    lengthInBits += 1

    // Simple field (unknownFlag)
    lengthInBits += 1

    // Simple field (lostFlag)
    lengthInBits += 1

    // Simple field (sequenceNumber)
    lengthInBits += 3

    return lengthInBits
}

func (m CEMIAdditionalInformationBusmonitorInfo) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIAdditionalInformationBusmonitorInfoParse(io *spi.ReadBuffer) (CEMIAdditionalInformationInitializer, error) {

    // Const Field (len)
    len, _lenErr := io.ReadUint8(8)
    if _lenErr != nil {
        return nil, errors.New("Error parsing 'len' field " + _lenErr.Error())
    }
    if len != CEMIAdditionalInformationBusmonitorInfo_LEN {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(CEMIAdditionalInformationBusmonitorInfo_LEN)) + " but got " + strconv.Itoa(int(len)))
    }

    // Simple Field (frameErrorFlag)
    frameErrorFlag, _frameErrorFlagErr := io.ReadBit()
    if _frameErrorFlagErr != nil {
        return nil, errors.New("Error parsing 'frameErrorFlag' field " + _frameErrorFlagErr.Error())
    }

    // Simple Field (bitErrorFlag)
    bitErrorFlag, _bitErrorFlagErr := io.ReadBit()
    if _bitErrorFlagErr != nil {
        return nil, errors.New("Error parsing 'bitErrorFlag' field " + _bitErrorFlagErr.Error())
    }

    // Simple Field (parityErrorFlag)
    parityErrorFlag, _parityErrorFlagErr := io.ReadBit()
    if _parityErrorFlagErr != nil {
        return nil, errors.New("Error parsing 'parityErrorFlag' field " + _parityErrorFlagErr.Error())
    }

    // Simple Field (unknownFlag)
    unknownFlag, _unknownFlagErr := io.ReadBit()
    if _unknownFlagErr != nil {
        return nil, errors.New("Error parsing 'unknownFlag' field " + _unknownFlagErr.Error())
    }

    // Simple Field (lostFlag)
    lostFlag, _lostFlagErr := io.ReadBit()
    if _lostFlagErr != nil {
        return nil, errors.New("Error parsing 'lostFlag' field " + _lostFlagErr.Error())
    }

    // Simple Field (sequenceNumber)
    sequenceNumber, _sequenceNumberErr := io.ReadUint8(3)
    if _sequenceNumberErr != nil {
        return nil, errors.New("Error parsing 'sequenceNumber' field " + _sequenceNumberErr.Error())
    }

    // Create the instance
    return NewCEMIAdditionalInformationBusmonitorInfo(frameErrorFlag, bitErrorFlag, parityErrorFlag, unknownFlag, lostFlag, sequenceNumber), nil
}

func (m CEMIAdditionalInformationBusmonitorInfo) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Const Field (len)
    _lenErr := io.WriteUint8(8, 1)
    if _lenErr != nil {
        return errors.New("Error serializing 'len' field " + _lenErr.Error())
    }

    // Simple Field (frameErrorFlag)
    frameErrorFlag := bool(m.FrameErrorFlag)
    _frameErrorFlagErr := io.WriteBit((bool) (frameErrorFlag))
    if _frameErrorFlagErr != nil {
        return errors.New("Error serializing 'frameErrorFlag' field " + _frameErrorFlagErr.Error())
    }

    // Simple Field (bitErrorFlag)
    bitErrorFlag := bool(m.BitErrorFlag)
    _bitErrorFlagErr := io.WriteBit((bool) (bitErrorFlag))
    if _bitErrorFlagErr != nil {
        return errors.New("Error serializing 'bitErrorFlag' field " + _bitErrorFlagErr.Error())
    }

    // Simple Field (parityErrorFlag)
    parityErrorFlag := bool(m.ParityErrorFlag)
    _parityErrorFlagErr := io.WriteBit((bool) (parityErrorFlag))
    if _parityErrorFlagErr != nil {
        return errors.New("Error serializing 'parityErrorFlag' field " + _parityErrorFlagErr.Error())
    }

    // Simple Field (unknownFlag)
    unknownFlag := bool(m.UnknownFlag)
    _unknownFlagErr := io.WriteBit((bool) (unknownFlag))
    if _unknownFlagErr != nil {
        return errors.New("Error serializing 'unknownFlag' field " + _unknownFlagErr.Error())
    }

    // Simple Field (lostFlag)
    lostFlag := bool(m.LostFlag)
    _lostFlagErr := io.WriteBit((bool) (lostFlag))
    if _lostFlagErr != nil {
        return errors.New("Error serializing 'lostFlag' field " + _lostFlagErr.Error())
    }

    // Simple Field (sequenceNumber)
    sequenceNumber := uint8(m.SequenceNumber)
    _sequenceNumberErr := io.WriteUint8(3, (sequenceNumber))
    if _sequenceNumberErr != nil {
        return errors.New("Error serializing 'sequenceNumber' field " + _sequenceNumberErr.Error())
    }

        return nil
    }
    return CEMIAdditionalInformationSerialize(io, m.CEMIAdditionalInformation, CastICEMIAdditionalInformation(m), ser)
}
