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
	log "github.com/sirupsen/logrus"
	"io"
)

// The data-structure of this message
type GroupObjectDescriptorRealisationType1 struct {
	DataPointer           uint8
	TransmitEnable        bool
	SegmentSelectorEnable bool
	WriteEnable           bool
	ReadEnable            bool
	CommunicationEnable   bool
	Priority              CEMIPriority
	ValueType             ComObjectValueType
	IGroupObjectDescriptorRealisationType1
}

// The corresponding interface
type IGroupObjectDescriptorRealisationType1 interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(io utils.WriteBuffer) error
	xml.Marshaler
}

func NewGroupObjectDescriptorRealisationType1(dataPointer uint8, transmitEnable bool, segmentSelectorEnable bool, writeEnable bool, readEnable bool, communicationEnable bool, priority CEMIPriority, valueType ComObjectValueType) *GroupObjectDescriptorRealisationType1 {
	return &GroupObjectDescriptorRealisationType1{DataPointer: dataPointer, TransmitEnable: transmitEnable, SegmentSelectorEnable: segmentSelectorEnable, WriteEnable: writeEnable, ReadEnable: readEnable, CommunicationEnable: communicationEnable, Priority: priority, ValueType: valueType}
}

func CastGroupObjectDescriptorRealisationType1(structType interface{}) *GroupObjectDescriptorRealisationType1 {
	castFunc := func(typ interface{}) *GroupObjectDescriptorRealisationType1 {
		if casted, ok := typ.(GroupObjectDescriptorRealisationType1); ok {
			return &casted
		}
		if casted, ok := typ.(*GroupObjectDescriptorRealisationType1); ok {
			return casted
		}
		return nil
	}
	return castFunc(structType)
}

func (m *GroupObjectDescriptorRealisationType1) GetTypeName() string {
	return "GroupObjectDescriptorRealisationType1"
}

func (m *GroupObjectDescriptorRealisationType1) LengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Simple field (dataPointer)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 1

	// Simple field (transmitEnable)
	lengthInBits += 1

	// Simple field (segmentSelectorEnable)
	lengthInBits += 1

	// Simple field (writeEnable)
	lengthInBits += 1

	// Simple field (readEnable)
	lengthInBits += 1

	// Simple field (communicationEnable)
	lengthInBits += 1

	// Simple field (priority)
	lengthInBits += 2

	// Simple field (valueType)
	lengthInBits += 8

	return lengthInBits
}

func (m *GroupObjectDescriptorRealisationType1) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func GroupObjectDescriptorRealisationType1Parse(io *utils.ReadBuffer) (*GroupObjectDescriptorRealisationType1, error) {

	// Simple Field (dataPointer)
	dataPointer, _dataPointerErr := io.ReadUint8(8)
	if _dataPointerErr != nil {
		return nil, errors.New("Error parsing 'dataPointer' field " + _dataPointerErr.Error())
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(1)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint8(0x1) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x1),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (transmitEnable)
	transmitEnable, _transmitEnableErr := io.ReadBit()
	if _transmitEnableErr != nil {
		return nil, errors.New("Error parsing 'transmitEnable' field " + _transmitEnableErr.Error())
	}

	// Simple Field (segmentSelectorEnable)
	segmentSelectorEnable, _segmentSelectorEnableErr := io.ReadBit()
	if _segmentSelectorEnableErr != nil {
		return nil, errors.New("Error parsing 'segmentSelectorEnable' field " + _segmentSelectorEnableErr.Error())
	}

	// Simple Field (writeEnable)
	writeEnable, _writeEnableErr := io.ReadBit()
	if _writeEnableErr != nil {
		return nil, errors.New("Error parsing 'writeEnable' field " + _writeEnableErr.Error())
	}

	// Simple Field (readEnable)
	readEnable, _readEnableErr := io.ReadBit()
	if _readEnableErr != nil {
		return nil, errors.New("Error parsing 'readEnable' field " + _readEnableErr.Error())
	}

	// Simple Field (communicationEnable)
	communicationEnable, _communicationEnableErr := io.ReadBit()
	if _communicationEnableErr != nil {
		return nil, errors.New("Error parsing 'communicationEnable' field " + _communicationEnableErr.Error())
	}

	// Simple Field (priority)
	priority, _priorityErr := CEMIPriorityParse(io)
	if _priorityErr != nil {
		return nil, errors.New("Error parsing 'priority' field " + _priorityErr.Error())
	}

	// Simple Field (valueType)
	valueType, _valueTypeErr := ComObjectValueTypeParse(io)
	if _valueTypeErr != nil {
		return nil, errors.New("Error parsing 'valueType' field " + _valueTypeErr.Error())
	}

	// Create the instance
	return NewGroupObjectDescriptorRealisationType1(dataPointer, transmitEnable, segmentSelectorEnable, writeEnable, readEnable, communicationEnable, priority, valueType), nil
}

func (m *GroupObjectDescriptorRealisationType1) Serialize(io utils.WriteBuffer) error {

	// Simple Field (dataPointer)
	dataPointer := uint8(m.DataPointer)
	_dataPointerErr := io.WriteUint8(8, (dataPointer))
	if _dataPointerErr != nil {
		return errors.New("Error serializing 'dataPointer' field " + _dataPointerErr.Error())
	}

	// Reserved Field (reserved)
	{
		_err := io.WriteUint8(1, uint8(0x1))
		if _err != nil {
			return errors.New("Error serializing 'reserved' field " + _err.Error())
		}
	}

	// Simple Field (transmitEnable)
	transmitEnable := bool(m.TransmitEnable)
	_transmitEnableErr := io.WriteBit((transmitEnable))
	if _transmitEnableErr != nil {
		return errors.New("Error serializing 'transmitEnable' field " + _transmitEnableErr.Error())
	}

	// Simple Field (segmentSelectorEnable)
	segmentSelectorEnable := bool(m.SegmentSelectorEnable)
	_segmentSelectorEnableErr := io.WriteBit((segmentSelectorEnable))
	if _segmentSelectorEnableErr != nil {
		return errors.New("Error serializing 'segmentSelectorEnable' field " + _segmentSelectorEnableErr.Error())
	}

	// Simple Field (writeEnable)
	writeEnable := bool(m.WriteEnable)
	_writeEnableErr := io.WriteBit((writeEnable))
	if _writeEnableErr != nil {
		return errors.New("Error serializing 'writeEnable' field " + _writeEnableErr.Error())
	}

	// Simple Field (readEnable)
	readEnable := bool(m.ReadEnable)
	_readEnableErr := io.WriteBit((readEnable))
	if _readEnableErr != nil {
		return errors.New("Error serializing 'readEnable' field " + _readEnableErr.Error())
	}

	// Simple Field (communicationEnable)
	communicationEnable := bool(m.CommunicationEnable)
	_communicationEnableErr := io.WriteBit((communicationEnable))
	if _communicationEnableErr != nil {
		return errors.New("Error serializing 'communicationEnable' field " + _communicationEnableErr.Error())
	}

	// Simple Field (priority)
	_priorityErr := m.Priority.Serialize(io)
	if _priorityErr != nil {
		return errors.New("Error serializing 'priority' field " + _priorityErr.Error())
	}

	// Simple Field (valueType)
	_valueTypeErr := m.ValueType.Serialize(io)
	if _valueTypeErr != nil {
		return errors.New("Error serializing 'valueType' field " + _valueTypeErr.Error())
	}

	return nil
}

func (m *GroupObjectDescriptorRealisationType1) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
			case "dataPointer":
				var data uint8
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.DataPointer = data
			case "transmitEnable":
				var data bool
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.TransmitEnable = data
			case "segmentSelectorEnable":
				var data bool
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.SegmentSelectorEnable = data
			case "writeEnable":
				var data bool
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.WriteEnable = data
			case "readEnable":
				var data bool
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ReadEnable = data
			case "communicationEnable":
				var data bool
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.CommunicationEnable = data
			case "priority":
				var data CEMIPriority
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.Priority = data
			case "valueType":
				var data ComObjectValueType
				if err := d.DecodeElement(&data, &tok); err != nil {
					return err
				}
				m.ValueType = data
			}
		}
	}
}

func (m *GroupObjectDescriptorRealisationType1) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	className := "org.apache.plc4x.java.knxnetip.readwrite.GroupObjectDescriptorRealisationType1"
	if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
		{Name: xml.Name{Local: "className"}, Value: className},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.DataPointer, xml.StartElement{Name: xml.Name{Local: "dataPointer"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.TransmitEnable, xml.StartElement{Name: xml.Name{Local: "transmitEnable"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.SegmentSelectorEnable, xml.StartElement{Name: xml.Name{Local: "segmentSelectorEnable"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.WriteEnable, xml.StartElement{Name: xml.Name{Local: "writeEnable"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ReadEnable, xml.StartElement{Name: xml.Name{Local: "readEnable"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.CommunicationEnable, xml.StartElement{Name: xml.Name{Local: "communicationEnable"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Priority, xml.StartElement{Name: xml.Name{Local: "priority"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ValueType, xml.StartElement{Name: xml.Name{Local: "valueType"}}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}
