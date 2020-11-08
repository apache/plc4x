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
    log "github.com/sirupsen/logrus"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

// The data-structure of this message
type APDUComplexAck struct {
    SegmentedMessage bool
    MoreFollows bool
    OriginalInvokeId uint8
    SequenceNumber *uint8
    ProposedWindowSize *uint8
    ServiceAck *BACnetServiceAck
    Parent *APDU
    IAPDUComplexAck
}

// The corresponding interface
type IAPDUComplexAck interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *APDUComplexAck) ApduType() uint8 {
    return 0x3
}


func (m *APDUComplexAck) InitializeParent(parent *APDU) {
}

func NewAPDUComplexAck(segmentedMessage bool, moreFollows bool, originalInvokeId uint8, sequenceNumber *uint8, proposedWindowSize *uint8, serviceAck *BACnetServiceAck, ) *APDU {
    child := &APDUComplexAck{
        SegmentedMessage: segmentedMessage,
        MoreFollows: moreFollows,
        OriginalInvokeId: originalInvokeId,
        SequenceNumber: sequenceNumber,
        ProposedWindowSize: proposedWindowSize,
        ServiceAck: serviceAck,
        Parent: NewAPDU(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastAPDUComplexAck(structType interface{}) APDUComplexAck {
    castFunc := func(typ interface{}) APDUComplexAck {
        if casted, ok := typ.(APDUComplexAck); ok {
            return casted
        }
        if casted, ok := typ.(*APDUComplexAck); ok {
            return *casted
        }
        if casted, ok := typ.(APDU); ok {
            return CastAPDUComplexAck(casted.Child)
        }
        if casted, ok := typ.(*APDU); ok {
            return CastAPDUComplexAck(casted.Child)
        }
        return APDUComplexAck{}
    }
    return castFunc(structType)
}

func (m *APDUComplexAck) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (segmentedMessage)
    lengthInBits += 1

    // Simple field (moreFollows)
    lengthInBits += 1

    // Reserved Field (reserved)
    lengthInBits += 2

    // Simple field (originalInvokeId)
    lengthInBits += 8

    // Optional Field (sequenceNumber)
    if m.SequenceNumber != nil {
        lengthInBits += 8
    }

    // Optional Field (proposedWindowSize)
    if m.ProposedWindowSize != nil {
        lengthInBits += 8
    }

    // Simple field (serviceAck)
    lengthInBits += m.ServiceAck.LengthInBits()

    return lengthInBits
}

func (m *APDUComplexAck) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func APDUComplexAckParse(io *utils.ReadBuffer) (*APDU, error) {

    // Simple Field (segmentedMessage)
    segmentedMessage, _segmentedMessageErr := io.ReadBit()
    if _segmentedMessageErr != nil {
        return nil, errors.New("Error parsing 'segmentedMessage' field " + _segmentedMessageErr.Error())
    }

    // Simple Field (moreFollows)
    moreFollows, _moreFollowsErr := io.ReadBit()
    if _moreFollowsErr != nil {
        return nil, errors.New("Error parsing 'moreFollows' field " + _moreFollowsErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(2)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0) {
            log.WithFields(log.Fields{
                "expected value": uint8(0),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (originalInvokeId)
    originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
    if _originalInvokeIdErr != nil {
        return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
    }

    // Optional Field (sequenceNumber) (Can be skipped, if a given expression evaluates to false)
    var sequenceNumber *uint8 = nil
    if segmentedMessage {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'sequenceNumber' field " + _err.Error())
        }

        sequenceNumber = &_val
    }

    // Optional Field (proposedWindowSize) (Can be skipped, if a given expression evaluates to false)
    var proposedWindowSize *uint8 = nil
    if segmentedMessage {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'proposedWindowSize' field " + _err.Error())
        }

        proposedWindowSize = &_val
    }

    // Simple Field (serviceAck)
    serviceAck, _serviceAckErr := BACnetServiceAckParse(io)
    if _serviceAckErr != nil {
        return nil, errors.New("Error parsing 'serviceAck' field " + _serviceAckErr.Error())
    }

    // Create a partially initialized instance
    _child := &APDUComplexAck{
        SegmentedMessage: segmentedMessage,
        MoreFollows: moreFollows,
        OriginalInvokeId: originalInvokeId,
        SequenceNumber: sequenceNumber,
        ProposedWindowSize: proposedWindowSize,
        ServiceAck: serviceAck,
        Parent: &APDU{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *APDUComplexAck) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (segmentedMessage)
    segmentedMessage := bool(m.SegmentedMessage)
    _segmentedMessageErr := io.WriteBit((segmentedMessage))
    if _segmentedMessageErr != nil {
        return errors.New("Error serializing 'segmentedMessage' field " + _segmentedMessageErr.Error())
    }

    // Simple Field (moreFollows)
    moreFollows := bool(m.MoreFollows)
    _moreFollowsErr := io.WriteBit((moreFollows))
    if _moreFollowsErr != nil {
        return errors.New("Error serializing 'moreFollows' field " + _moreFollowsErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(2, uint8(0))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (originalInvokeId)
    originalInvokeId := uint8(m.OriginalInvokeId)
    _originalInvokeIdErr := io.WriteUint8(8, (originalInvokeId))
    if _originalInvokeIdErr != nil {
        return errors.New("Error serializing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
    }

    // Optional Field (sequenceNumber) (Can be skipped, if the value is null)
    var sequenceNumber *uint8 = nil
    if m.SequenceNumber != nil {
        sequenceNumber = m.SequenceNumber
        _sequenceNumberErr := io.WriteUint8(8, *(sequenceNumber))
        if _sequenceNumberErr != nil {
            return errors.New("Error serializing 'sequenceNumber' field " + _sequenceNumberErr.Error())
        }
    }

    // Optional Field (proposedWindowSize) (Can be skipped, if the value is null)
    var proposedWindowSize *uint8 = nil
    if m.ProposedWindowSize != nil {
        proposedWindowSize = m.ProposedWindowSize
        _proposedWindowSizeErr := io.WriteUint8(8, *(proposedWindowSize))
        if _proposedWindowSizeErr != nil {
            return errors.New("Error serializing 'proposedWindowSize' field " + _proposedWindowSizeErr.Error())
        }
    }

    // Simple Field (serviceAck)
    _serviceAckErr := m.ServiceAck.Serialize(io)
    if _serviceAckErr != nil {
        return errors.New("Error serializing 'serviceAck' field " + _serviceAckErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *APDUComplexAck) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "segmentedMessage":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SegmentedMessage = data
            case "moreFollows":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MoreFollows = data
            case "originalInvokeId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.OriginalInvokeId = data
            case "sequenceNumber":
                var data *uint8
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.SequenceNumber = data
            case "proposedWindowSize":
                var data *uint8
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.ProposedWindowSize = data
            case "serviceAck":
                var dt *BACnetServiceAck
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.ServiceAck = dt
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

func (m *APDUComplexAck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.SegmentedMessage, xml.StartElement{Name: xml.Name{Local: "segmentedMessage"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MoreFollows, xml.StartElement{Name: xml.Name{Local: "moreFollows"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.OriginalInvokeId, xml.StartElement{Name: xml.Name{Local: "originalInvokeId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SequenceNumber, xml.StartElement{Name: xml.Name{Local: "sequenceNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ProposedWindowSize, xml.StartElement{Name: xml.Name{Local: "proposedWindowSize"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ServiceAck, xml.StartElement{Name: xml.Name{Local: "serviceAck"}}); err != nil {
        return err
    }
    return nil
}

