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
	"io"
	"reflect"
	"strings"
)

// Code generated by build-utils. DO NOT EDIT.

// The data-structure of this message
type S7PayloadUserDataItem struct {
	ReturnCode    DataTransportErrorCode
	TransportSize DataTransportSize
	SzlId         *SzlId
	SzlIndex      uint16
	Child         IS7PayloadUserDataItemChild
}

// The corresponding interface
type IS7PayloadUserDataItem interface {
	CpuFunctionType() uint8
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

type IS7PayloadUserDataItemParent interface {
	SerializeParent(io utils.WriteBuffer, child IS7PayloadUserDataItem, serializeChildFunction func() error) error
	GetTypeName() string
}

type IS7PayloadUserDataItemChild interface {
	Serialize(io utils.WriteBuffer) error
	InitializeParent(parent *S7PayloadUserDataItem, returnCode DataTransportErrorCode, transportSize DataTransportSize, szlId *SzlId, szlIndex uint16)
	GetTypeName() string
	IS7PayloadUserDataItem
}

func NewS7PayloadUserDataItem(returnCode DataTransportErrorCode, transportSize DataTransportSize, szlId *SzlId, szlIndex uint16) *S7PayloadUserDataItem {
	return &S7PayloadUserDataItem{ReturnCode: returnCode, TransportSize: transportSize, SzlId: szlId, SzlIndex: szlIndex}
}

func CastS7PayloadUserDataItem(structType interface{}) *S7PayloadUserDataItem {
	castFunc := func(typ interface{}) *S7PayloadUserDataItem {
		if casted, ok := typ.(S7PayloadUserDataItem); ok {
			return &casted
		}
		if casted, ok := typ.(*S7PayloadUserDataItem); ok {
			return casted
		}
		return nil
	}
	return castFunc(structType)
}

func (m *S7PayloadUserDataItem) GetTypeName() string {
	return "S7PayloadUserDataItem"
}

func (m *S7PayloadUserDataItem) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Enum Field (returnCode)
	lengthInBits += 8

	// Enum Field (transportSize)
	lengthInBits += 8

	// Implicit Field (dataLength)
	lengthInBits += 16

	// Simple field (szlId)
	lengthInBits += m.SzlId.LengthInBits()

	// Simple field (szlIndex)
	lengthInBits += 16

	// Length of sub-type elements will be added by sub-type...
	lengthInBits += m.Child.LengthInBits()

	return lengthInBits
}

func (m *S7PayloadUserDataItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7PayloadUserDataItemParse(io *utils.ReadBuffer, cpuFunctionType uint8) (*S7PayloadUserDataItem, error) {

	// Enum field (returnCode)
	returnCode, _returnCodeErr := DataTransportErrorCodeParse(io)
	if _returnCodeErr != nil {
		return nil, errors.New("Error parsing 'returnCode' field " + _returnCodeErr.Error())
	}

	// Enum field (transportSize)
	transportSize, _transportSizeErr := DataTransportSizeParse(io)
	if _transportSizeErr != nil {
		return nil, errors.New("Error parsing 'transportSize' field " + _transportSizeErr.Error())
	}

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _dataLengthErr := io.ReadUint16(16)
	if _dataLengthErr != nil {
		return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
	}

	// Simple Field (szlId)
	szlId, _szlIdErr := SzlIdParse(io)
	if _szlIdErr != nil {
		return nil, errors.New("Error parsing 'szlId' field " + _szlIdErr.Error())
	}

	// Simple Field (szlIndex)
	szlIndex, _szlIndexErr := io.ReadUint16(16)
	if _szlIndexErr != nil {
		return nil, errors.New("Error parsing 'szlIndex' field " + _szlIndexErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _parent *S7PayloadUserDataItem
	var typeSwitchError error
	switch {
	case cpuFunctionType == 0x04:
		_parent, typeSwitchError = S7PayloadUserDataItemCpuFunctionReadSzlRequestParse(io)
	case cpuFunctionType == 0x08:
		_parent, typeSwitchError = S7PayloadUserDataItemCpuFunctionReadSzlResponseParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Finish initializing
	_parent.Child.InitializeParent(_parent, returnCode, transportSize, szlId, szlIndex)
	return _parent, nil
}

func (m *S7PayloadUserDataItem) Serialize(io utils.WriteBuffer) error {
	return m.Child.Serialize(io)
}

func (m *S7PayloadUserDataItem) SerializeParent(io utils.WriteBuffer, child IS7PayloadUserDataItem, serializeChildFunction func() error) error {

	// Enum field (returnCode)
	returnCode := CastDataTransportErrorCode(m.ReturnCode)
	_returnCodeErr := returnCode.Serialize(io)
	if _returnCodeErr != nil {
		return errors.New("Error serializing 'returnCode' field " + _returnCodeErr.Error())
	}

	// Enum field (transportSize)
	transportSize := CastDataTransportSize(m.TransportSize)
	_transportSizeErr := transportSize.Serialize(io)
	if _transportSizeErr != nil {
		return errors.New("Error serializing 'transportSize' field " + _transportSizeErr.Error())
	}

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	dataLength := uint16(uint16(uint16(m.LengthInBytes())) - uint16(uint16(4)))
	_dataLengthErr := io.WriteUint16(16, (dataLength))
	if _dataLengthErr != nil {
		return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
	}

	// Simple Field (szlId)
	_szlIdErr := m.SzlId.Serialize(io)
	if _szlIdErr != nil {
		return errors.New("Error serializing 'szlId' field " + _szlIdErr.Error())
	}

	// Simple Field (szlIndex)
	szlIndex := uint16(m.SzlIndex)
	_szlIndexErr := io.WriteUint16(16, (szlIndex))
	if _szlIndexErr != nil {
		return errors.New("Error serializing 'szlIndex' field " + _szlIndexErr.Error())
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	_typeSwitchErr := serializeChildFunction()
	if _typeSwitchErr != nil {
		return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
	}

	return nil
}

func (m *S7PayloadUserDataItem) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
			case "returnCode":
				var data DataTransportErrorCode
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ReturnCode = data
			case "transportSize":
				var data DataTransportSize
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.TransportSize = data
			case "szlId":
				var data *SzlId
				if err := d.DecodeElement(data, &tok); err != nil {
					return err
				}
				m.SzlId = data
			case "szlIndex":
				var data uint16
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.SzlIndex = data
			default:
				switch start.Attr[0].Value {
				case "org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlRequest":
					var dt *S7PayloadUserDataItemCpuFunctionReadSzlRequest
					if m.Child != nil {
						dt = m.Child.(*S7PayloadUserDataItemCpuFunctionReadSzlRequest)
					}
					if err := d.DecodeElement(&dt, &tok); err != nil {
						return err
					}
					if m.Child == nil {
						dt.Parent = m
						m.Child = dt
					}
				case "org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlResponse":
					var dt *S7PayloadUserDataItemCpuFunctionReadSzlResponse
					if m.Child != nil {
						dt = m.Child.(*S7PayloadUserDataItemCpuFunctionReadSzlResponse)
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

func (m *S7PayloadUserDataItem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	className := reflect.TypeOf(m.Child).String()
	className = "org.apache.plc4x.java.s7.readwrite." + className[strings.LastIndex(className, ".")+1:]
	if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
		{Name: xml.Name{Local: "className"}, Value: className},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ReturnCode, xml.StartElement{Name: xml.Name{Local: "returnCode"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.TransportSize, xml.StartElement{Name: xml.Name{Local: "transportSize"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.SzlId, xml.StartElement{Name: xml.Name{Local: "szlId"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.SzlIndex, xml.StartElement{Name: xml.Name{Local: "szlIndex"}}); err != nil {
		return err
	}
	marshaller, ok := m.Child.(xml.Marshaler)
	if !ok {
		return errors.New("child is not castable to Marshaler")
	}
	if err := marshaller.MarshalXML(e, start); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}
