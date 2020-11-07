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
type S7ParameterSetupCommunication struct {
    MaxAmqCaller uint16
    MaxAmqCallee uint16
    PduLength uint16
    Parent *S7Parameter
    IS7ParameterSetupCommunication
}

// The corresponding interface
type IS7ParameterSetupCommunication interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7ParameterSetupCommunication) ParameterType() uint8 {
    return 0xF0
}

func (m *S7ParameterSetupCommunication) MessageType() uint8 {
    return 0
}


func (m *S7ParameterSetupCommunication) InitializeParent(parent *S7Parameter) {
}

func NewS7ParameterSetupCommunication(maxAmqCaller uint16, maxAmqCallee uint16, pduLength uint16, ) *S7Parameter {
    child := &S7ParameterSetupCommunication{
        MaxAmqCaller: maxAmqCaller,
        MaxAmqCallee: maxAmqCallee,
        PduLength: pduLength,
        Parent: NewS7Parameter(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7ParameterSetupCommunication(structType interface{}) S7ParameterSetupCommunication {
    castFunc := func(typ interface{}) S7ParameterSetupCommunication {
        if casted, ok := typ.(S7ParameterSetupCommunication); ok {
            return casted
        }
        if casted, ok := typ.(*S7ParameterSetupCommunication); ok {
            return *casted
        }
        if casted, ok := typ.(S7Parameter); ok {
            return CastS7ParameterSetupCommunication(casted.Child)
        }
        if casted, ok := typ.(*S7Parameter); ok {
            return CastS7ParameterSetupCommunication(casted.Child)
        }
        return S7ParameterSetupCommunication{}
    }
    return castFunc(structType)
}

func (m *S7ParameterSetupCommunication) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Reserved Field (reserved)
    lengthInBits += 8

    // Simple field (maxAmqCaller)
    lengthInBits += 16

    // Simple field (maxAmqCallee)
    lengthInBits += 16

    // Simple field (pduLength)
    lengthInBits += 16

    return lengthInBits
}

func (m *S7ParameterSetupCommunication) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7ParameterSetupCommunicationParse(io *utils.ReadBuffer) (*S7Parameter, error) {

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(8)
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

    // Simple Field (maxAmqCaller)
    maxAmqCaller, _maxAmqCallerErr := io.ReadUint16(16)
    if _maxAmqCallerErr != nil {
        return nil, errors.New("Error parsing 'maxAmqCaller' field " + _maxAmqCallerErr.Error())
    }

    // Simple Field (maxAmqCallee)
    maxAmqCallee, _maxAmqCalleeErr := io.ReadUint16(16)
    if _maxAmqCalleeErr != nil {
        return nil, errors.New("Error parsing 'maxAmqCallee' field " + _maxAmqCalleeErr.Error())
    }

    // Simple Field (pduLength)
    pduLength, _pduLengthErr := io.ReadUint16(16)
    if _pduLengthErr != nil {
        return nil, errors.New("Error parsing 'pduLength' field " + _pduLengthErr.Error())
    }

    // Create a partially initialized instance
    _child := &S7ParameterSetupCommunication{
        MaxAmqCaller: maxAmqCaller,
        MaxAmqCallee: maxAmqCallee,
        PduLength: pduLength,
        Parent: &S7Parameter{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7ParameterSetupCommunication) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(8, uint8(0x00))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (maxAmqCaller)
    maxAmqCaller := uint16(m.MaxAmqCaller)
    _maxAmqCallerErr := io.WriteUint16(16, (maxAmqCaller))
    if _maxAmqCallerErr != nil {
        return errors.New("Error serializing 'maxAmqCaller' field " + _maxAmqCallerErr.Error())
    }

    // Simple Field (maxAmqCallee)
    maxAmqCallee := uint16(m.MaxAmqCallee)
    _maxAmqCalleeErr := io.WriteUint16(16, (maxAmqCallee))
    if _maxAmqCalleeErr != nil {
        return errors.New("Error serializing 'maxAmqCallee' field " + _maxAmqCalleeErr.Error())
    }

    // Simple Field (pduLength)
    pduLength := uint16(m.PduLength)
    _pduLengthErr := io.WriteUint16(16, (pduLength))
    if _pduLengthErr != nil {
        return errors.New("Error serializing 'pduLength' field " + _pduLengthErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7ParameterSetupCommunication) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "maxAmqCaller":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MaxAmqCaller = data
            case "maxAmqCallee":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MaxAmqCallee = data
            case "pduLength":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.PduLength = data
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

func (m *S7ParameterSetupCommunication) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.MaxAmqCaller, xml.StartElement{Name: xml.Name{Local: "maxAmqCaller"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MaxAmqCallee, xml.StartElement{Name: xml.Name{Local: "maxAmqCallee"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.PduLength, xml.StartElement{Name: xml.Name{Local: "pduLength"}}); err != nil {
        return err
    }
    return nil
}

