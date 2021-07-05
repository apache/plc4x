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
	"strconv"
)

type PlcField struct {
	FieldType        FieldType
	StringLength     int32
	NumberOfElements uint32
	Datatype         model2.AdsDataType
}

func (m PlcField) GetTypeName() string {
	return m.FieldType.GetName()
}

func (m PlcField) GetQuantity() uint16 {
	return uint16(m.NumberOfElements)
}

func (m PlcField) GetNumberOfElements() uint32 {
	return m.NumberOfElements
}

func (m PlcField) GetDatatype() model2.AdsDataType {
	return m.Datatype
}

func (m PlcField) GetStringLength() int32 {
	return m.StringLength
}

func (m PlcField) GetAddressString() string {
	return fmt.Sprintf("%dx%05d%05d:%s", m.FieldType, m.StringLength, m.NumberOfElements, m.Datatype.String())
}

type AdsPlcField interface {
	GetDatatype() model2.AdsDataType
	GetStringLength() int32
	GetNumberOfElements() uint32
	model.PlcField
}

func castToAdsFieldFromPlcField(plcField model.PlcField) (AdsPlcField, error) {
	if adsField, ok := plcField.(AdsPlcField); ok {
		return adsField, nil
	}
	return nil, errors.Errorf("couldn't %T cast to AdsPlcField", plcField)
}

type DirectPlcField struct {
	IndexGroup  uint32
	IndexOffset uint32
	PlcField
}

func (m DirectPlcField) GetAddressString() string {
	return fmt.Sprintf("%dx%05d%05d%05d%05d:%s", m.FieldType, m.IndexGroup, m.IndexOffset, m.StringLength, m.NumberOfElements, m.Datatype.String())
}

func newDirectAdsPlcField(indexGroup uint32, indexOffset uint32, adsDataType model2.AdsDataType, stringLength int32, numberOfElements uint32) (model.PlcField, error) {
	fieldType := DirectAdsField
	if stringLength > 0 {
		fieldType = DirectAdsStringField
	}
	return DirectPlcField{
		IndexGroup:  indexGroup,
		IndexOffset: indexOffset,
		PlcField: PlcField{
			FieldType:        fieldType,
			StringLength:     stringLength,
			NumberOfElements: numberOfElements,
			Datatype:         adsDataType,
		},
	}, nil
}

func castToDirectAdsFieldFromPlcField(plcField model.PlcField) (DirectPlcField, error) {
	if adsField, ok := plcField.(DirectPlcField); ok {
		return adsField, nil
	}
	return DirectPlcField{}, errors.Errorf("couldn't %T cast to DirectPlcField", plcField)
}

func (m DirectPlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.IndexGroup, xml.StartElement{Name: xml.Name{Local: "indexGroup"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.IndexOffset, xml.StartElement{Name: xml.Name{Local: "indexOffset"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if m.StringLength > 0 {
		if err := e.EncodeElement(m.StringLength, xml.StartElement{Name: xml.Name{Local: "stringLength"}, Attr: []xml.Attr{
			{Name: xml.Name{Local: "dataType"}, Value: "int"},
			{Name: xml.Name{Local: "bitLength"}, Value: "64"},
		}}); err != nil {
			return err
		}
	}
	if err := e.EncodeElement(m.NumberOfElements, xml.StartElement{Name: xml.Name{Local: "numberOfElements"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Datatype.String(), xml.StartElement{Name: xml.Name{Local: "dataType"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "string"},
		{Name: xml.Name{Local: "bitLength"}, Value: strconv.Itoa(len(m.Datatype.String()) * 8)},
	}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}

func (m DirectPlcField) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	panic(name)
}

type SymbolicPlcField struct {
	SymbolicAddress string
	PlcField
}

func (m SymbolicPlcField) GetAddressString() string {
	return fmt.Sprintf("%dx%s%05d%05d:%s", m.FieldType, m.SymbolicAddress, m.StringLength, m.NumberOfElements, m.Datatype.String())
}

func newAdsSymbolicPlcField(symbolicAddress string, adsDataType model2.AdsDataType, stringLength int32, numberOfElements uint32) (model.PlcField, error) {
	fieldType := SymbolicAdsField
	if stringLength > 0 {
		fieldType = SymbolicAdsStringField
	}
	return SymbolicPlcField{
		SymbolicAddress: symbolicAddress,
		PlcField: PlcField{
			FieldType:        fieldType,
			StringLength:     stringLength,
			NumberOfElements: numberOfElements,
			Datatype:         adsDataType,
		},
	}, nil
}

func needsResolving(plcField model.PlcField) bool {
	switch plcField.(type) {
	case SymbolicPlcField:
		return true
	case DirectPlcField:
		return false
	default:
		return false
	}
}

func castToSymbolicPlcFieldFromPlcField(plcField model.PlcField) (SymbolicPlcField, error) {
	if adsField, ok := plcField.(SymbolicPlcField); ok {
		return adsField, nil
	}
	return SymbolicPlcField{}, errors.Errorf("couldn't cast %T to SymbolicPlcField", plcField)
}

func (m SymbolicPlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.SymbolicAddress, xml.StartElement{Name: xml.Name{Local: "symbolicAddress"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "string"},
		{Name: xml.Name{Local: "bitLength"}, Value: strconv.Itoa(len(m.SymbolicAddress) * 8)},
	}}); err != nil {
		return err
	}
	if m.StringLength > 0 {
		if err := e.EncodeElement(m.StringLength, xml.StartElement{Name: xml.Name{Local: "stringLength"}}); err != nil {
			return err
		}
	}
	if err := e.EncodeElement(m.NumberOfElements, xml.StartElement{Name: xml.Name{Local: "numberOfElements"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Datatype.String(), xml.StartElement{Name: xml.Name{Local: "dataType"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "string"},
		{Name: xml.Name{Local: "bitLength"}, Value: strconv.Itoa(len(m.Datatype.String()) * 8)},
	}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}

func (m DirectPlcField) banana(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.IndexGroup, xml.StartElement{Name: xml.Name{Local: "indexGroup"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.IndexOffset, xml.StartElement{Name: xml.Name{Local: "indexOffset"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if m.StringLength > 0 {
		if err := e.EncodeElement(m.StringLength, xml.StartElement{Name: xml.Name{Local: "stringLength"}, Attr: []xml.Attr{
			{Name: xml.Name{Local: "dataType"}, Value: "int"},
			{Name: xml.Name{Local: "bitLength"}, Value: "64"},
		}}); err != nil {
			return err
		}
	}
	if err := e.EncodeElement(m.NumberOfElements, xml.StartElement{Name: xml.Name{Local: "numberOfElements"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "int"},
		{Name: xml.Name{Local: "bitLength"}, Value: "64"},
	}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Datatype.String(), xml.StartElement{Name: xml.Name{Local: "dataType"}, Attr: []xml.Attr{
		{Name: xml.Name{Local: "dataType"}, Value: "string"},
		{Name: xml.Name{Local: "bitLength"}, Value: strconv.Itoa(len(m.Datatype.String()) * 8)},
	}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}
