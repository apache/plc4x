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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"io"
)

// The data-structure of this message
type LDataFrameAck struct {
	Parent *LDataFrame
	ILDataFrameAck
}

// The corresponding interface
type ILDataFrameAck interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *LDataFrameAck) NotAckFrame() bool {
	return false
}

func (m *LDataFrameAck) ExtendedFrame() bool {
	return false
}

func (m *LDataFrameAck) Polling() bool {
	return false
}

func (m *LDataFrameAck) InitializeParent(parent *LDataFrame, repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) {
	m.Parent.Repeated = repeated
	m.Parent.Priority = priority
	m.Parent.AcknowledgeRequested = acknowledgeRequested
	m.Parent.ErrorFlag = errorFlag
}

func NewLDataFrameAck(repeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool) *LDataFrame {
	child := &LDataFrameAck{
		Parent: NewLDataFrame(repeated, priority, acknowledgeRequested, errorFlag),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastLDataFrameAck(structType interface{}) *LDataFrameAck {
	castFunc := func(typ interface{}) *LDataFrameAck {
		if casted, ok := typ.(LDataFrameAck); ok {
			return &casted
		}
		if casted, ok := typ.(*LDataFrameAck); ok {
			return casted
		}
		if casted, ok := typ.(LDataFrame); ok {
			return CastLDataFrameAck(casted.Child)
		}
		if casted, ok := typ.(*LDataFrame); ok {
			return CastLDataFrameAck(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *LDataFrameAck) GetTypeName() string {
	return "LDataFrameAck"
}

func (m *LDataFrameAck) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	return lengthInBits
}

func (m *LDataFrameAck) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func LDataFrameAckParse(io *utils.ReadBuffer) (*LDataFrame, error) {

	// Create a partially initialized instance
	_child := &LDataFrameAck{
		Parent: &LDataFrame{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *LDataFrameAck) Serialize(io utils.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return m.Parent.SerializeParent(io, m, ser)
}

func (m *LDataFrameAck) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *LDataFrameAck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}
