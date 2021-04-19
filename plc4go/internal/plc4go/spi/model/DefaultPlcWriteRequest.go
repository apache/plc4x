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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/pkg/errors"
)

type DefaultPlcWriteRequestBuilder struct {
	writer       spi.PlcWriter
	fieldHandler spi.PlcFieldHandler
	valueHandler spi.PlcValueHandler
	queries      map[string]string
	queryNames   []string
	fields       map[string]model.PlcField
	fieldNames   []string
	values       map[string]interface{}
}

func NewDefaultPlcWriteRequestBuilder(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, writer spi.PlcWriter) *DefaultPlcWriteRequestBuilder {
	return &DefaultPlcWriteRequestBuilder{
		writer:       writer,
		fieldHandler: fieldHandler,
		valueHandler: valueHandler,
		queries:      map[string]string{},
		queryNames:   make([]string, 0),
		fields:       map[string]model.PlcField{},
		fieldNames:   make([]string, 0),
		values:       map[string]interface{}{},
	}
}

func (m *DefaultPlcWriteRequestBuilder) AddQuery(name string, query string, value interface{}) model.PlcWriteRequestBuilder {
	m.queryNames = append(m.queryNames, name)
	m.queries[name] = query
	m.values[name] = value
	return m
}

func (m *DefaultPlcWriteRequestBuilder) AddField(name string, field model.PlcField, value interface{}) model.PlcWriteRequestBuilder {
	m.fieldNames = append(m.fieldNames, name)
	m.fields[name] = field
	m.values[name] = value
	return m
}

func (m *DefaultPlcWriteRequestBuilder) Build() (model.PlcWriteRequest, error) {
	// Parse the queries as well as pro
	for _, name := range m.queryNames {
		query := m.queries[name]
		field, err := m.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing query: %s", query)
		}
		m.AddField(name, field, m.values[name])
	}

	// Process the values for fields.
	plcValues := make(map[string]values.PlcValue)
	for name, field := range m.fields {
		value, err := m.valueHandler.NewPlcValue(field, m.values[name])
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing value of type: %s", field.GetTypeName())
		}
		plcValues[name] = value
	}
	return NewDefaultPlcWriteRequest(m.fields, m.fieldNames, plcValues, m.writer), nil
}

type DefaultPlcWriteRequest struct {
	DefaultRequest
	values map[string]values.PlcValue
	writer spi.PlcWriter
}

func NewDefaultPlcWriteRequest(fields map[string]model.PlcField, fieldNames []string, values map[string]values.PlcValue, writer spi.PlcWriter) model.PlcWriteRequest {
	return DefaultPlcWriteRequest{NewDefaultRequest(fields, fieldNames), values, writer}
}

func (m DefaultPlcWriteRequest) Execute() <-chan model.PlcWriteRequestResult {
	return m.writer.Write(m)
}

func (m DefaultPlcWriteRequest) GetValue(name string) values.PlcValue {
	return m.values[name]
}

func (m DefaultPlcWriteRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcWriteRequest"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}
	for _, fieldName := range m.fieldNames {
		field := m.fields[fieldName]
		value := m.values[fieldName]
		if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
		if err := e.EncodeElement(field, xml.StartElement{Name: xml.Name{Local: "field"}}); err != nil {
			return err
		}
		switch value.(type) {
		case values2.PlcList:
			listValue, ok := value.(values2.PlcList)
			if !ok {
				return errors.New("couldn't cast PlcValue to PlcList")
			}
			for _, subValue := range listValue.Values {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
					return err
				}
				if !subValue.IsString() {
					return errors.New("value not serializable to string")
				}
				if err := e.EncodeToken(xml.CharData(subValue.GetString())); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "value"}}); err != nil {
					return err
				}
			}
		default:
			if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
				return err
			}
			if !value.IsString() {
				return errors.New("value not serializable to string")
			}
			if err := e.EncodeToken(xml.CharData(value.GetString())); err != nil {
				return err
			}
			if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "value"}}); err != nil {
				return err
			}
		}
		if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}}); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "fields"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcWriteRequest"}}); err != nil {
		return err
	}
	return nil
}
