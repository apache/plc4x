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
    "encoding/base64"
    "encoding/xml"
    "errors"
    "io"
    log "github.com/sirupsen/logrus"
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

// The data-structure of this message
type LDataFramePollingData struct {
    SourceAddress *KnxAddress
    TargetAddress []int8
    NumberExpectedPollData uint8
    Parent *LDataFrame
    ILDataFramePollingData
}

// The corresponding interface
type ILDataFramePollingData interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *LDataFramePollingData) NotAckFrame() bool {
    return true
}

func (m *LDataFramePollingData) ExtendedFrame() bool {
    return false
}

func (m *LDataFramePollingData) Polling() bool {
    return true
}


func (m *LDataFramePollingData) InitializeParent(parent *LDataFrame, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) {
    m.Parent.Repeated = repeated
    m.Parent.Priority = priority
    m.Parent.AcknowledgeRequested = acknowledgeRequested
    m.Parent.ErrorFlag = errorFlag
}

func NewLDataFramePollingData(sourceAddress *KnxAddress, targetAddress []int8, numberExpectedPollData uint8, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) *LDataFrame {
    child := &LDataFramePollingData{
        SourceAddress: sourceAddress,
        TargetAddress: targetAddress,
        NumberExpectedPollData: numberExpectedPollData,
        Parent: NewLDataFrame(repeated, priority, acknowledgeRequested, errorFlag),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastLDataFramePollingData(structType interface{}) *LDataFramePollingData {
    castFunc := func(typ interface{}) *LDataFramePollingData {
        if casted, ok := typ.(LDataFramePollingData); ok {
            return &casted
        }
        if casted, ok := typ.(*LDataFramePollingData); ok {
            return casted
        }
        if casted, ok := typ.(LDataFrame); ok {
            return CastLDataFramePollingData(casted.Child)
        }
        if casted, ok := typ.(*LDataFrame); ok {
            return CastLDataFramePollingData(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *LDataFramePollingData) GetTypeName() string {
    return "LDataFramePollingData"
}

func (m *LDataFramePollingData) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (sourceAddress)
    lengthInBits += m.SourceAddress.LengthInBits()

    // Array field
    if len(m.TargetAddress) > 0 {
        lengthInBits += 8 * uint16(len(m.TargetAddress))
    }

    // Reserved Field (reserved)
    lengthInBits += 4

    // Simple field (numberExpectedPollData)
    lengthInBits += 6

    return lengthInBits
}

func (m *LDataFramePollingData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func LDataFramePollingDataParse(io *utils.ReadBuffer) (*LDataFrame, error) {

    // Simple Field (sourceAddress)
    sourceAddress, _sourceAddressErr := KnxAddressParse(io)
    if _sourceAddressErr != nil {
        return nil, errors.New("Error parsing 'sourceAddress' field " + _sourceAddressErr.Error())
    }

    // Array field (targetAddress)
    // Count array
    targetAddress := make([]int8, uint16(2))
    for curItem := uint16(0); curItem < uint16(uint16(2)); curItem++ {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'targetAddress' field " + _err.Error())
        }
        targetAddress[curItem] = _item
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(4)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0x00) {
            log.WithFields(log.Fields{
                "expected value": uint8(0x00),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (numberExpectedPollData)
    numberExpectedPollData, _numberExpectedPollDataErr := io.ReadUint8(6)
    if _numberExpectedPollDataErr != nil {
        return nil, errors.New("Error parsing 'numberExpectedPollData' field " + _numberExpectedPollDataErr.Error())
    }

    // Create a partially initialized instance
    _child := &LDataFramePollingData{
        SourceAddress: sourceAddress,
        TargetAddress: targetAddress,
        NumberExpectedPollData: numberExpectedPollData,
        Parent: &LDataFrame{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *LDataFramePollingData) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (sourceAddress)
    _sourceAddressErr := m.SourceAddress.Serialize(io)
    if _sourceAddressErr != nil {
        return errors.New("Error serializing 'sourceAddress' field " + _sourceAddressErr.Error())
    }

    // Array Field (targetAddress)
    if m.TargetAddress != nil {
        for _, _element := range m.TargetAddress {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'targetAddress' field " + _elementErr.Error())
            }
        }
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(4, uint8(0x00))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (numberExpectedPollData)
    numberExpectedPollData := uint8(m.NumberExpectedPollData)
    _numberExpectedPollDataErr := io.WriteUint8(6, (numberExpectedPollData))
    if _numberExpectedPollDataErr != nil {
        return errors.New("Error serializing 'numberExpectedPollData' field " + _numberExpectedPollDataErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *LDataFramePollingData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "sourceAddress":
                var data *KnxAddress
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.SourceAddress = data
            case "targetAddress":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.TargetAddress = utils.ByteArrayToInt8Array(_decoded[0:_len])
            case "numberExpectedPollData":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.NumberExpectedPollData = data
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

func (m *LDataFramePollingData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.SourceAddress, xml.StartElement{Name: xml.Name{Local: "sourceAddress"}}); err != nil {
        return err
    }
    _encodedTargetAddress := make([]byte, base64.StdEncoding.EncodedLen(len(m.TargetAddress)))
    base64.StdEncoding.Encode(_encodedTargetAddress, utils.Int8ArrayToByteArray(m.TargetAddress))
    if err := e.EncodeElement(_encodedTargetAddress, xml.StartElement{Name: xml.Name{Local: "targetAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.NumberExpectedPollData, xml.StartElement{Name: xml.Name{Local: "numberExpectedPollData"}}); err != nil {
        return err
    }
    return nil
}

