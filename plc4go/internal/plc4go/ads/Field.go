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

package ads

import (
	"encoding/xml"
	"fmt"
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type PlcField struct {
	FieldType        FieldType
	IndexGroup       uint32
	IndexOffset      uint32
	SymbolicAddress  string
	StringLength     int32
	NumberOfElements int64
	Datatype         model2.AdsDataType
}

func (m PlcField) GetAddressString() string {
	return fmt.Sprintf("%dx%05d%05d%s%05d%05d:%s", m.FieldType, m.IndexGroup, m.IndexOffset, m.SymbolicAddress, m.StringLength, m.NumberOfElements, m.Datatype.String())
}

func (m PlcField) GetTypeName() string {
	return m.FieldType.GetName()
}

func (m PlcField) GetQuantity() uint16 {
	return 1
}

func NewAdsPlcField(fieldType FieldType, indexGroup uint32, indexOffset uint32, adsDataType model2.AdsDataType, stringLength int32, numberOfElements int64) (model.PlcField, error) {
	return PlcField{
		FieldType:        fieldType,
		IndexGroup:       indexGroup,
		IndexOffset:      indexOffset,
		SymbolicAddress:  "",
		StringLength:     stringLength,
		NumberOfElements: numberOfElements,
		Datatype:         adsDataType,
	}, nil
}

func NewAdsSymbolicPlcField(fieldType FieldType, symbolicAddress string, adsDataType model2.AdsDataType, stringLength int32, numberOfElements int64) (model.PlcField, error) {
	return PlcField{
		FieldType:        fieldType,
		IndexGroup:       0,
		IndexOffset:      0,
		SymbolicAddress:  symbolicAddress,
		StringLength:     stringLength,
		NumberOfElements: numberOfElements,
		Datatype:         adsDataType,
	}, nil
}

func CastToAdsFieldFromPlcField(plcField model.PlcField) (PlcField, error) {
	if adsField, ok := plcField.(PlcField); ok {
		return adsField, nil
	}
	return PlcField{}, errors.New("couldn't cast to AdsPlcField")
}

func (m PlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.IndexGroup, xml.StartElement{Name: xml.Name{Local: "indexGroup"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.IndexOffset, xml.StartElement{Name: xml.Name{Local: "indexOffset"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.SymbolicAddress, xml.StartElement{Name: xml.Name{Local: "symbolicAddress"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.StringLength, xml.StartElement{Name: xml.Name{Local: "stringLength"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.NumberOfElements, xml.StartElement{Name: xml.Name{Local: "numberOfElements"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Datatype.String(), xml.StartElement{Name: xml.Name{Local: "dataType"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}
