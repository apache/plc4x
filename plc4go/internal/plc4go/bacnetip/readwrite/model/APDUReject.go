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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type APDUReject struct {
    OriginalInvokeId uint8
    RejectReason uint8
    Parent *APDU
    IAPDUReject
}

// The corresponding interface
type IAPDUReject interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *APDUReject) ApduType() uint8 {
    return 0x6
}


func (m *APDUReject) InitializeParent(parent *APDU) {
}

func NewAPDUReject(originalInvokeId uint8, rejectReason uint8, ) *APDU {
    child := &APDUReject{
        OriginalInvokeId: originalInvokeId,
        RejectReason: rejectReason,
        Parent: NewAPDU(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastAPDUReject(structType interface{}) APDUReject {
    castFunc := func(typ interface{}) APDUReject {
        if casted, ok := typ.(APDUReject); ok {
            return casted
        }
        if casted, ok := typ.(*APDUReject); ok {
            return *casted
        }
        if casted, ok := typ.(APDU); ok {
            return CastAPDUReject(casted.Child)
        }
        if casted, ok := typ.(*APDU); ok {
            return CastAPDUReject(casted.Child)
        }
        return APDUReject{}
    }
    return castFunc(structType)
}

func (m *APDUReject) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Reserved Field (reserved)
    lengthInBits += 4

    // Simple field (originalInvokeId)
    lengthInBits += 8

    // Simple field (rejectReason)
    lengthInBits += 8

    return lengthInBits
}

func (m *APDUReject) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func APDURejectParse(io *utils.ReadBuffer) (*APDU, error) {

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

    // Simple Field (originalInvokeId)
    originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
    if _originalInvokeIdErr != nil {
        return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
    }

    // Simple Field (rejectReason)
    rejectReason, _rejectReasonErr := io.ReadUint8(8)
    if _rejectReasonErr != nil {
        return nil, errors.New("Error parsing 'rejectReason' field " + _rejectReasonErr.Error())
    }

    // Create a partially initialized instance
    _child := &APDUReject{
        OriginalInvokeId: originalInvokeId,
        RejectReason: rejectReason,
        Parent: &APDU{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *APDUReject) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(4, uint8(0x00))
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

    // Simple Field (rejectReason)
    rejectReason := uint8(m.RejectReason)
    _rejectReasonErr := io.WriteUint8(8, (rejectReason))
    if _rejectReasonErr != nil {
        return errors.New("Error serializing 'rejectReason' field " + _rejectReasonErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *APDUReject) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "originalInvokeId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.OriginalInvokeId = data
            case "rejectReason":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.RejectReason = data
            }
        }
    }
}

func (m *APDUReject) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.APDUReject"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.OriginalInvokeId, xml.StartElement{Name: xml.Name{Local: "originalInvokeId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.RejectReason, xml.StartElement{Name: xml.Name{Local: "rejectReason"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

