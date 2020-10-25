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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "reflect"
    "strconv"
)

// Constant values.
const S7Message_PROTOCOLID uint8 = 0x32

// The data-structure of this message
type S7Message struct {
    TpduReference uint16
    Parameter *IS7Parameter
    Payload *IS7Payload

}

// The corresponding interface
type IS7Message interface {
    spi.Message
    MessageType() uint8
    Serialize(io utils.WriteBuffer) error
}

type S7MessageInitializer interface {
    initialize(tpduReference uint16, parameter *IS7Parameter, payload *IS7Payload) spi.Message
}

func S7MessageMessageType(m IS7Message) uint8 {
    return m.MessageType()
}


func CastIS7Message(structType interface{}) IS7Message {
    castFunc := func(typ interface{}) IS7Message {
        if iS7Message, ok := typ.(IS7Message); ok {
            return iS7Message
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7Message(structType interface{}) S7Message {
    castFunc := func(typ interface{}) S7Message {
        if sS7Message, ok := typ.(S7Message); ok {
            return sS7Message
        }
        if sS7Message, ok := typ.(*S7Message); ok {
            return *sS7Message
        }
        return S7Message{}
    }
    return castFunc(structType)
}

func (m S7Message) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Const Field (protocolId)
    lengthInBits += 8

    // Discriminator Field (messageType)
    lengthInBits += 8

    // Reserved Field (reserved)
    lengthInBits += 16

    // Simple field (tpduReference)
    lengthInBits += 16

    // Implicit Field (parameterLength)
    lengthInBits += 16

    // Implicit Field (payloadLength)
    lengthInBits += 16

    // Length of sub-type elements will be added by sub-type...

    // Optional Field (parameter)
    if m.Parameter != nil {
        lengthInBits += (*m.Parameter).LengthInBits()
    }

    // Optional Field (payload)
    if m.Payload != nil {
        lengthInBits += (*m.Payload).LengthInBits()
    }

    return lengthInBits
}

func (m S7Message) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7MessageParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Const Field (protocolId)
    protocolId, _protocolIdErr := io.ReadUint8(8)
    if _protocolIdErr != nil {
        return nil, errors.New("Error parsing 'protocolId' field " + _protocolIdErr.Error())
    }
    if protocolId != S7Message_PROTOCOLID {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(S7Message_PROTOCOLID)) + " but got " + strconv.Itoa(int(protocolId)))
    }

    // Discriminator Field (messageType) (Used as input to a switch field)
    messageType, _messageTypeErr := io.ReadUint8(8)
    if _messageTypeErr != nil {
        return nil, errors.New("Error parsing 'messageType' field " + _messageTypeErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint16(16)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint16(0x0000) {
            log.WithFields(log.Fields{
                "expected value": uint16(0x0000),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (tpduReference)
    tpduReference, _tpduReferenceErr := io.ReadUint16(16)
    if _tpduReferenceErr != nil {
        return nil, errors.New("Error parsing 'tpduReference' field " + _tpduReferenceErr.Error())
    }

    // Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    parameterLength, _parameterLengthErr := io.ReadUint16(16)
    if _parameterLengthErr != nil {
        return nil, errors.New("Error parsing 'parameterLength' field " + _parameterLengthErr.Error())
    }

    // Implicit Field (payloadLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    payloadLength, _payloadLengthErr := io.ReadUint16(16)
    if _payloadLengthErr != nil {
        return nil, errors.New("Error parsing 'payloadLength' field " + _payloadLengthErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer S7MessageInitializer
    var typeSwitchError error
    switch {
    case messageType == 0x01:
        initializer, typeSwitchError = S7MessageRequestParse(io)
    case messageType == 0x02:
        initializer, typeSwitchError = S7MessageResponseParse(io)
    case messageType == 0x03:
        initializer, typeSwitchError = S7MessageResponseDataParse(io)
    case messageType == 0x07:
        initializer, typeSwitchError = S7MessageUserDataParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Optional Field (parameter) (Can be skipped, if a given expression evaluates to false)
    var parameter *IS7Parameter = nil
    if bool((parameterLength) > ((0))) {
        _message, _err := S7ParameterParse(io, messageType)
        if _err != nil {
            return nil, errors.New("Error parsing 'parameter' field " + _err.Error())
        }
        var _item IS7Parameter
        _item, _ok := _message.(IS7Parameter)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IS7Parameter")
        }
        parameter = &_item
    }

    // Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
    var payload *IS7Payload = nil
    if bool((payloadLength) > ((0))) {
        _message, _err := S7PayloadParse(io, messageType, (*parameter))
        if _err != nil {
            return nil, errors.New("Error parsing 'payload' field " + _err.Error())
        }
        var _item IS7Payload
        _item, _ok := _message.(IS7Payload)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IS7Payload")
        }
        payload = &_item
    }

    // Create the instance
    return initializer.initialize(tpduReference, parameter, payload), nil
}

func S7MessageSerialize(io utils.WriteBuffer, m S7Message, i IS7Message, childSerialize func() error) error {

    // Const Field (protocolId)
    _protocolIdErr := io.WriteUint8(8, 0x32)
    if _protocolIdErr != nil {
        return errors.New("Error serializing 'protocolId' field " + _protocolIdErr.Error())
    }

    // Discriminator Field (messageType) (Used as input to a switch field)
    messageType := uint8(i.MessageType())
    _messageTypeErr := io.WriteUint8(8, (messageType))
    if _messageTypeErr != nil {
        return errors.New("Error serializing 'messageType' field " + _messageTypeErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint16(16, uint16(0x0000))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (tpduReference)
    tpduReference := uint16(m.TpduReference)
    _tpduReferenceErr := io.WriteUint16(16, (tpduReference))
    if _tpduReferenceErr != nil {
        return errors.New("Error serializing 'tpduReference' field " + _tpduReferenceErr.Error())
    }

    // Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    parameterLength := uint16(utils.InlineIf(bool((m.Parameter) != (nil)), uint16((*m.Parameter).LengthInBytes()), uint16(uint16(0))))
    _parameterLengthErr := io.WriteUint16(16, (parameterLength))
    if _parameterLengthErr != nil {
        return errors.New("Error serializing 'parameterLength' field " + _parameterLengthErr.Error())
    }

    // Implicit Field (payloadLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    payloadLength := uint16(utils.InlineIf(bool((m.Payload) != (nil)), uint16((*m.Payload).LengthInBytes()), uint16(uint16(0))))
    _payloadLengthErr := io.WriteUint16(16, (payloadLength))
    if _payloadLengthErr != nil {
        return errors.New("Error serializing 'payloadLength' field " + _payloadLengthErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    // Optional Field (parameter) (Can be skipped, if the value is null)
    var parameter *IS7Parameter = nil
    if m.Parameter != nil {
        parameter = m.Parameter
        _parameterErr := CastIS7Parameter(*parameter).Serialize(io)
        if _parameterErr != nil {
            return errors.New("Error serializing 'parameter' field " + _parameterErr.Error())
        }
    }

    // Optional Field (payload) (Can be skipped, if the value is null)
    var payload *IS7Payload = nil
    if m.Payload != nil {
        payload = m.Payload
        _payloadErr := CastIS7Payload(*payload).Serialize(io)
        if _payloadErr != nil {
            return errors.New("Error serializing 'payload' field " + _payloadErr.Error())
        }
    }

    return nil
}

func (m *S7Message) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "tpduReference":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TpduReference = data
            case "parameter":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.s7.readwrite.S7ParameterSetupCommunication":
                        var dt *S7ParameterSetupCommunication
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Parameter = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarRequest":
                        var dt *S7ParameterReadVarRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Parameter = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarResponse":
                        var dt *S7ParameterReadVarResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Parameter = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarRequest":
                        var dt *S7ParameterWriteVarRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Parameter = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarResponse":
                        var dt *S7ParameterWriteVarResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Parameter = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7ParameterUserData":
                        var dt *S7ParameterUserData
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Parameter = dt
                    }
            case "payload":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.s7.readwrite.S7PayloadReadVarResponse":
                        var dt *S7PayloadReadVarResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Payload = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarRequest":
                        var dt *S7PayloadWriteVarRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Payload = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarResponse":
                        var dt *S7PayloadWriteVarResponse
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Payload = dt
                    case "org.apache.plc4x.java.s7.readwrite.S7PayloadUserData":
                        var dt *S7PayloadUserData
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Payload = dt
                    }
            }
        }
    }
}

func (m S7Message) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7Message"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.TpduReference, xml.StartElement{Name: xml.Name{Local: "tpduReference"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Parameter, xml.StartElement{Name: xml.Name{Local: "parameter"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Payload, xml.StartElement{Name: xml.Name{Local: "payload"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

