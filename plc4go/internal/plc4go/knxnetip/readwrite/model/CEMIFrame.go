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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    "reflect"
    "strings"
)

// The data-structure of this message
type CEMIFrame struct {
    Repeated bool
    Priority CEMIPriority
    AcknowledgeRequested bool
    ErrorFlag bool
    Child ICEMIFrameChild
    ICEMIFrame
    ICEMIFrameParent
}

// The corresponding interface
type ICEMIFrame interface {
    NotAckFrame() bool
    Polling() bool
    StandardFrame() bool
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type ICEMIFrameParent interface {
    SerializeParent(io utils.WriteBuffer, child ICEMIFrame, serializeChildFunction func() error) error
    GetTypeName() string
}

type ICEMIFrameChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *CEMIFrame, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool)
    GetTypeName() string
    ICEMIFrame
}

func NewCEMIFrame(repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) *CEMIFrame {
    return &CEMIFrame{Repeated: repeated, Priority: priority, AcknowledgeRequested: acknowledgeRequested, ErrorFlag: errorFlag}
}

func CastCEMIFrame(structType interface{}) *CEMIFrame {
    castFunc := func(typ interface{}) *CEMIFrame {
        if casted, ok := typ.(CEMIFrame); ok {
            return &casted
        }
        if casted, ok := typ.(*CEMIFrame); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *CEMIFrame) GetTypeName() string {
    return "CEMIFrame"
}

func (m *CEMIFrame) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (standardFrame)
    lengthInBits += 1

    // Discriminator Field (polling)
    lengthInBits += 1

    // Simple field (repeated)
    lengthInBits += 1

    // Discriminator Field (notAckFrame)
    lengthInBits += 1

    // Enum Field (priority)
    lengthInBits += 2

    // Simple field (acknowledgeRequested)
    lengthInBits += 1

    // Simple field (errorFlag)
    lengthInBits += 1

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *CEMIFrame) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIFrameParse(io *utils.ReadBuffer) (*CEMIFrame, error) {

    // Discriminator Field (standardFrame) (Used as input to a switch field)
    standardFrame, _standardFrameErr := io.ReadBit()
    if _standardFrameErr != nil {
        return nil, errors.New("Error parsing 'standardFrame' field " + _standardFrameErr.Error())
    }

    // Discriminator Field (polling) (Used as input to a switch field)
    polling, _pollingErr := io.ReadBit()
    if _pollingErr != nil {
        return nil, errors.New("Error parsing 'polling' field " + _pollingErr.Error())
    }

    // Simple Field (repeated)
    repeated, _repeatedErr := io.ReadBit()
    if _repeatedErr != nil {
        return nil, errors.New("Error parsing 'repeated' field " + _repeatedErr.Error())
    }

    // Discriminator Field (notAckFrame) (Used as input to a switch field)
    notAckFrame, _notAckFrameErr := io.ReadBit()
    if _notAckFrameErr != nil {
        return nil, errors.New("Error parsing 'notAckFrame' field " + _notAckFrameErr.Error())
    }

    // Enum field (priority)
    priority, _priorityErr := CEMIPriorityParse(io)
    if _priorityErr != nil {
        return nil, errors.New("Error parsing 'priority' field " + _priorityErr.Error())
    }

    // Simple Field (acknowledgeRequested)
    acknowledgeRequested, _acknowledgeRequestedErr := io.ReadBit()
    if _acknowledgeRequestedErr != nil {
        return nil, errors.New("Error parsing 'acknowledgeRequested' field " + _acknowledgeRequestedErr.Error())
    }

    // Simple Field (errorFlag)
    errorFlag, _errorFlagErr := io.ReadBit()
    if _errorFlagErr != nil {
        return nil, errors.New("Error parsing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *CEMIFrame
    var typeSwitchError error
    switch {
    case notAckFrame == false:
        _parent, typeSwitchError = CEMIFrameAckParse(io)
    case notAckFrame == true && standardFrame == true && polling == false:
        _parent, typeSwitchError = CEMIFrameDataParse(io)
    case notAckFrame == true && standardFrame == true && polling == true:
        _parent, typeSwitchError = CEMIFramePollingDataParse(io)
    case notAckFrame == true && standardFrame == false && polling == false:
        _parent, typeSwitchError = CEMIFrameDataExtParse(io)
    case notAckFrame == true && standardFrame == false && polling == true:
        _parent, typeSwitchError = CEMIFramePollingDataExtParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent, repeated, priority, acknowledgeRequested, errorFlag)
    return _parent, nil
}

func (m *CEMIFrame) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *CEMIFrame) SerializeParent(io utils.WriteBuffer, child ICEMIFrame, serializeChildFunction func() error) error {

    // Discriminator Field (standardFrame) (Used as input to a switch field)
    standardFrame := bool(child.StandardFrame())
    _standardFrameErr := io.WriteBit((standardFrame))
    if _standardFrameErr != nil {
        return errors.New("Error serializing 'standardFrame' field " + _standardFrameErr.Error())
    }

    // Discriminator Field (polling) (Used as input to a switch field)
    polling := bool(child.Polling())
    _pollingErr := io.WriteBit((polling))
    if _pollingErr != nil {
        return errors.New("Error serializing 'polling' field " + _pollingErr.Error())
    }

    // Simple Field (repeated)
    repeated := bool(m.Repeated)
    _repeatedErr := io.WriteBit((repeated))
    if _repeatedErr != nil {
        return errors.New("Error serializing 'repeated' field " + _repeatedErr.Error())
    }

    // Discriminator Field (notAckFrame) (Used as input to a switch field)
    notAckFrame := bool(child.NotAckFrame())
    _notAckFrameErr := io.WriteBit((notAckFrame))
    if _notAckFrameErr != nil {
        return errors.New("Error serializing 'notAckFrame' field " + _notAckFrameErr.Error())
    }

    // Enum field (priority)
    priority := CastCEMIPriority(m.Priority)
    _priorityErr := priority.Serialize(io)
    if _priorityErr != nil {
        return errors.New("Error serializing 'priority' field " + _priorityErr.Error())
    }

    // Simple Field (acknowledgeRequested)
    acknowledgeRequested := bool(m.AcknowledgeRequested)
    _acknowledgeRequestedErr := io.WriteBit((acknowledgeRequested))
    if _acknowledgeRequestedErr != nil {
        return errors.New("Error serializing 'acknowledgeRequested' field " + _acknowledgeRequestedErr.Error())
    }

    // Simple Field (errorFlag)
    errorFlag := bool(m.ErrorFlag)
    _errorFlagErr := io.WriteBit((errorFlag))
    if _errorFlagErr != nil {
        return errors.New("Error serializing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *CEMIFrame) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "repeated":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Repeated = data
            case "priority":
                var data CEMIPriority
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Priority = data
            case "acknowledgeRequested":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.AcknowledgeRequested = data
            case "errorFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ErrorFlag = data
                default:
                    switch start.Attr[0].Value {
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIFrameAck":
                            var dt *CEMIFrameAck
                            if m.Child != nil {
                                dt = m.Child.(*CEMIFrameAck)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIFrameData":
                            var dt *CEMIFrameData
                            if m.Child != nil {
                                dt = m.Child.(*CEMIFrameData)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIFramePollingData":
                            var dt *CEMIFramePollingData
                            if m.Child != nil {
                                dt = m.Child.(*CEMIFramePollingData)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIFrameDataExt":
                            var dt *CEMIFrameDataExt
                            if m.Child != nil {
                                dt = m.Child.(*CEMIFrameDataExt)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.CEMIFramePollingDataExt":
                            var dt *CEMIFramePollingDataExt
                            if m.Child != nil {
                                dt = m.Child.(*CEMIFramePollingDataExt)
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

func (m *CEMIFrame) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.knxnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Repeated, xml.StartElement{Name: xml.Name{Local: "repeated"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Priority, xml.StartElement{Name: xml.Name{Local: "priority"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.AcknowledgeRequested, xml.StartElement{Name: xml.Name{Local: "acknowledgeRequested"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ErrorFlag, xml.StartElement{Name: xml.Name{Local: "errorFlag"}}); err != nil {
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

