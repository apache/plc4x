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

package eip

import (
	"encoding/xml"
	readWrite "github.com/apache/plc4x/plc4go/internal/plc4go/eip/readwrite/model"
	"github.com/rs/zerolog/log"
)

type EIPPlcField interface {
	GetTag() string
	GetType() readWrite.CIPDataTypeCode
	GetElementNb() uint16
}

type PlcField struct {
	Tag       string
	Type      readWrite.CIPDataTypeCode
	ElementNb uint16
}

func (m PlcField) GetAddressString() string {
	return m.GetTag()
}

func (m PlcField) GetTypeName() string {
	return m.GetType().String()
}

func (m PlcField) GetQuantity() uint16 {
	return 1
}

func NewField(tag string, _type readWrite.CIPDataTypeCode, elementNb uint16) PlcField {
	return PlcField{
		Tag:       tag,
		Type:      _type,
		ElementNb: elementNb,
	}
}

func (m PlcField) GetTag() string {
	return m.Tag
}

func (m PlcField) GetType() readWrite.CIPDataTypeCode {
	return m.Type
}

func (m PlcField) GetElementNb() uint16 {
	return m.ElementNb
}

func (m PlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "EipField"}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.Tag, xml.StartElement{Name: xml.Name{Local: "node"}}); err != nil {
		return err
	}
	if m.Type != 0 {
		if err := e.EncodeElement(m.Type, xml.StartElement{Name: xml.Name{Local: "type"}}); err != nil {
			return err
		}
	}
	if err := e.EncodeElement(m.ElementNb, xml.StartElement{Name: xml.Name{Local: "elementNb"}}); err != nil {
		return err
	}
	// TODO: remove this from the spec
	if err := e.EncodeElement("java.lang.Object", xml.StartElement{Name: xml.Name{Local: "defaultJavaType"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "EipField"}}); err != nil {
		return err
	}
	return nil
}
