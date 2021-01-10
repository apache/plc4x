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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "reflect"
    "strings"
)

// The data-structure of this message
type CEMI struct {
    Child ICEMIChild
    ICEMI
    ICEMIParent
}

// The corresponding interface
type ICEMI interface {
    MessageCode() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type ICEMIParent interface {
    SerializeParent(io utils.WriteBuffer, child ICEMI, serializeChildFunction func() error) error
    GetTypeName() string
}

type ICEMIChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *CEMI)
    GetTypeName() string
    ICEMI
}

func NewCEMI() *CEMI {
    return &CEMI{}
}

func CastCEMI(structType interface{}) *CEMI {
    castFunc := func(typ interface{}) *CEMI {
        if casted, ok := typ.(CEMI); ok {
            return &casted
        }
        if casted, ok := typ.(*CEMI); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *CEMI) GetTypeName() string {
    return "CEMI"
}

func (m *CEMI) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (messageCode)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *CEMI) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIParse(io *utils.ReadBuffer, size uint8) (*CEMI, error) {

    // Discriminator Field (messageCode) (Used as input to a switch field)
    messageCode, _messageCodeErr := io.ReadUint8(8)
    if _messageCodeErr != nil {
        return nil, errors.New("Error parsing 'messageCode' field " + _messageCodeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *CEMI
    var typeSwitchError error
    switch {
    case messageCode == 0x2B:
        _parent, typeSwitchError = LBusmonIndParse(io)
    case messageCode == 0x11:
        _parent, typeSwitchError = LDataReqParse(io)
    case messageCode == 0x29:
        _parent, typeSwitchError = LDataIndParse(io)
    case messageCode == 0x2E:
        _parent, typeSwitchError = LDataConParse(io)
    case messageCode == 0x10:
        _parent, typeSwitchError = LRawReqParse(io)
    case messageCode == 0x2D:
        _parent, typeSwitchError = LRawIndParse(io)
    case messageCode == 0x2F:
        _parent, typeSwitchError = LRawConParse(io)
    case messageCode == 0x13:
        _parent, typeSwitchError = LPollDataReqParse(io)
    case messageCode == 0x25:
        _parent, typeSwitchError = LPollDataConParse(io)
    case messageCode == 0x41:
        _parent, typeSwitchError = TDataConnectedReqParse(io)
    case messageCode == 0x89:
        _parent, typeSwitchError = TDataConnectedIndParse(io)
    case messageCode == 0x4A:
        _parent, typeSwitchError = TDataIndividualReqParse(io)
    case messageCode == 0x94:
        _parent, typeSwitchError = TDataIndividualIndParse(io)
    case messageCode == 0xFC:
        _parent, typeSwitchError = MPropReadReqParse(io)
    case messageCode == 0xFB:
        _parent, typeSwitchError = MPropReadConParse(io)
    case messageCode == 0xF6:
        _parent, typeSwitchError = MPropWriteReqParse(io)
    case messageCode == 0xF5:
        _parent, typeSwitchError = MPropWriteConParse(io)
    case messageCode == 0xF7:
        _parent, typeSwitchError = MPropInfoIndParse(io)
    case messageCode == 0xF8:
        _parent, typeSwitchError = MFuncPropCommandReqParse(io)
    case messageCode == 0xF9:
        _parent, typeSwitchError = MFuncPropStateReadReqParse(io)
    case messageCode == 0xFA:
        _parent, typeSwitchError = MFuncPropConParse(io)
    case messageCode == 0xF1:
        _parent, typeSwitchError = MResetReqParse(io)
    case messageCode == 0xF0:
        _parent, typeSwitchError = MResetIndParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *CEMI) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *CEMI) SerializeParent(io utils.WriteBuffer, child ICEMI, serializeChildFunction func() error) error {

    // Discriminator Field (messageCode) (Used as input to a switch field)
    messageCode := uint8(child.MessageCode())
    _messageCodeErr := io.WriteUint8(8, (messageCode))
    if _messageCodeErr != nil {
        return errors.New("Error serializing 'messageCode' field " + _messageCodeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *CEMI) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
            default:
                switch start.Attr[0].Value {
                    case "org.apache.plc4x.java.knxnetip.readwrite.LBusmonInd":
                        var dt *LBusmonInd
                        if m.Child != nil {
                            dt = m.Child.(*LBusmonInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LDataReq":
                        var dt *LDataReq
                        if m.Child != nil {
                            dt = m.Child.(*LDataReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LDataInd":
                        var dt *LDataInd
                        if m.Child != nil {
                            dt = m.Child.(*LDataInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LDataCon":
                        var dt *LDataCon
                        if m.Child != nil {
                            dt = m.Child.(*LDataCon)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LRawReq":
                        var dt *LRawReq
                        if m.Child != nil {
                            dt = m.Child.(*LRawReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LRawInd":
                        var dt *LRawInd
                        if m.Child != nil {
                            dt = m.Child.(*LRawInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LRawCon":
                        var dt *LRawCon
                        if m.Child != nil {
                            dt = m.Child.(*LRawCon)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LPollDataReq":
                        var dt *LPollDataReq
                        if m.Child != nil {
                            dt = m.Child.(*LPollDataReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.LPollDataCon":
                        var dt *LPollDataCon
                        if m.Child != nil {
                            dt = m.Child.(*LPollDataCon)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.TDataConnectedReq":
                        var dt *TDataConnectedReq
                        if m.Child != nil {
                            dt = m.Child.(*TDataConnectedReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.TDataConnectedInd":
                        var dt *TDataConnectedInd
                        if m.Child != nil {
                            dt = m.Child.(*TDataConnectedInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.TDataIndividualReq":
                        var dt *TDataIndividualReq
                        if m.Child != nil {
                            dt = m.Child.(*TDataIndividualReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.TDataIndividualInd":
                        var dt *TDataIndividualInd
                        if m.Child != nil {
                            dt = m.Child.(*TDataIndividualInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MPropReadReq":
                        var dt *MPropReadReq
                        if m.Child != nil {
                            dt = m.Child.(*MPropReadReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MPropReadCon":
                        var dt *MPropReadCon
                        if m.Child != nil {
                            dt = m.Child.(*MPropReadCon)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MPropWriteReq":
                        var dt *MPropWriteReq
                        if m.Child != nil {
                            dt = m.Child.(*MPropWriteReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MPropWriteCon":
                        var dt *MPropWriteCon
                        if m.Child != nil {
                            dt = m.Child.(*MPropWriteCon)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MPropInfoInd":
                        var dt *MPropInfoInd
                        if m.Child != nil {
                            dt = m.Child.(*MPropInfoInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MFuncPropCommandReq":
                        var dt *MFuncPropCommandReq
                        if m.Child != nil {
                            dt = m.Child.(*MFuncPropCommandReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MFuncPropStateReadReq":
                        var dt *MFuncPropStateReadReq
                        if m.Child != nil {
                            dt = m.Child.(*MFuncPropStateReadReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MFuncPropCon":
                        var dt *MFuncPropCon
                        if m.Child != nil {
                            dt = m.Child.(*MFuncPropCon)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MResetReq":
                        var dt *MResetReq
                        if m.Child != nil {
                            dt = m.Child.(*MResetReq)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.MResetInd":
                        var dt *MResetInd
                        if m.Child != nil {
                            dt = m.Child.(*MResetInd)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                }
            }
        }
    }
}

func (m *CEMI) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.knxnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

