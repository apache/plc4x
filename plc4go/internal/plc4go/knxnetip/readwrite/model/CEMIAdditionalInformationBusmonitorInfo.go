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
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
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
    Parent *CEMIAdditionalInformation
    ICEMIAdditionalInformationBusmonitorInfo
}

// The corresponding interface
type ICEMIAdditionalInformationBusmonitorInfo interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *CEMIAdditionalInformationBusmonitorInfo) AdditionalInformationType() uint8 {
    return 0x03
}


func (m *CEMIAdditionalInformationBusmonitorInfo) InitializeParent(parent *CEMIAdditionalInformation) {
}

func NewCEMIAdditionalInformationBusmonitorInfo(frameErrorFlag bool, bitErrorFlag bool, parityErrorFlag bool, unknownFlag bool, lostFlag bool, sequenceNumber uint8, ) *CEMIAdditionalInformation {
    child := &CEMIAdditionalInformationBusmonitorInfo{
        FrameErrorFlag: frameErrorFlag,
        BitErrorFlag: bitErrorFlag,
        ParityErrorFlag: parityErrorFlag,
        UnknownFlag: unknownFlag,
        LostFlag: lostFlag,
        SequenceNumber: sequenceNumber,
        Parent: NewCEMIAdditionalInformation(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCEMIAdditionalInformationBusmonitorInfo(structType interface{}) CEMIAdditionalInformationBusmonitorInfo {
    castFunc := func(typ interface{}) CEMIAdditionalInformationBusmonitorInfo {
        if casted, ok := typ.(CEMIAdditionalInformationBusmonitorInfo); ok {
            return casted
        }
        if casted, ok := typ.(*CEMIAdditionalInformationBusmonitorInfo); ok {
            return *casted
        }
        if casted, ok := typ.(CEMIAdditionalInformation); ok {
            return CastCEMIAdditionalInformationBusmonitorInfo(casted.Child)
        }
        if casted, ok := typ.(*CEMIAdditionalInformation); ok {
            return CastCEMIAdditionalInformationBusmonitorInfo(casted.Child)
        }
        return CEMIAdditionalInformationBusmonitorInfo{}
    }
    return castFunc(structType)
}

func (m *CEMIAdditionalInformationBusmonitorInfo) LengthInBits() uint16 {
    lengthInBits := uint16(0)

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

func (m *CEMIAdditionalInformationBusmonitorInfo) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIAdditionalInformationBusmonitorInfoParse(io *utils.ReadBuffer) (*CEMIAdditionalInformation, error) {

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

    // Create a partially initialized instance
    _child := &CEMIAdditionalInformationBusmonitorInfo{
        FrameErrorFlag: frameErrorFlag,
        BitErrorFlag: bitErrorFlag,
        ParityErrorFlag: parityErrorFlag,
        UnknownFlag: unknownFlag,
        LostFlag: lostFlag,
        SequenceNumber: sequenceNumber,
        Parent: &CEMIAdditionalInformation{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *CEMIAdditionalInformationBusmonitorInfo) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Const Field (len)
    _lenErr := io.WriteUint8(8, 1)
    if _lenErr != nil {
        return errors.New("Error serializing 'len' field " + _lenErr.Error())
    }

    // Simple Field (frameErrorFlag)
    frameErrorFlag := bool(m.FrameErrorFlag)
    _frameErrorFlagErr := io.WriteBit((frameErrorFlag))
    if _frameErrorFlagErr != nil {
        return errors.New("Error serializing 'frameErrorFlag' field " + _frameErrorFlagErr.Error())
    }

    // Simple Field (bitErrorFlag)
    bitErrorFlag := bool(m.BitErrorFlag)
    _bitErrorFlagErr := io.WriteBit((bitErrorFlag))
    if _bitErrorFlagErr != nil {
        return errors.New("Error serializing 'bitErrorFlag' field " + _bitErrorFlagErr.Error())
    }

    // Simple Field (parityErrorFlag)
    parityErrorFlag := bool(m.ParityErrorFlag)
    _parityErrorFlagErr := io.WriteBit((parityErrorFlag))
    if _parityErrorFlagErr != nil {
        return errors.New("Error serializing 'parityErrorFlag' field " + _parityErrorFlagErr.Error())
    }

    // Simple Field (unknownFlag)
    unknownFlag := bool(m.UnknownFlag)
    _unknownFlagErr := io.WriteBit((unknownFlag))
    if _unknownFlagErr != nil {
        return errors.New("Error serializing 'unknownFlag' field " + _unknownFlagErr.Error())
    }

    // Simple Field (lostFlag)
    lostFlag := bool(m.LostFlag)
    _lostFlagErr := io.WriteBit((lostFlag))
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
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *CEMIAdditionalInformationBusmonitorInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "frameErrorFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.FrameErrorFlag = data
            case "bitErrorFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.BitErrorFlag = data
            case "parityErrorFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ParityErrorFlag = data
            case "unknownFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.UnknownFlag = data
            case "lostFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.LostFlag = data
            case "sequenceNumber":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SequenceNumber = data
            }
        }
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *CEMIAdditionalInformationBusmonitorInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.FrameErrorFlag, xml.StartElement{Name: xml.Name{Local: "frameErrorFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.BitErrorFlag, xml.StartElement{Name: xml.Name{Local: "bitErrorFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ParityErrorFlag, xml.StartElement{Name: xml.Name{Local: "parityErrorFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.UnknownFlag, xml.StartElement{Name: xml.Name{Local: "unknownFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.LostFlag, xml.StartElement{Name: xml.Name{Local: "lostFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SequenceNumber, xml.StartElement{Name: xml.Name{Local: "sequenceNumber"}}); err != nil {
        return err
    }
    return nil
}

