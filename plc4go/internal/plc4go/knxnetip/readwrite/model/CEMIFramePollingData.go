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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMIFramePollingData struct {
    Parent *CEMIFrame
    ICEMIFramePollingData
}

// The corresponding interface
type ICEMIFramePollingData interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *CEMIFramePollingData) NotAckFrame() bool {
    return true
}

func (m *CEMIFramePollingData) StandardFrame() bool {
    return true
}

func (m *CEMIFramePollingData) Polling() bool {
    return true
}


func (m *CEMIFramePollingData) InitializeParent(parent *CEMIFrame, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) {
    m.Parent.Repeated = repeated
    m.Parent.Priority = priority
    m.Parent.AcknowledgeRequested = acknowledgeRequested
    m.Parent.ErrorFlag = errorFlag
}

func NewCEMIFramePollingData(repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) *CEMIFrame {
    child := &CEMIFramePollingData{
        Parent: NewCEMIFrame(repeated, priority, acknowledgeRequested, errorFlag),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCEMIFramePollingData(structType interface{}) CEMIFramePollingData {
    castFunc := func(typ interface{}) CEMIFramePollingData {
        if casted, ok := typ.(CEMIFramePollingData); ok {
            return casted
        }
        if casted, ok := typ.(*CEMIFramePollingData); ok {
            return *casted
        }
        if casted, ok := typ.(CEMIFrame); ok {
            return CastCEMIFramePollingData(casted.Child)
        }
        if casted, ok := typ.(*CEMIFrame); ok {
            return CastCEMIFramePollingData(casted.Child)
        }
        return CEMIFramePollingData{}
    }
    return castFunc(structType)
}

func (m *CEMIFramePollingData) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *CEMIFramePollingData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func CEMIFramePollingDataParse(io *utils.ReadBuffer) (*CEMIFrame, error) {

    // Create a partially initialized instance
    _child := &CEMIFramePollingData{
        Parent: &CEMIFrame{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *CEMIFramePollingData) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *CEMIFramePollingData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *CEMIFramePollingData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

