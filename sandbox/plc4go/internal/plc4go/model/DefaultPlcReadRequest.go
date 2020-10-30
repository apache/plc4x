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
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/model"
)

type DefaultPlcReadRequestBuilder struct {
	reader       spi.PlcReader
	fieldHandler spi.PlcFieldHandler
	queries      map[string]string
	model.PlcReadRequestBuilder
}

func NewDefaultPlcReadRequestBuilder(fieldHandler spi.PlcFieldHandler, reader spi.PlcReader) *DefaultPlcReadRequestBuilder {
	return &DefaultPlcReadRequestBuilder{
		reader:       reader,
		fieldHandler: fieldHandler,
		queries:      map[string]string{},
	}
}

func (m *DefaultPlcReadRequestBuilder) AddItem(name string, query string) {
	m.queries[name] = query
}

func (m *DefaultPlcReadRequestBuilder) Build() (model.PlcReadRequest, error) {
	fields := make(map[string]model.PlcField)
	for name := range m.queries {
		query := m.queries[name]
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.New("Error parsing query: " + query + ". Got error: " + err.Error())
		}
		fields[name] = field
	}
	return DefaultPlcReadRequest{
		fields: fields,
		reader: m.reader,
	}, nil
}

type DefaultPlcReadRequest struct {
	fields map[string]model.PlcField
	reader spi.PlcReader
	model.PlcReadRequest
}

func (m DefaultPlcReadRequest) Execute() <-chan model.PlcReadRequestResult {
	return m.reader.Read(m)
}

func (m DefaultPlcReadRequest) GetFieldNames() []string {
	fieldNames := []string{}
	for name := range m.fields {
		fieldNames = append(fieldNames, name)
	}
	return fieldNames
}

func (m DefaultPlcReadRequest) GetField(name string) model.PlcField {
	if field, ok := m.fields[name]; ok {
		return field
	}
	return nil
}

func (m DefaultPlcReadRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcReadRequest"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}
	for fieldName, field := range m.fields {
		if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
		if err := e.EncodeElement(field, xml.StartElement{Name: xml.Name{Local: "field"}}); err != nil {
			return err
		}
		if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcReadRequest"}}); err != nil {
		return err
	}
	return nil
}
