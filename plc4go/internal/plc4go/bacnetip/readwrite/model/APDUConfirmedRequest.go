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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    log "github.com/sirupsen/logrus"
    "io"
)

// The data-structure of this message
type APDUConfirmedRequest struct {
    SegmentedMessage bool
    MoreFollows bool
    SegmentedResponseAccepted bool
    MaxSegmentsAccepted uint8
    MaxApduLengthAccepted uint8
    InvokeId uint8
    SequenceNumber *uint8
    ProposedWindowSize *uint8
    ServiceRequest *BACnetConfirmedServiceRequest
    Parent *APDU
    IAPDUConfirmedRequest
}

// The corresponding interface
type IAPDUConfirmedRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *APDUConfirmedRequest) ApduType() uint8 {
    return 0x0
}


func (m *APDUConfirmedRequest) InitializeParent(parent *APDU) {
}

func NewAPDUConfirmedRequest(segmentedMessage bool, moreFollows bool, segmentedResponseAccepted bool, maxSegmentsAccepted uint8, maxApduLengthAccepted uint8, invokeId uint8, sequenceNumber *uint8, proposedWindowSize *uint8, serviceRequest *BACnetConfirmedServiceRequest) *APDU {
    child := &APDUConfirmedRequest{
        SegmentedMessage: segmentedMessage,
        MoreFollows: moreFollows,
        SegmentedResponseAccepted: segmentedResponseAccepted,
        MaxSegmentsAccepted: maxSegmentsAccepted,
        MaxApduLengthAccepted: maxApduLengthAccepted,
        InvokeId: invokeId,
        SequenceNumber: sequenceNumber,
        ProposedWindowSize: proposedWindowSize,
        ServiceRequest: serviceRequest,
        Parent: NewAPDU(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastAPDUConfirmedRequest(structType interface{}) *APDUConfirmedRequest {
    castFunc := func(typ interface{}) *APDUConfirmedRequest {
        if casted, ok := typ.(APDUConfirmedRequest); ok {
            return &casted
        }
        if casted, ok := typ.(*APDUConfirmedRequest); ok {
            return casted
        }
        if casted, ok := typ.(APDU); ok {
            return CastAPDUConfirmedRequest(casted.Child)
        }
        if casted, ok := typ.(*APDU); ok {
            return CastAPDUConfirmedRequest(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *APDUConfirmedRequest) GetTypeName() string {
    return "APDUConfirmedRequest"
}

func (m *APDUConfirmedRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (segmentedMessage)
    lengthInBits += 1

    // Simple field (moreFollows)
    lengthInBits += 1

    // Simple field (segmentedResponseAccepted)
    lengthInBits += 1

    // Reserved Field (reserved)
    lengthInBits += 2

    // Simple field (maxSegmentsAccepted)
    lengthInBits += 3

    // Simple field (maxApduLengthAccepted)
    lengthInBits += 4

    // Simple field (invokeId)
    lengthInBits += 8

    // Optional Field (sequenceNumber)
    if m.SequenceNumber != nil {
        lengthInBits += 8
    }

    // Optional Field (proposedWindowSize)
    if m.ProposedWindowSize != nil {
        lengthInBits += 8
    }

    // Simple field (serviceRequest)
    lengthInBits += m.ServiceRequest.LengthInBits()

    return lengthInBits
}

func (m *APDUConfirmedRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func APDUConfirmedRequestParse(io *utils.ReadBuffer, apduLength uint16) (*APDU, error) {

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

    // Simple Field (segmentedResponseAccepted)
    segmentedResponseAccepted, _segmentedResponseAcceptedErr := io.ReadBit()
    if _segmentedResponseAcceptedErr != nil {
        return nil, errors.New("Error parsing 'segmentedResponseAccepted' field " + _segmentedResponseAcceptedErr.Error())
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

    // Simple Field (maxSegmentsAccepted)
    maxSegmentsAccepted, _maxSegmentsAcceptedErr := io.ReadUint8(3)
    if _maxSegmentsAcceptedErr != nil {
        return nil, errors.New("Error parsing 'maxSegmentsAccepted' field " + _maxSegmentsAcceptedErr.Error())
    }

    // Simple Field (maxApduLengthAccepted)
    maxApduLengthAccepted, _maxApduLengthAcceptedErr := io.ReadUint8(4)
    if _maxApduLengthAcceptedErr != nil {
        return nil, errors.New("Error parsing 'maxApduLengthAccepted' field " + _maxApduLengthAcceptedErr.Error())
    }

    // Simple Field (invokeId)
    invokeId, _invokeIdErr := io.ReadUint8(8)
    if _invokeIdErr != nil {
        return nil, errors.New("Error parsing 'invokeId' field " + _invokeIdErr.Error())
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

    // Simple Field (serviceRequest)
    serviceRequest, _serviceRequestErr := BACnetConfirmedServiceRequestParse(io, uint16(apduLength) - uint16(uint16(uint16(uint16(3)) + uint16(uint16(utils.InlineIf(segmentedMessage, uint16(uint16(2)), uint16(uint16(0))))))))
    if _serviceRequestErr != nil {
        return nil, errors.New("Error parsing 'serviceRequest' field " + _serviceRequestErr.Error())
    }

    // Create a partially initialized instance
    _child := &APDUConfirmedRequest{
        SegmentedMessage: segmentedMessage,
        MoreFollows: moreFollows,
        SegmentedResponseAccepted: segmentedResponseAccepted,
        MaxSegmentsAccepted: maxSegmentsAccepted,
        MaxApduLengthAccepted: maxApduLengthAccepted,
        InvokeId: invokeId,
        SequenceNumber: sequenceNumber,
        ProposedWindowSize: proposedWindowSize,
        ServiceRequest: serviceRequest,
        Parent: &APDU{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *APDUConfirmedRequest) Serialize(io utils.WriteBuffer) error {
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

    // Simple Field (segmentedResponseAccepted)
    segmentedResponseAccepted := bool(m.SegmentedResponseAccepted)
    _segmentedResponseAcceptedErr := io.WriteBit((segmentedResponseAccepted))
    if _segmentedResponseAcceptedErr != nil {
        return errors.New("Error serializing 'segmentedResponseAccepted' field " + _segmentedResponseAcceptedErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(2, uint8(0))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (maxSegmentsAccepted)
    maxSegmentsAccepted := uint8(m.MaxSegmentsAccepted)
    _maxSegmentsAcceptedErr := io.WriteUint8(3, (maxSegmentsAccepted))
    if _maxSegmentsAcceptedErr != nil {
        return errors.New("Error serializing 'maxSegmentsAccepted' field " + _maxSegmentsAcceptedErr.Error())
    }

    // Simple Field (maxApduLengthAccepted)
    maxApduLengthAccepted := uint8(m.MaxApduLengthAccepted)
    _maxApduLengthAcceptedErr := io.WriteUint8(4, (maxApduLengthAccepted))
    if _maxApduLengthAcceptedErr != nil {
        return errors.New("Error serializing 'maxApduLengthAccepted' field " + _maxApduLengthAcceptedErr.Error())
    }

    // Simple Field (invokeId)
    invokeId := uint8(m.InvokeId)
    _invokeIdErr := io.WriteUint8(8, (invokeId))
    if _invokeIdErr != nil {
        return errors.New("Error serializing 'invokeId' field " + _invokeIdErr.Error())
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

    // Simple Field (serviceRequest)
    _serviceRequestErr := m.ServiceRequest.Serialize(io)
    if _serviceRequestErr != nil {
        return errors.New("Error serializing 'serviceRequest' field " + _serviceRequestErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *APDUConfirmedRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "segmentedResponseAccepted":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SegmentedResponseAccepted = data
            case "maxSegmentsAccepted":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MaxSegmentsAccepted = data
            case "maxApduLengthAccepted":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MaxApduLengthAccepted = data
            case "invokeId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.InvokeId = data
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
            case "serviceRequest":
                var dt *BACnetConfirmedServiceRequest
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.ServiceRequest = dt
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

func (m *APDUConfirmedRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.SegmentedMessage, xml.StartElement{Name: xml.Name{Local: "segmentedMessage"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MoreFollows, xml.StartElement{Name: xml.Name{Local: "moreFollows"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SegmentedResponseAccepted, xml.StartElement{Name: xml.Name{Local: "segmentedResponseAccepted"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MaxSegmentsAccepted, xml.StartElement{Name: xml.Name{Local: "maxSegmentsAccepted"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MaxApduLengthAccepted, xml.StartElement{Name: xml.Name{Local: "maxApduLengthAccepted"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.InvokeId, xml.StartElement{Name: xml.Name{Local: "invokeId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SequenceNumber, xml.StartElement{Name: xml.Name{Local: "sequenceNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ProposedWindowSize, xml.StartElement{Name: xml.Name{Local: "proposedWindowSize"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ServiceRequest, xml.StartElement{Name: xml.Name{Local: "serviceRequest"}}); err != nil {
        return err
    }
    return nil
}

