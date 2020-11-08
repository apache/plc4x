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
    "io"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMIFrameAck struct {
    Parent *CEMIFrame
    ICEMIFrameAck
}

// The corresponding interface
type ICEMIFrameAck interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *CEMIFrameAck) NotAckFrame() bool {
    return false
}

func (m *CEMIFrameAck) StandardFrame() bool {
    return false
}

func (m *CEMIFrameAck) Polling() bool {
    return false
}


func (m *CEMIFrameAck) InitializeParent(parent *CEMIFrame, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) {
    m.Parent.Repeated = repeated
    m.Parent.Priority = priority
    m.Parent.AcknowledgeRequested = acknowledgeRequested
    m.Parent.ErrorFlag = errorFlag
}

func NewCEMIFrameAck(repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) *CEMIFrame {
    child := &CEMIFrameAck{
        Parent: NewCEMIFrame(repeated, priority, acknowledgeRequested, errorFlag),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCEMIFrameAck(structType interface{}) CEMIFrameAck {
    castFunc := func(typ interface{}) CEMIFrameAck {
        if casted, ok := typ.(CEMIFrameAck); ok {
            return casted
        }
        if casted, ok := typ.(*CEMIFrameAck); ok {
            return *casted
        }
        if casted, ok := typ.(CEMIFrame); ok {
            return CastCEMIFrameAck(casted.Child)
        }
        if casted, ok := typ.(*CEMIFrame); ok {
            return CastCEMIFrameAck(casted.Child)
        }
        return CEMIFrameAck{}
    }
    return castFunc(structType)
}

func (m *CEMIFrameAck) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *CEMIFrameAck) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIFrameAckParse(io *utils.ReadBuffer) (*CEMIFrame, error) {

    // Create a partially initialized instance
    _child := &CEMIFrameAck{
        Parent: &CEMIFrame{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *CEMIFrameAck) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *CEMIFrameAck) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
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

func (m *CEMIFrameAck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

