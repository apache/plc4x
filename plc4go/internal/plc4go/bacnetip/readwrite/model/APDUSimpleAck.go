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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

// The data-structure of this message
type APDUSimpleAck struct {
    OriginalInvokeId uint8
    ServiceChoice uint8
    Parent *APDU
    IAPDUSimpleAck
}

// The corresponding interface
type IAPDUSimpleAck interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *APDUSimpleAck) ApduType() uint8 {
    return 0x2
}


func (m *APDUSimpleAck) InitializeParent(parent *APDU) {
}

func NewAPDUSimpleAck(originalInvokeId uint8, serviceChoice uint8, ) *APDU {
    child := &APDUSimpleAck{
        OriginalInvokeId: originalInvokeId,
        ServiceChoice: serviceChoice,
        Parent: NewAPDU(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastAPDUSimpleAck(structType interface{}) *APDUSimpleAck {
    castFunc := func(typ interface{}) *APDUSimpleAck {
        if casted, ok := typ.(APDUSimpleAck); ok {
            return &casted
        }
        if casted, ok := typ.(*APDUSimpleAck); ok {
            return casted
        }
        if casted, ok := typ.(APDU); ok {
            return CastAPDUSimpleAck(casted.Child)
        }
        if casted, ok := typ.(*APDU); ok {
            return CastAPDUSimpleAck(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *APDUSimpleAck) GetTypeName() string {
    return "APDUSimpleAck"
}

func (m *APDUSimpleAck) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Reserved Field (reserved)
    lengthInBits += 4

    // Simple field (originalInvokeId)
    lengthInBits += 8

    // Simple field (serviceChoice)
    lengthInBits += 8

    return lengthInBits
}

func (m *APDUSimpleAck) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func APDUSimpleAckParse(io *utils.ReadBuffer) (*APDU, error) {

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(4)
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

    // Simple Field (serviceChoice)
    serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
    if _serviceChoiceErr != nil {
        return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Create a partially initialized instance
    _child := &APDUSimpleAck{
        OriginalInvokeId: originalInvokeId,
        ServiceChoice: serviceChoice,
        Parent: &APDU{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *APDUSimpleAck) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(4, uint8(0))
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

    // Simple Field (serviceChoice)
    serviceChoice := uint8(m.ServiceChoice)
    _serviceChoiceErr := io.WriteUint8(8, (serviceChoice))
    if _serviceChoiceErr != nil {
        return errors.New("Error serializing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *APDUSimpleAck) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
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
            case "serviceChoice":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ServiceChoice = data
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

func (m *APDUSimpleAck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.OriginalInvokeId, xml.StartElement{Name: xml.Name{Local: "originalInvokeId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ServiceChoice, xml.StartElement{Name: xml.Name{Local: "serviceChoice"}}); err != nil {
        return err
    }
    return nil
}

